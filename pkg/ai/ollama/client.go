// Package ollama 提供 Ollama 本地模型客户端实现
package ollama

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Cshiyuan/Gomarkdown2image/pkg/ai"
	"github.com/ollama/ollama/api"
)

// Provider Ollama 提供器实现
type Provider struct {
	client  *api.Client
	config  *ai.Config
	timeout time.Duration
}

// New 创建 Ollama Provider 实例(包外可见)
func New(cfg *ai.Config) (ai.Provider, error) {
	return newOllamaProvider(cfg)
}

// newOllamaProvider 创建 Ollama Provider 实例(包内使用)
func newOllamaProvider(cfg *ai.Config) (*Provider, error) {
	if err := ai.ValidateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// 设置 Ollama 服务地址
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434" // 默认地址
	}

	// 创建客户端
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, ai.NewError(ai.ErrorTypeNetwork, "failed to create Ollama client", err)
	}

	// 设置自定义 BaseURL(如果提供)
	if cfg.BaseURL != "" {
		// Ollama API 客户端会自动从环境变量或默认地址创建
		// 如需自定义,可以通过环境变量 OLLAMA_HOST 设置
	}

	// 设置超时
	timeout := time.Duration(cfg.Timeout) * time.Second
	if cfg.Timeout == 0 {
		timeout = 60 * time.Second // Ollama 本地模型可能需要更长时间
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

	// 构建系统提示词
	systemPrompt := req.System
	if systemPrompt == "" && p.config.Prompts != nil {
		systemPrompt = p.config.Prompts.DefaultSystem
	}

	// 构建完整提示词
	fullPrompt := req.Prompt
	if systemPrompt != "" {
		fullPrompt = systemPrompt + "\n\n" + req.Prompt
	}

	// 构建生成请求
	generateReq := &api.GenerateRequest{
		Model:  p.config.Model,
		Prompt: fullPrompt,
	}

	// 设置可选参数
	if req.Temperature > 0 {
		generateReq.Options = map[string]interface{}{
			"temperature": req.Temperature,
		}
	}

	// 调用 Ollama API
	var responseText strings.Builder
	var totalTokens int

	err := p.client.Generate(ctx, generateReq, func(resp api.GenerateResponse) error {
		responseText.WriteString(resp.Response)
		if resp.Done {
			totalTokens = resp.EvalCount + resp.PromptEvalCount
		}
		return nil
	})

	if err != nil {
		return nil, p.handleError(err)
	}

	// 构建响应
	response := &ai.GenerateResponse{
		Content:      responseText.String(),
		TokensUsed:   totalTokens,
		FinishReason: "stop",
		Metadata:     make(map[string]string),
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

	// 构建系统提示词
	systemPrompt := req.System
	if systemPrompt == "" && p.config.Prompts != nil {
		systemPrompt = p.config.Prompts.DefaultSystem
	}

	// 构建完整提示词
	fullPrompt := req.Prompt
	if systemPrompt != "" {
		fullPrompt = systemPrompt + "\n\n" + req.Prompt
	}

	// 构建生成请求
	generateReq := &api.GenerateRequest{
		Model:  p.config.Model,
		Prompt: fullPrompt,
	}

	// 设置可选参数
	if req.Temperature > 0 {
		generateReq.Options = map[string]interface{}{
			"temperature": req.Temperature,
		}
	}

	// 创建流式响应通道
	streamChan := make(chan ai.StreamChunk, 10)

	// 启动 goroutine 处理流式响应
	go func() {
		defer cancel()
		defer close(streamChan)

		// 调用 Ollama 流式 API
		err := p.client.Generate(ctx, generateReq, func(resp api.GenerateResponse) error {
			// 发送内容块
			if resp.Response != "" {
				streamChan <- ai.StreamChunk{
					Content: resp.Response,
					Done:    resp.Done,
					Error:   nil,
				}
			}

			// 如果完成,发送完成信号
			if resp.Done {
				streamChan <- ai.StreamChunk{
					Content: "",
					Done:    true,
					Error:   nil,
				}
			}

			return nil
		})

		if err != nil {
			streamChan <- ai.StreamChunk{
				Content: "",
				Done:    true,
				Error:   p.handleError(err),
			}
		}
	}()

	return streamChan, nil
}

// Name 返回 Provider 名称
func (p *Provider) Name() string {
	return "ollama"
}

// Close 关闭连接和清理资源
func (p *Provider) Close() error {
	// Ollama 客户端无需显式关闭
	return nil
}

// handleError 处理 Ollama API 错误并转换为统一的 AI 错误
func (p *Provider) handleError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	errMsgLower := strings.ToLower(errMsg)

	// 连接错误
	if strings.Contains(errMsgLower, "connection refused") ||
		strings.Contains(errMsgLower, "dial tcp") ||
		strings.Contains(errMsgLower, "no such host") {
		return ai.NewError(ai.ErrorTypeNetwork, "cannot connect to Ollama service", err)
	}

	// 模型不存在
	if strings.Contains(errMsgLower, "model") && strings.Contains(errMsgLower, "not found") {
		return ai.NewError(ai.ErrorTypeInvalidReq, "model not found", err)
	}

	// 超时
	if strings.Contains(errMsgLower, "timeout") ||
		strings.Contains(errMsgLower, "deadline exceeded") ||
		strings.Contains(errMsgLower, "context deadline") {
		return ai.NewError(ai.ErrorTypeTimeout, "request timeout", err)
	}

	// 网络错误
	if strings.Contains(errMsgLower, "connection") ||
		strings.Contains(errMsgLower, "network") {
		return ai.NewError(ai.ErrorTypeNetwork, "network error", err)
	}

	// 默认为未知错误
	return ai.NewError(ai.ErrorTypeUnknown, "unknown error", err)
}
