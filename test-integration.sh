#!/bin/bash

# Auto-Visualizer 集成测试脚本

echo "🧪 Auto-Visualizer 集成测试"
echo "=========================================="
echo ""

# 设置环境变量
export ENABLE_AUTO_VISUALIZER=true
export AUTO_VISUALIZER_PORT=8090

echo "✅ 环境变量已设置:"
echo "   ENABLE_AUTO_VISUALIZER=$ENABLE_AUTO_VISUALIZER"
echo "   AUTO_VISUALIZER_PORT=$AUTO_VISUALIZER_PORT"
echo ""

# 进入测试目录
cd examples/test-integration

echo "📦 安装依赖..."
go mod tidy 2>&1 | grep -v "go: finding" || true
echo ""

echo "🚀 启动测试程序..."
echo "   (程序将在后台运行 10 秒)"
echo ""

# 在后台运行测试程序
timeout 10s go run main.go &
PID=$!

# 等待服务器启动
echo "⏳ 等待服务器启动..."
sleep 3

# 检查进程是否还在运行
if ! kill -0 $PID 2>/dev/null; then
    echo "❌ 测试程序启动失败"
    exit 1
fi

echo ""
echo "🔍 测试端点..."
echo ""

# 测试健康检查
echo "1️⃣  测试健康检查端点:"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8090/health)
if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ 健康检查通过 (HTTP $HTTP_CODE)"
else
    echo "   ❌ 健康检查失败 (HTTP $HTTP_CODE)"
fi
echo ""

# 测试追踪 API
echo "2️⃣  测试追踪 API:"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8090/api/method-traces)
if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ 追踪 API 可访问 (HTTP $HTTP_CODE)"
    
    # 获取追踪数据
    TRACE_COUNT=$(curl -s http://localhost:8090/api/method-traces | grep -o '"total":[0-9]*' | grep -o '[0-9]*')
    echo "   📊 已记录 $TRACE_COUNT 条追踪"
else
    echo "   ❌ 追踪 API 失败 (HTTP $HTTP_CODE)"
fi
echo ""

# 测试主页
echo "3️⃣  测试 Web UI:"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8090/)
if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ Web UI 可访问 (HTTP $HTTP_CODE)"
else
    echo "   ❌ Web UI 失败 (HTTP $HTTP_CODE)"
fi
echo ""

# 清理
echo "🧹 清理..."
kill $PID 2>/dev/null || true
wait $PID 2>/dev/null || true

echo ""
echo "=========================================="
echo "✅ 集成测试完成！"
echo ""
echo "💡 如果所有测试都通过，说明集成成功！"
echo ""
echo "🌐 在您的项目中使用:"
echo "   1. go get github.com/Ryan-myp/auto-visualizer-service@latest"
echo "   2. export ENABLE_AUTO_VISUALIZER=true"
echo "   3. import _ \"github.com/Ryan-myp/auto-visualizer-service\""
echo ""

