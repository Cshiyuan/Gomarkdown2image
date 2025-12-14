# 快速开始指南

## 1️⃣ 编译项目

```bash
# 编译 API 服务
go build -o markdown2image-api ./cmd/api

# 编译 CLI 工具 (可选)
go build -o markdown2image ./cmd/markdown2image
```

## 2️⃣ 启动 API 服务

```bash
# 默认端口 8080
./markdown2image-api

# 或指定端口
PORT=3000 ./markdown2image-api
```

看到以下输出说明启动成功:
```
🚀 Gomarkdown2image API 服务启动中...
📡 监听端口: 8080
🌍 访问地址: http://localhost:8080
💚 健康检查: http://localhost:8080/health
```

## 3️⃣ 测试 API

### 方式 1: 使用 curl (JSON)

```bash
curl -X POST http://localhost:8080/api/convert \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# Hello World\n\nThis is **bold** and *italic* text.",
    "theme": "dark",
    "imageFormat": "png"
  }' \
  --output output.png
```

### 方式 2: 使用 curl (文件上传)

```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@examples/basic.md" \
  -F "theme=light" \
  -F "imageFormat=webp" \
  --output output.webp
```

### 方式 3: 使用测试脚本

```bash
# 运行完整测试套件
./examples/api-test.sh
```

## 4️⃣ 在代码中调用

### Python

```python
import requests

response = requests.post(
    'http://localhost:8080/api/convert',
    json={
        'markdown': '# Hello from Python',
        'theme': 'dark',
        'imageFormat': 'png'
    }
)

with open('output.png', 'wb') as f:
    f.write(response.content)
```

### JavaScript (Node.js)

```javascript
const axios = require('axios');
const fs = require('fs');

async function convert() {
  const response = await axios.post('http://localhost:8080/api/convert', {
    markdown: '# Hello from Node.js',
    theme: 'dark',
    imageFormat: 'png'
  }, {
    responseType: 'arraybuffer'
  });

  fs.writeFileSync('output.png', response.data);
}

convert();
```

### JavaScript (浏览器)

```javascript
async function convertMarkdown() {
  const response = await fetch('http://localhost:8080/api/convert', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      markdown: '# Hello from Browser',
      theme: 'dark',
      imageFormat: 'png'
    })
  });

  const blob = await response.blob();
  const url = URL.createObjectURL(blob);

  // 显示图片
  const img = document.createElement('img');
  img.src = url;
  document.body.appendChild(img);
}
```

## 5️⃣ 可用参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `markdown` | string | (必需) | Markdown 内容 |
| `title` | string | "Markdown to Image" | 页面标题 |
| `theme` | string | "light" | 主题 (light/dark) |
| `width` | integer | 1200 | 页面宽度 |
| `fontSize` | integer | 16 | 字体大小 |
| `fontFamily` | string | "Arial, sans-serif" | 字体族 |
| `imageFormat` | string | "png" | 图片格式 (png/jpeg/webp) |
| `imageQuality` | integer | 90 | 图片质量 (1-100) |
| `devicePixelRatio` | number | 1.0 | 设备像素比 |

## 6️⃣ 查看完整文档

- **API 文档**: [docs/API.md](docs/API.md)
- **实现说明**: [docs/IMPLEMENTATION.md](docs/IMPLEMENTATION.md)
- **项目 README**: [README.md](README.md)

## 常见问题

**Q: 首次启动很慢?**
A: Rod 首次运行会下载 Chromium (~150MB),后续启动会很快。

**Q: 支持哪些 Markdown 语法?**
A: 完整支持 CommonMark + GitHub Flavored Markdown (表格、删除线、任务列表、代码高亮等)。

**Q: 如何在生产环境部署?**
A: 参考 [docs/API.md](docs/API.md) 的部署章节,推荐使用 Docker 或 Systemd。

---

**需要帮助?** 查看完整文档或提交 Issue!
