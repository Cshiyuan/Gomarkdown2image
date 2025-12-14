#!/bin/bash
# API 测试脚本

API_URL="http://localhost:8080"
OUTPUT_DIR="./testdata/output"

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

echo "🧪 开始测试 Gomarkdown2image API..."
echo ""

# 测试 1: 健康检查
echo "1️⃣ 测试健康检查端点..."
curl -s "$API_URL/health" | python3 -m json.tool
echo ""

# 测试 2: JSON 转换 (PNG)
echo "2️⃣ 测试 JSON 转换 (PNG 格式)..."
curl -X POST "$API_URL/api/convert" \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# API 测试\n\n这是通过 **JSON** 方式提交的 Markdown。\n\n```go\nfunc main() {\n    fmt.Println(\"Hello API\")\n}\n```",
    "theme": "light",
    "width": 1200,
    "imageFormat": "png"
  }' \
  --output "$OUTPUT_DIR/json-test.png" \
  -w "\n✅ PNG 图片已保存: %{size_download} 字节\n\n"

# 测试 3: JSON 转换 (WebP, Dark 主题)
echo "3️⃣ 测试 JSON 转换 (WebP 格式 + 暗色主题)..."
curl -X POST "$API_URL/api/convert" \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# Dark Theme Test\n\nThis is a **WebP** format test with dark theme.\n\n- Item 1\n- Item 2\n- Item 3",
    "theme": "dark",
    "width": 1400,
    "fontSize": 18,
    "imageFormat": "webp",
    "imageQuality": 95
  }' \
  --output "$OUTPUT_DIR/json-dark-webp.webp" \
  -w "\n✅ WebP 图片已保存: %{size_download} 字节\n\n"

# 测试 4: 文件上传 (如果存在示例文件)
if [ -f "examples/basic.md" ]; then
  echo "4️⃣ 测试文件上传端点..."
  curl -X POST "$API_URL/api/upload" \
    -F "file=@examples/basic.md" \
    -F "theme=light" \
    -F "imageFormat=png" \
    --output "$OUTPUT_DIR/upload-test.png" \
    -w "\n✅ 上传转换成功: %{size_download} 字节\n\n"
else
  echo "4️⃣ ⏭️  跳过文件上传测试 (未找到 examples/basic.md)"
  echo ""
fi

# 测试 5: 错误处理 (空内容)
echo "5️⃣ 测试错误处理 (空 Markdown)..."
curl -X POST "$API_URL/api/convert" \
  -H "Content-Type: application/json" \
  -d '{"markdown": ""}' \
  2>/dev/null | python3 -m json.tool
echo ""

# 测试 6: 错误处理 (无效主题)
echo "6️⃣ 测试错误处理 (无效主题)..."
curl -X POST "$API_URL/api/convert" \
  -H "Content-Type: application/json" \
  -d '{"markdown": "# Test", "theme": "invalid"}' \
  2>/dev/null | python3 -m json.tool
echo ""

# 总结
echo "✅ 测试完成!"
echo ""
echo "📁 输出文件位置: $OUTPUT_DIR"
ls -lh "$OUTPUT_DIR"
