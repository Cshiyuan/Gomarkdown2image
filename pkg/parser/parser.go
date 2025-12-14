// Package parser 提供 Markdown 到 HTML 的转换功能
//
// 本包使用 Goldmark 解析器,支持 CommonMark 标准和代码语法高亮
package parser

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

// Parser Markdown 解析器接口
type Parser interface {
	// Parse 将 Markdown 转换为 HTML
	Parse(markdown []byte) ([]byte, error)
}

// GoldmarkParser 基于 Goldmark 的解析器实现
type GoldmarkParser struct {
	md goldmark.Markdown
}

// NewGoldmarkParser 创建新的 Goldmark 解析器
//
// 特性:
//   - 支持 CommonMark 标准
//   - 支持 GFM 扩展 (表格、删除线、自动链接等)
//   - 支持代码语法高亮 (使用 Chroma)
func NewGoldmarkParser() *GoldmarkParser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,         // GitHub Flavored Markdown
			extension.Typographer, // 智能标点符号
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"), // 代码高亮主题
				highlighting.WithFormatOptions(
					html.WithLineNumbers(true), // 显示行号
					html.WithClasses(true),     // 使用 CSS 类
				),
			),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(), // 硬换行
			goldmarkhtml.WithXHTML(),     // 使用 XHTML 标签
			goldmarkhtml.WithUnsafe(),    // 允许原始 HTML
		),
	)

	return &GoldmarkParser{md: md}
}

// Parse 将 Markdown 文本转换为 HTML
//
// 参数:
//   - markdown: Markdown 文本字节数组
//
// 返回:
//   - []byte: HTML 字节数组
//   - error: 转换错误(如有)
func (p *GoldmarkParser) Parse(markdown []byte) ([]byte, error) {
	var buf bytes.Buffer

	if err := p.md.Convert(markdown, &buf); err != nil {
		return nil, fmt.Errorf("failed to parse markdown: %w", err)
	}

	return buf.Bytes(), nil
}

// ParseToString 将 Markdown 文本转换为 HTML 字符串
//
// 这是 Parse 方法的便捷包装
func (p *GoldmarkParser) ParseToString(markdown string) (string, error) {
	html, err := p.Parse([]byte(markdown))
	if err != nil {
		return "", err
	}
	return string(html), nil
}
