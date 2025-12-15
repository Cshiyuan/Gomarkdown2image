package utils

import (
	"testing"

	"github.com/Cshiyuan/Gomarkdown2image/pkg/renderer"
)

func TestParseImageFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    renderer.ImageFormat
		wantErr bool
	}{
		{
			name:    "png小写",
			input:   "png",
			want:    renderer.FormatPNG,
			wantErr: false,
		},
		{
			name:    "PNG大写",
			input:   "PNG",
			want:    renderer.FormatPNG,
			wantErr: false,
		},
		{
			name:    "jpeg",
			input:   "jpeg",
			want:    renderer.FormatJPEG,
			wantErr: false,
		},
		{
			name:    "jpg",
			input:   "jpg",
			want:    renderer.FormatJPEG,
			wantErr: false,
		},
		{
			name:    "JPEG大写",
			input:   "JPEG",
			want:    renderer.FormatJPEG,
			wantErr: false,
		},
		{
			name:    "webp",
			input:   "webp",
			want:    renderer.FormatWebP,
			wantErr: false,
		},
		{
			name:    "WebP混合大小写",
			input:   "WebP",
			want:    renderer.FormatWebP,
			wantErr: false,
		},
		{
			name:    "带空格的格式",
			input:   "  png  ",
			want:    renderer.FormatPNG,
			wantErr: false,
		},
		{
			name:    "不支持的格式",
			input:   "gif",
			want:    "",
			wantErr: true,
		},
		{
			name:    "空字符串",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "无效格式",
			input:   "invalid",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseImageFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseImageFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseImageFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseImageFormatOrDefault(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  renderer.ImageFormat
	}{
		{
			name:  "有效格式png",
			input: "png",
			want:  renderer.FormatPNG,
		},
		{
			name:  "有效格式jpeg",
			input: "jpeg",
			want:  renderer.FormatJPEG,
		},
		{
			name:  "有效格式webp",
			input: "webp",
			want:  renderer.FormatWebP,
		},
		{
			name:  "无效格式返回默认",
			input: "gif",
			want:  renderer.FormatPNG,
		},
		{
			name:  "空字符串返回默认",
			input: "",
			want:  renderer.FormatPNG,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseImageFormatOrDefault(tt.input)
			if got != tt.want {
				t.Errorf("ParseImageFormatOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetContentType(t *testing.T) {
	tests := []struct {
		name   string
		format renderer.ImageFormat
		want   string
	}{
		{
			name:   "PNG格式",
			format: renderer.FormatPNG,
			want:   "image/png",
		},
		{
			name:   "JPEG格式",
			format: renderer.FormatJPEG,
			want:   "image/jpeg",
		},
		{
			name:   "WebP格式",
			format: renderer.FormatWebP,
			want:   "image/webp",
		},
		{
			name:   "未知格式默认PNG",
			format: "unknown",
			want:   "image/png",
		},
		{
			name:   "空格式默认PNG",
			format: "",
			want:   "image/png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetContentType(tt.format)
			if got != tt.want {
				t.Errorf("GetContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 基准测试
func BenchmarkParseImageFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseImageFormat("png")
	}
}

func BenchmarkGetContentType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetContentType(renderer.FormatPNG)
	}
}
