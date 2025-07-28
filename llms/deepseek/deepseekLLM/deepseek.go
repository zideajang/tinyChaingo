package deepseekLLM

import (
	"context"
	"fmt"
	"sync"

	"github.com/zideajang/langChaingo/llms/deepseek/internal/deepseekclient" // Import the new deepseekclient
)

// DeepSeekLLM 结构体封装了 DeepSeek 客户端和模型配置。
type DeepSeekLLM struct {
	client *deepseekclient.Client
	model  string //模型名称，提供选择的是 deepseek-chat /deepseek-reason
}

// Option 类型定义了用于配置 DeepSeekLLM 实例的函数选项。
type Option func(*DeepSeekLLM)

// New 函数是 DeepSeekLLM 的构造函数，用于创建和初始化一个 DeepSeekLLM 实例。
// 它支持使用函数选项模式进行配置，并会初始化内部的 DeepSeek 客户端。
func New(opts ...Option) (*DeepSeekLLM, error) {
	// llm 的指针
	llm := &DeepSeekLLM{
		model: deepseekclient.DefaultChatModel, // 设置默认模型
	}

	// 应用所有传入的选项来更新 LLM 配置
	for _, opt := range opts {
		opt(llm)
	}

	// 初始化 DeepSeek API 客户端
	// 注意：deepseekclient.New() 现在从配置文件中读取API密钥，所以不需要参数
	client, err := deepseekclient.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create DeepSeek client: %w", err)
	}
	llm.client = client
	return llm, nil
}

// WithModel 是一个选项函数，用于设置 DeepSeekLLM 实例的模型名称。
func WithModel(model string) Option {
	return func(llm *DeepSeekLLM) {
		llm.model = model
	}
}

// Call 方法实现了 llms.LLM 接口的 Call 方法，用于向 DeepSeek 模型发送单个提示。
func (l *DeepSeekLLM) Call(prompt string) (string, error) {
	ctx := context.Background()

	// 构建 DeepSeek 聊天请求
	req := &deepseekclient.ChatRequest{
		Model: l.model,
		Messages: []deepseekclient.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false, // 强制为非流式响应
	}

	// 调用内部 DeepSeek 客户端的 Chat 方法
	resp, err := l.client.Chat(ctx, req)
	if err != nil {
		return "", fmt.Errorf("DeepSeek Chat failed: %w", err)
	}

	return resp.Content, nil
}

// Generate 方法实现了 llms.LLM 接口的 Generate 方法，用于向 DeepSeek 模型批量发送提示。
// 它通过并发 Goroutine 来提高效率。
func (l *DeepSeekLLM) Generate(prompts []string) ([]string, error) {
	completions := make([]string, len(prompts)) // 存储所有完成的文本
	var wg sync.WaitGroup                       // 用于等待所有 Goroutine 完成
	errs := make(chan error, len(prompts))      // 缓冲通道，用于收集并发 Goroutine 中的错误

	for i, prompt := range prompts {
		wg.Add(1) // 增加 WaitGroup 计数器

		go func(i int, p string) {
			defer wg.Done() // Goroutine 完成时，减少 WaitGroup 计数器
			ctx := context.Background()

			req := &deepseekclient.ChatRequest{
				Model: l.model,
				Messages: []deepseekclient.Message{
					{
						Role:    "user",
						Content: p,
					},
				},
				Stream: false, // 强制为非流式响应
			}

			resp, err := l.client.Chat(ctx, req)
			if err != nil {
				errs <- fmt.Errorf("DeepSeek Generate for prompt %d failed: %w", i, err)
				return
			}
			completions[i] = resp.Content
		}(i, prompt)
	}

	wg.Wait()   // 等待所有 Goroutine 完成执行
	close(errs) // 关闭错误通道，表示不会再有错误发送

	var allErrors []error
	// 从错误通道中收集所有错误
	for err := range errs {
		allErrors = append(allErrors, err)
	}

	if len(allErrors) > 0 {
		return nil, fmt.Errorf("multiple errors during DeepSeek Generate: %v", allErrors)
	}

	return completions, nil
}
