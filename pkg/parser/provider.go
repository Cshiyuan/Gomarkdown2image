package parser

import (
	"context"
	"fmt"

	"github.com/Cshiyuan/Gomarkdown2image/pkg/ai"
	"github.com/Cshiyuan/Gomarkdown2image/pkg/ai/factory"
)

// ParserProvider 解析器提供器接口
//
// 该接口定义了创建 Parser 的工厂方法,支持不同的解析策略:
//   - 传统模式: 使用 Goldmark 直接解析
//   - AI 增强模式: 先通过 AI 增强/优化内容,再解析
type ParserProvider interface {
	// CreateParser 创建 Parser 实例
	CreateParser() (Parser, error)

	// Name 返回提供器名称
	Name() string
}

// ProviderType 提供器类型
type ProviderType string

const (
	// ProviderTypeTraditional 传统 Goldmark 解析器
	ProviderTypeTraditional ProviderType = "traditional"

	// ProviderTypeAI AI 增强解析器
	ProviderTypeAI ProviderType = "ai"
)

// ProviderConfig 提供器配置
type ProviderConfig struct {
	// Type 提供器类型 (traditional 或 ai)
	Type ProviderType

	// AIConfig AI 配置 (仅当 Type=ai 时使用)
	AIConfig *ai.Config

	// AIPromptTemplate 自定义提示词模板名称 (可选)
	// 可用值: "enhance", "translate", "format" 等
	AIPromptTemplate string

	// AIPromptData 提示词模板数据 (可选)
	// 例如: map[string]interface{}{"TargetLang": "英文"}
	AIPromptData map[string]interface{}

	// CustomPrompt 自定义提示词 (可选,覆盖模板)
	CustomPrompt string
}

// GoldmarkProvider 传统 Goldmark 解析器提供器
type GoldmarkProvider struct{}

// NewGoldmarkProvider 创建 Goldmark 提供器
func NewGoldmarkProvider() *GoldmarkProvider {
	return &GoldmarkProvider{}
}

// CreateParser 创建 GoldmarkParser 实例
func (p *GoldmarkProvider) CreateParser() (Parser, error) {
	return NewGoldmarkParser(), nil
}

// Name 返回提供器名称
func (p *GoldmarkProvider) Name() string {
	return "goldmark"
}

// AIParserProvider AI 增强解析器提供器
type AIParserProvider struct {
	aiProvider       ai.Provider
	promptTemplate   string
	promptData       map[string]interface{}
	customPrompt     string
	fallbackProvider ParserProvider
}

// NewAIParserProvider 创建 AI 增强解析器提供器
//
// 参数:
//   - aiConfig: AI 配置
//   - promptTemplate: 提示词模板名称 (如 "enhance")
//   - promptData: 模板数据 (可选)
//   - customPrompt: 自定义提示词 (可选,覆盖模板)
//
// 返回:
//   - *AIParserProvider: AI 提供器实例
//   - error: 创建错误
func NewAIParserProvider(aiConfig *ai.Config, promptTemplate string, promptData map[string]interface{}, customPrompt string) (*AIParserProvider, error) {
	// 创建 AI Provider
	aiProvider, err := factory.NewProvider(aiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI provider: %w", err)
	}

	// 创建降级提供器 (AI 失败时使用)
	fallbackProvider := NewGoldmarkProvider()

	return &AIParserProvider{
		aiProvider:       aiProvider,
		promptTemplate:   promptTemplate,
		promptData:       promptData,
		customPrompt:     customPrompt,
		fallbackProvider: fallbackProvider,
	}, nil
}

// CreateParser 创建 AI 增强的 Parser 实例
func (p *AIParserProvider) CreateParser() (Parser, error) {
	return &AIParser{
		aiProvider:       p.aiProvider,
		promptTemplate:   p.promptTemplate,
		promptData:       p.promptData,
		customPrompt:     p.customPrompt,
		fallbackParser:   NewGoldmarkParser(),
		enableFallback:   true,
	}, nil
}

// Name 返回提供器名称
func (p *AIParserProvider) Name() string {
	return fmt.Sprintf("ai-%s", p.aiProvider.Name())
}

// Close 关闭 AI Provider 连接
func (p *AIParserProvider) Close() error {
	if p.aiProvider != nil {
		return p.aiProvider.Close()
	}
	return nil
}

// AIParser AI 增强的 Markdown 解析器
type AIParser struct {
	aiProvider       ai.Provider
	promptTemplate   string
	promptData       map[string]interface{}
	customPrompt     string
	fallbackParser   Parser
	enableFallback   bool
}

// Parse 使用 AI 增强 Markdown 内容,然后转换为 HTML
//
// 工作流程:
//  1. 调用 AI 服务增强/优化 Markdown 内容
//  2. 使用 Goldmark 解析增强后的 Markdown
//  3. 如果 AI 失败且启用降级,直接使用 Goldmark 解析原始内容
func (p *AIParser) Parse(markdown []byte) ([]byte, error) {
	// 第 1 步: 使用 AI 增强内容
	enhancedMarkdown, err := p.enhanceWithAI(string(markdown))
	if err != nil {
		// AI 失败,检查是否启用降级
		if p.enableFallback && p.fallbackParser != nil {
			// 降级到传统解析
			return p.fallbackParser.Parse(markdown)
		}
		return nil, fmt.Errorf("AI enhancement failed: %w", err)
	}

	// 第 2 步: 使用 Goldmark 解析增强后的内容
	parser := NewGoldmarkParser()
	return parser.Parse([]byte(enhancedMarkdown))
}

// enhanceWithAI 使用 AI 增强 Markdown 内容
func (p *AIParser) enhanceWithAI(markdown string) (string, error) {
	ctx := context.Background()

	// 构建提示词
	var prompt string
	var err error

	if p.customPrompt != "" {
		// 使用自定义提示词
		prompt = p.customPrompt + "\n\n" + markdown
	} else if p.promptTemplate != "" {
		// 使用模板渲染提示词
		data := p.promptData
		if data == nil {
			data = make(map[string]interface{})
		}
		data["Content"] = markdown

		prompt, err = ai.RenderPrompt(p.promptTemplate, data)
		if err != nil {
			return "", fmt.Errorf("failed to render prompt: %w", err)
		}
	} else {
		// 使用默认增强提示词
		prompt, err = ai.RenderPrompt("enhance", map[string]interface{}{
			"Content": markdown,
		})
		if err != nil {
			return "", fmt.Errorf("failed to render default prompt: %w", err)
		}
	}

	// 调用 AI 生成
	req := &ai.GenerateRequest{
		Prompt:      prompt,
		MaxTokens:   8192,
		Temperature: 0.7,
		Context:     ctx,
	}

	resp, err := p.aiProvider.Generate(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

// NewProvider 根据配置创建对应的 ParserProvider
//
// 工厂函数,根据 ProviderConfig.Type 自动选择实现
func NewProvider(cfg *ProviderConfig) (ParserProvider, error) {
	if cfg == nil {
		// 默认使用传统模式
		return NewGoldmarkProvider(), nil
	}

	switch cfg.Type {
	case ProviderTypeTraditional:
		return NewGoldmarkProvider(), nil

	case ProviderTypeAI:
		if cfg.AIConfig == nil {
			return nil, fmt.Errorf("AI config is required for AI provider")
		}
		return NewAIParserProvider(
			cfg.AIConfig,
			cfg.AIPromptTemplate,
			cfg.AIPromptData,
			cfg.CustomPrompt,
		)

	default:
		return nil, fmt.Errorf("unsupported provider type: %s", cfg.Type)
	}
}
