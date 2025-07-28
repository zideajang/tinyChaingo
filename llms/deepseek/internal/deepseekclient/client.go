package deepseekclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os" // For reading config file

	// For path manipulation
	"gopkg.in/yaml.v3" // For parsing YAML config
)

// --- Constants ---
const (
	// DefaultChatModel 是DeepSeek客户端使用的默认模型。
	DefaultChatModel = "deepseek-chat"
	// chatAPIPath 是DeepSeek API中用于聊天补全的路由。
	chatAPIPath = "/chat/completions"
	// DefaultBaseURL 是DeepSeek服务的默认基础URL。
	DefaultBaseURL = "https://api.deepseek.com"
	// configFilePath 是DeepSeek API密钥的配置文件路径。
	configFilePath = "D:/config.yaml" // Adjust this path as needed
	// deepseekAPIKeyEnvVar 是在config.yaml中查找DeepSeek API Key的关键字。
	// deepseekAPIKeyConfigKey = "DEEPSEEK_API_KEY"
)

// --- Errors ---
// ErrEmptyResponse 表示DeepSeek模型返回的内容为空。
var ErrEmptyResponse = errors.New("empty response from DeepSeek")

// ErrAPIKeyNotFound 表示在配置文件中未找到DeepSeek API Key。
var ErrAPIKeyNotFound = errors.New("DeepSeek API Key not found in config file")

// --- Client Structure ---

// Client 表示与DeepSeek API交互的客户端。
type Client struct {
	apikey  string // DeepSeek API密钥
	baseURL string // DeepSeek服务的基准URL
}

// --- Config Structure for YAML Parsing ---
type Config struct {
	DeepSeekAPIKey string `yaml:"DEEPSEEK_API_KEY"`
}

// --- Client Constructor ---

// New 创建并返回一个新的DeepSeek Client实例。
// 它会从指定路径的config.yaml文件中读取API密钥。
func New() (*Client, error) {
	// Read API key from config file
	// 读取 yaml 文件
	configBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file at %s: %w", configFilePath, err)
	}

	var config Config
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	if config.DeepSeekAPIKey == "" {
		return nil, ErrAPIKeyNotFound
	}

	c := &Client{
		apikey:  config.DeepSeekAPIKey, //直接把 API 写道这里，
		baseURL: DefaultBaseURL,
	}
	return c, nil
}

// --- Request and Response Payloads (OpenAI-compatible) ---

// Message 结构体表示聊天中的一条消息。
type Message struct {
	Role    string `json:"role"`    // 消息发送者的角色 (e.g., "user", "system", "assistant")
	Content string `json:"content"` // 消息的文本内容
}

// ChatRequest 结构体定义了发送到DeepSeek API的聊天请求体。
// 它与OpenAI的API请求体兼容。
type ChatRequest struct {
	Model    string    `json:"model"`    // 要使用的DeepSeek模型名称
	Messages []Message `json:"messages"` // 聊天消息列表
	Stream   bool      `json:"stream"`   // 是否以流式方式获取响应 (false表示获取完整响应)
}

// ChatChoice 结构体表示聊天补全的一个选项。
type ChatChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 结构体表示本次API调用的token使用情况。
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// DeepSeekChatResponsePayload 结构体用于解析DeepSeek API返回的完整JSON响应。
// 它与OpenAI的API响应体兼容。
type DeepSeekChatResponsePayload struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Usage   Usage        `json:"usage"`
}

// ChatResponse 结构体是DeepSeek客户端向外部暴露的简化聊天响应。
// 它只包含最重要的信息：LLM生成的内容。
type ChatResponse struct {
	Content string // LLM生成的内容
}

// this bind(object)
// self.doChat()
// type contains
//
// --- Internal HTTP Request Method ---

// doChat 是一个内部方法，负责向DeepSeek API发送实际的HTTP聊天请求。
// ctx 用于管理请求的生命周期和超时。
// payload 包含要发送到DeepSeek的聊天请求数据。
func (c *Client) doChat(ctx context.Context, payload *ChatRequest) (*DeepSeekChatResponsePayload, error) {
	// 如果没有指定模型，则使用默认模型。
	if payload.Model == "" {
		payload.Model = DefaultChatModel
	}

	// 客户端期望DeepSeek API返回一个完整的、一次性的响应，而不是流式传输。
	payload.Stream = false

	// 将请求体转换为JSON字节数组。
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	// 将字节数组包装成io.Reader，以便http.NewRequestWithContext使用。
	body := bytes.NewReader(payloadBytes)

	// 构建完整的请求URL。
	url := c.baseURL + chatAPIPath

	// 创建新的HTTP POST请求。
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// 设置请求头，指定内容类型为JSON和授权信息。
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apikey) // Add API Key

	// 发送HTTP请求。
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer r.Body.Close() // 确保响应体在使用后关闭

	// 检查HTTP响应状态码。
	if r.StatusCode != http.StatusOK {
		respBodyBytes, _ := io.ReadAll(r.Body) // Read body for detailed error
		return nil, fmt.Errorf("DeepSeek API request failed with status %d: %s", r.StatusCode, string(respBodyBytes))
	}

	// 声明一个变量来存储解析后的DeepSeek API响应。
	var response DeepSeekChatResponsePayload
	// 将HTTP响应体中的JSON数据解码到结构体中。
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DeepSeek API response: %w", err)
	}
	return &response, nil
}

// --- Public Chat Method ---

// Chat 方法是DeepSeek客户端的公共入口点，用于发送聊天请求。
// 它调用内部的doChat方法并处理返回的响应。
func (c *Client) Chat(ctx context.Context, r *ChatRequest) (*ChatResponse, error) {
	resp, err := c.doChat(ctx, r)
	if err != nil {
		return nil, err
	}
	// 检查DeepSeek的响应是否包含有效的消息内容。
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return nil, ErrEmptyResponse
	}
	// 返回一个简化的ChatResponse。
	return &ChatResponse{
		Content: resp.Choices[0].Message.Content,
	}, nil
}
