package parser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "标题H1",
			input:   "# Hello World",
			want:    "<h1",
			wantErr: false,
		},
		{
			name:    "标题H2",
			input:   "## Subtitle",
			want:    "<h2",
			wantErr: false,
		},
		{
			name:    "段落",
			input:   "This is a paragraph.",
			want:    "<p>This is a paragraph.</p>",
			wantErr: false,
		},
		{
			name:    "粗体",
			input:   "**bold text**",
			want:    "<strong>bold text</strong>",
			wantErr: false,
		},
		{
			name:    "斜体",
			input:   "*italic text*",
			want:    "<em>italic text</em>",
			wantErr: false,
		},
		{
			name:    "代码块",
			input:   "```go\nfunc main() {}\n```",
			want:    "<pre",
			wantErr: false,
		},
		{
			name:    "行内代码",
			input:   "`code`",
			want:    "<code>code</code>",
			wantErr: false,
		},
		{
			name:    "链接",
			input:   "[Go](https://go.dev)",
			want:    `<a href="https://go.dev"`,
			wantErr: false,
		},
		{
			name:    "无序列表",
			input:   "- Item 1\n- Item 2",
			want:    "<ul>",
			wantErr: false,
		},
		{
			name:    "有序列表",
			input:   "1. First\n2. Second",
			want:    "<ol>",
			wantErr: false,
		},
		{
			name:    "引用块",
			input:   "> This is a quote",
			want:    "<blockquote>",
			wantErr: false,
		},
		{
			name:    "表格(GFM)",
			input:   "| Header |\n| --- |\n| Cell |",
			want:    "<table>",
			wantErr: false,
		},
		{
			name:    "删除线(GFM)",
			input:   "~~strikethrough~~",
			want:    "<del>strikethrough</del>",
			wantErr: false,
		},
		{
			name:    "任务列表(GFM)",
			input:   "- [ ] Unchecked\n- [x] Checked",
			want:    `<input`,
			wantErr: false,
		},
		{
			name:    "空输入",
			input:   "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "多行混合内容",
			input:   "# Title\n\nParagraph with **bold** and *italic*.\n\n```\ncode\n```",
			want:    "<h1",
			wantErr: false,
		},
	}

	parser := NewGoldmarkParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseToString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.ParseToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != "" && !strings.Contains(got, tt.want) {
				t.Errorf("Parser.ParseToString() output doesn't contain expected substring\ngot: %s\nwant substring: %s", got, tt.want)
			}
		})
	}
}

func TestWrapHTML(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		template  *HTMLTemplate
		wantParts []string
		wantErr   bool
	}{
		{
			name:    "默认模板",
			content: "<p>Hello</p>",
			template: &HTMLTemplate{
				Title:      "Test",
				Theme:      "light",
				Width:      1200,
				FontSize:   16,
				FontFamily: "Arial",
			},
			wantParts: []string{
				"<!DOCTYPE html>",
				`<html lang="zh-CN">`,
				"<title>Test</title>",
				"<p>Hello</p>",
				"</html>",
			},
			wantErr: false,
		},
		{
			name:    "暗色主题",
			content: "<h1>Title</h1>",
			template: &HTMLTemplate{
				Title: "Dark Theme",
				Theme: "dark",
				Width: 1400,
			},
			wantParts: []string{
				"<title>Dark Theme</title>",
				"<h1>Title</h1>",
			},
			wantErr: false,
		},
		{
			name:    "自定义CSS",
			content: "<div>Custom</div>",
			template: &HTMLTemplate{
				Title:     "Custom",
				CustomCSS: "body { color: red; }",
			},
			wantParts: []string{
				"body { color: red; }",
				"<div>Custom</div>",
			},
			wantErr: false,
		},
		{
			name:    "空内容",
			content: "",
			template: &HTMLTemplate{
				Title: "Empty",
			},
			wantParts: []string{
				"<!DOCTYPE html>",
				"<title>Empty</title>",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WrapHTML(tt.content, tt.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("WrapHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, part := range tt.wantParts {
				if !strings.Contains(got, part) {
					t.Errorf("WrapHTML() output doesn't contain expected part: %s", part)
				}
			}
		})
	}
}

func TestDefaultTemplate(t *testing.T) {
	tmpl := DefaultTemplate()

	if tmpl.Title == "" {
		t.Error("DefaultTemplate() Title should not be empty")
	}
	if tmpl.Theme != "light" {
		t.Errorf("DefaultTemplate() Theme = %v, want %v", tmpl.Theme, "light")
	}
	if tmpl.Width <= 0 {
		t.Error("DefaultTemplate() Width should be positive")
	}
	if tmpl.FontSize <= 0 {
		t.Error("DefaultTemplate() FontSize should be positive")
	}
	if tmpl.FontFamily == "" {
		t.Error("DefaultTemplate() FontFamily should not be empty")
	}
}

// 基准测试
func BenchmarkParse(b *testing.B) {
	parser := NewGoldmarkParser()
	markdown := `# Benchmark Test

This is a **benchmark** test with *various* elements:

- List item 1
- List item 2

` + "```go\nfunc main() {}\n```"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.ParseToString(markdown)
	}
}

func BenchmarkWrapHTML(b *testing.B) {
	content := "<h1>Test</h1><p>Paragraph</p>"
	tmpl := DefaultTemplate()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WrapHTML(content, tmpl)
	}
}
