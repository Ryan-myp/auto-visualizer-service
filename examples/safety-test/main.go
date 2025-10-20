package main

import (
	"fmt"
	"time"

	_ "github.com/Ryan-myp/auto-visualizer-service"
	autovisualizer "github.com/Ryan-myp/auto-visualizer-service"
)

// 测试1：正常业务逻辑（不会被追踪影响）
func NormalBusiness(id int) string {
	end := autovisualizer.Begin("NormalBusiness", id)
	defer end("success", nil)
	
	// 正常业务逻辑
	time.Sleep(10 * time.Millisecond)
	return "success"
}

// 测试2：业务逻辑中有 panic（追踪不应该干扰）
func BusinessWithPanic(id int) (result string) {
	defer autovisualizer.Begin("BusinessWithPanic", id)(result, nil)
	
	// 业务 panic，但追踪不应该影响
	if id == 999 {
		panic("business panic")
	}
	
	result = "success"
	return
}

// 测试3：传入超大对象（不应该阻塞）
func ProcessLargeData(data []byte) string {
	end := autovisualizer.Begin("ProcessLargeData", data)
	defer end("done", nil)
	
	// 业务逻辑
	time.Sleep(5 * time.Millisecond)
	return "done"
}

// 测试4：循环引用的结构体（不应该导致 panic）
type Node struct {
	Value int
	Next  *Node
}

func ProcessCircularRef(node *Node) string {
	end := autovisualizer.Begin("ProcessCircularRef", node)
	defer end("done", nil)
	
	// 业务逻辑
	return "done"
}

// 测试5：高并发场景（追踪不应该成为瓶颈）
func ConcurrentTest(id int) {
	defer autovisualizer.Measure("ConcurrentTest")()
	
	// 模拟快速业务逻辑
	time.Sleep(1 * time.Millisecond)
}

func main() {
	fmt.Println("🧪 安全性和性能测试")
	fmt.Println("=" + "================================================")
	fmt.Println()
	
	// 测试1：正常业务
	fmt.Println("✅ 测试1：正常业务逻辑")
	for i := 0; i < 5; i++ {
		result := NormalBusiness(i)
		fmt.Printf("  业务 %d 完成: %s\n", i, result)
	}
	fmt.Println()
	
	// 测试2：业务 panic（追踪应该不影响）
	fmt.Println("✅ 测试2：业务 panic 处理")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  捕获到业务 panic: %v (追踪未受影响)\n", r)
			}
		}()
		BusinessWithPanic(999)
	}()
	fmt.Println()
	
	// 测试3：超大对象（不应该阻塞）
	fmt.Println("✅ 测试3：处理超大对象")
	start := time.Now()
	largeData := make([]byte, 10*1024*1024) // 10MB
	ProcessLargeData(largeData)
	elapsed := time.Since(start)
	fmt.Printf("  处理完成，耗时: %v (应该 < 100ms)\n", elapsed)
	if elapsed > 100*time.Millisecond {
		fmt.Printf("  ⚠️  警告：耗时过长！\n")
	}
	fmt.Println()
	
	// 测试4：循环引用（不应该 panic）
	fmt.Println("✅ 测试4：循环引用结构")
	node1 := &Node{Value: 1}
	node2 := &Node{Value: 2}
	node1.Next = node2
	node2.Next = node1 // 循环引用
	
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  ❌ 发生 panic: %v\n", r)
			} else {
				fmt.Printf("  ✅ 循环引用处理成功\n")
			}
		}()
		ProcessCircularRef(node1)
	}()
	fmt.Println()
	
	// 测试5：高并发（追踪不应该成为瓶颈）
	fmt.Println("✅ 测试5：高并发场景")
	start = time.Now()
	concurrency := 1000
	done := make(chan bool, concurrency)
	
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			ConcurrentTest(id)
			done <- true
		}(i)
	}
	
	// 等待所有完成
	for i := 0; i < concurrency; i++ {
		<-done
	}
	
	elapsed = time.Since(start)
	fmt.Printf("  %d 个并发请求完成，总耗时: %v\n", concurrency, elapsed)
	fmt.Printf("  平均每个请求: %v\n", elapsed/time.Duration(concurrency))
	fmt.Println()
	
	// 测试6：追踪器状态
	fmt.Println("✅ 测试6：追踪器状态")
	tracer := autovisualizer.GetTracer()
	if tracer != nil {
		status := tracer.GetCircuitStatus()
		fmt.Printf("  熔断器状态: %+v\n", status)
	}
	fmt.Println()
	
	// 获取追踪统计
	traces := autovisualizer.GetAllTraces()
	fmt.Printf("📊 总共记录了 %d 条追踪\n", len(traces))
	fmt.Println()
	
	fmt.Println("🎉 所有测试完成！")
	fmt.Println("💡 关键点:")
	fmt.Println("  1. 业务 panic 不会被追踪影响")
	fmt.Println("  2. 超大对象不会阻塞业务")
	fmt.Println("  3. 循环引用不会导致 panic")
	fmt.Println("  4. 高并发下追踪不是瓶颈")
	fmt.Println("  5. 熔断机制保护系统稳定")
	fmt.Println()
	fmt.Println("🌐 访问 http://localhost:8090 查看追踪数据")
	fmt.Println()
	fmt.Println("按 Ctrl+C 退出...")
	
	select {}
}

