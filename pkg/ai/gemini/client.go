// Package gemini 提供 Google Gemini API 客户端实现
package gemini

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Cshiyuan/Gomarkdown2image/pkg/ai"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Provider Gemini API 提供器实现
type Provider struct {
	client  *genai.Client
	config  *ai.Config
	timeout time.Duration
}

// New 创建 Gemini Provider 实例(包外可见)
func New(cfg *ai.Config) (ai.Provider, error) {
	return newGeminiProvider(cfg)
}

// newGeminiProvider 创建 Gemini Provider 实例(包内使用)
func newGeminiProvider(cfg *ai.Config) (*Provider, error) {
	if err := ai.ValidateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	ctx := context.Background()

	// 构建客户端选项
	opts := []option.ClientOption{
		option.WithAPIKey(cfg.APIKey),
	}

	// 创建客户端
	client, err := genai.NewClient(ctx, opts...)
	if err != nil {
		return nil, ai.NewError(ai.ErrorTypeAuth, "failed to create Gemini client", err)
	}

	// 设置超时
	timeout := time.Duration(cfg.Timeout) * time.Second
	if cfg.Timeout == 0 {
		timeout = 30 * time.Second // 默认 30 秒
	}

	return &Provider{
		client:  client,
		config:  cfg,
		timeout: timeout,
	}, nil
}

// Generate 生成文本响应
func (p *Provider) Generate(ctx context.Context, req *ai.GenerateRequest) (*ai.GenerateResponse, error) {
	if err := ai.ValidateGenerateRequest(req); err != nil {
		return nil, ai.NewError(ai.ErrorTypeInvalidReq, "invalid request", err)
	}

	// 设置超时上下文
	if req.Context != nil {
		ctx = req.Context
	}
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	// 获取模型
	model := p.client.GenerativeModel(p.config.Model)

	// 设置生成参数
	if req.MaxTokens > 0 {
		model.SetMaxOutputTokens(int32(req.MaxTokens))
	}

	if req.Temperature > 0 {
		model.SetTemperature(float32(req.Temperature))
	}

	// 设置系统指令
	systemPrompt := req.System
	if systemPrompt == "" && p.config.Prompts != nil {
		systemPrompt = p.config.Prompts.DefaultSystem
	}
	if systemPrompt != "" {
		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(systemPrompt)},
		}
	}

	// 生成内容
	resp, err := model.GenerateContent(ctx, genai.Text(req.Prompt))
	if err != nil {
		return nil, p.handleError(err)
	}

	// 提取响应内容
	if resp == nil || len(resp.Candidates) == 0 {
		return nil, ai.NewError(ai.ErrorTypeServerError, "empty response from Gemini", nil)
	}

	candidate := resp.Candidates[0]
	var content strings.Builder

	if candidate.Content != nil {
		for _, part := range candidate.Content.Parts {
			if text, ok := part.(genai.Text); ok {
				content.WriteString(string(text))
			}
		}
	}

	// 构建响应
	response := &ai.GenerateResponse{
		Content:      content.String(),
		FinishReason: string(candidate.FinishReason),
		Metadata:     make(map[string]string),
	}

	// 添加 token 使用信息
	if resp.UsageMetadata != nil {
		response.TokensUsed = int(resp.UsageMetadata.TotalTokenCount)
		response.Metadata["prompt_tokens"] = fmt.Sprintf("%d", resp.UsageMetadata.PromptTokenCount)
		response.Metadata["candidates_tokens"] = fmt.Sprintf("%d", resp.UsageMetadata.CandidatesTokenCount)
	}
	response.Metadata["model"] = p.config.Model

	return response, nil
}

// GenerateStream 流式生成文本响应
func (p *Provider) GenerateStream(ctx context.Context, req *ai.GenerateRequest) (<-chan ai.StreamChunk, error) {
	if err := ai.ValidateGenerateRequest(req); err != nil {
		return nil, ai.NewError(ai.ErrorTypeInvalidReq, "invalid request", err)
	}

	// 设置超时上下文
	if req.Context != nil {
		ctx = req.Context
	}
	ctx, cancel := context.WithTimeout(ctx, p.timeout)

	// 获取模型
	model := p.client.GenerativeModel(p.config.Model)

	// 设置生成参数
	if req.MaxTokens > 0 {
		model.SetMaxOutputTokens(int32(req.MaxTokens))
	}

	if req.Temperature > 0 {
		model.SetTemperature(float32(req.Temperature))
	}

	// 设置系统指令
	systemPrompt := req.System
	if systemPrompt == "" && p.config.Prompts != nil {
		systemPrompt = p.config.Prompts.DefaultSystem
	}
	if systemPrompt != "" {
		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(systemPrompt)},
		}
	}

	// 创建流式响应通道
	streamChan := make(chan ai.StreamChunk, 10)

	// 启动 goroutine 处理流式响应
	go func() {
		defer cancel()
		defer close(streamChan)

		// 调用 Gemini 流式 API
		iter := model.GenerateContentStream(ctx, genai.Text(req.Prompt))

		// 处理流式响应
		for {
			resp, err := iter.Next()
			if err != nil {
				// 检查是否是正常结束
				if err.Error() == "iterator is done" || strings.Contains(err.Error(), "EOF") {
					break
				}
				streamChan <- ai.StreamChunk{
					Content: "",
					Done:    true,
					Error:   p.handleError(err),
				}
				return
			}

			// 提取文本内容
			if resp != nil && len(resp.Candidates) > 0 {
				candidate := resp.Candidates[0]
				if candidate.Content != nil {
					for _, part := range candidate.Content.Parts {
						if text, ok := part.(genai.Text); ok {
							streamChan <- ai.StreamChunk{
								Content: string(text),
								Done:    false,
								Error:   nil,
							}
						}
					}
				}
			}
		}

		// 发送完成信号
		streamChan <- ai.StreamChunk{
			Content: "",
			Done:    true,
			Error:   nil,
		}
	}()

	return streamChan, nil
}

// Name 返回 Provider 名称
func (p *Provider) Name() string {
	return "gemini"
}

// Close 关闭连接和清理资源
func (p *Provider) Close() error {
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}

// handleError 处理 Gemini API 错误并转换为统一的 AI 错误
func (p *Provider) handleError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	errMsgLower := strings.ToLower(errMsg)

	// 认证错误
	if strings.Contains(errMsgLower, "authentication") ||
		strings.Contains(errMsgLower, "api key") ||
		strings.Contains(errMsgLower, "unauthorized") ||
		strings.Contains(errMsgLower, "401") {
		return ai.NewError(ai.ErrorTypeAuth, "authentication failed", err)
	}

	// 速率限制
	if strings.Contains(errMsgLower, "rate limit") ||
		strings.Contains(errMsgLower, "quota") ||
		strings.Contains(errMsgLower, "too many requests") ||
		strings.Contains(errMsgLower, "429") {
		return ai.NewError(ai.ErrorTypeRateLimit, "rate limit exceeded", err)
	}

	// 超时
	if strings.Contains(errMsgLower, "timeout") ||
		strings.Contains(errMsgLower, "deadline exceeded") ||
		strings.Contains(errMsgLower, "context deadline") {
		return ai.NewError(ai.ErrorTypeTimeout, "request timeout", err)
	}

	// 网络错误
	if strings.Contains(errMsgLower, "connection") ||
		strings.Contains(errMsgLower, "network") ||
		strings.Contains(errMsgLower, "dial") {
		return ai.NewError(ai.ErrorTypeNetwork, "network error", err)
	}

	// 服务器错误
	if strings.Contains(errMsgLower, "500") ||
		strings.Contains(errMsgLower, "502") ||
		strings.Contains(errMsgLower, "503") ||
		strings.Contains(errMsgLower, "504") ||
		strings.Contains(errMsgLower, "server error") ||
		strings.Contains(errMsgLower, "internal error") {
		return ai.NewError(ai.ErrorTypeServerError, "server error", err)
	}

	// 无效请求
	if strings.Contains(errMsgLower, "invalid") ||
		strings.Contains(errMsgLower, "bad request") ||
		strings.Contains(errMsgLower, "400") {
		return ai.NewError(ai.ErrorTypeInvalidReq, "invalid request", err)
	}

	// 默认为未知错误
	return ai.NewError(ai.ErrorTypeUnknown, "unknown error", err)
}
