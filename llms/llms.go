package llms

type LLM interface {
	Call(prompt string) (string, error)
	Generate(prompts []string) ([]string, error)
}
