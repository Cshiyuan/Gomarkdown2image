package ai

import "fmt"

// ErrorType 错误类型
type ErrorType string

const (
	// ErrorTypeAuth 认证错误 (API Key 无效等)
	ErrorTypeAuth ErrorType = "auth"

	// ErrorTypeRateLimit 速率限制错误
	ErrorTypeRateLimit ErrorType = "rate_limit"

	// ErrorTypeInvalidReq 无效请求错误 (参数错误等)
	ErrorTypeInvalidReq ErrorType = "invalid_req"

	// ErrorTypeServerError 服务器错误 (5xx)
	ErrorTypeServerError ErrorType = "server_error"

	// ErrorTypeTimeout 超时错误
	ErrorTypeTimeout ErrorType = "timeout"

	// ErrorTypeNetwork 网络错误
	ErrorTypeNetwork ErrorType = "network"

	// ErrorTypeUnknown 未知错误
	ErrorTypeUnknown ErrorType = "unknown"
)

// Error AI 服务错误
type Error struct {
	// Type 错误类型
	Type ErrorType

	// Message 错误消息
	Message string

	// StatusCode HTTP 状态码 (如适用)
	StatusCode int

	// Retryable 是否可重试
	Retryable bool

	// Original 原始错误
	Original error

	// Metadata 额外元数据
	Metadata map[string]string
}

// Error 实现 error 接口
func (e *Error) Error() string {
	if e.Original != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Original)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Unwrap 支持 errors.Unwrap
func (e *Error) Unwrap() error {
	return e.Original
}

// NewError 创建新的 AI 错误
func NewError(errType ErrorType, msg string, original error) *Error {
	return &Error{
		Type:      errType,
		Message:   msg,
		Retryable: isRetryable(errType),
		Original:  original,
		Metadata:  make(map[string]string),
	}
}

// NewErrorWithCode 创建带 HTTP 状态码的错误
func NewErrorWithCode(errType ErrorType, msg string, statusCode int, original error) *Error {
	return &Error{
		Type:       errType,
		Message:    msg,
		StatusCode: statusCode,
		Retryable:  isRetryable(errType),
		Original:   original,
		Metadata:   make(map[string]string),
	}
}

// isRetryable 判断错误类型是否可重试
func isRetryable(errType ErrorType) bool {
	switch errType {
	case ErrorTypeRateLimit, ErrorTypeServerError, ErrorTypeTimeout, ErrorTypeNetwork:
		return true
	default:
		return false
	}
}

// IsRetryable 判断错误是否可重试
func IsRetryable(err error) bool {
	if aiErr, ok := err.(*Error); ok {
		return aiErr.Retryable
	}
	return false
}
