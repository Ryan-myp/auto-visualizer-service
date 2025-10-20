package main

import (
	"fmt"
	"time"

	// 只需要这一行导入，就能自动启动追踪！
	_ "github.com/Ryan-myp/auto-visualizer-service"
	autovisualizer "github.com/Ryan-myp/auto-visualizer-service"
)

// 业务函数1：处理订单
func ProcessOrder(orderID string, amount float64) (string, error) {
	// 添加追踪 - 方式1：使用 Begin
	end := autovisualizer.Begin("ProcessOrder", orderID, amount)
	
	var result string
	var err error
	
	defer func() {
		end(result, err)
	}()
	
	// 业务逻辑
	fmt.Printf("处理订单: %s, 金额: %.2f\n", orderID, amount)
	time.Sleep(100 * time.Millisecond)
	
	// 调用其他方法
	ValidateOrder(orderID)
	SaveOrder(orderID, amount)
	
	result = fmt.Sprintf("订单 %s 处理成功", orderID)
	return result, nil
}

// 业务函数2：验证订单
func ValidateOrder(orderID string) bool {
	// 添加追踪 - 方式2：使用 Measure（只测量时间）
	defer autovisualizer.Measure("ValidateOrder")()
	
	fmt.Printf("  验证订单: %s\n", orderID)
	time.Sleep(50 * time.Millisecond)
	return true
}

// 业务函数3：保存订单
func SaveOrder(orderID string, amount float64) {
	// 添加追踪 - 方式3：使用 TraceMethod
	defer autovisualizer.TraceMethod("SaveOrder")()
	
	fmt.Printf("  保存订单: %s\n", orderID)
	time.Sleep(30 * time.Millisecond)
}

// 计算折扣
func CalculateDiscount(amount float64, rate float64) float64 {
	defer autovisualizer.Measure("CalculateDiscount")()
	
	time.Sleep(10 * time.Millisecond)
	return amount * rate
}

func main() {
	fmt.Println("🚀 Auto-Visualizer 快速开始示例")
	fmt.Println("=" + "================================================")
	fmt.Println()
	
	// 直接调用业务函数，追踪会自动记录
	fmt.Println("📦 处理订单...")
	result, _ := ProcessOrder("ORD-001", 999.99)
	fmt.Printf("✅ %s\n\n", result)
	
	// 再处理一个订单
	fmt.Println("📦 处理订单...")
	result, _ = ProcessOrder("ORD-002", 1299.50)
	fmt.Printf("✅ %s\n\n", result)
	
	// 计算折扣
	fmt.Println("💰 计算折扣...")
	discount := CalculateDiscount(1000.0, 0.1)
	fmt.Printf("✅ 折扣金额: %.2f\n\n", discount)
	
	// 打印追踪信息
	traces := autovisualizer.GetAllTraces()
	fmt.Printf("📊 已记录 %d 条追踪\n\n", len(traces))
	
	fmt.Println("🌐 打开浏览器访问以下地址查看可视化:")
	fmt.Println("   主页:     http://localhost:8090")
	fmt.Println("   追踪列表: http://localhost:8090/api/method-traces")
	fmt.Println("   调用树:   http://localhost:8090/api/method-traces/tree")
	fmt.Println()
	fmt.Println("💡 提示: 设置环境变量 ENABLE_AUTO_VISUALIZER=true 启用追踪")
	fmt.Println()
	fmt.Println("按 Ctrl+C 退出...")
	
	// 保持程序运行
	select {}
}

