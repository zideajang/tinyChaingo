# tinyChain

## LLM 供应商支持
- 支持 ollama LLM 平台上提供模型
- 支持 deepseek 系列模型



## 关于如何调用 deepSeek 系列模型
首先需要做一些准备工作，也就是去官网申请 deepseek API key，然后可以保存到一个地方。便于在实例化客户端(client)提供 API key

### 客户端

可以将申请的 API key 保存一个指定位置的 yaml 文件中
然后在 `llms/deepseek/internal/deepseek/client.go`

### 两种方式去加载 API key
1. 按照 tinyChaingo 提供通过读取保存 API key yaml 文件来获取 API key

```go
configFilePath = "<路径>/config.yaml" 
```
将官网申请的 deepseek API 放到 config.yaml 文件中 
```yaml
DEEPSEEK_API_KEY: <deepske apikey>
```

```go
type Config struct {
	DeepSeekAPIKey string `yaml:"DEEPSEEK_API_KEY"`
}
```
可以通过修改 `DEEPSEEK_API_KEY` 这个值来自定义的键值用于在 `config.yaml` 文件中放置的 deepseek API

1. 直接加载 API key 

```go
c := &Client{
		apikey:  config.DeepSeekAPIKey, //直接把申请的 API key 写
		baseURL: DefaultBaseURL,
	}
```
也就是实例化 `Client` 时，直接替换调用 `config.DeepSeekAPIKey` 即可


```go

// 实例化的 deepseekLLM 模型 
llm, err := deepseekLLM.New(deepseekLLM.WithModel("deepseek-chat"))

if err != nil {
    log.Fatalf("Failed to create DeepSeek LLM: %v", err)
}

var genericLLM llms.LLM = llm

prompt := "用中文告诉我，天空为什么是蓝的？"

completion, err := genericLLM.Call(prompt)
if err != nil {
    log.Fatalf("Failed to get DeepSeek completion: %v", err)
}
fmt.Printf("DeepSeek Completion: %s\n", completion)

```
`deepseekLLM` 提供 `New` 方法，然后`deepseekLLM.WithModel(<模型名称>)` 只支持 `deepseek-chat` 