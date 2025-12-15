package ai

import (
	"context"
	"fmt"
)

// Provider 定义了 AI 服务提供器的统一接口
//
// 该接口抽象了不同 AI 后端(Claude, Ollama 等)的实现细节,
// 提供一致的调用方式。
type Provider interface {
	// Generate 生成文本响应
	//
	// 参数:
	//   - ctx: 请求上下文(用于超时控制和取消)
	//   - req: 生成请求
	//
	// 返回:
	//   - *GenerateResponse: 生成响应
	//   - error: 错误信息
	Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)

	// GenerateStream 流式生成文本响应
	//
	// 参数:
	//   - ctx: 请求上下文
	//   - req: 生成请求
	//
	// 返回:
	//   - <-chan StreamChunk: 流式响应通道
	//   - error: 错误信息(初始化错误,流式错误通过 StreamChunk.Error 返回)
	GenerateStream(ctx context.Context, req *GenerateRequest) (<-chan StreamChunk, error)

	// Name 返回 Provider 名称
	Name() string

	// Close 关闭连接和清理资源
	Close() error
}

// ValidateConfig 验证 AI 配置是否有效
func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// 验证 Provider 类型
	if cfg.Provider != ProviderGemini && cfg.Provider != ProviderOllama {
		return fmt.Errorf("invalid provider type: %s", cfg.Provider)
	}

	// Gemini 特定验证
	if cfg.Provider == ProviderGemini {
		if cfg.APIKey == "" {
			return fmt.Errorf("API key is required for Gemini provider")
		}
	}

	// 验证模型名称
	if cfg.Model == "" {
		return fmt.Errorf("model name cannot be empty")
	}

	// 验证超时时间
	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout must be non-negative")
	}

	// 验证重试次数
	if cfg.MaxRetries < 0 {
		return fmt.Errorf("max retries must be non-negative")
	}

	return nil
}

// ValidateGenerateRequest 验证生成请求是否有效
func ValidateGenerateRequest(req *GenerateRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.Prompt == "" {
		return fmt.Errorf("prompt cannot be empty")
	}

	if req.MaxTokens < 0 {
		return fmt.Errorf("max tokens must be non-negative")
	}

	if req.Temperature < 0 || req.Temperature > 1 {
		return fmt.Errorf("temperature must be between 0 and 1")
	}

	return nil
}
