package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	// 导入 auto-visualizer
	_ "github.com/Ryan-myp/auto-visualizer-service"
	autovisualizer "github.com/Ryan-myp/auto-visualizer-service"
)

func TestMethod1(name string) string {
	end := autovisualizer.Begin("TestMethod1", name)
	defer end("success", nil)
	
	time.Sleep(50 * time.Millisecond)
	return "success"
}

func TestMethod2(id int) {
	defer autovisualizer.Measure("TestMethod2")()
	time.Sleep(30 * time.Millisecond)
}

func main() {
	fmt.Println("🚀 测试 Auto-Visualizer 集成")
	fmt.Println("=" + "================================================")
	
	// 检查环境变量
	if os.Getenv("ENABLE_AUTO_VISUALIZER") != "true" {
		fmt.Println("⚠️  警告: ENABLE_AUTO_VISUALIZER 未设置为 true")
		fmt.Println("请运行: export ENABLE_AUTO_VISUALIZER=true")
		fmt.Println()
	}
	
	// 检查可视化器是否启用
	if autovisualizer.IsEnabled() {
		fmt.Println("✅ Auto-Visualizer 已启用")
	} else {
		fmt.Println("❌ Auto-Visualizer 未启用")
		fmt.Println("请设置环境变量: export ENABLE_AUTO_VISUALIZER=true")
		os.Exit(1)
	}
	
	// 等待服务器启动
	fmt.Println("⏳ 等待 Web 服务器启动...")
	time.Sleep(2 * time.Second)
	
	// 检查端口是否监听
	port := os.Getenv("AUTO_VISUALIZER_PORT")
	if port == "" {
		port = "8090"
	}
	
	url := fmt.Sprintf("http://localhost:%s/health", port)
	fmt.Printf("🔍 检查健康检查端点: %s\n", url)
	
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ 无法连接到 Web 服务器: %v\n", err)
		fmt.Printf("端口 %s 可能没有被监听\n", port)
		fmt.Println()
		fmt.Println("可能的原因:")
		fmt.Println("1. Web 服务器启动失败")
		fmt.Println("2. 端口被占用")
		fmt.Println("3. init 函数没有正确执行")
		os.Exit(1)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Printf("✅ Web 服务器正常运行 (端口: %s)\n", port)
	} else {
		fmt.Printf("⚠️  Web 服务器响应异常: %d\n", resp.StatusCode)
	}
	
	fmt.Println()
	fmt.Println("📝 执行测试方法...")
	
	// 执行一些测试方法
	for i := 0; i < 3; i++ {
		TestMethod1(fmt.Sprintf("test-%d", i))
		TestMethod2(i)
	}
	
	fmt.Println("✅ 测试方法执行完成")
	fmt.Println()
	
	// 检查追踪数据
	traces := autovisualizer.GetAllTraces()
	fmt.Printf("📊 已记录 %d 条追踪\n", len(traces))
	
	if len(traces) > 0 {
		fmt.Println("✅ 追踪功能正常工作")
	} else {
		fmt.Println("⚠️  没有记录到追踪数据")
	}
	
	fmt.Println()
	fmt.Printf("🌐 访问 Web UI: http://localhost:%s\n", port)
	fmt.Printf("🌐 查看追踪数据: http://localhost:%s/api/method-traces\n", port)
	fmt.Println()
	fmt.Println("✅ 集成测试完成！")
	fmt.Println()
	fmt.Println("按 Ctrl+C 退出...")
	
	// 保持运行
	select {}
}

