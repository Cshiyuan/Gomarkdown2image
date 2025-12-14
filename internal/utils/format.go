package utils

import (
	"fmt"
	"strings"

	"github.com/Cshiyuan/Gomarkdown2image/pkg/renderer"
)

// ParseImageFormat 解析图片格式字符串
// 统一的实现,消除代码重复
func ParseImageFormat(format string) (renderer.ImageFormat, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "png":
		return renderer.FormatPNG, nil
	case "jpeg", "jpg":
		return renderer.FormatJPEG, nil
	case "webp":
		return renderer.FormatWebP, nil
	default:
		return "", fmt.Errorf("不支持的图片格式: %s (支持: png, jpeg, webp)", format)
	}
}

// ParseImageFormatOrDefault 解析图片格式,失败时返回默认格式
func ParseImageFormatOrDefault(format string) renderer.ImageFormat {
	result, err := ParseImageFormat(format)
	if err != nil {
		return renderer.FormatPNG
	}
	return result
}

// GetContentType 根据图片格式获取 MIME 类型
func GetContentType(format renderer.ImageFormat) string {
	switch format {
	case renderer.FormatJPEG:
		return "image/jpeg"
	case renderer.FormatWebP:
		return "image/webp"
	default:
		return "image/png"
	}
}
