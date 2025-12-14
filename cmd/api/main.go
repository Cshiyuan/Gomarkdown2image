package main

import (
	"Gomarkdown2image/pkg/handlers"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	version = "0.1.0"
)

func main() {
	// 设置 Gin 模式 (可通过环境变量 GIN_MODE 控制)
	// 生产环境: export GIN_MODE=release
	// 开发环境: export GIN_MODE=debug (默认)

	// 创建路由器
	router := gin.New()

	// 应用中间件
	router.Use(handlers.ErrorRecovery())    // 错误恢复
	router.Use(handlers.SetupCORS())        // CORS 跨域
	router.Use(gin.Logger())                // Gin 内置日志
	router.Use(handlers.RequestLogger())    // 自定义日志

	// 设置最大上传文件大小 (10MB)
	router.MaxMultipartMemory = 10 << 20 // 10MB

	// 健康检查端点
	router.GET("/health", handlers.HealthCheckHandler)

	// API 路由组
	api := router.Group("/api")
	{
		// POST /api/convert - JSON 方式转换 Markdown
		api.POST("/convert", handlers.ConvertHandler)

		// POST /api/upload - 文件上传方式转换 Markdown
		api.POST("/upload", handlers.UploadHandler)
	}

	// 根路径欢迎信息
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "Gomarkdown2image API",
			"version": version,
			"endpoints": gin.H{
				"health":  "GET /health",
				"convert": "POST /api/convert",
				"upload":  "POST /api/upload",
			},
			"docs": "https://github.com/yourusername/Gomarkdown2image",
		})
	})

	// 获取端口 (默认 8080,可通过环境变量 PORT 指定)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务
	fmt.Printf("\n🚀 Gomarkdown2image API 服务启动中...\n")
	fmt.Printf("📡 监听端口: %s\n", port)
	fmt.Printf("🌍 访问地址: http://localhost:%s\n", port)
	fmt.Printf("💚 健康检查: http://localhost:%s/health\n", port)
	fmt.Printf("\n可用端点:\n")
	fmt.Printf("  POST http://localhost:%s/api/convert - JSON 转换\n", port)
	fmt.Printf("  POST http://localhost:%s/api/upload  - 文件上传\n", port)
	fmt.Printf("\n按 Ctrl+C 停止服务\n\n")

	if err := router.Run(":" + port); err != nil {
		fmt.Fprintf(os.Stderr, "❌ 服务启动失败: %v\n", err)
		os.Exit(1)
	}
}
