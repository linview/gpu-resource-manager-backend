#!/bin/bash

# GPU资源管理系统 API 测试脚本

BASE_URL="http://localhost:8080"

echo "=== GPU资源管理系统 API 测试 ==="
echo "基础URL: $BASE_URL"
echo

# 1. 健康检查
echo "1. 健康检查"
curl -s "$BASE_URL/health" | jq .
echo

# 2. 获取GPU列表
echo "2. 获取GPU列表"
curl -s "$BASE_URL/api/v1/gpus" | jq .
echo

# 3. 获取GPU详情
echo "3. 获取GPU详情"
curl -s "$BASE_URL/api/v1/gpus/gpu-001" | jq .
echo

# 4. 获取GPU状态
echo "4. 获取GPU状态"
curl -s "$BASE_URL/api/v1/gpus/gpu-001/status" | jq .
echo

# 5. 获取GPU指标
echo "5. 获取GPU指标"
curl -s "$BASE_URL/api/v1/gpus/gpu-001/metrics" | jq .
echo

# 6. 更新GPU配置
echo "6. 更新GPU配置"
curl -s -X POST "$BASE_URL/api/v1/gpus/gpu-001/config" \
  -H "Content-Type: application/json" \
  -d '{
    "power_limit": 200,
    "memory_clock": 8000,
    "driver_config": {
      "driver_version": "470.82.01",
      "cuda_version": "11.4"
    }
  }' | jq .
echo

# 7. 获取GPU配置
echo "7. 获取GPU配置"
curl -s "$BASE_URL/api/v1/gpus/gpu-001/config" | jq .
echo

# 8. 获取服务器列表
echo "8. 获取服务器列表"
curl -s "$BASE_URL/api/v1/servers" | jq .
echo

# 9. 创建资源分配
echo "9. 创建资源分配"
curl -s -X POST "$BASE_URL/api/v1/allocations" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "server_id": "server456",
    "purpose": "AI训练"
  }' | jq .
echo

# 10. 获取分配列表
echo "10. 获取分配列表"
curl -s "$BASE_URL/api/v1/allocations" | jq .
echo

# 11. 创建部署工作流
echo "11. 创建部署工作流"
curl -s -X POST "$BASE_URL/api/v1/workflows/deploy" \
  -H "Content-Type: application/json" \
  -d '{
    "allocation_id": "alloc789",
    "server_id": "server456"
  }' | jq .
echo

# 12. 获取工作流列表
echo "12. 获取工作流列表"
curl -s "$BASE_URL/api/v1/workflows" | jq .
echo

# 13. 获取事件列表
echo "13. 获取事件列表"
curl -s "$BASE_URL/api/v1/events" | jq .
echo

# 14. 获取告警列表
echo "14. 获取告警列表"
curl -s "$BASE_URL/api/v1/alerts" | jq .
echo

echo "=== 测试完成 ==="
