package utils

import (
	"strings"
	"testing"
)

func TestValidateCustomCSS(t *testing.T) {
	tests := []struct {
		name    string
		css     string
		wantErr bool
	}{
		// 有效的 CSS
		{
			name:    "有效CSS",
			css:     "body { color: red; background: blue; }",
			wantErr: false,
		},
		{
			name:    "空CSS",
			css:     "",
			wantErr: false,
		},
		{
			name:    "复杂但安全的CSS",
			css:     ".container { padding: 20px; margin: auto; } h1 { font-size: 24px; }",
			wantErr: false,
		},

		// XSS 注入尝试
		{
			name:    "闭合style标签",
			css:     "</style><script>alert('xss')</script>",
			wantErr: true,
		},
		{
			name:    "注入script标签",
			css:     "<script>alert(1)</script>",
			wantErr: true,
		},
		{
			name:    "注入iframe",
			css:     "<iframe src='http://evil.com'></iframe>",
			wantErr: true,
		},
		{
			name:    "JavaScript协议",
			css:     "background-image: url('javascript:alert(1)');",
			wantErr: true,
		},
		{
			name:    "data URI",
			css:     "background-image: url('data:text/html,<script>alert(1)</script>');",
			wantErr: true,
		},
		{
			name:    "事件处理器-onerror",
			css:     "div { onerror: alert(1); }",
			wantErr: true,
		},
		{
			name:    "事件处理器-onload",
			css:     "img { onload: alert(1); }",
			wantErr: true,
		},
		{
			name:    "混合大小写绕过",
			css:     "</STYLE><ScRiPt>alert('xss')</ScRiPt>",
			wantErr: true,
		},

		// 大小限制
		{
			name:    "超过100KB限制",
			css:     strings.Repeat("a", 100001),
			wantErr: true,
		},

		// 引号不配对
		{
			name:    "双引号不配对",
			css:     `body { content: "hello; }`,
			wantErr: true,
		},
		{
			name:    "单引号不配对",
			css:     `body { content: 'hello; }`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCustomCSS(tt.css)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCustomCSS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
