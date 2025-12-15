#!/bin/bash

# AI 增强 Markdown 转换示例脚本
# 展示如何使用 API 进行 AI 增强的 Markdown 转换

API_URL="http://localhost:8080/api/convert"

echo "=== Gomarkdown2image AI 功能示例 ==="
echo ""

# 示例 1: 使用 Gemini 增强内容
echo "1. 使用 Gemini AI 增强 Markdown 内容..."
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# AI Test\nThis is a simple test document.",
    "parserMode": "ai",
    "aiProvider": "gemini",
    "aiModel": "gemini-2.0-flash-exp",
    "aiApiKey": "YOUR_GEMINI_API_KEY",
    "aiPromptTemplate": "enhance",
    "theme": "dark",
    "width": 1200
  }' \
  --output gemini-enhanced.png

echo "✓ 图片已保存为 gemini-enhanced.png"
echo ""

# 示例 2: 使用 Ollama 本地模型
echo "2. 使用 Ollama 本地模型增强内容..."
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# Local AI Test\nThis uses local Ollama model.",
    "parserMode": "ai",
    "aiProvider": "ollama",
    "aiModel": "llama3.2",
    "aiEndpoint": "http://localhost:11434",
    "aiPromptTemplate": "enhance",
    "theme": "light"
  }' \
  --output ollama-enhanced.png

echo "✓ 图片已保存为 ollama-enhanced.png"
echo ""

# 示例 3: 自定义提示词
echo "3. 使用自定义提示词..."
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# Custom Prompt Test\nOriginal content here.",
    "parserMode": "ai",
    "aiProvider": "gemini",
    "aiApiKey": "YOUR_GEMINI_API_KEY",
    "aiCustomPrompt": "请将以下内容改写为技术文档风格,并添加更多细节:",
    "theme": "light"
  }' \
  --output custom-prompt.png

echo "✓ 图片已保存为 custom-prompt.png"
echo ""

# 示例 4: 传统模式 (不使用 AI)
echo "4. 传统模式 (无 AI 增强)..."
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# Traditional Mode\nNo AI enhancement.",
    "parserMode": "traditional",
    "theme": "light"
  }' \
  --output traditional.png

echo "✓ 图片已保存为 traditional.png"
echo ""

echo "=== 所有示例完成 ==="
echo ""
echo "注意:"
echo "1. 使用 Gemini 需要有效的 API Key (https://ai.google.dev/)"
echo "2. 使用 Ollama 需要本地运行 Ollama 服务 (ollama serve)"
echo "3. AI 模式转换时间较长 (10-30秒),请耐心等待"
echo "4. 建议先使用 Ollama 本地模型测试,无需 API Key"
