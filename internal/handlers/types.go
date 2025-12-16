package handlers

// RequestParams 统一的请求参数接口
// 用于消除 buildConvertOptions 和 buildConvertOptionsFromForm 的代码重复
type RequestParams interface {
	GetTitle() string
	GetTheme() string
	GetCustomCSS() string
	GetWidth() int
	GetFontSize() int
	GetFontFamily() string
	GetImageFormat() string
	GetImageQuality() int
	GetDevicePixelRatio() float64
	GetParserMode() string
	GetAIProvider() string
	GetAIModel() string
	GetAIAPIKey() string
	GetAIEndpoint() string
	GetAIPromptTemplate() string
	GetAICustomPrompt() string
}

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

	// AI 增强选项 (新增)
	ParserMode       string `json:"parserMode,omitempty" binding:"omitempty,oneof=traditional ai"` // 解析器模式
	AIProvider       string `json:"aiProvider,omitempty" binding:"omitempty,oneof=gemini ollama"`  // AI 提供器
	AIModel          string `json:"aiModel,omitempty"`                                             // AI 模型名称
	AIAPIKey         string `json:"aiApiKey,omitempty"`                                            // AI API 密钥
	AIEndpoint       string `json:"aiEndpoint,omitempty"`                                          // AI 服务端点
	AIPromptTemplate string `json:"aiPromptTemplate,omitempty"`                                    // 提示词模板
	AICustomPrompt   string `json:"aiCustomPrompt,omitempty"`                                      // 自定义提示词
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

	// AI 增强选项 (新增)
	ParserMode       string `form:"parserMode" binding:"omitempty,oneof=traditional ai"`
	AIProvider       string `form:"aiProvider" binding:"omitempty,oneof=gemini ollama"`
	AIModel          string `form:"aiModel"`
	AIAPIKey         string `form:"aiApiKey"`
	AIEndpoint       string `form:"aiEndpoint"`
	AIPromptTemplate string `form:"aiPromptTemplate"`
	AICustomPrompt   string `form:"aiCustomPrompt"`
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

// ===== ConvertRequest 实现 RequestParams 接口 =====

func (r *ConvertRequest) GetTitle() string             { return r.Title }
func (r *ConvertRequest) GetTheme() string             { return r.Theme }
func (r *ConvertRequest) GetCustomCSS() string         { return r.CustomCSS }
func (r *ConvertRequest) GetWidth() int                { return r.Width }
func (r *ConvertRequest) GetFontSize() int             { return r.FontSize }
func (r *ConvertRequest) GetFontFamily() string        { return r.FontFamily }
func (r *ConvertRequest) GetImageFormat() string       { return r.ImageFormat }
func (r *ConvertRequest) GetImageQuality() int         { return r.ImageQuality }
func (r *ConvertRequest) GetDevicePixelRatio() float64 { return r.DevicePixelRatio }
func (r *ConvertRequest) GetParserMode() string        { return r.ParserMode }
func (r *ConvertRequest) GetAIProvider() string        { return r.AIProvider }
func (r *ConvertRequest) GetAIModel() string           { return r.AIModel }
func (r *ConvertRequest) GetAIAPIKey() string          { return r.AIAPIKey }
func (r *ConvertRequest) GetAIEndpoint() string        { return r.AIEndpoint }
func (r *ConvertRequest) GetAIPromptTemplate() string  { return r.AIPromptTemplate }
func (r *ConvertRequest) GetAICustomPrompt() string    { return r.AICustomPrompt }

// ===== UploadRequest 实现 RequestParams 接口 =====

func (r *UploadRequest) GetTitle() string             { return r.Title }
func (r *UploadRequest) GetTheme() string             { return r.Theme }
func (r *UploadRequest) GetCustomCSS() string         { return r.CustomCSS }
func (r *UploadRequest) GetWidth() int                { return r.Width }
func (r *UploadRequest) GetFontSize() int             { return r.FontSize }
func (r *UploadRequest) GetFontFamily() string        { return r.FontFamily }
func (r *UploadRequest) GetImageFormat() string       { return r.ImageFormat }
func (r *UploadRequest) GetImageQuality() int         { return r.ImageQuality }
func (r *UploadRequest) GetDevicePixelRatio() float64 { return r.DevicePixelRatio }
func (r *UploadRequest) GetParserMode() string        { return r.ParserMode }
func (r *UploadRequest) GetAIProvider() string        { return r.AIProvider }
func (r *UploadRequest) GetAIModel() string           { return r.AIModel }
func (r *UploadRequest) GetAIAPIKey() string          { return r.AIAPIKey }
func (r *UploadRequest) GetAIEndpoint() string        { return r.AIEndpoint }
func (r *UploadRequest) GetAIPromptTemplate() string  { return r.AIPromptTemplate }
func (r *UploadRequest) GetAICustomPrompt() string    { return r.AICustomPrompt }
