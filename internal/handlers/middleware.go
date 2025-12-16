package handlers

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS 配置 CORS 中间件
func SetupCORS() gin.HandlerFunc {
	// 从环境变量读取允许的域名
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		// 默认值：开发环境
		allowedOrigins = "http://localhost:3000,http://localhost:8080"
	}

	origins := strings.Split(allowedOrigins, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false, // 除非必要，否则关闭
		MaxAge:           12 * time.Hour,
	})
}

// RequestLogger 自定义请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// 使用 Gin 的日志格式 (修复: 使用 strconv.Itoa 而非 string(rune()))
		gin.DefaultWriter.Write([]byte(
			time.Now().Format("2006/01/02 - 15:04:05") +
				" | " + method +
				" | " + path +
				" | " + clientIP +
				" | " + latency.String() +
				" | " + strconv.Itoa(statusCode) + "\n",
		))
	}
}

// ErrorRecovery 错误恢复中间件
func ErrorRecovery() gin.HandlerFunc {
	return gin.Recovery()
}

// HealthCheckHandler 健康检查处理器
// @Summary 健康检查
// @Description 返回 API 服务健康状态
// @Produce json
// @Success 200 {object} APIResponse
// @Router /health [get]
func HealthCheckHandler(c *gin.Context) {
	c.JSON(200, APIResponse{
		Success: true,
		Message: "服务运行正常",
		Data: gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		},
	})
}
