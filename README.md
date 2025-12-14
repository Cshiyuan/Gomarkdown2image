# Gomarkdown2image

一个强大的 Markdown 转图片工具,支持代码语法高亮和多种输出格式。

## ✨ 特性

- ✅ **Markdown 解析**: 基于 Goldmark,完全支持 CommonMark 标准
- ✅ **代码高亮**: 集成 Chroma,支持多种编程语言的语法高亮
- ✅ **无头浏览器渲染**: 使用 Rod 进行高质量 HTML 渲染
- ✅ **多格式输出**: 支持 PNG, JPEG, WebP 格式
- ✅ **自定义样式**: 支持亮色/暗色主题,可自定义字体和样式
- ✅ **GFM 扩展**: 支持表格、删除线、任务列表等 GitHub 风格特性
- 🚧 **AI 增强**: (计划中) 支持 AI 内容润色和增强

## 🚀 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/Gomarkdown2image.git
cd Gomarkdown2image

# 构建
go build -o markdown2image ./cmd/markdown2image

# 或安装到 $GOPATH/bin
go install ./cmd/markdown2image
```

### 基础用法

```bash
# 基本转换
./markdown2image -input examples/basic.md -output output.png

# 使用暗色主题
./markdown2image -input doc.md -output doc.png -theme dark

# 指定图片格式和质量
./markdown2image -input doc.md -output doc.jpg -format jpeg -quality 95

# 自定义宽度和字体大小
./markdown2image -input doc.md -output doc.png -width 1920 -font-size 18
```

## 📋 命令行参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-input` | string | (必需) | 输入的 Markdown 文件路径 |
| `-output` | string | (必需) | 输出的图片文件路径 |
| `-title` | string | "Markdown to Image" | 页面标题 |
| `-theme` | string | "light" | 主题 (light, dark) |
| `-width` | int | 1200 | 页面宽度(像素) |
| `-font-size` | int | 16 | 字体大小(px) |
| `-font-family` | string | "Arial, sans-serif" | 字体族 |
| `-format` | string | "png" | 图片格式 (png, jpeg, webp) |
| `-quality` | int | 90 | 图片质量 1-100 (仅 JPEG/WebP) |
| `-dpr` | float | 1.0 | 设备像素比 (用于高清屏) |
| `-version` | bool | false | 显示版本信息 |

## 📖 示例

### 示例 1: 技术文档转图片

```bash
./markdown2image \
  -input examples/technical-doc.md \
  -output technical-doc.png \
  -theme light \
  -width 1400 \
  -font-size 16
```

### 示例 2: 代码分享卡片

```bash
./markdown2image \
  -input code-snippet.md \
  -output code-card.png \
  -theme dark \
  -width 800 \
  -font-size 14 \
  -dpr 2.0
```

### 示例 3: 博客文章预览图

```bash
./markdown2image \
  -input blog-post.md \
  -output preview.jpg \
  -format jpeg \
  -quality 90 \
  -width 1200
```

## 🏗️ 架构

```
Markdown 输入
    ↓
[Parser] - Goldmark 解析器
    ↓
HTML 内容
    ↓
[Template] - 应用样式和主题
    ↓
完整 HTML 文档
    ↓
[Renderer] - Rod 无头浏览器
    ↓
图片输出 (PNG/JPEG/WebP)
```

## 📁 项目结构

```
Gomarkdown2image/
├── cmd/
│   └── markdown2image/      # 命令行入口
│       └── main.go
├── pkg/
│   ├── parser/              # Markdown → HTML
│   │   ├── parser.go        # Goldmark 解析器
│   │   └── template.go      # HTML 模板
│   ├── renderer/            # HTML → 图片
│   │   └── renderer.go      # Rod 渲染器
│   └── converter/           # 核心转换器
│       └── converter.go     # 协调 Parser 和 Renderer
├── examples/                # 示例文件
│   └── basic.md
├── testdata/                # 测试数据
│   ├── input/
│   └── output/
└── README.md
```

## 🔧 开发

### 运行测试

```bash
go test ./...
```

### 代码格式化

```bash
go fmt ./...
```

### 静态分析

```bash
go vet ./...
```

## 🛠️ 技术栈

- **Markdown 解析**: [Goldmark](https://github.com/yuin/goldmark) - CommonMark 兼容的 Markdown 解析器
- **代码高亮**: [Chroma](https://github.com/alecthomas/chroma) - 语法高亮库
- **无头浏览器**: [Rod](https://github.com/go-rod/rod) - 高性能浏览器自动化工具
- **Go 版本**: 1.25.1+

## 🚧 路线图

- [x] Markdown → HTML 转换 (支持 CommonMark + GFM)
- [x] 代码语法高亮
- [x] 无头浏览器渲染
- [x] 多格式输出 (PNG, JPEG, WebP)
- [x] 自定义样式和主题
- [ ] AI 内容增强 (Claude API / Ollama)
- [ ] 自定义 CSS 模板
- [ ] 批量转换
- [ ] Web UI 界面
- [ ] Docker 镜像

## 🤝 贡献

欢迎贡献!请随时提交 Pull Request。

## 📄 许可证

MIT License

## 🔗 相关项目

- [Goldmark](https://github.com/yuin/goldmark) - Go Markdown 解析器
- [Rod](https://github.com/go-rod/rod) - Go 浏览器自动化
- [Chroma](https://github.com/alecthomas/chroma) - 语法高亮

---

**生成工具**: Gomarkdown2image v0.1.0
**作者**: Your Name
**GitHub**: https://github.com/yourusername/Gomarkdown2image
