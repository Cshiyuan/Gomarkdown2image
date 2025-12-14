package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Cshiyuan/Gomarkdown2image/internal/utils"
	"github.com/Cshiyuan/Gomarkdown2image/pkg/converter"
)

const (
	version = "0.1.0"
)

func main() {
	// 定义命令行参数
	var (
		input       = flag.String("input", "", "输入的 Markdown 文件路径 (必需)")
		output      = flag.String("output", "", "输出的图片文件路径 (必需)")
		title       = flag.String("title", "Markdown to Image", "页面标题")
		theme       = flag.String("theme", "light", "主题 (light, dark)")
		width       = flag.Int("width", 1200, "页面宽度(像素)")
		fontSize    = flag.Int("font-size", 16, "字体大小(px)")
		fontFamily  = flag.String("font-family", "Arial, sans-serif", "字体族")
		format      = flag.String("format", "png", "图片格式 (png, jpeg, webp)")
		quality     = flag.Int("quality", 90, "图片质量 1-100 (仅 JPEG/WebP)")
		dpr         = flag.Float64("dpr", 1.0, "设备像素比")
		showVersion = flag.Bool("version", false, "显示版本信息")
	)

	flag.Parse()

	// 显示版本
	if *showVersion {
		fmt.Printf("markdown2image version %s\n", version)
		os.Exit(0)
	}

	// 验证必需参数
	if *input == "" {
		fmt.Fprintln(os.Stderr, "错误: 必须指定输入文件 (-input)")
		flag.Usage()
		os.Exit(1)
	}

	if *output == "" {
		fmt.Fprintln(os.Stderr, "错误: 必须指定输出文件 (-output)")
		flag.Usage()
		os.Exit(1)
	}

	// 验证输入文件存在
	if _, err := os.Stat(*input); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "错误: 输入文件不存在: %s\n", *input)
		os.Exit(1)
	}

	// 解析图片格式
	imageFormat, err := utils.ParseImageFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	// 验证图片质量
	if err := utils.ValidateQuality(*quality); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	// 创建输出目录(如果不存在)
	if err := os.MkdirAll(filepath.Dir(*output), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法创建输出目录: %v\n", err)
		os.Exit(1)
	}

	// 创建转换器
	fmt.Println("正在初始化转换器...")
	conv, err := converter.NewConverter()
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法创建转换器: %v\n", err)
		os.Exit(1)
	}
	defer conv.Close()

	// 配置转换选项
	opts := &converter.ConvertOptions{
		Title:            *title,
		Theme:            *theme,
		Width:            *width,
		FontSize:         *fontSize,
		FontFamily:       *fontFamily,
		ImageFormat:      imageFormat,
		ImageQuality:     *quality,
		FullPage:         true,
		DevicePixelRatio: *dpr,
	}

	// 执行转换
	fmt.Printf("正在转换 %s...\n", *input)
	if err := conv.ConvertFile(*input, *output, opts); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 转换失败: %v\n", err)
		os.Exit(1)
	}

	// 获取输出文件大小
	stat, err := os.Stat(*output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: 无法读取输出文件信息: %v\n", err)
	}

	fmt.Printf("✅ 转换成功!\n")
	fmt.Printf("   输入: %s\n", *input)
	fmt.Printf("   输出: %s\n", *output)
	if stat != nil {
		fmt.Printf("   大小: %.2f KB\n", float64(stat.Size())/1024.0)
	}
	fmt.Printf("   格式: %s\n", imageFormat)
	fmt.Printf("   尺寸: %dpx (宽度)\n", *width)
}
