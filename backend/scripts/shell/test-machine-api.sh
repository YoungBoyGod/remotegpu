#!/bin/bash
# 测试设备添加 API
# 使用方法: ./test-machine-api.sh

set -e

HOST="${API_HOST:-http://localhost:8080}"
USERNAME="${ADMIN_USER:-admin}"
PASSWORD="${ADMIN_PASS:-admin123}"

echo "=== 1. 登录获取 Token ==="
TOKEN=$(curl -s -X POST "$HOST/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" \
  | jq -r '.data.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "登录失败"
  exit 1
fi
echo "Token: ${TOKEN:0:20}..."

echo ""
echo "=== 2. 获取机器列表 ==="
curl -s "$HOST/api/v1/admin/machines" \
  -H "Authorization: Bearer $TOKEN" | jq .
