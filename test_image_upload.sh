#!/bin/bash

# 图片上传测试脚本

echo "=== 图书图片上传测试 ==="
echo ""

# 1. 登录获取管理员token
echo "1. 正在登录获取管理员token..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@lib.com","password":"12345678"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ 登录失败！"
  echo "响应: $LOGIN_RESPONSE"
  exit 1
fi

echo "✅ 登录成功！"
echo "Token: ${TOKEN:0:50}..."
echo ""

# 2. 检查是否有图书
echo "2. 检查图书列表..."
BOOKS_RESPONSE=$(curl -s -X POST http://localhost:8080/api/book/search \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\",\"page\":1,\"limit\":1}")

BOOK_ID=$(echo $BOOKS_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$BOOK_ID" ]; then
  echo "⚠️  没有找到图书，请先添加一本图书"
  echo "响应: $BOOKS_RESPONSE"
  exit 1
fi

echo "✅ 找到图书 ID: $BOOK_ID"
echo ""

# 3. 检查图片文件
if [ -z "$1" ]; then
  echo "❌ 请提供图片文件路径"
  echo "用法: ./test_image_upload.sh <图片文件路径>"
  echo "示例: ./test_image_upload.sh ~/Downloads/test.jpg"
  exit 1
fi

IMAGE_FILE="$1"

if [ ! -f "$IMAGE_FILE" ]; then
  echo "❌ 图片文件不存在: $IMAGE_FILE"
  exit 1
fi

FILE_SIZE=$(stat -f%z "$IMAGE_FILE" 2>/dev/null || stat -c%s "$IMAGE_FILE" 2>/dev/null)
FILE_SIZE_MB=$(echo "scale=2; $FILE_SIZE / 1024 / 1024" | bc)

echo "3. 准备上传图片..."
echo "   文件: $IMAGE_FILE"
echo "   大小: ${FILE_SIZE_MB} MB"
echo ""

if (( $(echo "$FILE_SIZE_MB > 5" | bc -l) )); then
  echo "⚠️  警告: 图片大小超过 5MB，可能会上传失败"
  echo ""
fi

# 4. 上传图片
echo "4. 正在上传图片到图书 ID: $BOOK_ID..."
UPLOAD_RESPONSE=$(curl -s -X POST http://localhost:8080/api/book/uploadCover \
  -F "token=$TOKEN" \
  -F "book_id=$BOOK_ID" \
  -F "image=@$IMAGE_FILE")

echo "上传响应:"
echo "$UPLOAD_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$UPLOAD_RESPONSE"
echo ""

# 5. 检查上传结果
if echo "$UPLOAD_RESPONSE" | grep -q '"code":0'; then
  IMAGE_URL=$(echo $UPLOAD_RESPONSE | grep -o '"image_url":"[^"]*' | cut -d'"' -f4)
  echo "✅ 上传成功！"
  echo "图片 URL: $IMAGE_URL"
  echo ""
  echo "5. 验证图片URL..."
  echo "   在浏览器中打开: $IMAGE_URL"
  echo ""
  
  # 6. 查看图书详情
  echo "6. 查看图书详情（包含图片URL）..."
  DETAIL_RESPONSE=$(curl -s -X POST http://localhost:8080/api/book/detail \
    -H "Content-Type: application/json" \
    -d "{\"token\":\"$TOKEN\",\"id\":$BOOK_ID}")
  
  echo "$DETAIL_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$DETAIL_RESPONSE"
else
  echo "❌ 上传失败！"
  exit 1
fi

echo ""
echo "=== 测试完成 ==="

