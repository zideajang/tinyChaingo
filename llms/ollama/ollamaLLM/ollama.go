package ollamaLLM

import (
	"context"
	"fmt"
	"sync"

	"github.com/zideajang/langChaingo/llms/ollama/internal/ollamaclient"
)

type OllamaLLM struct {
	// 包含一个客户端和模型名称
	client *ollamaclient.Client
	model  string
}

// Option 的切片
// New 函数用于创建返回一个 llm 的结构体指针
func New(opts ...Option) (*OllamaLLM, error) {

	// 初始化 llm 这里 llm 时 OllamaLLM* llm
	llm := &OllamaLLM{
		model: ollamaclient.DefaultChatModel,
	}

	// 循环 function(*llm) 然后在函数内部去更新 llm
	for _, opt := range opts {
		opt(llm)
	}

	// 初始化 client
	client, err := ollamaclient.New("")
	if err != nil {
		return nil, err
	}
	llm.client = client
	return llm, nil
}

type Option func(*OllamaLLM)

func WithModel(model string) Option {
	return func(llm *OllamaLLM) {
		llm.model = model
	}
}

// llm 上方法
func (l *OllamaLLM) Call(prompt string) (string, error) {
	ctx := context.Background()

	// 构建聊天请求
	req := &ollamaclient.ChatRequest{
		Model: l.model, //当前llma
		Messages: []ollamaclient.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	resp, err := l.client.Chat(ctx, req)
	if err != nil {
		return "", fmt.Errorf("ollama Chat failed: %w", err)
	}

	return resp.Content, nil
}

func (l *OllamaLLM) Generate(prompts []string) ([]string, error) {
	// 用于存储所有完成的文本
	completions := make([]string, len(prompts))
	// WaitGroup 用于等待所有 Goroutine 完成
	var wg sync.WaitGroup
	// 缓冲通道，用于收集并发 Goroutine 中的错误
	errs := make(chan error, len(prompts))

	for i, prompt := range prompts {
		wg.Add(1) // 增加 WaitGroup 计数器

		// 启动一个新的 Goroutine
		go func(i int, p string) {
			defer wg.Done() // Goroutine 完成时，减少 WaitGroup 计数器
			ctx := context.Background()

			req := &ollamaclient.ChatRequest{
				Model: l.model,
				Messages: []ollamaclient.Message{
					{
						Role:    "user",
						Content: p,
					},
				},
				Stream: false,
			}

			resp, err := l.client.Chat(ctx, req)
			if err != nil {
				errs <- fmt.Errorf("ollama Generate for prompt %d failed: %w", i, err)
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
		return nil, fmt.Errorf("multiple errors during ollama Generate: %v", allErrors)
	}

	return completions, nil
}
