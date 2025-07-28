package main

import (
	"fmt"
	"log"

	"github.com/zideajang/langChaingo/llms"
	"github.com/zideajang/langChaingo/llms/deepseek/deepseekLLM"
	"github.com/zideajang/langChaingo/llms/ollama/ollamaLLM"
)

func main() {
	// --- Test Call method ---
	// fmt.Println("--- Testing Chat method ---")
	// testCallMethod() //单独一条提示词的发起询问

	// fmt.Println("\n--- Testing Generate method ---")
	// testGenerateMethod() //批量提示词发起询问

	// --- Test DeepSeek Call 调用方法 ---
	fmt.Println("\n--- Testing DeepSeek Call method(单次请求大语言模型) ---")
	testDeepSeekCallMethod()

	fmt.Println("\n--- Testing DeepSeek Generate method(批量地请求大语言模型) ---")
	testDeepSeekGenerateMethod()
}

func testCallMethod() {
	// 初始化一个 ollamaLLM 模型
	llm, err := ollamaLLM.New(ollamaLLM.WithModel("qwen3:8b"))
	if err != nil {
		log.Fatalf("创建 Ollama LLM 时发生错误: %v", err)
	}

	var genericLLM llms.LLM = llm

	prompt := "天空为什么是蓝色"
	fmt.Printf("Calling LLM with prompt: \"%s\"\n", prompt)
	completion, err := genericLLM.Call(prompt) // Using the interface method
	if err != nil {
		log.Fatalf("Failed to get completion: %v", err)
	}
	fmt.Printf("Completion: %s\n", completion)
}

func testGenerateMethod() {
	llm, err := ollamaLLM.New(ollamaLLM.WithModel("qwen3:8b"))
	if err != nil {
		log.Fatalf("Failed to create Ollama LLM: %v", err)
	}

	var genericLLM llms.LLM = llm

	prompts := []string{
		"辽宁的省会是名称",
		" 2+3 等于几",
	}
	fmt.Printf("对于每一个 prompt 生成结果: %v\n", prompts)
	completions, err := genericLLM.Generate(prompts) // Using the interface method
	if err != nil {
		log.Fatalf("Failed to generate completions: %v", err)
	}

	for i, comp := range completions {
		fmt.Printf("Prompt %d: \"%s\"\n", i+1, prompts[i])
		fmt.Printf("Completion %d: %s\n", i+1, comp)
	}
}

// New functions for DeepSeek
func testDeepSeekCallMethod() {
	// 模型用 deepseek-chat
	// 调用方法就是实例化模型`deepseekLLM.New`
	// 通过 deepseekLLM.WithModel("<填写模型名称>") deepseek-chat/deepseek-reason
	// 具体模型调用参见 deepseek 官网提供了哪些模型
	llm, err := deepseekLLM.New(deepseekLLM.WithModel("deepseek-chat")) // Use deepseek-chat model
	//
	if err != nil {
		log.Fatalf("创建 DeepSeek LLM 时发生错误: %v", err)
	}

	// 向上转型，所有 deepseekLLM 也实现 (LLM 的接口)，所以可以向上转型
	var genericLLM llms.LLM = llm
	// prompt
	prompt := "用中文告诉我，天空为什么是蓝的？"
	fmt.Printf("调用 DeepSeek LLM 模型提供: \"%s\"\n", prompt)
	completion, err := genericLLM.Call(prompt)
	if err != nil {
		log.Fatalf("Failed to get DeepSeek completion: %v", err)
	}
	fmt.Printf("DeepSeek Completion: %s\n", completion)
}

func testDeepSeekGenerateMethod() {
	llm, err := deepseekLLM.New(deepseekLLM.WithModel("deepseek-chat"))
	if err != nil {
		log.Fatalf("Failed to create DeepSeek LLM: %v", err)
	}

	var genericLLM llms.LLM = llm

	prompts := []string{
		"请用中文介绍一下Go语言的特点。",
		"请问猫科动物有哪些？",
	}
	fmt.Printf("对于每一个 DeepSeek prompt 生成结果: %v\n", prompts)
	completions, err := genericLLM.Generate(prompts)
	if err != nil {
		log.Fatalf("Failed to generate DeepSeek completions: %v", err)
	}

	for i, comp := range completions {
		fmt.Printf("DeepSeek Prompt %d: \"%s\"\n", i+1, prompts[i])
		fmt.Printf("DeepSeek Completion %d: %s\n", i+1, comp)
	}
}
