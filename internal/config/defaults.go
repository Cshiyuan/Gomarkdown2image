package config

import "github.com/Cshiyuan/Gomarkdown2image/pkg/renderer"

// 默认配置值常量
const (
	// DefaultTitle 默认页面标题
	DefaultTitle = "Markdown to Image"

	// DefaultTheme 默认主题
	DefaultTheme = "light"

	// DefaultWidth 默认页面宽度 (像素)
	DefaultWidth = 1200

	// DefaultFontSize 默认字体大小 (px)
	DefaultFontSize = 16

	// DefaultFontFamily 默认字体族
	DefaultFontFamily = "Arial, sans-serif"

	// DefaultQuality 默认图片质量
	DefaultQuality = 90

	// DefaultDevicePixelRatio 默认设备像素比
	DefaultDevicePixelRatio = 1.0
)

// DefaultImageFormat 返回默认图片格式
func DefaultImageFormat() renderer.ImageFormat {
	return renderer.FormatPNG
}
