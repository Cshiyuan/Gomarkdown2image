package utils

import (
	"fmt"
	"strings"

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

// ValidateCustomCSS 验证自定义 CSS，防止 XSS 注入
func ValidateCustomCSS(css string) error {
	if css == "" {
		return nil // 空 CSS 是允许的
	}

	// 检查 CSS 长度限制（防止 DoS）
	const maxCSSLength = 100000 // 100KB
	if len(css) > maxCSSLength {
		return fmt.Errorf("自定义 CSS 过大（最大 %d 字节）", maxCSSLength)
	}

	// 转换为小写以进行不区分大小写的检查
	lowerCSS := strings.ToLower(css)

	// 禁止的模式列表（防止 XSS 注入）
	forbiddenPatterns := []string{
		"</style>",    // 闭合 style 标签
		"<style",      // 嵌套 style 标签
		"<script",     // 注入脚本
		"<iframe",     // 注入 iframe
		"javascript:", // JavaScript 协议
		"data:",       // data URI（可能包含脚本）
		"vbscript:",   // VBScript 协议
		"<object",     // 对象标签
		"<embed",      // 嵌入标签
		"onerror",     // 事件处理器
		"onload",      // 事件处理器
		"onclick",     // 事件处理器
	}

	for _, pattern := range forbiddenPatterns {
		if strings.Contains(lowerCSS, pattern) {
			return fmt.Errorf("无效的 CSS：包含禁止的模式 '%s'", pattern)
		}
	}

	// 检查是否包含不配对的引号或括号（可能的注入尝试）
	if strings.Count(css, `"`)%2 != 0 || strings.Count(css, `'`)%2 != 0 {
		return fmt.Errorf("无效的 CSS：引号不配对")
	}

	return nil
}
