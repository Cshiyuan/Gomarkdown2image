package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Cshiyuan/Gomarkdown2image/internal/config"
	"github.com/Cshiyuan/Gomarkdown2image/internal/utils"
	"github.com/Cshiyuan/Gomarkdown2image/pkg/converter"
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

	// 验证 Markdown 内容大小
	if len(req.Markdown) > config.MaxMarkdownSize {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "CONTENT_TOO_LARGE",
				Message: "Markdown 内容过大",
				Details: fmt.Sprintf("最大支持 %d MB", config.MaxMarkdownSize/(1024*1024)),
			},
		})
		return
	}

	// 验证自定义 CSS (防止 XSS 注入)
	if err := utils.ValidateCustomCSS(req.CustomCSS); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "INVALID_CUSTOM_CSS",
				Message: "自定义 CSS 验证失败",
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
	contentType := utils.GetContentType(opts.ImageFormat)
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

	// 验证文件大小
	if file.Size > config.MaxFileUploadSize {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "FILE_TOO_LARGE",
				Message: "文件过大",
				Details: fmt.Sprintf("最大支持 %d MB", config.MaxFileUploadSize/(1024*1024)),
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

	// 验证自定义 CSS (防止 XSS 注入)
	if err := utils.ValidateCustomCSS(formReq.CustomCSS); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "INVALID_CUSTOM_CSS",
				Message: "自定义 CSS 验证失败",
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
	contentType := utils.GetContentType(opts.ImageFormat)
	c.Data(http.StatusOK, contentType, imageData)
}

// buildConvertOptionsFromParams 从 RequestParams 接口构建 ConvertOptions
// 这是统一的构建函数,消除了 buildConvertOptions 和 buildConvertOptionsFromForm 的代码重复
func buildConvertOptionsFromParams(params RequestParams) *converter.ConvertOptions {
	opts := converter.DefaultConvertOptions()

	// HTML 模板选项
	if v := params.GetTitle(); v != "" {
		opts.Title = v
	}
	if v := params.GetTheme(); v != "" {
		opts.Theme = v
	}
	if v := params.GetCustomCSS(); v != "" {
		opts.CustomCSS = v
	}
	if v := params.GetWidth(); v > 0 {
		opts.Width = v
	}
	if v := params.GetFontSize(); v > 0 {
		opts.FontSize = v
	}
	if v := params.GetFontFamily(); v != "" {
		opts.FontFamily = v
	}

	// 图像渲染选项
	if v := params.GetImageFormat(); v != "" {
		opts.ImageFormat = utils.ParseImageFormatOrDefault(v)
	}
	if v := params.GetImageQuality(); v > 0 {
		opts.ImageQuality = v
	}
	if v := params.GetDevicePixelRatio(); v > 0 {
		opts.DevicePixelRatio = v
	}

	// AI 增强选项
	if v := params.GetParserMode(); v != "" {
		opts.ParserMode = v
	}
	if v := params.GetAIProvider(); v != "" {
		opts.AIProvider = v
	}
	if v := params.GetAIModel(); v != "" {
		opts.AIModel = v
	}
	if v := params.GetAIAPIKey(); v != "" {
		opts.AIAPIKey = v
	}
	if v := params.GetAIEndpoint(); v != "" {
		opts.AIEndpoint = v
	}
	if v := params.GetAIPromptTemplate(); v != "" {
		opts.AIPromptTemplate = v
	}
	if v := params.GetAICustomPrompt(); v != "" {
		opts.AICustomPrompt = v
	}

	return opts
}

// buildConvertOptions 从 ConvertRequest 构建 ConvertOptions
func buildConvertOptions(req *ConvertRequest) *converter.ConvertOptions {
	return buildConvertOptionsFromParams(req)
}

// buildConvertOptionsFromForm 从 UploadRequest 构建 ConvertOptions
func buildConvertOptionsFromForm(req *UploadRequest) *converter.ConvertOptions {
	return buildConvertOptionsFromParams(req)
}
