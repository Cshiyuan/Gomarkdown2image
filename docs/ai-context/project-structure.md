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

**阶段**: AI 增强功能完成 (v0.2.0)

**已完成:**
- ✅ 核心转换系统 (Parser, Renderer, Converter)
- ✅ CLI 工具 (完整命令行参数)
- ✅ HTTP API 服务 (Gin 框架,JSON + 文件上传)
- ✅ 多格式输出 (PNG, JPEG, WebP)
- ✅ 主题系统 (light, dark)
- ✅ 代码语法高亮 (Chroma)
- ✅ 完整文档 (API 文档,实现说明,快速开始)
- ✅ **代码重组** (internal/ 架构,消除重复代码)
- ✅ **单元测试** (18.3% 覆盖率,utils 包 100%)
- ✅ **AI 增强功能** (Gemini + Ollama 双后端支持) 🆕

**下一步:**
- 提高测试覆盖率 (目标 80%)
- API 性能优化 (连接池,缓存)
- 自定义 CSS 模板
- AI 功能扩展 (更多提示词模板)

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
| **AI SDK (Gemini)** | `github.com/google/generative-ai-go` | v0.20.1 | Google Gemini API 集成 🆕 |
| **AI SDK (Ollama)** | `github.com/ollama/ollama` | v0.13.3 | 本地 AI 模型支持 🆕 |

### 支持库

| 组件 | 包名 | 版本 | 用途 |
|------|------|------|------|
| **参数验证** | `github.com/go-playground/validator/v10` | v10.29.0 | 请求参数验证 |
| **正则表达式** | `github.com/dlclark/regexp2` | v1.11.5 | Chroma 依赖 |
| **Rod 工具库** | `github.com/ysmood/*` | - | Rod 运行时支持 |
| **Google API** | `google.golang.org/api` | v0.257.0 | Google 服务基础设施 🆕 |
| **重试 HTTP** | `github.com/hashicorp/go-retryablehttp` | v0.7.7 | HTTP 重试机制 🆕 |

---

## 实际项目文件树

```
github.com/Cshiyuan/Gomarkdown2image/
├── cmd/
│   ├── markdown2image/           # CLI 命令行工具
│   │   └── main.go               # CLI 主程序 (使用 utils 包)
│   └── api/                      # HTTP API 服务
│       └── main.go               # API 服务入口 (Gin 路由,中间件配置)
│
├── pkg/                          # 公共库 (可被外部导入)
│   ├── ai/                       # AI 服务抽象层 🆕
│   │   ├── provider.go           # Provider 接口定义
│   │   ├── types.go              # 核心数据类型 (Config, Request, Response)
│   │   ├── errors.go             # AI 错误处理和分类
│   │   ├── prompts.go            # 提示词模板系统 (5 个内置模板)
│   │   ├── factory/              # Provider 工厂 (解决循环依赖)
│   │   │   └── factory.go        # NewProvider() 工厂函数
│   │   ├── gemini/               # Google Gemini 客户端
│   │   │   └── client.go         # Gemini API 实现
│   │   └── ollama/               # Ollama 本地模型客户端
│   │       └── client.go         # Ollama API 实现
│   │
│   ├── parser/                   # Markdown → HTML 转换
│   │   ├── parser.go             # GoldmarkParser (Parse, ParseToString)
│   │   ├── parser_test.go        # 单元测试 (16 个测试用例,89.3% 覆盖率)
│   │   ├── provider.go           # Parser Provider 抽象层 (AI + 传统) 🆕
│   │   └── template.go           # HTML 模板系统 (WrapHTML, CSS 生成)
│   │
│   ├── renderer/                 # HTML → 图片渲染
│   │   └── renderer.go           # RodRenderer (RenderToImage, RenderToFile)
│   │
│   └── converter/                # 端到端转换协调
│       └── converter.go          # DefaultConverter (Convert, ConvertFile)
│
├── internal/                     # 内部实现 (不可被外部导入)
│   ├── config/                   # 配置管理 (单一真相来源)
│   │   ├── defaults.go           # 默认配置值 (DefaultTitle, DefaultTheme, DefaultWidth 等)
│   │   ├── limits.go             # 限制常量 (MaxMarkdownSize, MinWidth, MaxQuality 等)
│   │   └── ai.go                 # AI 相关配置 (Provider, Model, Timeout 等) 🆕
│   │
│   ├── utils/                    # 工具函数 (消除代码重复)
│   │   ├── format.go             # 图片格式解析 (ParseImageFormat, GetContentType)
│   │   ├── format_test.go        # 格式测试 (11 个测试用例,100% 覆盖率)
│   │   ├── validation.go         # 参数验证 + XSS 防护 (6 个验证函数)
│   │   │                         # - ValidateQuality/Width/FontSize/DevicePixelRatio
│   │   │                         # - ValidateTheme
│   │   │                         # - ValidateCustomCSS() (12 个禁止模式 XSS 防护) 🆕
│   │   ├── validation_test.go    # 验证测试 (40+ 个测试用例,100% 覆盖率)
│   │   └── validation_css_test.go # CustomCSS XSS 防护测试 (14 个安全测试用例) 🆕
│   │
│   └── handlers/                 # HTTP 处理器 (应用层,非公共 API)
│       ├── types.go              # 请求/响应数据结构 + RequestParams 接口 (17 个 getter) 🆕
│       ├── convert.go            # 转换端点 (ConvertHandler, UploadHandler)
│       ├── convert_test.go       # 接口实现测试 (43 个子测试,参数映射验证) 🆕
│       └── middleware.go         # 中间件 (CORS, 日志, 错误恢复, 健康检查)
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
│   ├── api-test.sh               # API 测试脚本
│   └── ai-example.sh             # AI 功能示例脚本 🆕
│
├── testdata/                     # 测试数据
│   ├── input/                    # 测试输入
│   └── output/                   # 生成的图片
│
├── CLAUDE.md                     # 主 AI 上下文
├── README.md                     # 用户文档
├── QUICKSTART.md                 # 快速开始指南
├── go.mod                        # Go 模块定义 (github.com/Cshiyuan/Gomarkdown2image)
├── go.sum                        # 依赖校验和
├── coverage.out                  # 测试覆盖率报告 (18.3%)
├── markdown2image                # CLI 可执行文件
└── markdown2image-api            # API 可执行文件 (39MB)
```

---

## 核心架构设计

### 实现架构

#### 传统模式
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

#### AI 增强模式 🆕
```
Markdown 输入
    ↓
[AI Parser] 提示词模板 + AI Provider (Gemini/Ollama)
    ↓ (AI 增强的 Markdown)
[Parser] Goldmark 解析器
    ↓ (HTML 内容)
[Template] 应用样式和主题
    ↓ (完整 HTML 文档)
[Renderer] Rod 无头浏览器
    ↓ (截图)
图像输出 (PNG/JPEG/WebP)
```

**AI 模式特性:**
- 双后端支持: Gemini (云端) 和 Ollama (本地)
- 5 个内置提示词模板: enhance, translate, format, explain_code, summarize
- 自定义提示词支持
- 自动降级: AI 失败时自动回退到传统模式
- 错误分类: 7 种错误类型 (auth, rate_limit, invalid_req, server_error, timeout, network, unknown)

### 组件实现

**AI 服务层 (pkg/ai/)** 🆕
- **实现**: Provider 模式 + 双后端支持 (Gemini, Ollama)
- **功能**: AI 内容增强,提示词模板系统,错误处理和重试,自动降级
- **文件**:
  - provider.go (Provider 接口)
  - types.go (Config, Request, Response 数据类型)
  - errors.go (错误分类和处理)
  - prompts.go (5 个内置提示词模板)
  - factory/factory.go (Provider 工厂)
  - gemini/client.go (Gemini 客户端)
  - ollama/client.go (Ollama 客户端)

**Parser (pkg/parser/)**
- **实现**: GoldmarkParser + HTMLTemplate + AI Parser Provider 🆕
- **功能**: Markdown → HTML,GFM 扩展,Chroma 代码高亮,主题系统,AI 增强
- **文件**:
  - parser.go (GoldmarkParser 解析)
  - template.go (HTML 模板和 CSS)
  - provider.go (Parser Provider 抽象层,支持传统/AI 双模式) 🆕

**Renderer (pkg/renderer/)**
- **实现**: RodRenderer (基于无头浏览器)
- **功能**: HTML → 图片,全页截图,多格式输出,自定义视口
- **文件**: renderer.go

**Converter (pkg/converter/)**
- **实现**: DefaultConverter (协调 Parser 和 Renderer)
- **功能**: 端到端转换,统一配置管理,文件操作封装
- **文件**: converter.go

**Handlers (internal/handlers/)**
- **实现**: Gin HTTP 处理器
- **功能**: JSON 转换端点,文件上传端点,CORS 中间件,参数验证
- **文件**: types.go (数据结构), convert.go (端点), middleware.go (中间件)
- **注意**: 移至 internal/ 因为是应用层代码,不应作为公共 API

**Config (internal/config/)**
- **实现**: 配置常量和默认值管理
- **功能**: 单一真相来源,消除硬编码常量
- **文件**:
  - defaults.go (默认配置)
  - limits.go (限制常量)
  - ai.go (AI 相关配置: Provider, Model, Timeout 等) 🆕

**Utils (internal/utils/)**
- **实现**: 通用工具函数
- **功能**: 格式解析,参数验证 (消除代码重复)
- **文件**: format.go (图片格式), validation.go (参数验证)
- **测试覆盖率**: 100%

---

## AI 服务架构 🆕

### 设计理念

AI 增强功能采用 **Provider Pattern** 实现,核心设计原则:
- **可插拔**: 支持多个 AI 后端,统一接口
- **可靠性**: 自动降级,AI 失败时回退到传统模式
- **灵活性**: 内置模板 + 自定义提示词
- **安全性**: 错误分类,不暴露敏感信息

### 双后端架构

**Gemini (云端 AI)**
- **提供器**: Google Generative AI
- **模型**: gemini-2.0-flash-exp (默认)
- **优势**: 强大的生成能力,云端推理,无需本地资源
- **需求**: API Key (从 https://ai.google.dev/ 获取)
- **适用**: 生产环境,需要高质量内容增强

**Ollama (本地 AI)**
- **提供器**: Ollama 本地服务
- **模型**: llama3.2, qwen2.5 等 (可选)
- **优势**: 隐私保护,无网络依赖,无 API 成本
- **需求**: 本地运行 Ollama 服务 (ollama serve)
- **适用**: 开发环境,隐私敏感场景

### 提示词模板系统

**5 个内置模板:**

1. **enhance** (默认)
   - 用途: 内容润色和增强
   - 特性: 保持原意,改善表达,添加细节

2. **translate**
   - 用途: 多语言翻译
   - 参数: TargetLang (目标语言)
   - 特性: 保留格式,准确翻译

3. **format**
   - 用途: 格式优化
   - 特性: 标题层级,列表结构,代码块标注

4. **explain_code**
   - 用途: 代码解释
   - 参数: Language (编程语言)
   - 特性: 添加注释,解释逻辑

5. **summarize**
   - 用途: 内容摘要
   - 特性: 提取关键点,生成概要

**自定义提示词:**
- 完全自定义提示词内容
- 最大长度: 10000 字符
- 支持模板变量插值

### 错误处理策略

**7 种错误分类:**
1. `auth` - 认证失败 (API Key 错误)
2. `rate_limit` - 速率限制 (超出配额)
3. `invalid_req` - 无效请求 (参数错误)
4. `server_error` - 服务器错误 (5xx)
5. `timeout` - 超时错误
6. `network` - 网络错误
7. `unknown` - 未知错误

**自动降级机制:**
```
AI 模式请求 → AI Provider 处理
    ↓ (失败)
错误分类和记录
    ↓
自动切换到传统模式
    ↓
Traditional Parser 处理 → 返回结果
```

### 集成方式

**CLI 工具集成:**
```bash
# 使用 Gemini
markdown2image -input doc.md -output doc.png \
  -parser-mode ai \
  -ai-provider gemini \
  -ai-api-key YOUR_KEY \
  -ai-template enhance

# 使用 Ollama
markdown2image -input doc.md -output doc.png \
  -parser-mode ai \
  -ai-provider ollama \
  -ai-endpoint http://localhost:11434
```

**HTTP API 集成:**
```bash
# JSON 请求
curl -X POST http://localhost:8080/api/convert \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# Test",
    "parserMode": "ai",
    "aiProvider": "gemini",
    "aiModel": "gemini-2.0-flash-exp",
    "aiApiKey": "YOUR_KEY",
    "aiPromptTemplate": "enhance"
  }'
```

**Go 代码集成:**
```go
// 创建 AI 配置
aiConfig := &ai.Config{
    Provider:   ai.ProviderGemini,
    APIKey:     "YOUR_KEY",
    Model:      "gemini-2.0-flash-exp",
    Timeout:    30,
    MaxRetries: 3,
}

// 创建 Parser Provider
providerConfig := &parser.ProviderConfig{
    Type:             parser.ProviderTypeAI,
    AIConfig:         aiConfig,
    AIPromptTemplate: "enhance",
}

provider, _ := parser.NewProvider(providerConfig)
p, _ := provider.CreateParser()

// 使用 Parser
html, _ := p.ParseToString([]byte("# Test"))
```

---

## 接口设计

### AI Provider 接口 🆕
```go
// AI Provider 统一接口
type Provider interface {
    Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
    GenerateStream(ctx context.Context, req *GenerateRequest) (<-chan StreamChunk, error)
    Name() string
    Close() error
}

// AI 配置
type Config struct {
    Provider   ProviderType              // gemini 或 ollama
    APIKey     string                    // API 密钥 (Gemini)
    BaseURL    string                    // 服务端点 (Ollama)
    Model      string                    // 模型名称
    Timeout    int                       // 超时时间 (秒)
    MaxRetries int                       // 最大重试次数
    Prompts    *PromptConfig            // 提示词配置
}

// AI 请求
type GenerateRequest struct {
    Prompt      string                  // 用户提示词
    System      string                  // 系统提示词
    MaxTokens   int                     // 最大 token 数
    Temperature float64                 // 温度参数
}
```

### Parser Provider 接口 🆕
```go
// Parser Provider 抽象层
type ParserProvider interface {
    CreateParser() (Parser, error)
    Name() string
}

// Provider 配置
type ProviderConfig struct {
    Type             ProviderType           // traditional 或 ai
    AIConfig         *ai.Config            // AI 配置
    AIPromptTemplate string                // 提示词模板名称
    AIPromptData     map[string]interface{} // 模板数据
    CustomPrompt     string                // 自定义提示词
}
```

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

### 阶段 3.6: 代码质量优化 ✅ 完成 (2025-12-15)
- [x] Go 模块路径标准化 (github.com/Cshiyuan/Gomarkdown2image)
- [x] internal/ 架构重组 (config, utils, handlers)
- [x] 消除代码重复 (统一格式解析和验证)
- [x] 单元测试套件 (18.3% 覆盖率,utils 包 100%)
- [x] 文档更新 (反映新架构)

### 阶段 4: AI 增强功能 ✅ 完成 (2025-12-15)
- [x] AI Provider 抽象层设计
- [x] Gemini API 集成 (Google Generative AI)
- [x] Ollama 本地模型集成
- [x] Parser Provider 架构 (传统/AI 双模式)
- [x] 提示词模板系统 (5 个内置模板)
- [x] 自定义提示词支持
- [x] AI 错误处理和自动降级
- [x] HTTP API 扩展 (7 个 AI 参数)
- [x] AI 使用示例脚本
- [x] 文档更新 (AI 架构说明)

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
- 表驱动测试 (struct slices with test cases)
- 至少 80% 覆盖率 (当前 18.3%)
- 性能关键路径需要基准测试
- 测试文件命名: `*_test.go`

**当前测试状态**:
- `internal/utils/format_test.go`: 100% 覆盖率 (11 个测试用例 + 2 个基准测试)
- `internal/utils/validation_test.go`: 100% 覆盖率 (40 个测试用例 + 3 个基准测试)
- `pkg/parser/parser_test.go`: 89.3% 覆盖率 (16 个测试用例 + 2 个基准测试)
- 总体覆盖率: 18.3% (需要为 renderer, converter, handlers 添加测试)

---

## HTTP API 端点

### 可用端点
- `GET /health` - 健康检查
- `GET /` - 服务信息
- `POST /api/convert` - JSON 方式 Markdown 转换
- `POST /api/upload` - 文件上传方式转换

### API 参数 (17 个)
- **HTML 样式**: title, theme, width, fontSize, fontFamily, customCss
- **图片配置**: imageFormat (png/jpeg/webp), imageQuality (1-100), devicePixelRatio (0.5-4.0)
- **AI 增强** 🆕:
  - parserMode (traditional/ai)
  - aiProvider (gemini/ollama)
  - aiModel (模型名称)
  - aiApiKey (API 密钥)
  - aiEndpoint (服务端点)
  - aiPromptTemplate (提示词模板)
  - aiCustomPrompt (自定义提示词)
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

**文档版本**: 2025-12-15
**项目阶段**: AI 增强功能完成 (v0.2.0)
**代码库状态**: 生产就绪
**针对**: AI 代理优化 - 快速导航和技术参考
