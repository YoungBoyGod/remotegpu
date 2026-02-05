#!/bin/bash
# 添加测试设备
# 使用方法: ./add-test-machine.sh <ip> <name>

HOST="${API_HOST:-http://localhost:8080}"
TOKEN="$API_TOKEN"

IP="${1:-192.168.1.100}"
NAME="${2:-test-gpu-01}"

if [ -z "$TOKEN" ]; then
  echo "请先设置 API_TOKEN 环境变量"
  echo "export API_TOKEN=your_token"
  exit 1
fi

echo "添加设备: $NAME ($IP)"
curl -s -X POST "$HOST/api/v1/admin/machines" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"$NAME\",
    \"hostname\": \"$NAME\",
    \"region\": \"default\",
    \"ip_address\": \"$IP\",
    \"ssh_port\": 22,
    \"ssh_username\": \"root\"
  }" | jq .
