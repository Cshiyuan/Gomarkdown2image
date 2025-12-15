package utils

import (
	"testing"

	"github.com/Cshiyuan/Gomarkdown2image/internal/config"
)

func TestValidateQuality(t *testing.T) {
	tests := []struct {
		name    string
		quality int
		wantErr bool
	}{
		{"最小有效值", config.MinQuality, false},
		{"最大有效值", config.MaxQuality, false},
		{"中间值", 50, false},
		{"默认值", 90, false},
		{"低于最小值", config.MinQuality - 1, true},
		{"高于最大值", config.MaxQuality + 1, true},
		{"零值", 0, true},
		{"负数", -10, true},
		{"超大值", 200, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuality(tt.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQuality(%d) error = %v, wantErr %v", tt.quality, err, tt.wantErr)
			}
		})
	}
}

func TestValidateWidth(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		wantErr bool
	}{
		{"最小有效宽度", config.MinWidth, false},
		{"最大有效宽度", config.MaxWidth, false},
		{"默认宽度", 1200, false},
		{"低于最小值", config.MinWidth - 1, true},
		{"高于最大值", config.MaxWidth + 1, true},
		{"零值", 0, true},
		{"负数", -100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWidth(tt.width)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWidth(%d) error = %v, wantErr %v", tt.width, err, tt.wantErr)
			}
		})
	}
}

func TestValidateFontSize(t *testing.T) {
	tests := []struct {
		name     string
		fontSize int
		wantErr  bool
	}{
		{"最小字体", config.MinFontSize, false},
		{"最大字体", config.MaxFontSize, false},
		{"默认字体", 16, false},
		{"中等字体", 24, false},
		{"低于最小值", config.MinFontSize - 1, true},
		{"高于最大值", config.MaxFontSize + 1, true},
		{"零值", 0, true},
		{"负数", -5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFontSize(tt.fontSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFontSize(%d) error = %v, wantErr %v", tt.fontSize, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDevicePixelRatio(t *testing.T) {
	tests := []struct {
		name    string
		dpr     float64
		wantErr bool
	}{
		{"最小DPR", config.MinDevicePixelRatio, false},
		{"最大DPR", config.MaxDevicePixelRatio, false},
		{"默认DPR", 1.0, false},
		{"高清屏", 2.0, false},
		{"超高清", 3.0, false},
		{"低于最小值", 0.4, true},
		{"高于最大值", 4.1, true},
		{"零值", 0.0, true},
		{"负数", -1.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDevicePixelRatio(tt.dpr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDevicePixelRatio(%v) error = %v, wantErr %v", tt.dpr, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTheme(t *testing.T) {
	tests := []struct {
		name    string
		theme   string
		wantErr bool
	}{
		{"亮色主题", "light", false},
		{"暗色主题", "dark", false},
		{"无效主题", "blue", true},
		{"空字符串", "", true},
		{"大写主题", "LIGHT", true},
		{"混合大小写", "Light", true},
		{"不存在的主题", "solarized", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTheme(tt.theme)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTheme(%q) error = %v, wantErr %v", tt.theme, err, tt.wantErr)
			}
		})
	}
}

// 基准测试
func BenchmarkValidateQuality(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateQuality(90)
	}
}

func BenchmarkValidateWidth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateWidth(1200)
	}
}

func BenchmarkValidateTheme(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateTheme("light")
	}
}
