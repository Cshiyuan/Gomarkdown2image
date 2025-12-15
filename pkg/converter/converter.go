// Package converter 提供 Markdown 到图片的端到端转换功能
//
// 本包协调 Parser 和 Renderer,提供统一的转换接口
package converter

import (
	"fmt"
	"os"

	"github.com/Cshiyuan/Gomarkdown2image/pkg/ai"
	"github.com/Cshiyuan/Gomarkdown2image/pkg/parser"
	"github.com/Cshiyuan/Gomarkdown2image/pkg/renderer"
)

// Converter Markdown 到图片转换器接口
type Converter interface {
	// Convert 将 Markdown 转换为图片
	Convert(markdown []byte, opts *ConvertOptions) ([]byte, error)

	// ConvertFile 转换 Markdown 文件为图片文件
	ConvertFile(inputPath string, outputPath string, opts *ConvertOptions) error

	// Close 关闭转换器,释放资源
	Close() error
}

// ConvertOptions 转换选项
type ConvertOptions struct {
	// HTML 模板选项
	Title      string // 页面标题
	Theme      string // 主题 (light, dark)
	CustomCSS  string // 自定义 CSS
	Width      int    // 页面宽度
	FontSize   int    // 字体大小
	FontFamily string // 字体族

	// 渲染选项
	ImageFormat      renderer.ImageFormat // 图片格式
	ImageQuality     int                  // 图片质量
	FullPage         bool                 // 全页截图
	DevicePixelRatio float64              // 设备像素比

	// AI 增强选项 (新增)
	ParserMode       string // 解析器模式: "traditional" (默认) 或 "ai"
	AIProvider       string // AI 提供器: "gemini" 或 "ollama"
	AIModel          string // AI 模型名称
	AIAPIKey         string // AI API 密钥 (Gemini 需要)
	AIEndpoint       string // AI 服务端点 (Ollama 使用)
	AIPromptTemplate string // 提示词模板: "enhance", "translate", "format" 等
	AICustomPrompt   string // 自定义提示词 (覆盖模板)
	AIPromptData     map[string]interface{} // 提示词模板数据
}

// DefaultConvertOptions 返回默认转换选项
func DefaultConvertOptions() *ConvertOptions {
	return &ConvertOptions{
		Title:            "Markdown to Image",
		Theme:            "light",
		Width:            1200,
		FontSize:         16,
		FontFamily:       "Arial, sans-serif",
		ImageFormat:      renderer.FormatPNG,
		ImageQuality:     90,
		FullPage:         true,
		DevicePixelRatio: 1.0,
		// AI 默认值
		ParserMode:       "traditional", // 默认使用传统模式
		AIProvider:       "gemini",
		AIModel:          "gemini-2.0-flash-exp",
		AIEndpoint:       "http://localhost:11434",
		AIPromptTemplate: "enhance",
		AIPromptData:     make(map[string]interface{}),
	}
}

// DefaultConverter 默认转换器实现
type DefaultConverter struct {
	parser   parser.Parser
	renderer renderer.Renderer
}

// NewConverter 创建新的转换器
//
// 返回:
//   - Converter: 转换器实例
//   - error: 初始化错误(如有)
func NewConverter() (Converter, error) {
	// 创建 Parser
	p := parser.NewGoldmarkParser()

	// 创建 Renderer
	r, err := renderer.NewRodRenderer()
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}

	return &DefaultConverter{
		parser:   p,
		renderer: r,
	}, nil
}

// Convert 将 Markdown 字节数组转换为图片字节数组
//
// 工作流程:
//  1. 根据 ParserMode 创建对应的 Parser (传统/AI)
//  2. Markdown → HTML (使用 Parser)
//  3. HTML → 完整 HTML 文档 (应用模板)
//  4. HTML 文档 → 图片 (使用 Renderer)
//
// 参数:
//   - markdown: Markdown 文本字节数组
//   - opts: 转换选项
//
// 返回:
//   - []byte: 图片字节数组
//   - error: 转换错误(如有)
func (c *DefaultConverter) Convert(markdown []byte, opts *ConvertOptions) ([]byte, error) {
	if opts == nil {
		opts = DefaultConvertOptions()
	}

	// 步骤 1: 根据 ParserMode 创建 Parser
	var currentParser parser.Parser
	var err error

	if opts.ParserMode == "ai" {
		// 创建 AI Parser
		currentParser, err = c.createAIParser(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to create AI parser: %w", err)
		}

		// 如果是 AI Parser Provider,需要在完成后关闭
		if closer, ok := currentParser.(interface{ Close() error }); ok {
			defer closer.Close()
		}
	} else {
		// 使用传统 Parser
		currentParser = c.parser
	}

	// 步骤 2: 解析 Markdown → HTML
	htmlContent, err := currentParser.Parse(markdown)
	if err != nil {
		return nil, fmt.Errorf("failed to parse markdown: %w", err)
	}

	// 步骤 3: 包装为完整的 HTML 文档
	tmpl := &parser.HTMLTemplate{
		Title:      opts.Title,
		Theme:      opts.Theme,
		CustomCSS:  opts.CustomCSS,
		Width:      opts.Width,
		FontSize:   opts.FontSize,
		FontFamily: opts.FontFamily,
	}

	fullHTML, err := parser.WrapHTML(string(htmlContent), tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap HTML: %w", err)
	}

	// 步骤 4: 渲染 HTML → 图片
	renderOpts := &renderer.RenderOptions{
		Width:            opts.Width,
		Height:           0, // 自动高度
		Format:           opts.ImageFormat,
		Quality:          opts.ImageQuality,
		FullPage:         opts.FullPage,
		DevicePixelRatio: opts.DevicePixelRatio,
	}

	imageData, err := c.renderer.RenderToImage(fullHTML, renderOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to render image: %w", err)
	}

	return imageData, nil
}

// createAIParser 根据配置创建 AI Parser
func (c *DefaultConverter) createAIParser(opts *ConvertOptions) (parser.Parser, error) {
	// 构建 AI 配置
	aiConfig := &ai.Config{
		Provider:   ai.ProviderType(opts.AIProvider),
		APIKey:     opts.AIAPIKey,
		BaseURL:    opts.AIEndpoint,
		Model:      opts.AIModel,
		Timeout:    30,
		MaxRetries: 3,
	}

	// 创建 Parser Provider 配置
	providerConfig := &parser.ProviderConfig{
		Type:             parser.ProviderTypeAI,
		AIConfig:         aiConfig,
		AIPromptTemplate: opts.AIPromptTemplate,
		AIPromptData:     opts.AIPromptData,
		CustomPrompt:     opts.AICustomPrompt,
	}

	// 创建 Provider
	provider, err := parser.NewProvider(providerConfig)
	if err != nil {
		return nil, err
	}

	// 创建 Parser
	return provider.CreateParser()
}

// ConvertFile 转换 Markdown 文件为图片文件
//
// 参数:
//   - inputPath: Markdown 文件路径
//   - outputPath: 输出图片文件路径
//   - opts: 转换选项
//
// 返回:
//   - error: 转换或文件操作错误(如有)
func (c *DefaultConverter) ConvertFile(inputPath string, outputPath string, opts *ConvertOptions) error {
	// 读取 Markdown 文件
	markdown, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// 转换为图片
	imageData, err := c.Convert(markdown, opts)
	if err != nil {
		return fmt.Errorf("failed to convert: %w", err)
	}

	// 写入图片文件
	if err := os.WriteFile(outputPath, imageData, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// Close 关闭转换器,释放资源
func (c *DefaultConverter) Close() error {
	if c.renderer != nil {
		return c.renderer.Close()
	}
	return nil
}
