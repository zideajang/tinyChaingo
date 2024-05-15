package agent

type AgentConfig {
	Type string
	MaxRetries int
	EnableMemory bool
}

type Agent struct {
	Name string
	Config AgentConfig
}

type Option func(*LLMAgent)


func WithTimeout(timeout time.Duration) Option {
    return func(s *MyService) {
        s.Timeout = timeout
    }
}