package parser

import (
	"fmt"
	"html/template"
	"strings"
)

// HTMLTemplate HTML 模板配置
type HTMLTemplate struct {
	Title      string // 页面标题
	Theme      string // 主题名称 (light, dark)
	CustomCSS  string // 自定义 CSS
	Width      int    // 页面宽度
	FontSize   int    // 字体大小
	FontFamily string // 字体族
}

// DefaultTemplate 返回默认模板配置
func DefaultTemplate() *HTMLTemplate {
	return &HTMLTemplate{
		Title:      "Markdown to Image",
		Theme:      "light",
		Width:      1200,
		FontSize:   16,
		FontFamily: "Arial, sans-serif",
	}
}

// WrapHTML 将 HTML 内容包装成完整的 HTML 文档
//
// 参数:
//   - content: HTML 内容(由 Parser 生成)
//   - tmpl: HTML 模板配置
//
// 返回:
//   - string: 完整的 HTML 文档
//   - error: 模板渲染错误(如有)
func WrapHTML(content string, tmpl *HTMLTemplate) (string, error) {
	if tmpl == nil {
		tmpl = DefaultTemplate()
	}

	// 构建完整的 HTML 文档
	htmlDoc := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        %s
    </style>
    %s
</head>
<body>
    <div class="container">
        %s
    </div>
</body>
</html>`,
		template.HTMLEscapeString(tmpl.Title),
		generateBaseCSS(tmpl),
		tmpl.CustomCSS,
		content,
	)

	return htmlDoc, nil
}

// generateBaseCSS 生成基础 CSS 样式
func generateBaseCSS(tmpl *HTMLTemplate) string {
	var css strings.Builder

	// 基础样式
	css.WriteString(fmt.Sprintf(`
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: %s;
            font-size: %dpx;
            line-height: 1.6;
            color: %s;
            background-color: %s;
            padding: 40px 20px;
        }

        .container {
            max-width: %dpx;
            margin: 0 auto;
            padding: 40px;
            background: %s;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
`,
		tmpl.FontFamily,
		tmpl.FontSize,
		getTextColor(tmpl.Theme),
		getBackgroundColor(tmpl.Theme),
		tmpl.Width,
		getContainerColor(tmpl.Theme),
	))

	// Markdown 元素样式
	css.WriteString(`
        h1, h2, h3, h4, h5, h6 {
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }

        h1 { font-size: 2em; border-bottom: 1px solid #eaecef; padding-bottom: 8px; }
        h2 { font-size: 1.5em; border-bottom: 1px solid #eaecef; padding-bottom: 8px; }
        h3 { font-size: 1.25em; }
        h4 { font-size: 1em; }
        h5 { font-size: 0.875em; }
        h6 { font-size: 0.85em; color: #6a737d; }

        p {
            margin-bottom: 16px;
        }

        a {
            color: #0366d6;
            text-decoration: none;
        }

        a:hover {
            text-decoration: underline;
        }

        ul, ol {
            margin-bottom: 16px;
            padding-left: 2em;
        }

        li {
            margin-bottom: 8px;
        }

        blockquote {
            padding: 0 1em;
            color: #6a737d;
            border-left: 4px solid #dfe2e5;
            margin-bottom: 16px;
        }

        code {
            padding: 2px 6px;
            font-family: 'Courier New', Courier, monospace;
            font-size: 0.9em;
            background-color: rgba(27,31,35,0.05);
            border-radius: 3px;
        }

        pre {
            padding: 16px;
            overflow: auto;
            font-size: 0.9em;
            line-height: 1.45;
            background-color: #f6f8fa;
            border-radius: 6px;
            margin-bottom: 16px;
        }

        pre code {
            display: block;
            padding: 0;
            background: transparent;
            border-radius: 0;
        }

        table {
            border-collapse: collapse;
            width: 100%;
            margin-bottom: 16px;
        }

        th, td {
            padding: 12px;
            border: 1px solid #dfe2e5;
            text-align: left;
        }

        th {
            background-color: #f6f8fa;
            font-weight: 600;
        }

        hr {
            height: 1px;
            margin: 24px 0;
            background-color: #e1e4e8;
            border: 0;
        }

        img {
            max-width: 100%;
            height: auto;
            display: block;
            margin: 16px 0;
        }

        /* Chroma 代码高亮样式 */
        .chroma {
            background-color: #272822;
            color: #f8f8f2;
            padding: 16px;
            border-radius: 6px;
            overflow-x: auto;
        }

        .chroma .line {
            display: block;
        }

        .chroma .ln {
            margin-right: 12px;
            color: #75715e;
        }
    `)

	return css.String()
}

// 主题颜色辅助函数
func getTextColor(theme string) string {
	if theme == "dark" {
		return "#e6edf3"
	}
	return "#24292f"
}

func getBackgroundColor(theme string) string {
	if theme == "dark" {
		return "#0d1117"
	}
	return "#ffffff"
}

func getContainerColor(theme string) string {
	if theme == "dark" {
		return "#161b22"
	}
	return "#ffffff"
}
