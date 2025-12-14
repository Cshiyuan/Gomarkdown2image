package handlers

import (
	"Gomarkdown2image/pkg/converter"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ConvertHandler 处理 JSON 方式的 Markdown 转换
// @Summary 转换 Markdown 为图片
// @Description 接收 JSON 格式的 Markdown 内容,返回生成的图片
// @Accept json
// @Produce image/png,image/jpeg,image/webp
// @Param request body ConvertRequest true "转换请求"
// @Success 200 {file} binary "生成的图片"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/convert [post]
func ConvertHandler(c *gin.Context) {
	var req ConvertRequest

	// 绑定并验证 JSON 请求
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "INVALID_REQUEST",
				Message: "请求参数验证失败",
				Details: err.Error(),
			},
		})
		return
	}

	// 验证 Markdown 内容大小 (限制 10MB)
	const maxMarkdownSize = 10 * 1024 * 1024
	if len(req.Markdown) > maxMarkdownSize {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "CONTENT_TOO_LARGE",
				Message: "Markdown 内容过大",
				Details: fmt.Sprintf("最大支持 %d MB", maxMarkdownSize/(1024*1024)),
			},
		})
		return
	}

	// 创建转换器
	conv, err := converter.NewConverter()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "CONVERTER_INIT_FAILED",
				Message: "转换器初始化失败",
				Details: err.Error(),
			},
		})
		return
	}
	defer conv.Close()

	// 构建转换选项
	opts := buildConvertOptions(&req)

	// 执行转换
	imageData, err := conv.Convert([]byte(req.Markdown), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "CONVERSION_FAILED",
				Message: "Markdown 转换失败",
				Details: err.Error(),
			},
		})
		return
	}

	// 返回图片
	contentType := getContentType(string(opts.ImageFormat))
	c.Data(http.StatusOK, contentType, imageData)
}

// UploadHandler 处理文件上传方式的 Markdown 转换
// @Summary 上传 Markdown 文件并转换为图片
// @Description 接收 Markdown 文件上传,返回生成的图片
// @Accept multipart/form-data
// @Produce image/png,image/jpeg,image/webp
// @Param file formData file true "Markdown 文件"
// @Param theme formData string false "主题 (light/dark)"
// @Param width formData int false "页面宽度"
// @Param fontSize formData int false "字体大小"
// @Param imageFormat formData string false "图片格式 (png/jpeg/webp)"
// @Param imageQuality formData int false "图片质量 (1-100)"
// @Success 200 {file} binary "生成的图片"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/upload [post]
func UploadHandler(c *gin.Context) {
	var formReq UploadRequest

	// 绑定并验证表单参数
	if err := c.ShouldBind(&formReq); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "INVALID_FORM",
				Message: "表单参数验证失败",
				Details: err.Error(),
			},
		})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "NO_FILE_UPLOADED",
				Message: "未找到上传文件",
				Details: "请在表单中上传文件 (字段名: file)",
			},
		})
		return
	}

	// 验证文件大小 (限制 10MB)
	const maxFileSize = 10 * 1024 * 1024
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "FILE_TOO_LARGE",
				Message: "文件过大",
				Details: fmt.Sprintf("最大支持 %d MB", maxFileSize/(1024*1024)),
			},
		})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "FILE_READ_FAILED",
				Message: "文件读取失败",
				Details: err.Error(),
			},
		})
		return
	}
	defer src.Close()

	// 读取文件内容
	markdownData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "FILE_READ_FAILED",
				Message: "文件内容读取失败",
				Details: err.Error(),
			},
		})
		return
	}

	// 创建转换器
	conv, err := converter.NewConverter()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "CONVERTER_INIT_FAILED",
				Message: "转换器初始化失败",
				Details: err.Error(),
			},
		})
		return
	}
	defer conv.Close()

	// 构建转换选项 (从表单参数)
	opts := buildConvertOptionsFromForm(&formReq)

	// 执行转换
	imageData, err := conv.Convert(markdownData, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "CONVERSION_FAILED",
				Message: "Markdown 转换失败",
				Details: err.Error(),
			},
		})
		return
	}

	// 返回图片
	contentType := getContentType(string(opts.ImageFormat))
	c.Data(http.StatusOK, contentType, imageData)
}

// buildConvertOptions 从 ConvertRequest 构建 ConvertOptions
func buildConvertOptions(req *ConvertRequest) *converter.ConvertOptions {
	opts := converter.DefaultConvertOptions()

	// 应用用户指定的选项
	if req.Title != "" {
		opts.Title = req.Title
	}
	if req.Theme != "" {
		opts.Theme = req.Theme
	}
	if req.CustomCSS != "" {
		opts.CustomCSS = req.CustomCSS
	}
	if req.Width > 0 {
		opts.Width = req.Width
	}
	if req.FontSize > 0 {
		opts.FontSize = req.FontSize
	}
	if req.FontFamily != "" {
		opts.FontFamily = req.FontFamily
	}
	if req.ImageFormat != "" {
		opts.ImageFormat = ParseImageFormat(req.ImageFormat)
	}
	if req.ImageQuality > 0 {
		opts.ImageQuality = req.ImageQuality
	}
	if req.DevicePixelRatio > 0 {
		opts.DevicePixelRatio = req.DevicePixelRatio
	}

	return opts
}

// buildConvertOptionsFromForm 从 UploadRequest 构建 ConvertOptions
func buildConvertOptionsFromForm(req *UploadRequest) *converter.ConvertOptions {
	opts := converter.DefaultConvertOptions()

	// 应用用户指定的选项
	if req.Title != "" {
		opts.Title = req.Title
	}
	if req.Theme != "" {
		opts.Theme = req.Theme
	}
	if req.CustomCSS != "" {
		opts.CustomCSS = req.CustomCSS
	}
	if req.Width > 0 {
		opts.Width = req.Width
	}
	if req.FontSize > 0 {
		opts.FontSize = req.FontSize
	}
	if req.FontFamily != "" {
		opts.FontFamily = req.FontFamily
	}
	if req.ImageFormat != "" {
		opts.ImageFormat = ParseImageFormat(req.ImageFormat)
	}
	if req.ImageQuality > 0 {
		opts.ImageQuality = req.ImageQuality
	}
	if req.DevicePixelRatio > 0 {
		opts.DevicePixelRatio = req.DevicePixelRatio
	}

	return opts
}

// getContentType 根据图片格式获取 Content-Type
func getContentType(format string) string {
	switch format {
	case "jpeg", "jpg":
		return "image/jpeg"
	case "webp":
		return "image/webp"
	default:
		return "image/png"
	}
}
