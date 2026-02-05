#!/bin/bash
# API 完整测试脚本
# 使用方法: ./test-api.sh

set -e

HOST="${API_HOST:-http://localhost:8080}"
USERNAME="${ADMIN_USER:-admin}"
PASSWORD="${ADMIN_PASS:-admin@123}"

echo "=========================================="
echo "RemoteGPU API 测试"
echo "=========================================="
echo "Host: $HOST"
echo ""

# 1. 健康检查
echo "=== 1. 健康检查 ==="
curl -s "$HOST/api/v1/health"
echo ""
echo ""

# 2. 登录
echo "=== 2. 登录 ==="
LOGIN_RESP=$(curl -s -X POST "$HOST/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")

echo "$LOGIN_RESP" | jq .
TOKEN=$(echo "$LOGIN_RESP" | jq -r '.data.access_token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "登录失败!"
  exit 1
fi
echo "Token: ${TOKEN:0:30}..."
echo ""

# 3. 获取机器列表
echo "=== 3. 获取机器列表 ==="
curl -s "$HOST/api/v1/admin/machines" \
  -H "Authorization: Bearer $TOKEN" | jq '.data.total'
echo ""

echo "=========================================="
echo "测试完成!"
echo "=========================================="
