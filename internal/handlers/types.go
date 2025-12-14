package handlers

// ConvertRequest 表示 /api/convert 端点的请求体
type ConvertRequest struct {
	// Markdown 内容 (必需)
	Markdown string `json:"markdown" binding:"required"`

	// HTML 模板选项
	Title      string `json:"title,omitempty"`                                      // 页面标题
	Theme      string `json:"theme,omitempty" binding:"omitempty,oneof=light dark"` // 主题: light/dark
	CustomCSS  string `json:"customCss,omitempty"`                                  // 自定义 CSS
	Width      int    `json:"width,omitempty" binding:"omitempty,min=200,max=4000"` // 页面宽度
	FontSize   int    `json:"fontSize,omitempty" binding:"omitempty,min=8,max=72"`  // 字体大小
	FontFamily string `json:"fontFamily,omitempty"`                                 // 字体族

	// 图像渲染选项
	ImageFormat      string  `json:"imageFormat,omitempty" binding:"omitempty,oneof=png jpeg webp"` // 图片格式
	ImageQuality     int     `json:"imageQuality,omitempty" binding:"omitempty,min=1,max=100"`      // 图片质量 (1-100)
	DevicePixelRatio float64 `json:"devicePixelRatio,omitempty" binding:"omitempty,min=0.5,max=4"`  // 设备像素比
}

// UploadRequest 表示 /api/upload 端点的表单字段
type UploadRequest struct {
	// 所有字段都通过表单 (multipart/form-data) 提交
	// 文件字段名: "file"
	Title      string `form:"title"`
	Theme      string `form:"theme" binding:"omitempty,oneof=light dark"`
	Width      int    `form:"width" binding:"omitempty,min=200,max=4000"`
	FontSize   int    `form:"fontSize" binding:"omitempty,min=8,max=72"`
	FontFamily string `form:"fontFamily"`
	CustomCSS  string `form:"customCss"`

	ImageFormat      string  `form:"imageFormat" binding:"omitempty,oneof=png jpeg webp"`
	ImageQuality     int     `form:"imageQuality" binding:"omitempty,min=1,max=100"`
	DevicePixelRatio float64 `form:"devicePixelRatio" binding:"omitempty,min=0.5,max=4"`
}

// APIResponse 统一的 API 响应格式
type APIResponse struct {
	Success bool        `json:"success"`           // 是否成功
	Message string      `json:"message,omitempty"` // 成功消息
	Data    interface{} `json:"data,omitempty"`    // 响应数据
	Error   *APIError   `json:"error,omitempty"`   // 错误信息
}

// APIError API 错误详情
type APIError struct {
	Code    string `json:"code"`              // 错误代码
	Message string `json:"message"`           // 错误消息
	Details string `json:"details,omitempty"` // 错误详情
}

// ConvertResponse 转换成功时的响应数据
type ConvertResponse struct {
	Format string `json:"format"` // 图片格式
	Size   int    `json:"size"`   // 文件大小 (字节)
	Width  int    `json:"width"`  // 图片宽度 (像素)
}
