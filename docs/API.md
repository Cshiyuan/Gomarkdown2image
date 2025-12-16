# Gomarkdown2image API 文档

## 概述

Gomarkdown2image API 提供了将 Markdown 文本转换为图片的 HTTP 接口,支持多种输出格式、自定义样式选项,以及 **AI 增强功能**。

**版本**: v0.2.1
**基础 URL**: `http://localhost:8080`

**核心功能**:
- ✅ 传统 Markdown → 图片转换 (Goldmark + Chroma 语法高亮)
- ✅ **AI 增强模式** (Gemini + Ollama 双后端支持) 🆕
- ✅ 5 种内置提示词模板 (润色、翻译、格式化、代码解释、总结)
- ✅ 多种输出格式 (PNG, JPEG, WebP)
- ✅ 主题系统 (light, dark)
- ✅ 完整的安全防护 (XSS, 输入验证, 并发安全)

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

**基础参数**:

| 参数 | 类型 | 必需 | 默认值 | 说明 | 验证规则 |
|------|------|------|--------|------|----------|
| `markdown` | string | ✅ | - | Markdown 内容 | 最大 10MB |
| `title` | string | ❌ | "Markdown to Image" | 页面标题 | - |
| `theme` | string | ❌ | "light" | 主题 | `light` 或 `dark` |
| `customCss` | string | ❌ | "" | 自定义 CSS | 最大 100KB,XSS 防护 |
| `width` | integer | ❌ | 1200 | 页面宽度(px) | 200-4000 |
| `fontSize` | integer | ❌ | 16 | 字体大小(px) | 8-72 |
| `fontFamily` | string | ❌ | "Arial, sans-serif" | 字体族 | CSS font-family |
| `imageFormat` | string | ❌ | "png" | 图片格式 | `png`, `jpeg`, `webp` |
| `imageQuality` | integer | ❌ | 90 | 图片质量 | 1-100 (仅 JPEG/WebP) |
| `devicePixelRatio` | number | ❌ | 1.0 | 设备像素比 | 0.5-4.0 |

**AI 增强参数** 🆕:

| 参数 | 类型 | 必需 | 默认值 | 说明 | 验证规则 |
|------|------|------|--------|------|----------|
| `parserMode` | string | ❌ | "traditional" | 解析器模式 | `traditional` 或 `ai` |
| `aiProvider` | string | ❌ | "gemini" | AI 提供器 | `gemini` 或 `ollama` |
| `aiModel` | string | ❌ | "gemini-2.0-flash-exp" | AI 模型名称 | 依赖提供器 |
| `aiApiKey` | string | ❌ | - | AI API 密钥 | Gemini 必需 |
| `aiEndpoint` | string | ❌ | "http://localhost:11434" | AI 服务端点 | Ollama 使用 |
| `aiPromptTemplate` | string | ❌ | "enhance" | 提示词模板 | 见下方模板列表 |
| `aiCustomPrompt` | string | ❌ | - | 自定义提示词 | 覆盖模板 |

**可用的提示词模板**:
- `enhance`: 润色和优化 Markdown 内容,提升可读性
- `translate`: 翻译文档 (需配合 `aiPromptData` 指定目标语言)
- `format`: 格式化和美化 Markdown 结构
- `explain_code`: 为代码块添加解释和注释
- `summarize`: 生成文档摘要和关键要点

#### 请求示例

**1. 传统模式 (不使用 AI)**:

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

**2. AI 增强模式 - 使用 Gemini 润色内容** 🆕:

```bash
curl -X POST http://localhost:8080/api/convert \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# Project\n\nThis project is good.",
    "parserMode": "ai",
    "aiProvider": "gemini",
    "aiModel": "gemini-2.0-flash-exp",
    "aiApiKey": "YOUR_GEMINI_API_KEY",
    "aiPromptTemplate": "enhance",
    "theme": "light",
    "imageFormat": "png"
  }' \
  --output enhanced.png
```

**3. AI 增强模式 - 使用 Ollama 本地翻译** 🆕:

```bash
curl -X POST http://localhost:8080/api/convert \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# Hello World\n\nThis is a technical document.",
    "parserMode": "ai",
    "aiProvider": "ollama",
    "aiModel": "llama3.2",
    "aiEndpoint": "http://localhost:11434",
    "aiPromptTemplate": "translate",
    "theme": "dark",
    "imageFormat": "webp"
  }' \
  --output translated.webp
```

**4. AI 增强模式 - 自定义提示词** 🆕:

```json
{
  "markdown": "# API Documentation\n\n## Overview\nThis is an API...",
  "parserMode": "ai",
  "aiProvider": "gemini",
  "aiApiKey": "YOUR_KEY",
  "aiCustomPrompt": "请将以下 Markdown 文档转换为更专业的技术文档格式,添加必要的图表说明和代码注释:",
  "theme": "light",
  "imageFormat": "png"
}
```

**5. 完整参数示例**:

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

**基础字段**:

| 字段名 | 类型 | 必需 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `file` | file | ✅ | - | Markdown 文件 (最大 10MB) |
| `title` | string | ❌ | "Markdown to Image" | 页面标题 |
| `theme` | string | ❌ | "light" | 主题 (`light`/`dark`) |
| `width` | integer | ❌ | 1200 | 页面宽度 |
| `fontSize` | integer | ❌ | 16 | 字体大小 |
| `fontFamily` | string | ❌ | "Arial, sans-serif" | 字体族 |
| `customCss` | string | ❌ | "" | 自定义 CSS (最大 100KB) |
| `imageFormat` | string | ❌ | "png" | 图片格式 |
| `imageQuality` | integer | ❌ | 90 | 图片质量 |
| `devicePixelRatio` | number | ❌ | 1.0 | 设备像素比 |

**AI 增强字段** 🆕:

| 字段名 | 类型 | 必需 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `parserMode` | string | ❌ | "traditional" | 解析器模式 (`traditional`/`ai`) |
| `aiProvider` | string | ❌ | "gemini" | AI 提供器 (`gemini`/`ollama`) |
| `aiModel` | string | ❌ | "gemini-2.0-flash-exp" | AI 模型名称 |
| `aiApiKey` | string | ❌ | - | AI API 密钥 |
| `aiEndpoint` | string | ❌ | "http://localhost:11434" | AI 服务端点 |
| `aiPromptTemplate` | string | ❌ | "enhance" | 提示词模板 |
| `aiCustomPrompt` | string | ❌ | - | 自定义提示词 |

#### 请求示例

**1. 基础上传 (传统模式)**:

```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@document.md" \
  --output output.png
```

**2. 带参数上传**:

```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@document.md" \
  -F "theme=dark" \
  -F "width=1600" \
  -F "imageFormat=webp" \
  -F "imageQuality=95" \
  --output output.webp
```

**3. AI 增强模式上传 - Gemini 润色** 🆕:

```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@document.md" \
  -F "parserMode=ai" \
  -F "aiProvider=gemini" \
  -F "aiApiKey=YOUR_GEMINI_API_KEY" \
  -F "aiPromptTemplate=enhance" \
  -F "theme=light" \
  -F "imageFormat=png" \
  --output enhanced.png
```

**4. AI 增强模式上传 - Ollama 翻译** 🆕:

```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@technical-doc.md" \
  -F "parserMode=ai" \
  -F "aiProvider=ollama" \
  -F "aiModel=llama3.2" \
  -F "aiPromptTemplate=translate" \
  -F "theme=dark" \
  --output translated.png
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

**通用错误**:

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

**AI 相关错误** 🆕:

| 错误代码 | HTTP 状态 | 说明 |
|----------|-----------|------|
| `AI_AUTH_FAILED` | 400 | AI API 认证失败 (无效的 API Key) |
| `AI_RATE_LIMIT` | 429 | AI 服务速率限制 |
| `AI_TIMEOUT` | 504 | AI 服务超时 |
| `AI_NETWORK_ERROR` | 502 | AI 服务网络错误 |
| `AI_SERVER_ERROR` | 500 | AI 服务内部错误 |
| `AI_MODEL_NOT_FOUND` | 400 | AI 模型不存在 |
| `AI_INVALID_PROMPT` | 400 | 无效的提示词模板 |

---

## 使用示例

### Python 示例

**1. 传统模式 - JSON 转换**:

```python
import requests

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

**2. AI 增强模式 - Gemini 润色** 🆕:

```python
import requests
import os

# 从环境变量获取 API Key
gemini_api_key = os.getenv('GEMINI_API_KEY', 'YOUR_API_KEY')

response = requests.post(
    'http://localhost:8080/api/convert',
    json={
        'markdown': '# My Project\n\nThis is a simple project.',
        'parserMode': 'ai',
        'aiProvider': 'gemini',
        'aiModel': 'gemini-2.0-flash-exp',
        'aiApiKey': gemini_api_key,
        'aiPromptTemplate': 'enhance',
        'theme': 'light',
        'imageFormat': 'png'
    }
)

if response.status_code == 200:
    with open('enhanced.png', 'wb') as f:
        f.write(response.content)
    print('✅ AI 增强转换成功!')
else:
    error = response.json()
    print(f'❌ 转换失败: {error}')
```

**3. AI 增强模式 - Ollama 本地翻译** 🆕:

```python
import requests

response = requests.post(
    'http://localhost:8080/api/convert',
    json={
        'markdown': '# Technical Document\n\nThis document explains...',
        'parserMode': 'ai',
        'aiProvider': 'ollama',
        'aiModel': 'llama3.2',
        'aiEndpoint': 'http://localhost:11434',
        'aiPromptTemplate': 'translate',
        'theme': 'dark',
        'imageFormat': 'webp'
    }
)

if response.status_code == 200:
    with open('translated.webp', 'wb') as f:
        f.write(response.content)
    print('✅ AI 翻译转换成功!')
```

**4. 文件上传 - 传统模式**:

```python
import requests

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
        print('✅ 文件上传转换成功!')
```

**5. 文件上传 - AI 增强模式** 🆕:

```python
import requests
import os

gemini_api_key = os.getenv('GEMINI_API_KEY')

with open('document.md', 'rb') as f:
    files = {'file': f}
    data = {
        'parserMode': 'ai',
        'aiProvider': 'gemini',
        'aiApiKey': gemini_api_key,
        'aiPromptTemplate': 'enhance',
        'theme': 'light',
        'imageFormat': 'png'
    }
    response = requests.post(
        'http://localhost:8080/api/upload',
        files=files,
        data=data
    )

    if response.status_code == 200:
        with open('enhanced.png', 'wb') as out:
            out.write(response.content)
        print('✅ AI 增强文件转换成功!')
```

### JavaScript (Node.js) 示例

**1. 传统模式 - JSON 转换**:

```javascript
const axios = require('axios');
const fs = require('fs');

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

convertMarkdown();
```

**2. AI 增强模式 - Gemini 润色** 🆕:

```javascript
const axios = require('axios');
const fs = require('fs');

async function convertWithAI() {
  try {
    const response = await axios.post('http://localhost:8080/api/convert', {
      markdown: '# My Project\n\nThis project is good.',
      parserMode: 'ai',
      aiProvider: 'gemini',
      aiModel: 'gemini-2.0-flash-exp',
      aiApiKey: process.env.GEMINI_API_KEY,
      aiPromptTemplate: 'enhance',
      theme: 'light',
      imageFormat: 'png'
    }, {
      responseType: 'arraybuffer'
    });

    fs.writeFileSync('enhanced.png', response.data);
    console.log('✅ AI 增强转换成功!');
  } catch (error) {
    console.error('❌ 转换失败:', error.response?.data || error.message);
  }
}

convertWithAI();
```

**3. AI 增强模式 - Ollama 本地处理** 🆕:

```javascript
const axios = require('axios');
const fs = require('fs');

async function convertWithOllama() {
  const response = await axios.post('http://localhost:8080/api/convert', {
    markdown: '# Technical Document\n\nThis is a technical guide...',
    parserMode: 'ai',
    aiProvider: 'ollama',
    aiModel: 'llama3.2',
    aiEndpoint: 'http://localhost:11434',
    aiPromptTemplate: 'format',
    theme: 'dark',
    imageFormat: 'webp'
  }, {
    responseType: 'arraybuffer'
  });

  fs.writeFileSync('formatted.webp', response.data);
  console.log('✅ Ollama AI 格式化成功!');
}

convertWithOllama();
```

**4. 文件上传 - 传统模式**:

```javascript
const axios = require('axios');
const fs = require('fs');
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
  console.log('✅ 文件上传转换成功!');
}

uploadMarkdown();
```

**5. 文件上传 - AI 增强模式** 🆕:

```javascript
const axios = require('axios');
const fs = require('fs');
const FormData = require('form-data');

async function uploadWithAI() {
  const form = new FormData();
  form.append('file', fs.createReadStream('document.md'));
  form.append('parserMode', 'ai');
  form.append('aiProvider', 'gemini');
  form.append('aiApiKey', process.env.GEMINI_API_KEY);
  form.append('aiPromptTemplate', 'enhance');
  form.append('theme', 'light');
  form.append('imageFormat', 'png');

  const response = await axios.post(
    'http://localhost:8080/api/upload',
    form,
    {
      headers: form.getHeaders(),
      responseType: 'arraybuffer'
    }
  );

  fs.writeFileSync('enhanced.png', response.data);
  console.log('✅ AI 增强文件转换成功!');
}

uploadWithAI();
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

**基础配置**:

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | 8080 | 服务监听端口 |
| `GIN_MODE` | debug | Gin 运行模式 (`debug`/`release`) |
| `ALLOWED_ORIGINS` | * | CORS 允许的源 (生产环境应指定具体域名) |

**AI 服务配置** 🆕:

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `GEMINI_API_KEY` | - | Google Gemini API 密钥 |
| `OLLAMA_ENDPOINT` | http://localhost:11434 | Ollama 服务端点 |
| `AI_TIMEOUT` | 60 | AI 请求超时时间 (秒) |
| `AI_MAX_RETRIES` | 3 | AI 请求最大重试次数 |

### 生产环境配置

**基础配置**:

```bash
# 设置为 release 模式 (关闭调试日志)
export GIN_MODE=release

# 指定端口
export PORT=3000

# 配置 CORS
export ALLOWED_ORIGINS="https://yourdomain.com,https://app.yourdomain.com"

# 启动服务
./markdown2image-api
```

**AI 增强配置** 🆕:

```bash
# Gemini API 配置
export GEMINI_API_KEY="your-gemini-api-key-here"

# Ollama 本地配置 (如果使用)
export OLLAMA_ENDPOINT="http://localhost:11434"

# AI 超时和重试
export AI_TIMEOUT=60
export AI_MAX_RETRIES=3

# 启动服务
./markdown2image-api
```

### 获取 Gemini API Key

1. 访问 [Google AI Studio](https://ai.google.dev/)
2. 登录 Google 账号
3. 点击 "Get API Key"
4. 创建新的 API Key
5. 复制 API Key 并设置到环境变量

### 安装 Ollama (本地 AI)

```bash
# macOS
brew install ollama

# Linux
curl -fsSL https://ollama.com/install.sh | sh

# 启动 Ollama 服务
ollama serve

# 下载模型
ollama pull llama3.2
ollama pull mistral
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

### 基础功能

**Q: 为什么首次转换很慢?**
A: 首次运行时 Rod 需要下载 Chromium 浏览器 (~150MB),后续会复用已下载的浏览器。

**Q: 支持哪些 Markdown 语法?**
A: 支持 CommonMark 标准 + GitHub Flavored Markdown (GFM),包括表格、删除线、任务列表、代码高亮等。

**Q: 如何自定义样式?**
A: 使用 `customCss` 参数注入自定义 CSS:
```json
{
  "markdown": "# Title",
  "customCss": ".container { background: #f0f0f0; padding: 40px; }"
}
```

**Q: 支持数学公式吗?**
A: 当前版本不支持 LaTeX/MathJax,后续版本会添加。

**Q: 如何提高转换性能?**
A: 见"性能优化建议"章节,使用连接池、缓存和限流。

### AI 增强功能 🆕

**Q: AI 增强模式和传统模式有什么区别?**
A:
- **传统模式**: 直接使用 Goldmark 解析 Markdown,速度快,结果可预测
- **AI 增强模式**: 先使用 AI 优化内容(润色、翻译、格式化等),再转换为图片,质量更高但耗时更长

**Q: Gemini 和 Ollama 应该选哪个?**
A:
- **Gemini**: 云端 AI,速度快,质量高,需要 API Key 和网络连接,有速率限制
- **Ollama**: 本地 AI,完全离线,无速率限制,免费,但需要本地硬件资源(CPU/内存)

**Q: 如何获取 Gemini API Key?**
A: 访问 [Google AI Studio](https://ai.google.dev/),登录后点击 "Get API Key" 创建。免费版有速率限制,付费版无限制。

**Q: Ollama 需要什么硬件配置?**
A:
- **最低**: 8GB RAM,支持 7B 参数模型 (如 llama3.2)
- **推荐**: 16GB+ RAM,支持更大模型 (如 llama3:70b)
- GPU 可选,但会显著加速

**Q: AI 增强模式失败会怎样?**
A: API 实现了自动降级机制。AI 服务失败时,会自动回退到传统 Goldmark 模式,确保转换总能成功。

**Q: AI 增强需要多长时间?**
A:
- **Gemini**: 3-10 秒 (取决于内容长度和网络)
- **Ollama**: 5-30 秒 (取决于模型大小和硬件)
- **传统模式**: < 1 秒

**Q: 可以自定义 AI 提示词吗?**
A: 可以。使用 `aiCustomPrompt` 参数覆盖默认模板:
```json
{
  "markdown": "# My Doc",
  "parserMode": "ai",
  "aiCustomPrompt": "请将以下文档转换为技术白皮书格式,添加专业术语解释:"
}
```

**Q: 支持哪些提示词模板?**
A: 内置 5 个模板:
- `enhance`: 润色和优化内容
- `translate`: 多语言翻译
- `format`: 格式化和美化
- `explain_code`: 代码解释
- `summarize`: 内容总结

**Q: AI 模式会修改我的原始 Markdown 吗?**
A: 不会。AI 只在内存中处理增强后的内容,原始文件/数据完全不受影响。

**Q: AI 模式是否安全?敏感信息会泄露吗?**
A:
- **Gemini**: 数据会发送到 Google 服务器,请勿用于敏感内容
- **Ollama**: 完全本地处理,数据不会离开你的机器,适合敏感内容

**Q: 如何调试 AI 模式的问题?**
A:
1. 检查环境变量 (GEMINI_API_KEY / OLLAMA_ENDPOINT)
2. 查看 API 返回的错误信息 (包含详细错误类型)
3. 尝试降低 `aiModel` 复杂度 (如使用 gemini-1.5-flash 而非 gemini-2.0-flash-exp)
4. 启用服务器日志 (`GIN_MODE=debug`)

---

## 更新日志

### v0.2.1 (2025-12-16) 🆕
- ✅ 生产级安全加固
  - Renderer panic 防护 (defer/recover 模式)
  - AI prompts 并发安全 (sync.RWMutex)
  - XSS 防护 (CustomCSS 验证,12 个禁止模式)
  - CORS 安全配置 (环境变量支持)
  - 超时控制优化 (30s 总超时 + 分级策略)
- ✅ 代码质量优化
  - RequestParams 接口抽象 (消除 116 行重复代码)
  - 接口设计模式应用 (ISP + DIP)
  - 单元测试覆盖扩展 (70+ 测试用例)

### v0.2.0 (2025-12-15) 🆕
- ✅ **AI 增强功能完整实现**
  - Google Gemini API 集成 (google/generative-ai-go v0.20.1)
  - Ollama 本地模型支持 (ollama/ollama v0.13.3)
  - Parser Provider 可插拔架构 (传统/AI 双模式)
  - 5 种内置提示词模板 (enhance, translate, format, explain_code, summarize)
  - AI 错误处理和自动降级机制
- ✅ HTTP API 扩展
  - 7 个 AI 增强参数 (parserMode, aiProvider, aiModel 等)
  - 完整的 AI 错误代码系统
- ✅ 文档更新
  - API 文档 AI 功能章节
  - AI 架构说明和使用示例

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
git clone https://github.com/Cshiyuan/Gomarkdown2image.git
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

- GitHub: https://github.com/Cshiyuan/Gomarkdown2image
- Issues: https://github.com/Cshiyuan/Gomarkdown2image/issues
