package config

// AI 相关默认配置
const (
	// DefaultAIProvider 默认 AI 提供器
	DefaultAIProvider = "gemini"

	// DefaultAIModel 默认 AI 模型
	DefaultAIModel = "gemini-2.0-flash-exp"

	// DefaultOllamaEndpoint 默认 Ollama 端点
	DefaultOllamaEndpoint = "http://localhost:11434"

	// DefaultAIPromptTemplate 默认提示词模板
	DefaultAIPromptTemplate = "enhance"

	// DefaultAITimeout AI 请求默认超时(秒)
	DefaultAITimeout = 60

	// MaxAIPromptLength AI 提示词最大长度
	MaxAIPromptLength = 10000
)
