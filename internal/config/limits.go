package config

// 文件大小限制常量
const (
	// MaxMarkdownSize Markdown 内容最大大小 (10MB)
	MaxMarkdownSize = 10 * 1024 * 1024

	// MaxFileUploadSize 文件上传最大大小 (10MB)
	MaxFileUploadSize = 10 * 1024 * 1024

	// MaxMultipartMemory 多部分表单内存限制 (10MB)
	MaxMultipartMemory = 10 << 20
)

// 参数范围限制常量
const (
	// MinWidth 最小页面宽度 (像素)
	MinWidth = 200

	// MaxWidth 最大页面宽度 (像素)
	MaxWidth = 4000

	// MinFontSize 最小字体大小 (px)
	MinFontSize = 8

	// MaxFontSize 最大字体大小 (px)
	MaxFontSize = 72

	// MinQuality 最小图片质量
	MinQuality = 1

	// MaxQuality 最大图片质量
	MaxQuality = 100

	// MinDevicePixelRatio 最小设备像素比
	MinDevicePixelRatio = 0.5

	// MaxDevicePixelRatio 最大设备像素比
	MaxDevicePixelRatio = 4.0
)
