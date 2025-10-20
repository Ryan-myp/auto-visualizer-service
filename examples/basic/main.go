package main

import (
	"fmt"
	"log"
	"time"

	// 导入auto-visualizer包 - 自动启动插件
	_ "github.com/Ryan-myp/auto-visualizer-service"
	autovisualizer "github.com/Ryan-myp/auto-visualizer-service"
)

// 模拟业务服务
type OrderService struct{}

// ProcessOrder 处理订单 - 会被自动拦截
func (s *OrderService) ProcessOrder(orderID string, userID int64, amount float64) error {
	log.Printf("📦 开始处理订单: %s, 用户: %d, 金额: %.2f", orderID, userID, amount)

	// 如果插件启用，拦截这个方法调用
	traceID := autovisualizer.InterceptCall("ProcessOrder", map[string]interface{}{
		"orderID": orderID,
		"userID":  userID,
		"amount":  amount,
	})

	// 模拟业务处理
	time.Sleep(500 * time.Millisecond)

	// 记录执行结果
	autovisualizer.RecordResult(traceID, map[string]interface{}{
		"success": true,
		"orderID": orderID,
		"status":  "completed",
	}, nil)

	log.Printf("✅ 订单处理完成: %s", orderID)
	return nil
}

// CreateCampaign 创建广告活动 - 会被自动拦截
func (s *OrderService) CreateCampaign(campaignName string, budget float64, platform string) error {
	log.Printf("📢 创建广告活动: %s, 预算: %.2f, 平台: %s", campaignName, budget, platform)

	// 拦截方法调用
	traceID := autovisualizer.InterceptCall("CreateCampaign", map[string]interface{}{
		"campaignName": campaignName,
		"budget":       budget,
		"platform":     platform,
	})

	// 模拟业务处理
	time.Sleep(800 * time.Millisecond)

	// 记录执行结果
	autovisualizer.RecordResult(traceID, map[string]interface{}{
		"success":    true,
		"campaignID": fmt.Sprintf("camp_%d", time.Now().Unix()),
		"status":     "active",
	}, nil)

	log.Printf("✅ 广告活动创建完成: %s", campaignName)
	return nil
}

func main() {
	fmt.Println("🚀 Auto-Visualizer 使用示例")
	fmt.Println("================================")

	// 检查插件是否启用
	if !autovisualizer.IsEnabled() {
		fmt.Println("💡 插件未启用，请设置环境变量:")
		fmt.Println("   export ENABLE_AUTO_VISUALIZER=true")
		fmt.Println("   go run main.go")
		return
	}

	fmt.Println("✅ Auto-Visualizer插件已启用")
	fmt.Printf("🌐 访问可视化界面: http://localhost:8090\n\n")

	// 注册自定义方法拦截器
	autovisualizer.RegisterMethod("ProcessOrder", "订单处理流程")
	autovisualizer.RegisterMethod("CreateCampaign", "广告创建流程")

	// 创建业务服务
	orderService := &OrderService{}

	// 模拟业务调用
	fmt.Println("🎯 开始模拟业务调用...")

	// 处理几个订单
	go func() {
		for i := 1; i <= 3; i++ {
			orderID := fmt.Sprintf("order_%d", i)
			userID := int64(1000 + i)
			amount := float64(100 + i*50)

			orderService.ProcessOrder(orderID, userID, amount)
			time.Sleep(2 * time.Second)
		}
	}()

	// 创建几个广告活动
	go func() {
		campaigns := []struct {
			name     string
			budget   float64
			platform string
		}{
			{"春节促销活动", 5000.0, "facebook"},
			{"情人节特惠", 3000.0, "google"},
			{"三八节大促", 4000.0, "tiktok"},
		}

		for _, campaign := range campaigns {
			orderService.CreateCampaign(campaign.name, campaign.budget, campaign.platform)
			time.Sleep(3 * time.Second)
		}
	}()

	fmt.Println("💡 业务调用已启动，所有方法调用都会被自动拦截和记录")
	fmt.Println("📊 访问 http://localhost:8090 查看实时执行记录")
	fmt.Println("⏰ 服务运行中，按 Ctrl+C 退出...")

	// 保持服务运行
	select {}
}
