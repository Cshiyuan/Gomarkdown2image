package utils

import (
	"fmt"

	"github.com/Cshiyuan/Gomarkdown2image/internal/config"
)

// ValidateQuality 验证图片质量参数
func ValidateQuality(quality int) error {
	if quality < config.MinQuality || quality > config.MaxQuality {
		return fmt.Errorf("图片质量必须在 %d-%d 之间", config.MinQuality, config.MaxQuality)
	}
	return nil
}

// ValidateWidth 验证页面宽度参数
func ValidateWidth(width int) error {
	if width < config.MinWidth || width > config.MaxWidth {
		return fmt.Errorf("页面宽度必须在 %d-%d 之间", config.MinWidth, config.MaxWidth)
	}
	return nil
}

// ValidateFontSize 验证字体大小参数
func ValidateFontSize(fontSize int) error {
	if fontSize < config.MinFontSize || fontSize > config.MaxFontSize {
		return fmt.Errorf("字体大小必须在 %d-%d 之间", config.MinFontSize, config.MaxFontSize)
	}
	return nil
}

// ValidateDevicePixelRatio 验证设备像素比参数
func ValidateDevicePixelRatio(dpr float64) error {
	if dpr < config.MinDevicePixelRatio || dpr > config.MaxDevicePixelRatio {
		return fmt.Errorf("设备像素比必须在 %.1f-%.1f 之间", config.MinDevicePixelRatio, config.MaxDevicePixelRatio)
	}
	return nil
}

// ValidateTheme 验证主题参数
func ValidateTheme(theme string) error {
	validThemes := []string{"light", "dark"}
	for _, valid := range validThemes {
		if theme == valid {
			return nil
		}
	}
	return fmt.Errorf("无效的主题: %s (支持: light, dark)", theme)
}
