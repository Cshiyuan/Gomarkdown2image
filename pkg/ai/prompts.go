package ai

import (
	"bytes"
	"fmt"
	"text/template"
)

// PromptTemplate 提示词模板
type PromptTemplate struct {
	name     string
	template *template.Template
}

// DefaultPrompts 默认提示词模板集合
var DefaultPrompts = map[string]string{
	// enhance 增强和润色 Markdown 文档
	"enhance": `你是一个专业的文档优化助手。请润色以下 Markdown 内容,使其更清晰、流畅和专业。

要求:
1. 优化语言表达,使文档更易读
2. 保持技术准确性,不改变原意
3. 保留所有代码块内容(不修改代码)
4. 保留 Markdown 格式结构(标题、列表、表格等)
5. 修正明显的语法和拼写错误
6. 直接返回优化后的 Markdown,不要添加任何说明文字

原始内容:
{{.Content}}`,

	// translate 翻译 Markdown 文档
	"translate": `你是一个专业的技术文档翻译助手。请将以下 Markdown 文档翻译成{{.TargetLang}}。

要求:
1. 保持 Markdown 格式完整(标题、列表、代码块、表格等)
2. 代码块内的代码不翻译
3. 专业术语使用准确的{{.TargetLang}}表达
4. 保持文档的技术准确性
5. 翻译要自然流畅,符合{{.TargetLang}}表达习惯
6. 直接返回翻译后的 Markdown,不要添加任何说明文字

原始内容:
{{.Content}}`,

	// explain_code 解释代码功能
	"explain_code": `你是一个代码分析专家。请解释以下{{.Language}}代码的功能和实现逻辑。

要求:
1. 用中文简洁清晰地解释代码的主要功能
2. 指出关键技术点和算法思路
3. 如有复杂逻辑,请分步骤说明
4. 指出潜在的优化点(如有)

代码:
` + "```{{.Language}}" + `
{{.Code}}
` + "```" + ``,

	// format 格式化和美化 Markdown
	"format": `你是一个 Markdown 格式化专家。请优化以下 Markdown 文档的格式。

要求:
1. 统一标题层级(确保逻辑结构清晰)
2. 规范化列表格式(缩进、符号一致)
3. 美化表格对齐
4. 优化代码块语言标记
5. 调整段落间距,提高可读性
6. 直接返回格式化后的 Markdown,不要添加任何说明文字

原始内容:
{{.Content}}`,

	// summarize 总结文档内容
	"summarize": `你是一个技术文档总结专家。请为以下 Markdown 文档生成简洁的摘要。

要求:
1. 提取文档的核心内容和关键信息
2. 总结控制在 3-5 个要点
3. 使用 Markdown 列表格式
4. 保持技术准确性
5. 语言简洁专业

原始内容:
{{.Content}}`,
}

// RenderPrompt 渲染提示词模板
//
// 参数:
//   - templateName: 模板名称 (如 "enhance", "translate")
//   - data: 模板数据 (如 map[string]interface{}{"Content": "...", "TargetLang": "英文"})
//
// 返回:
//   - string: 渲染后的提示词
//   - error: 渲染错误
func RenderPrompt(templateName string, data map[string]interface{}) (string, error) {
	// 获取模板字符串
	tmplStr, ok := DefaultPrompts[templateName]
	if !ok {
		return "", fmt.Errorf("template not found: %s", templateName)
	}

	// 解析模板
	tmpl, err := template.New(templateName).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// 渲染模板
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

// GetPromptTemplate 获取提示词模板字符串
func GetPromptTemplate(templateName string) (string, bool) {
	tmplStr, ok := DefaultPrompts[templateName]
	return tmplStr, ok
}

// SetPromptTemplate 设置自定义提示词模板
func SetPromptTemplate(templateName string, templateStr string) {
	DefaultPrompts[templateName] = templateStr
}

// ListPromptTemplates 列出所有可用的提示词模板名称
func ListPromptTemplates() []string {
	names := make([]string, 0, len(DefaultPrompts))
	for name := range DefaultPrompts {
		names = append(names, name)
	}
	return names
}
