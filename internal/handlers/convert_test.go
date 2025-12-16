package handlers

import (
	"testing"

	"github.com/Cshiyuan/Gomarkdown2image/pkg/converter"
)

// TestBuildConvertOptions 测试从 ConvertRequest 构建选项
func TestBuildConvertOptions(t *testing.T) {
	req := &ConvertRequest{
		Markdown:         "# Test",
		Title:            "测试标题",
		Theme:            "dark",
		CustomCSS:        "body { color: red; }",
		Width:            1400,
		FontSize:         18,
		FontFamily:       "Arial",
		ImageFormat:      "jpeg",
		ImageQuality:     85,
		DevicePixelRatio: 2.0,
		ParserMode:       "ai",
		AIProvider:       "gemini",
		AIModel:          "gemini-2.0-flash-exp",
		AIAPIKey:         "test-key",
		AIEndpoint:       "https://api.example.com",
		AIPromptTemplate: "enhance",
		AICustomPrompt:   "自定义提示词",
	}

	opts := buildConvertOptions(req)

	// 验证所有字段都正确映射
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Title", opts.Title, "测试标题"},
		{"Theme", opts.Theme, "dark"},
		{"CustomCSS", opts.CustomCSS, "body { color: red; }"},
		{"Width", opts.Width, 1400},
		{"FontSize", opts.FontSize, 18},
		{"FontFamily", opts.FontFamily, "Arial"},
		{"ImageQuality", opts.ImageQuality, 85},
		{"DevicePixelRatio", opts.DevicePixelRatio, 2.0},
		{"ParserMode", opts.ParserMode, "ai"},
		{"AIProvider", opts.AIProvider, "gemini"},
		{"AIModel", opts.AIModel, "gemini-2.0-flash-exp"},
		{"AIAPIKey", opts.AIAPIKey, "test-key"},
		{"AIEndpoint", opts.AIEndpoint, "https://api.example.com"},
		{"AIPromptTemplate", opts.AIPromptTemplate, "enhance"},
		{"AICustomPrompt", opts.AICustomPrompt, "自定义提示词"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s: got %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

// TestBuildConvertOptionsFromForm 测试从 UploadRequest 构建选项
func TestBuildConvertOptionsFromForm(t *testing.T) {
	req := &UploadRequest{
		Title:            "上传测试",
		Theme:            "light",
		CustomCSS:        ".container { padding: 20px; }",
		Width:            1200,
		FontSize:         16,
		FontFamily:       "Helvetica",
		ImageFormat:      "png",
		ImageQuality:     90,
		DevicePixelRatio: 1.5,
		ParserMode:       "traditional",
		AIProvider:       "ollama",
		AIModel:          "llama3.2",
		AIAPIKey:         "",
		AIEndpoint:       "http://localhost:11434",
		AIPromptTemplate: "translate",
		AICustomPrompt:   "",
	}

	opts := buildConvertOptionsFromForm(req)

	// 验证所有字段都正确映射
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Title", opts.Title, "上传测试"},
		{"Theme", opts.Theme, "light"},
		{"CustomCSS", opts.CustomCSS, ".container { padding: 20px; }"},
		{"Width", opts.Width, 1200},
		{"FontSize", opts.FontSize, 16},
		{"FontFamily", opts.FontFamily, "Helvetica"},
		{"ImageQuality", opts.ImageQuality, 90},
		{"DevicePixelRatio", opts.DevicePixelRatio, 1.5},
		{"ParserMode", opts.ParserMode, "traditional"},
		{"AIProvider", opts.AIProvider, "ollama"},
		{"AIModel", opts.AIModel, "llama3.2"},
		{"AIEndpoint", opts.AIEndpoint, "http://localhost:11434"},
		{"AIPromptTemplate", opts.AIPromptTemplate, "translate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s: got %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

// TestBuildConvertOptionsDefaults 测试默认值处理
func TestBuildConvertOptionsDefaults(t *testing.T) {
	// 空请求应该返回默认值
	req := &ConvertRequest{
		Markdown: "# Test",
	}

	opts := buildConvertOptions(req)
	defaults := converter.DefaultConvertOptions()

	// 验证未设置的字段使用默认值
	if opts.Title != defaults.Title {
		t.Errorf("Title: got %v, want %v", opts.Title, defaults.Title)
	}
	if opts.Theme != defaults.Theme {
		t.Errorf("Theme: got %v, want %v", opts.Theme, defaults.Theme)
	}
	if opts.Width != defaults.Width {
		t.Errorf("Width: got %v, want %v", opts.Width, defaults.Width)
	}
	if opts.ParserMode != defaults.ParserMode {
		t.Errorf("ParserMode: got %v, want %v", opts.ParserMode, defaults.ParserMode)
	}
}

// TestRequestParamsInterface 验证两个请求类型都实现了 RequestParams 接口
func TestRequestParamsInterface(t *testing.T) {
	var _ RequestParams = (*ConvertRequest)(nil)
	var _ RequestParams = (*UploadRequest)(nil)
}
