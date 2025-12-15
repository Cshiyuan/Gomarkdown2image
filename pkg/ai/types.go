// Package ai 提供统一的 AI 服务接口,支持多种 AI 后端(Claude, Ollama 等)
package ai

import (
	"context"
)

// ProviderType AI 提供器类型
type ProviderType string

const (
	// ProviderGemini Google Gemini API
	ProviderGemini ProviderType = "gemini"
	// ProviderOllama 本地 Ollama 服务
	ProviderOllama ProviderType = "ollama"
)

// Config AI 服务通用配置
type Config struct {
	// Provider 提供器类型 (gemini 或 ollama)
	Provider ProviderType

	// APIKey API 密钥 (仅 Gemini 需要)
	APIKey string

	// BaseURL 自定义 API 基础 URL (可选)
	BaseURL string

	// Model 模型名称
	// Gemini: "gemini-2.0-flash-exp", "gemini-1.5-pro" 等
	// Ollama: "llama3.2", "mistral" 等
	Model string

	// Timeout 超时时间(秒)
	Timeout int

	// MaxRetries 最大重试次数
	MaxRetries int

	// Prompts 提示词配置
	Prompts *PromptConfig

	// Extra 额外配置,用于特定提供器的自定义选项
	Extra map[string]string
}

// PromptConfig 提示词配置
type PromptConfig struct {
	// DefaultSystem 默认系统提示词
	DefaultSystem string

	// Templates 提示词模板集合
	Templates map[string]string
}

// GenerateRequest AI 生成请求
type GenerateRequest struct {
	// Prompt 用户提示词
	Prompt string

	// System 系统提示词 (可选,覆盖 Config 中的默认值)
	System string

	// MaxTokens 最大生成 token 数
	MaxTokens int

	// Temperature 温度参数 (0.0-1.0),控制随机性
	Temperature float64

	// Metadata 额外元数据
	Metadata map[string]string

	// Context 请求上下文
	Context context.Context
}

// GenerateResponse AI 生成响应
type GenerateResponse struct {
	// Content 生成的内容
	Content string

	// TokensUsed 使用的 token 数量
	TokensUsed int

	// FinishReason 结束原因 ("stop", "length", "content_filter" 等)
	FinishReason string

	// Metadata 响应元数据
	Metadata map[string]string
}

// StreamChunk 流式响应块
type StreamChunk struct {
	// Content 内容片段
	Content string

	// Done 是否完成
	Done bool

	// Error 错误信息
	Error error
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Provider:   ProviderGemini,
		Model:      "gemini-2.0-flash-exp",
		Timeout:    30,
		MaxRetries: 3,
		Prompts:    DefaultPromptConfig(),
		Extra:      make(map[string]string),
	}
}

// DefaultPromptConfig 返回默认提示词配置
func DefaultPromptConfig() *PromptConfig {
	return &PromptConfig{
		DefaultSystem: "你是一个专业的 Markdown 文档处理助手,擅长优化和增强文档内容。",
		Templates:     make(map[string]string),
	}
}
