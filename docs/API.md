# Gomarkdown2image API 文档

## 概述

Gomarkdown2image API 提供了将 Markdown 文本转换为图片的 HTTP 接口,支持多种输出格式和自定义样式选项。

**版本**: v0.1.0
**基础 URL**: `http://localhost:8080`

---

## 快速开始

### 启动服务

```bash
# 方式 1: 运行可执行文件
./markdown2image-api

# 方式 2: 从源码运行
go run cmd/api/main.go

# 方式 3: 指定端口
PORT=3000 ./markdown2image-api
```

服务启动后,访问 `http://localhost:8080` 查看可用端点。

---

## API 端点

### 1. 健康检查

**端点**: `GET /health`
**描述**: 检查服务运行状态

**响应示例**:
```json
{
  "success": true,
  "message": "服务运行正常",
  "data": {
    "status": "healthy",
    "timestamp": 1765706173
  }
}
```

---

### 2. JSON 转换 (推荐)

**端点**: `POST /api/convert`
**Content-Type**: `application/json`
**描述**: 接收 JSON 格式的 Markdown 内容,返回生成的图片

#### 请求参数

| 参数 | 类型 | 必需 | 默认值 | 说明 | 验证规则 |
|------|------|------|--------|------|----------|
| `markdown` | string | ✅ | - | Markdown 内容 | 最大 10MB |
| `title` | string | ❌ | "Markdown to Image" | 页面标题 | - |
| `theme` | string | ❌ | "light" | 主题 | `light` 或 `dark` |
| `customCss` | string | ❌ | "" | 自定义 CSS | - |
| `width` | integer | ❌ | 1200 | 页面宽度(px) | 200-4000 |
| `fontSize` | integer | ❌ | 16 | 字体大小(px) | 8-72 |
| `fontFamily` | string | ❌ | "Arial, sans-serif" | 字体族 | CSS font-family |
| `imageFormat` | string | ❌ | "png" | 图片格式 | `png`, `jpeg`, `webp` |
| `imageQuality` | integer | ❌ | 90 | 图片质量 | 1-100 (仅 JPEG/WebP) |
| `devicePixelRatio` | number | ❌ | 1.0 | 设备像素比 | 0.5-4.0 |

#### 请求示例

```bash
curl -X POST http://localhost:8080/api/convert \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# Hello World\n\nThis is **bold** and this is *italic*.",
    "theme": "dark",
    "width": 1400,
    "fontSize": 18,
    "imageFormat": "png"
  }' \
  --output output.png
```

```json
{
  "markdown": "# 标题\n\n这是一段**粗体**文本。\n\n```python\nprint('Hello')\n```",
  "theme": "dark",
  "width": 1400,
  "fontSize": 18,
  "fontFamily": "Georgia, serif",
  "imageFormat": "webp",
  "imageQuality": 95,
  "devicePixelRatio": 2.0
}
```

#### 响应

**成功 (200 OK)**:
- **Content-Type**: `image/png` / `image/jpeg` / `image/webp`
- **Body**: 二进制图片数据

**失败 (4xx/5xx)**:
```json
{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "请求参数验证失败",
    "details": "Key: 'ConvertRequest.Markdown' Error:Field validation for 'Markdown' failed on the 'required' tag"
  }
}
```

---

### 3. 文件上传转换

**端点**: `POST /api/upload`
**Content-Type**: `multipart/form-data`
**描述**: 上传 Markdown 文件并转换为图片

#### 表单字段

| 字段名 | 类型 | 必需 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `file` | file | ✅ | - | Markdown 文件 (最大 10MB) |
| `title` | string | ❌ | "Markdown to Image" | 页面标题 |
| `theme` | string | ❌ | "light" | 主题 (`light`/`dark`) |
| `width` | integer | ❌ | 1200 | 页面宽度 |
| `fontSize` | integer | ❌ | 16 | 字体大小 |
| `fontFamily` | string | ❌ | "Arial, sans-serif" | 字体族 |
| `customCss` | string | ❌ | "" | 自定义 CSS |
| `imageFormat` | string | ❌ | "png" | 图片格式 |
| `imageQuality` | integer | ❌ | 90 | 图片质量 |
| `devicePixelRatio` | number | ❌ | 1.0 | 设备像素比 |

#### 请求示例

```bash
# 基础上传
curl -X POST http://localhost:8080/api/upload \
  -F "file=@document.md" \
  --output output.png

# 带参数上传
curl -X POST http://localhost:8080/api/upload \
  -F "file=@document.md" \
  -F "theme=dark" \
  -F "width=1600" \
  -F "imageFormat=webp" \
  -F "imageQuality=95" \
  --output output.webp
```

#### HTML 表单示例

```html
<!DOCTYPE html>
<html>
<body>
  <form action="http://localhost:8080/api/upload" method="POST" enctype="multipart/form-data">
    <input type="file" name="file" accept=".md,.markdown" required>
    <select name="theme">
      <option value="light">亮色主题</option>
      <option value="dark">暗色主题</option>
    </select>
    <select name="imageFormat">
      <option value="png">PNG</option>
      <option value="jpeg">JPEG</option>
      <option value="webp">WebP</option>
    </select>
    <button type="submit">转换</button>
  </form>
</body>
</html>
```

---

## 错误代码

| 错误代码 | HTTP 状态 | 说明 |
|----------|-----------|------|
| `INVALID_REQUEST` | 400 | 请求参数验证失败 |
| `CONTENT_TOO_LARGE` | 400 | Markdown 内容过大 (>10MB) |
| `NO_FILE_UPLOADED` | 400 | 未找到上传文件 |
| `FILE_TOO_LARGE` | 400 | 文件过大 (>10MB) |
| `INVALID_FORM` | 400 | 表单参数验证失败 |
| `CONVERTER_INIT_FAILED` | 500 | 转换器初始化失败 |
| `CONVERSION_FAILED` | 500 | Markdown 转换失败 |
| `FILE_READ_FAILED` | 500 | 文件读取失败 |

---

## 使用示例

### Python 示例

```python
import requests

# JSON 转换
response = requests.post(
    'http://localhost:8080/api/convert',
    json={
        'markdown': '# Hello from Python\n\nThis is **bold** text.',
        'theme': 'dark',
        'imageFormat': 'png'
    }
)

if response.status_code == 200:
    with open('output.png', 'wb') as f:
        f.write(response.content)
    print('✅ 转换成功!')
else:
    print('❌ 转换失败:', response.json())
```

```python
# 文件上传
with open('document.md', 'rb') as f:
    files = {'file': f}
    data = {'theme': 'light', 'imageFormat': 'webp'}
    response = requests.post(
        'http://localhost:8080/api/upload',
        files=files,
        data=data
    )

    if response.status_code == 200:
        with open('output.webp', 'wb') as out:
            out.write(response.content)
```

### JavaScript (Node.js) 示例

```javascript
const axios = require('axios');
const fs = require('fs');

// JSON 转换
async function convertMarkdown() {
  const response = await axios.post('http://localhost:8080/api/convert', {
    markdown: '# Hello from Node.js\n\nThis is **bold** text.',
    theme: 'dark',
    imageFormat: 'png'
  }, {
    responseType: 'arraybuffer'
  });

  fs.writeFileSync('output.png', response.data);
  console.log('✅ 转换成功!');
}

// 文件上传
const FormData = require('form-data');
async function uploadMarkdown() {
  const form = new FormData();
  form.append('file', fs.createReadStream('document.md'));
  form.append('theme', 'dark');
  form.append('imageFormat', 'webp');

  const response = await axios.post(
    'http://localhost:8080/api/upload',
    form,
    {
      headers: form.getHeaders(),
      responseType: 'arraybuffer'
    }
  );

  fs.writeFileSync('output.webp', response.data);
}

convertMarkdown();
```

### JavaScript (浏览器 Fetch API) 示例

```javascript
// JSON 转换
async function convertMarkdown() {
  const response = await fetch('http://localhost:8080/api/convert', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      markdown: '# Hello from Browser\n\nThis is **bold** text.',
      theme: 'dark',
      imageFormat: 'png'
    })
  });

  if (response.ok) {
    const blob = await response.blob();
    const url = URL.createObjectURL(blob);

    // 显示图片
    const img = document.createElement('img');
    img.src = url;
    document.body.appendChild(img);

    // 或下载
    const a = document.createElement('a');
    a.href = url;
    a.download = 'output.png';
    a.click();
  }
}

// 文件上传
async function uploadMarkdown(fileInput) {
  const formData = new FormData();
  formData.append('file', fileInput.files[0]);
  formData.append('theme', 'dark');
  formData.append('imageFormat', 'png');

  const response = await fetch('http://localhost:8080/api/upload', {
    method: 'POST',
    body: formData
  });

  if (response.ok) {
    const blob = await response.blob();
    const url = URL.createObjectURL(blob);
    // 处理图片...
  }
}
```

---

## 配置

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | 8080 | 服务监听端口 |
| `GIN_MODE` | debug | Gin 运行模式 (`debug`/`release`) |

### 生产环境配置

```bash
# 设置为 release 模式 (关闭调试日志)
export GIN_MODE=release

# 指定端口
export PORT=3000

# 启动服务
./markdown2image-api
```

---

## CORS 配置

默认配置允许所有来源 (`AllowOrigins: ["*"]`),生产环境建议修改:

编辑 `pkg/handlers/middleware.go`:
```go
func SetupCORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"https://yourdomain.com"}, // 指定域名
        AllowMethods:     []string{"GET", "POST"},
        AllowHeaders:     []string{"Content-Type"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

---

## 性能优化建议

### 1. 复用转换器实例

当前每次请求都创建新的转换器,生产环境建议使用连接池:

```go
// 全局转换器池 (待实现)
var converterPool *sync.Pool

func init() {
    converterPool = &sync.Pool{
        New: func() interface{} {
            conv, _ := converter.NewConverter()
            return conv
        },
    }
}
```

### 2. 缓存结果

对相同 Markdown 内容进行缓存,避免重复渲染:

```go
// 使用 Redis 或内存缓存
cache.Set(md5(markdown), imageData, 1*time.Hour)
```

### 3. 限流保护

```bash
# 安装限流中间件
go get github.com/ulule/limiter/v3
```

---

## 部署

### Docker 部署 (推荐)

```dockerfile
# Dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o markdown2image-api ./cmd/api

FROM alpine:latest
RUN apk add --no-cache chromium
WORKDIR /app
COPY --from=builder /app/markdown2image-api .
EXPOSE 8080
CMD ["./markdown2image-api"]
```

```bash
# 构建镜像
docker build -t markdown2image-api .

# 运行容器
docker run -d -p 8080:8080 markdown2image-api
```

### Systemd 服务

```ini
# /etc/systemd/system/markdown2image-api.service
[Unit]
Description=Markdown to Image API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/markdown2image
ExecStart=/opt/markdown2image/markdown2image-api
Restart=always
Environment="GIN_MODE=release"
Environment="PORT=8080"

[Install]
WantedBy=multi-user.target
```

```bash
# 启动服务
sudo systemctl start markdown2image-api
sudo systemctl enable markdown2image-api
```

---

## 常见问题

### Q: 为什么首次转换很慢?
A: 首次运行时 Rod 需要下载 Chromium 浏览器 (~150MB),后续会复用已下载的浏览器。

### Q: 支持哪些 Markdown 语法?
A: 支持 CommonMark 标准 + GitHub Flavored Markdown (GFM),包括表格、删除线、任务列表、代码高亮等。

### Q: 如何自定义样式?
A: 使用 `customCss` 参数注入自定义 CSS:
```json
{
  "markdown": "# Title",
  "customCss": ".container { background: #f0f0f0; padding: 40px; }"
}
```

### Q: 支持数学公式吗?
A: 当前版本不支持 LaTeX/MathJax,后续版本会添加。

### Q: 如何提高转换性能?
A: 见"性能优化建议"章节,使用连接池、缓存和限流。

---

## 更新日志

### v0.1.0 (2025-12-14)
- ✅ 初始版本发布
- ✅ 支持 JSON 和文件上传两种转换方式
- ✅ 支持 PNG/JPEG/WebP 三种输出格式
- ✅ 支持 light/dark 主题
- ✅ 完整的参数验证和错误处理
- ✅ CORS 和日志中间件

---

## 贡献指南

欢迎提交 Issue 和 Pull Request!

**开发环境**:
```bash
git clone https://github.com/yourusername/Gomarkdown2image.git
cd Gomarkdown2image
go mod download
go run cmd/api/main.go
```

**测试**:
```bash
go test ./...
```

---

## 许可证

MIT License

---

## 联系方式

- GitHub: https://github.com/yourusername/Gomarkdown2image
- Issues: https://github.com/yourusername/Gomarkdown2image/issues
