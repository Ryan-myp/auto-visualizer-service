#!/bin/bash

echo "🚀 Auto-Visualizer 完整演示"
echo "=========================================="
echo ""

# 设置环境变量
export ENABLE_AUTO_VISUALIZER=true
export AUTO_VISUALIZER_PORT=8090
export AUTO_VISUALIZER_SERVICE_NAME=demo-service

echo "✅ 环境变量已设置:"
echo "   ENABLE_AUTO_VISUALIZER=$ENABLE_AUTO_VISUALIZER"
echo "   AUTO_VISUALIZER_PORT=$AUTO_VISUALIZER_PORT"
echo "   AUTO_VISUALIZER_SERVICE_NAME=$AUTO_VISUALIZER_SERVICE_NAME"
echo ""

# 进入演示目录
cd examples/full-demo

echo "📦 安装依赖..."
go mod tidy 2>&1 | grep -v "go: finding" || true
echo ""

echo "🎬 启动演示程序..."
echo ""
echo "=========================================="
echo ""

# 运行演示
go run main.go

echo ""
echo "演示结束"

