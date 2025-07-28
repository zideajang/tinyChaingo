package ollamaclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// --- Constants ---
const (
	// DefaultChatModel 是Ollama客户端使用的默认模型。
	DefaultChatModel = "qwen3:8b"
	// chatAPIPath 是Ollama API中用于聊天补全的路由。
	chatAPIPath = "/api/chat"
	// DefaultBaseURL 是Ollama服务的默认基础URL。
	DefaultBaseURL = "http://localhost:11434"
)

// --- Errors ---
// ErrEmptyResponse 表示Ollama模型返回的内容为空。
var ErrEmptyResponse = errors.New("empty response")

// --- Client Structure ---

// Client 表示与Ollama API交互的客户端。
type Client struct {
	// apikey 存储API密钥，对于Ollama通常为空或不适用。
	apikey string
	// baseURL 存储Ollama服务的基准URL。
	baseURL string
}

// --- Client Constructor ---

// New 创建并返回一个新的Ollama Client实例。
// apikey 参数目前对Ollama服务通常不使用，但保留以备将来兼容性。
func New(apikey string) (*Client, error) {
	c := &Client{
		apikey:  apikey,
		baseURL: DefaultBaseURL, // 使用常量设置默认URL
	}
	return c, nil
}

// --- Request and Response Payloads ---

// Message 结构体表示聊天中的一条消息。
type Message struct {
	Role    string `json:"role"`    // 消息发送者的角色 (e.g., "user", "system", "assistant")
	Content string `json:"content"` // 消息的文本内容
}

// ChatRequest 结构体定义了发送到Ollama API的聊天请求体。
type ChatRequest struct {
	Model    string    `json:"model"`    // 要使用的Ollama模型名称
	Messages []Message `json:"messages"` // 聊天消息列表
	Stream   bool      `json:"stream"`   // 是否以流式方式获取响应 (false表示获取完整响应)
}

// ollamaChatResponsePayload 结构体用于解析Ollama API返回的完整JSON响应。
// 它包含了模型生成的消息以及各种性能指标。
type ollamaChatResponsePayload struct {
	Model         string  `json:"model"`
	CreatedAt     string  `json:"created_at"`
	Message       Message `json:"message"` // LLM 生成的回复消息
	Done          bool    `json:"done"`
	TotalDuration int64   `json:"total_duration"`
	LoadDuration  int64   `json:"load_duration"`
	PromptEvalCount int `json:"prompt_eval_count"`
	PromptEvalDuration int64 `json:"prompt_eval_duration"`
	EvalCount     int     `json:"eval_count"`
	EvalDuration  int64   `json:"eval_duration"`
}

// ChatResponse 结构体是Ollama客户端向外部暴露的简化聊天响应。
// 它只包含最重要的信息：LLM生成的内容。
type ChatResponse struct {
	Content string // LLM生成的内容
}

// --- Internal HTTP Request Method ---

// doChat 是一个内部方法，负责向Ollama API发送实际的HTTP聊天请求。
// ctx 用于管理请求的生命周期和超时。
// payload 包含要发送到Ollama的聊天请求数据。
func (c *Client) doChat(ctx context.Context, payload *ChatRequest) (*ollamaChatResponsePayload, error) {
	// 如果没有指定模型，则使用默认模型。
	if payload.Model == "" {
		payload.Model = DefaultChatModel
	}

	// 客户端期望Ollama API返回一个完整的、一次性的响应，而不是流式传输。
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
	
	// 设置请求头，指定内容类型为JSON。
	req.Header.Set("Content-Type", "application/json")

	// 发送HTTP请求。
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer r.Body.Close() // 确保响应体在使用后关闭

	// 检查HTTP响应状态码。
	if r.StatusCode != http.StatusOK {
		var errRes map[string]interface{}
		// 尝试解析错误响应体。
		json.NewDecoder(r.Body).Decode(&errRes)
		return nil, fmt.Errorf("ollama API request failed with status %d: %v", r.StatusCode, errRes)
	}

	// 声明一个变量来存储解析后的Ollama API响应。
	var response ollamaChatResponsePayload
	// 将HTTP响应体中的JSON数据解码到结构体中。
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ollama API response: %w", err)
	}
	return &response, nil
}

// --- Public Chat Method ---

// Chat 方法是Ollama客户端的公共入口点，用于发送聊天请求。
// 它调用内部的doChat方法并处理返回的响应。
func (c *Client) Chat(ctx context.Context, r *ChatRequest) (*ChatResponse, error) {
	resp, err := c.doChat(ctx, r)
	if err != nil {
		return nil, err
	}
	// 检查Ollama的响应消息内容是否为空。
	if resp.Message.Content == "" {
		return nil, ErrEmptyResponse
	}
	// 返回一个简化的ChatResponse。
	return &ChatResponse{
		Content: resp.Message.Content,
	}, nil
}