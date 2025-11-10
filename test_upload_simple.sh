#!/bin/bash

# 简单的图片上传测试脚本
# 使用方法: ./test_upload_simple.sh <邮箱> <密码> <图书ID>

EMAIL="${1:-your-email@example.com}"
PASSWORD="${2:-your-password}"
BOOK_ID="${3:-1}"
IMAGE_FILE="/Users/cat/Desktop/book-manage/docs/troubleshooting/test.png"

echo "=== 图片上传测试 ==="
echo "邮箱: $EMAIL"
echo "图书ID: $BOOK_ID"
echo "图片: $IMAGE_FILE"
echo ""

# 1. 登录
echo "1. 正在登录..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESPONSE | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('data', {}).get('token', ''))" 2>/dev/null)

if [ -z "$TOKEN" ]; then
  echo "❌ 登录失败！"
  echo "响应: $LOGIN_RESPONSE"
  echo ""
  echo "请检查："
  echo "1. 邮箱和密码是否正确"
  echo "2. 账户是否存在"
  echo "3. 服务是否正在运行"
  exit 1
fi

echo "✅ 登录成功！"
echo ""

# 2. 检查图片文件
if [ ! -f "$IMAGE_FILE" ]; then
  echo "❌ 图片文件不存在: $IMAGE_FILE"
  exit 1
fi

FILE_SIZE=$(stat -f%z "$IMAGE_FILE" 2>/dev/null || stat -c%s "$IMAGE_FILE" 2>/dev/null)
FILE_SIZE_MB=$(echo "scale=2; $FILE_SIZE / 1024 / 1024" | bc 2>/dev/null || echo "0")

echo "2. 图片信息："
echo "   文件: $IMAGE_FILE"
echo "   大小: ${FILE_SIZE_MB} MB"
echo ""

# 3. 上传图片
echo "3. 正在上传图片..."
UPLOAD_RESPONSE=$(curl -s -X POST http://localhost:8080/api/book/uploadCover \
  -F "token=$TOKEN" \
  -F "book_id=$BOOK_ID" \
  -F "image=@$IMAGE_FILE")

echo "上传响应:"
echo "$UPLOAD_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$UPLOAD_RESPONSE"
echo ""

# 4. 检查结果
if echo "$UPLOAD_RESPONSE" | grep -q '"code":0'; then
  IMAGE_URL=$(echo $UPLOAD_RESPONSE | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('data', {}).get('image_url', ''))" 2>/dev/null)
  echo "✅ 上传成功！"
  echo ""
  echo "图片 URL: $IMAGE_URL"
  echo ""
  echo "4. 验证图片URL（在浏览器中打开）："
  echo "   $IMAGE_URL"
  echo ""
  
  # 5. 查看图书详情
  echo "5. 查看图书详情..."
  DETAIL_RESPONSE=$(curl -s -X POST http://localhost:8080/api/book/detail \
    -H "Content-Type: application/json" \
    -d "{\"token\":\"$TOKEN\",\"id\":$BOOK_ID}")
  
  echo "$DETAIL_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$DETAIL_RESPONSE"
else
  echo "❌ 上传失败！"
  echo ""
  echo "可能的原因："
  echo "1. 不是管理员账户（需要管理员权限）"
  echo "2. 图书ID不存在"
  echo "3. 图片格式或大小不符合要求"
  echo "4. R2服务未配置或配置错误"
fi

echo ""
echo "=== 测试完成 ==="

