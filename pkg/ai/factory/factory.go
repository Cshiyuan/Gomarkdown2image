// Package factory 提供 AI Provider 的工厂函数,避免循环导入
package factory

import (
	"fmt"

	"github.com/Cshiyuan/Gomarkdown2image/pkg/ai"
	"github.com/Cshiyuan/Gomarkdown2image/pkg/ai/gemini"
	"github.com/Cshiyuan/Gomarkdown2image/pkg/ai/ollama"
)

// NewProvider 根据配置创建对应的 AI Provider
//
// 这是一个工厂函数,根据 Config.Provider 字段自动选择实现:
//   - ProviderGemini: 创建 Google Gemini API 客户端
//   - ProviderOllama: 创建 Ollama 本地客户端
//
// 参数:
//   - cfg: AI 配置
//
// 返回:
//   - ai.Provider: AI 提供器实例
//   - error: 创建错误
func NewProvider(cfg *ai.Config) (ai.Provider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	switch cfg.Provider {
	case ai.ProviderGemini:
		return gemini.New(cfg)
	case ai.ProviderOllama:
		return ollama.New(cfg)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}
