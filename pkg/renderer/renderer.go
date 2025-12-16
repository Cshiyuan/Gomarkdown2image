// Package renderer 提供 HTML 到图片的渲染功能
//
// 本包使用 Rod 无头浏览器将 HTML 渲染为图片
package renderer

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// Renderer HTML 渲染器接口
type Renderer interface {
	// RenderToImage 将 HTML 渲染为图片
	RenderToImage(html string, opts *RenderOptions) ([]byte, error)

	// RenderToFile 将 HTML 渲染为图片并保存到文件
	RenderToFile(html string, outputPath string, opts *RenderOptions) error

	// Close 关闭渲染器,释放资源
	Close() error
}

// RenderOptions 渲染选项
type RenderOptions struct {
	Width            int         // 视口宽度(默认 1200)
	Height           int         // 视口高度(默认 0,表示自动高度)
	Format           ImageFormat // 图片格式(默认 PNG)
	Quality          int         // 图片质量 1-100 (仅 JPEG 有效,默认 90)
	FullPage         bool        // 是否全页截图(默认 true)
	DevicePixelRatio float64     // 设备像素比(默认 1.0)
}

// ImageFormat 图片格式
type ImageFormat string

const (
	FormatPNG  ImageFormat = "png"
	FormatJPEG ImageFormat = "jpeg"
	FormatWebP ImageFormat = "webp"
)

// DefaultRenderOptions 返回默认渲染选项
func DefaultRenderOptions() *RenderOptions {
	return &RenderOptions{
		Width:            1200,
		Height:           0,
		Format:           FormatPNG,
		Quality:          90,
		FullPage:         true,
		DevicePixelRatio: 1.0,
	}
}

// RodRenderer 基于 Rod 的渲染器实现
type RodRenderer struct {
	browser *rod.Browser
}

// NewRodRenderer 创建新的 Rod 渲染器
//
// 特性:
//   - 自动管理浏览器实例
//   - 支持全页截图
//   - 支持多种图片格式
func NewRodRenderer() (*RodRenderer, error) {
	// 启动无头浏览器 (使用 defer/recover 捕获 panic)
	var browser *rod.Browser
	var panicErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				panicErr = fmt.Errorf("failed to connect to browser: %v", r)
			}
		}()
		browser = rod.New().Timeout(10 * time.Second).MustConnect()
	}()

	if panicErr != nil {
		return nil, panicErr
	}

	return &RodRenderer{
		browser: browser,
	}, nil
}

// RenderToImage 将 HTML 渲染为图片字节数组
//
// 参数:
//   - html: 完整的 HTML 文档
//   - opts: 渲染选项
//
// 返回:
//   - []byte: 图片字节数组
//   - error: 渲染错误(如有)
func (r *RodRenderer) RenderToImage(html string, opts *RenderOptions) ([]byte, error) {
	if opts == nil {
		opts = DefaultRenderOptions()
	}

	// 创建带超时的上下文 (30 秒总超时)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建新页面 (使用 Page 而非 MustPage,避免 panic)
	page, err := r.browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	defer page.Close()

	// 设置页面上下文为带超时的 context
	page = page.Context(ctx)

	// 设置视口大小
	err = page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:             opts.Width,
		Height:            opts.Height,
		DeviceScaleFactor: opts.DevicePixelRatio,
		Mobile:            false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set viewport: %w", err)
	}

	// 注入 HTML 内容
	err = page.SetDocumentContent(html)
	if err != nil {
		return nil, fmt.Errorf("failed to set document content: %w", err)
	}

	// 等待页面加载完成
	err = page.WaitLoad()
	if err != nil {
		return nil, fmt.Errorf("failed to wait for page load: %w", err)
	}

	// 等待页面 idle (使用更短的超时,失败不影响主流程)
	// 使用 5 秒超时,如果失败仅记录警告
	idleCtx, idleCancel := context.WithTimeout(ctx, 5*time.Second)
	defer idleCancel()
	if err := page.Context(idleCtx).WaitIdle(1 * time.Second); err != nil {
		// WaitIdle 失败不应阻止截图,仅在调试时记录
		// 生产环境可以移除或使用日志库
		// log.Printf("warning: page idle wait failed (continuing): %v", err)
	}

	// 构建截图参数
	screenshotOpts := &proto.PageCaptureScreenshot{
		FromSurface: true,
	}

	// 设置图片格式
	switch opts.Format {
	case FormatPNG:
		screenshotOpts.Format = proto.PageCaptureScreenshotFormatPng
	case FormatJPEG:
		screenshotOpts.Format = proto.PageCaptureScreenshotFormatJpeg
		quality := opts.Quality
		screenshotOpts.Quality = &quality
	case FormatWebP:
		screenshotOpts.Format = proto.PageCaptureScreenshotFormatWebp
		quality := opts.Quality
		screenshotOpts.Quality = &quality
	default:
		return nil, fmt.Errorf("unsupported image format: %s", opts.Format)
	}

	// 全页截图
	if opts.FullPage {
		return page.Screenshot(true, screenshotOpts)
	}

	// 视口截图
	return page.Screenshot(false, screenshotOpts)
}

// RenderToFile 将 HTML 渲染为图片并保存到文件
//
// 参数:
//   - html: 完整的 HTML 文档
//   - outputPath: 输出文件路径
//   - opts: 渲染选项
//
// 返回:
//   - error: 渲染或文件写入错误(如有)
func (r *RodRenderer) RenderToFile(html string, outputPath string, opts *RenderOptions) error {
	// 渲染为图片
	imageData, err := r.RenderToImage(html, opts)
	if err != nil {
		return fmt.Errorf("failed to render image: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(outputPath, imageData, 0644); err != nil {
		return fmt.Errorf("failed to write image file: %w", err)
	}

	return nil
}

// Close 关闭渲染器,释放资源
func (r *RodRenderer) Close() error {
	if r.browser != nil {
		r.browser.Close()
	}
	return nil
}
