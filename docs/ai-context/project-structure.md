# Gomarkdown2image 项目结构

这是 Gomarkdown2image 项目的完整技术栈和项目结构文档,针对 AI 代理优化。

## 项目元信息

- **项目名称**: Gomarkdown2image
- **类型**: Go CLI 工具 + HTTP API 服务
- **主要用途**: Markdown 文档转换为图像
- **Go 版本**: 1.25.1
- **主要技术**: Goldmark + Rod (无头浏览器) + Gin (Web 框架)
- **架构模式**: 三层架构 (Parser → Converter → Renderer)

---

## 项目当前状态

**阶段**: HTTP API 完成 (v0.1.0)

**已完成:**
- ✅ 核心转换系统 (Parser, Renderer, Converter)
- ✅ CLI 工具 (完整命令行参数)
- ✅ HTTP API 服务 (Gin 框架,JSON + 文件上传)
- ✅ 多格式输出 (PNG, JPEG, WebP)
- ✅ 主题系统 (light, dark)
- ✅ 代码语法高亮 (Chroma)
- ✅ 完整文档 (API 文档,实现说明,快速开始)

**下一步:**
- API 性能优化 (连接池,缓存)
- AI 增强功能 (Claude API / Ollama)
- 自定义 CSS 模板

---

## 实际技术栈

### 核心库 (已集成)

| 组件 | 包名 | 版本 | 用途 |
|------|------|------|------|
| **Markdown 解析** | `github.com/yuin/goldmark` | v1.7.13 | CommonMark + GFM 解析 |
| **代码高亮** | `github.com/alecthomas/chroma/v2` | v2.20.0 | 多语言语法高亮 |
| **高亮扩展** | `github.com/yuin/goldmark-highlighting/v2` | v2.0.0 | Goldmark 集成 |
| **无头浏览器** | `github.com/go-rod/rod` | v0.116.2 | HTML 渲染为图片 |
| **Web 框架** | `github.com/gin-gonic/gin` | v1.11.0 | HTTP API 服务 |
| **CORS 中间件** | `github.com/gin-contrib/cors` | v1.7.6 | 跨域资源共享 |

### 支持库

| 组件 | 包名 | 版本 | 用途 |
|------|------|------|------|
| **参数验证** | `github.com/go-playground/validator/v10` | v10.29.0 | 请求参数验证 |
| **正则表达式** | `github.com/dlclark/regexp2` | v1.11.5 | Chroma 依赖 |
| **Rod 工具库** | `github.com/ysmood/*` | - | Rod 运行时支持 |

---

## 实际项目文件树

```
Gomarkdown2image/
├── cmd/
│   ├── markdown2image/           # CLI 命令行工具
│   │   └── main.go               # CLI 主程序
│   └── api/                      # HTTP API 服务
│       └── main.go               # API 服务入口 (Gin 路由,中间件配置)
│
├── pkg/                          # 公共库
│   ├── parser/                   # Markdown → HTML 转换
│   │   ├── parser.go             # GoldmarkParser (Parse, ParseToString)
│   │   └── template.go           # HTML 模板系统 (WrapHTML, CSS 生成)
│   │
│   ├── renderer/                 # HTML → 图片渲染
│   │   └── renderer.go           # RodRenderer (RenderToImage, RenderToFile)
│   │
│   ├── converter/                # 端到端转换协调
│   │   └── converter.go          # DefaultConverter (Convert, ConvertFile)
│   │
│   └── handlers/                 # HTTP 处理器
│       ├── types.go              # 请求/响应数据结构 (ConvertRequest, UploadRequest)
│       ├── convert.go            # 转换端点 (ConvertHandler, UploadHandler)
│       └── middleware.go         # 中间件 (CORS, 日志, 错误恢复, 健康检查)
│
├── internal/                     # 内部实现 (预留)
│   ├── config/                   # 配置管理
│   └── utils/                    # 工具函数
│
├── docs/                         # 文档
│   ├── ai-context/               # AI 上下文文档
│   │   ├── project-structure.md  # 本文档
│   │   └── docs-overview.md      # 文档概览
│   ├── API.md                    # HTTP API 完整文档
│   └── IMPLEMENTATION.md         # 实现说明
│
├── examples/                     # 示例文件
│   ├── basic.md                  # 基础功能示例
│   ├── technical-doc.md          # 技术文档示例
│   └── api-test.sh               # API 测试脚本
│
├── testdata/                     # 测试数据
│   ├── input/                    # 测试输入
│   └── output/                   # 生成的图片
│
├── CLAUDE.md                     # 主 AI 上下文
├── README.md                     # 用户文档
├── QUICKSTART.md                 # 快速开始指南
├── go.mod                        # Go 模块定义
├── go.sum                        # 依赖校验和
├── markdown2image                # CLI 可执行文件
└── markdown2image-api            # API 可执行文件 (39MB)
```

---

## 核心架构设计

### 实现架构

```
Markdown 输入
    ↓
[Parser] Goldmark 解析器
    ↓ (HTML 内容)
[Template] 应用样式和主题
    ↓ (完整 HTML 文档)
[Renderer] Rod 无头浏览器
    ↓ (截图)
图像输出 (PNG/JPEG/WebP)
```

### 组件实现

**Parser (pkg/parser/)**
- **实现**: GoldmarkParser + HTMLTemplate
- **功能**: Markdown → HTML,GFM 扩展,Chroma 代码高亮,主题系统
- **文件**: parser.go (解析), template.go (模板和 CSS)

**Renderer (pkg/renderer/)**
- **实现**: RodRenderer (基于无头浏览器)
- **功能**: HTML → 图片,全页截图,多格式输出,自定义视口
- **文件**: renderer.go

**Converter (pkg/converter/)**
- **实现**: DefaultConverter (协调 Parser 和 Renderer)
- **功能**: 端到端转换,统一配置管理,文件操作封装
- **文件**: converter.go

**Handlers (pkg/handlers/)**
- **实现**: Gin HTTP 处理器
- **功能**: JSON 转换端点,文件上传端点,CORS 中间件,参数验证
- **文件**: types.go (数据结构), convert.go (端点), middleware.go (中间件)

---

## 接口设计

### Parser 接口
```go
type Parser interface {
    Parse(input []byte) (ast.Node, error)
}
```

### Converter 接口
```go
type Converter interface {
    Convert(ast ast.Node, opts *Options) (*Layout, error)
}

type Options struct {
    Width          int
    Height         int
    BackgroundColor color.Color
    Theme          string
    FontFamily     string
    FontSize       int
    Padding        int
    OutputFormat   Format
}
```

### Renderer 接口
```go
type Renderer interface {
    Render(layout *Layout, output io.Writer) error
    SetFormat(format Format) error
}

type Format int
const (
    FormatPNG Format = iota
    FormatJPEG
    FormatWebP
    FormatSVG
)
```

---

## 数据流详解

### 步骤 1: 解析
```
Markdown 文本 → Parser.Parse() → AST
```

### 步骤 2: 转换
```
AST + Options → Converter.Convert() → Layout 树
- 应用样式配置
- 计算文本宽度和换行
- 计算元素位置
- 处理分页 (如需要)
```

### 步骤 3: 渲染
```
Layout 树 → Renderer.Render() → Image 文件
- 绘制背景
- 渲染文本 (字体处理)
- 渲染图形元素
- 输出指定格式
```

---

## 开发工作流

### 构建
```bash
# 构建 CLI 工具
go build -o markdown2image ./cmd/markdown2image

# 构建 API 服务
go build -o markdown2image-api ./cmd/api

# 优化发布版本
go build -ldflags="-s -w" -o markdown2image ./cmd/markdown2image
go build -ldflags="-s -w" -o markdown2image-api ./cmd/api
```

### 运行

#### CLI 工具
```bash
# 基础用法
./markdown2image -input examples/basic.md -output output.png

# 暗色主题
./markdown2image -input doc.md -output doc.png -theme dark

# 自定义宽度和格式
./markdown2image -input doc.md -output doc.jpg -format jpeg -width 1400
```

#### HTTP API 服务
```bash
# 启动服务 (默认端口 8080)
./markdown2image-api

# 指定端口
PORT=3000 ./markdown2image-api

# 生产模式
GIN_MODE=release ./markdown2image-api
```

#### API 调用示例
```bash
# JSON 转换
curl -X POST http://localhost:8080/api/convert \
  -H "Content-Type: application/json" \
  -d '{"markdown":"# Test","theme":"dark","imageFormat":"png"}' \
  --output output.png

# 文件上传
curl -X POST http://localhost:8080/api/upload \
  -F "file=@document.md" \
  -F "theme=light" \
  --output output.png

# 健康检查
curl http://localhost:8080/health
```

### 依赖管理
```bash
# 查看依赖
go list -m all

# 清理依赖
go mod tidy

# 验证依赖
go mod verify
```

---

## 实现路线图

### 阶段 1: 项目初始化 ✅ 完成
- [x] Go 模块初始化
- [x] AI 上下文文档
- [x] 创建目录结构
- [x] 创建 README

### 阶段 2: MVP 功能 ✅ 完成
- [x] Markdown → HTML Parser (Goldmark + GFM)
- [x] HTML → 图片 Renderer (Rod 无头浏览器)
- [x] Converter 协调器
- [x] 多格式输出 (PNG, JPEG, WebP)
- [x] CLI 工具
- [x] 代码语法高亮 (Chroma)

### 阶段 3: 样式系统 ✅ 完成
- [x] HTML 模板系统
- [x] 主题配置 (light/dark)
- [x] 代码块样式
- [x] 表格和列表渲染

### 阶段 3.5: HTTP API 服务 ✅ 完成 (2025-12-14)
- [x] Gin 框架集成
- [x] POST /api/convert (JSON 转换)
- [x] POST /api/upload (文件上传)
- [x] 参数验证 (10 个配置参数)
- [x] CORS 中间件
- [x] 完整 API 文档

### 阶段 4: AI 增强 🚧 规划中
- [ ] Claude API 集成
- [ ] Ollama 本地模型
- [ ] 内容增强功能

### 阶段 5: 高级特性 📋 待定
- [ ] 自定义 CSS 模板
- [ ] 配置文件支持
- [ ] 批量转换
- [ ] 性能优化

---

## Go 编码标准

### 命名约定
- **包名**: 小写,无下划线 (例如: `parser`)
- **文件名**: 小写,下划线分隔 (例如: `image_renderer.go`)
- **类型**: PascalCase (例如: `ImageRenderer`)
- **函数**: PascalCase (公开) 或 camelCase (私有)
- **接口**: `-er` 结尾 (例如: `Parser`, `Renderer`)

### 文档要求
- 每个公开类型/函数必须有文档注释
- 包级文档 (doc.go 或主文件)
- 使用示例测试 (`Example*`)

### 测试标准
- 表驱动测试
- 至少 80% 覆盖率
- 性能关键路径需要基准测试

---

## HTTP API 端点

### 可用端点
- `GET /health` - 健康检查
- `GET /` - 服务信息
- `POST /api/convert` - JSON 方式 Markdown 转换
- `POST /api/upload` - 文件上传方式转换

### API 参数 (10 个)
- **HTML 样式**: title, theme, width, fontSize, fontFamily, customCss
- **图片配置**: imageFormat (png/jpeg/webp), imageQuality (1-100), devicePixelRatio (0.5-4.0)
- **验证**: 自动参数验证,最大文件大小 10MB

### 中间件
- CORS (跨域资源共享)
- Logger (请求日志)
- Recovery (错误恢复)
- Validator (参数验证)

---

## 相关文档

- **[CLAUDE.md](/CLAUDE.md)** - 主 AI 上下文和架构文档
- **[docs-overview.md](/docs/ai-context/docs-overview.md)** - 文档架构导航
- **[API.md](/docs/API.md)** - HTTP API 完整文档
- **[IMPLEMENTATION.md](/docs/IMPLEMENTATION.md)** - 实现说明
- **[QUICKSTART.md](/QUICKSTART.md)** - 快速开始指南

---

**文档版本**: 2025-12-14
**项目阶段**: HTTP API 完成 (v0.1.0)
**代码库状态**: 生产就绪
**针对**: AI 代理优化 - 快速导航和技术参考
