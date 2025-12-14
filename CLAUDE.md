# Gomarkdown2image - AI 上下文文档

## 1. 项目概览

- **愿景**: 创建一个高质量的 Markdown 到图像转换工具,支持丰富的样式和多种输出格式
- **当前阶段**: MVP 完成 - 核心功能已实现并可用
- **关键架构**: Markdown → HTML → 图片 (使用无头浏览器渲染)
- **开发策略**: 迭代开发,优先实现核心功能,后续扩展 AI 增强和高级特性

---

## 2. 项目状态: HTTP API 完成 (v0.1.0)

**已实现功能:**
- ✅ Markdown → HTML 转换 (Goldmark + GFM 扩展)
- ✅ 代码语法高亮 (Chroma)
- ✅ HTML → 图片渲染 (Rod 无头浏览器)
- ✅ CLI 工具 (完整命令行参数支持)
- ✅ 多格式输出 (PNG, JPEG, WebP)
- ✅ 主题系统 (light, dark)
- ✅ **HTTP API 服务** (Gin 框架,支持 JSON + 文件上传) 🆕

**下一步:**
- 实现 AI 增强功能 (Claude API / Ollama)
- 添加自定义 CSS 模板支持
- 性能优化和批量转换
- API 性能优化 (连接池、缓存)

---

## 3. 项目结构

### 当前目录结构

```
Gomarkdown2image/
├── cmd/
│   ├── markdown2image/      # CLI 命令行工具
│   │   └── main.go          # CLI 主程序
│   └── api/                 # HTTP API 服务 🆕
│       └── main.go          # API 服务入口
│
├── pkg/                      # 公共库代码
│   ├── parser/               # Markdown → HTML 转换
│   │   ├── parser.go         # Goldmark 解析器实现
│   │   └── template.go       # HTML 模板和样式系统
│   │
│   ├── renderer/             # HTML → 图片渲染
│   │   └── renderer.go       # Rod 无头浏览器渲染器
│   │
│   ├── converter/            # 转换器协调层
│   │   └── converter.go      # 端到端转换逻辑
│   │
│   └── handlers/             # HTTP 处理器 🆕
│       ├── types.go          # 请求/响应数据结构
│       ├── convert.go        # 转换端点 (JSON + 上传)
│       └── middleware.go     # 中间件 (CORS, 日志, 恢复)
│
├── docs/                     # 文档 🆕
│   ├── ai-context/           # AI 上下文文档
│   │   ├── project-structure.md
│   │   └── docs-overview.md
│   ├── API.md                # HTTP API 完整文档 🆕
│   └── IMPLEMENTATION.md     # 实现说明 🆕
│
├── examples/                 # 示例文件
│   ├── basic.md              # 基础功能示例
│   ├── technical-doc.md      # 技术文档示例
│   └── api-test.sh           # API 测试脚本 🆕
│
├── testdata/                 # 测试数据
│   ├── input/
│   └── output/
│
├── go.mod                    # Go 模块定义
├── go.sum                    # 依赖校验和
├── README.md                 # 用户文档 (已更新)
├── QUICKSTART.md             # 快速开始指南 🆕
├── CLAUDE.md                 # 本文档
├── markdown2image            # CLI 可执行文件
└── markdown2image-api        # API 可执行文件 🆕
```

---

## 4. 核心架构设计

### 4.1 实现架构

```
Markdown 输入
    ↓
[Parser] Goldmark 解析器
    ↓ (HTML 内容)
[Template] 应用样式和主题
    ↓ (完整 HTML 文档)
[Renderer] Rod 无头浏览器
    ↓ (截图)
图片输出 (PNG/JPEG/WebP)
```

### 4.2 组件实现

#### Parser (`pkg/parser/`)
**实现**: 基于 Goldmark 的 Markdown → HTML 转换器

**核心功能**:
- Goldmark 解析器 (CommonMark 兼容)
- GFM 扩展 (表格、删除线、自动链接、任务列表)
- Chroma 代码语法高亮 (Monokai 主题)
- HTML 模板系统 (支持 light/dark 主题)

**关键文件**:
- `parser.go`: GoldmarkParser 实现,包含 Parse() 和 ParseToString() 方法
- `template.go`: WrapHTML() 函数,生成完整 HTML 文档和 CSS 样式

#### Renderer (`pkg/renderer/`)
**实现**: 基于 Rod 的 HTML → 图片渲染器

**核心功能**:
- Rod 无头浏览器自动化
- 全页截图支持
- 多格式输出 (PNG, JPEG, WebP)
- 自定义视口和设备像素比

**关键文件**:
- `renderer.go`: RodRenderer 实现,包含 RenderToImage() 和 RenderToFile() 方法

**渲染选项**:
```go
type RenderOptions struct {
    Width            int          // 视口宽度 (默认 1200)
    Height           int          // 自动高度
    Format           ImageFormat  // 图片格式
    Quality          int          // 质量 (JPEG/WebP)
    FullPage         bool         // 全页截图
    DevicePixelRatio float64      // 像素比 (高清屏)
}
```

#### Converter (`pkg/converter/`)
**实现**: 端到端转换协调器

**核心功能**:
- 协调 Parser 和 Renderer
- 统一的配置管理
- 文件到文件的转换接口

**关键文件**:
- `converter.go`: DefaultConverter 实现,包含 Convert() 和 ConvertFile() 方法

**转换流程**:
```go
1. Markdown → HTML (Parser.Parse)
2. HTML → 完整文档 (WrapHTML + 模板)
3. HTML 文档 → 图片 (Renderer.RenderToImage)
```

### 4.3 CLI 工具 (`cmd/markdown2image/`)

**实现**: 基于 Go 标准库 `flag` 的命令行工具

**支持参数**:
- `-input`: 输入 Markdown 文件 (必需)
- `-output`: 输出图片文件 (必需)
- `-theme`: 主题 (light/dark)
- `-width`: 页面宽度
- `-font-size`: 字体大小
- `-format`: 图片格式 (png/jpeg/webp)
- `-quality`: 图片质量 (1-100)
- `-dpr`: 设备像素比

**使用示例**:
```bash
./markdown2image -input doc.md -output doc.png -theme dark -width 1400
```

---

## 5. 技术栈和依赖

### 5.1 核心依赖 (已集成)

| 依赖 | 版本 | 用途 |
|------|------|------|
| `github.com/yuin/goldmark` | v1.7.13 | Markdown 解析 (CommonMark + GFM) |
| `github.com/alecthomas/chroma/v2` | v2.20.0 | 代码语法高亮 |
| `github.com/yuin/goldmark-highlighting/v2` | v2.0.0 | Goldmark 高亮扩展 |
| `github.com/go-rod/rod` | v0.116.2 | 无头浏览器自动化 |

### 5.2 支持库

| 依赖 | 版本 | 用途 |
|------|------|------|
| `github.com/dlclark/regexp2` | v1.11.5 | Chroma 正则表达式 |
| `github.com/ysmood/*` | - | Rod 依赖库 (goob, gson, got, fetchup, leakless) |

### 5.3 依赖管理

```bash
# 查看依赖
go list -m all

# 更新依赖
go get -u ./...

# 清理未使用依赖
go mod tidy

# 验证依赖
go mod verify
```

### 5.4 浏览器依赖

Rod 首次运行时会自动下载 Chromium:
- **位置**: `~/.cache/rod/browser/chromium-*`
- **大小**: 约 150MB
- **版本**: Chromium 1321438 (自动管理)

---

## 6. 编码标准与 Go 最佳实践

### 6.1 通用指令
- 遵循 Go 官方风格指南 ([Effective Go](https://go.dev/doc/effective_go))
- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 进行静态分析
- 遵循 KISS、YAGNI 和 DRY 原则
- 优先使用标准库,然后是成熟的第三方库

### 6.2 命名约定
- **包名**: 小写单词,无下划线 (例如: `parser`, `renderer`)
- **文件名**: 小写,使用下划线分隔 (例如: `image_renderer.go`)
- **类型**: PascalCase (例如: `ImageRenderer`, `ConvertOptions`)
- **函数/方法**: PascalCase (公开) 或 camelCase (私有) (例如: `Parse()`, `calculateWidth()`)
- **常量**: PascalCase 或 UPPER_SNAKE_CASE (例如: `DefaultWidth` 或 `MAX_IMAGE_SIZE`)
- **接口**: 以 `-er` 结尾 (例如: `Parser`, `Renderer`, `Converter`)

### 6.3 接口设计原则
- **小接口**: 优先定义小而专注的接口 (例如: `io.Reader`, `io.Writer`)
- **隐式实现**: 无需显式声明实现接口
- **接口隔离**: 客户端不应依赖不需要的方法
- **接受接口,返回结构体**: 函数参数使用接口,返回具体类型

```go
// 良好示例
func ProcessMarkdown(input io.Reader, output io.Writer) error {
    // 接受接口,提高可测试性
}

// 不推荐
func ProcessMarkdown(input *os.File, output *os.File) error {
    // 过于具体,难以测试
}
```

### 6.4 错误处理
- **显式错误处理**: 不忽略错误,明确处理每个错误
- **错误包装**: 使用 `fmt.Errorf` 和 `%w` 包装错误
- **自定义错误**: 对领域特定错误定义自定义类型

```go
// 错误包装
if err != nil {
    return fmt.Errorf("failed to parse markdown: %w", err)
}

// 自定义错误
type ParseError struct {
    Line   int
    Column int
    Msg    string
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("parse error at %d:%d: %s", e.Line, e.Column, e.Msg)
}
```

### 6.5 并发安全
- **不可变性**: 优先使用不可变数据结构
- **互斥锁**: 使用 `sync.Mutex` 保护共享状态
- **通道**: 使用 `chan` 进行 goroutine 间通信
- **避免数据竞争**: 使用 `go run -race` 检测数据竞争

### 6.6 文档要求
- **包文档**: 每个包需要包级文档 (在 `doc.go` 或包的主文件中)
- **公开类型/函数**: 必须有文档注释
- **示例**: 为关键功能提供示例代码 (使用 `Example` 测试)

```go
// Package parser 提供 Markdown 解析功能
//
// 本包支持 CommonMark 标准,并可通过扩展支持额外语法
package parser

// Parse 解析 Markdown 文本并返回 AST
//
// 参数:
//   - input: Markdown 文本字节数组
//
// 返回:
//   - ast.Node: 抽象语法树根节点
//   - error: 解析错误 (如有)
func Parse(input []byte) (ast.Node, error) {
    // 实现
}
```

### 6.7 测试标准
- **单元测试**: 每个公开函数都应有测试
- **表驱动测试**: 使用表驱动模式测试多个场景
- **测试覆盖率**: 目标至少 80% 覆盖率
- **基准测试**: 为性能关键路径提供基准测试

```go
// 表驱动测试示例
func TestParse(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    ast.NodeType
        wantErr bool
    }{
        {"heading", "# Title", ast.Heading, false},
        {"paragraph", "Hello", ast.Paragraph, false},
        {"empty", "", ast.Document, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Parse([]byte(tt.input))
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got.Type() != tt.want {
                t.Errorf("Parse() = %v, want %v", got.Type(), tt.want)
            }
        })
    }
}
```

---

## 7. 开发命令参考

### 7.1 依赖管理
```bash
# 初始化模块 (已完成)
go mod init Gomarkdown2image

# 添加依赖
go get <package>

# 下载所有依赖
go mod download

# 清理未使用的依赖
go mod tidy

# 验证依赖
go mod verify
```

### 7.2 构建与运行
```bash
# 构建可执行文件
go build -o markdown2image ./cmd/markdown2image

# 运行 (开发模式)
go run ./cmd/markdown2image -input input.md -output output.png

# 交叉编译 (例如: macOS 编译 Linux 版本)
GOOS=linux GOARCH=amd64 go build -o markdown2image-linux ./cmd/markdown2image

# 生成优化的发布版本
go build -ldflags="-s -w" -o markdown2image ./cmd/markdown2image
```

### 7.3 测试
```bash
# 运行所有测试
go test ./...

# 详细输出
go test -v ./...

# 测试覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 基准测试
go test -bench=. -benchmem ./...

# 数据竞争检测
go test -race ./...
```

### 7.4 代码质量
```bash
# 格式化代码
go fmt ./...
# 或使用 gofmt
gofmt -s -w .

# 静态分析
go vet ./...

# 使用 golangci-lint (需要安装)
golangci-lint run

# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 7.5 依赖可视化
```bash
# 查看依赖树
go mod graph

# 为什么需要某个依赖
go mod why <package>
```

---

## 8. 实现路线图

### 阶段 1: 项目初始化 ✅ 完成
- [x] 创建 `go.mod`
- [x] 编写 AI 上下文文档 (`CLAUDE.md`)
- [x] 创建项目目录结构
- [x] 创建 `README.md`

### 阶段 2: MVP 功能 ✅ 完成
- [x] 实现 Markdown → HTML Parser (Goldmark)
- [x] 实现 HTML → 图片 Renderer (Rod)
- [x] 实现 Converter 核心逻辑
- [x] 支持完整 Markdown 语法 (CommonMark + GFM)
- [x] 代码语法高亮 (Chroma)
- [x] 多格式输出 (PNG, JPEG, WebP)
- [x] CLI 工具实现

### 阶段 3: 样式系统 ✅ 完成
- [x] HTML 模板系统
- [x] 亮色/暗色主题
- [x] 代码块样式 (Monokai 主题)
- [x] 表格渲染
- [x] 引用块和列表样式

### 阶段 3.5: HTTP API 服务 ✅ 完成 (2025-12-14)
- [x] Gin 框架集成
- [x] POST /api/convert 端点 (JSON 方式)
- [x] POST /api/upload 端点 (文件上传方式)
- [x] 请求参数验证 (10 个配置参数)
- [x] CORS 中间件
- [x] 日志和错误恢复中间件
- [x] 健康检查端点
- [x] 完整 API 文档
- [x] 测试脚本和示例

### 阶段 4: AI 增强功能 🚧 规划中
- [ ] Claude API 集成
- [ ] Ollama 本地模型支持
- [ ] AI 内容润色和增强
- [ ] 多语言翻译
- [ ] 代码解释功能

### 阶段 5: 高级特性和优化 📋 待定
- [ ] 自定义 CSS 模板
- [ ] 配置文件支持 (YAML/JSON)
- [ ] 批量转换
- [ ] 性能优化 (大文件处理)
- [ ] Docker 镜像
- [ ] Web UI 界面

---

## 9. 开发注意事项

### 9.1 图像格式支持
- **优先级**: PNG (无损,广泛支持) > JPEG (有损,文件小) > WebP (现代格式) > SVG (矢量)
- **默认**: PNG
- **实现顺序**: PNG → JPEG → WebP → SVG

### 9.2 字体处理
- **跨平台兼容性**:
  - 内置默认字体 (嵌入到二进制)
  - 支持系统字体路径查找
  - 提供字体 fallback 机制
- **推荐字体**:
  - 西文: Roboto、Open Sans
  - 中文: Noto Sans CJK、Source Han Sans
  - 等宽: JetBrains Mono、Fira Code

### 9.3 样式配置
- **配置层级**: 默认配置 < 配置文件 < 命令行参数
- **主题系统**: 支持预定义主题 (Light、Dark、Solarized 等)
- **自定义选项**:
  - 字体族和大小
  - 颜色方案 (背景、前景、代码块、链接等)
  - 边距和间距
  - 图像尺寸

### 9.4 性能优化
- **大文件处理**:
  - 分页渲染 (单个图像 vs 多页)
  - 流式处理 AST
  - 延迟加载字体
- **内存管理**:
  - 复用图像缓冲区
  - 及时释放 AST 节点
- **并发**:
  - 考虑并行渲染多页
  - 字体加载可异步

### 9.5 错误处理
- **文件 I/O**: 明确的文件读写错误信息
- **解析错误**: 指出错误位置 (行号、列号)
- **渲染错误**: 字体缺失、颜色格式错误等
- **优雅降级**: 遇到不支持的语法时,回退到基础渲染

---

## 10. 安全和质量保证

### 10.1 输入验证
- 验证 Markdown 文件大小 (防止 OOM)
- 验证输出路径 (防止路径遍历)
- 验证配置参数 (字体大小、图像尺寸等)

### 10.2 依赖安全
- 定期更新依赖 (`go get -u ./...`)
- 使用 `go mod verify` 验证依赖完整性
- 关注依赖的安全公告

### 10.3 代码审查
- 使用 `go vet` 和 `golangci-lint` 静态分析
- 编写全面的单元测试
- 进行代码审查 (如果是团队协作)

---

## 11. 任务完成后检查清单

完成任何开发任务后,执行以下检查:

### 1. 代码质量
```bash
# 格式化
go fmt ./...

# 静态分析
go vet ./...

# Linting (如果已安装 golangci-lint)
golangci-lint run
```

### 2. 测试
```bash
# 运行测试
go test ./...

# 覆盖率检查
go test -cover ./...

# 数据竞争检测 (如有并发代码)
go test -race ./...
```

### 3. 构建验证
```bash
# 确保可以成功构建
go build ./...
```

### 4. 依赖管理
```bash
# 清理未使用的依赖
go mod tidy

# 验证依赖
go mod verify
```

---

## 12. 相关资源

### 官方文档
- [Go 官方网站](https://go.dev/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go 代码审查指南](https://github.com/golang/go/wiki/CodeReviewComments)

### 核心依赖文档
- [Goldmark](https://github.com/yuin/goldmark) - Markdown 解析器
- [gg](https://github.com/fogleman/gg) - 2D 图形库
- [Freetype](https://github.com/golang/freetype) - 字体渲染

### 工具
- [golangci-lint](https://golangci-lint.run/) - Go linters 聚合工具
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [Viper](https://github.com/spf13/viper) - 配置管理

---

**文档版本**: 2025-12-12
**项目阶段**: 初始化 - 架构规划
**Go 版本**: 1.25.1
**维护者**: AI 代理 + 开发团队
