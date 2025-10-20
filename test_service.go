package main

import (
	"fmt"
	"log"
	"time"

	// 导入auto-visualizer包 - 自动启动插件
	_ "github.com/shopee/auto-visualizer"
	"github.com/shopee/auto-visualizer/config"
)

// 模拟AdMgmt业务方法
func AppPublishOpsAdsRun(campaignTaskId int64, para string, userId int64) error {
	log.Printf("📢 AdMgmt业务方法执行: AppPublishOpsAdsRun(taskId=%d, para=%s, userId=%d)",
		campaignTaskId, para, userId)

	// 模拟业务处理时间
	time.Sleep(500 * time.Millisecond)

	return nil
}

func CreateCampaign(name string, budget float64, platform string) error {
	log.Printf("📢 创建广告活动: %s, 预算: %.2f, 平台: %s", name, budget, platform)

	time.Sleep(300 * time.Millisecond)

	return nil
}

func main() {
	fmt.Println("🔌 Auto-Visualizer 独立服务测试")
	fmt.Println("==================================")

	fmt.Println("✅ Auto-Visualizer独立服务已启用")
	fmt.Printf("🌐 访问可视化界面: http://localhost:8090\n\n")

	// 注册业务方法
	config.RegisterInterceptor("AppPublishOpsAdsRun", "AdMgmt广告发布流程")
	config.RegisterInterceptor("CreateCampaign", "广告活动创建")

	fmt.Println("🎯 开始模拟业务调用...")

	// 模拟AdMgmt调用
	go func() {
		for i := 1; i <= 5; i++ {
			taskId := int64(12340 + i)
			para := fmt.Sprintf(`{"campaign_name":"春节促销%d","platform":"facebook","budget":%d}`, i, 1000*i)
			userId := int64(888880 + i)

			AppPublishOpsAdsRun(taskId, para, userId)
			time.Sleep(2 * time.Second)
		}
	}()

	// 模拟其他业务调用
	go func() {
		campaigns := []struct {
			name     string
			budget   float64
			platform string
		}{
			{"情人节特惠", 3000.0, "google"},
			{"三八节大促", 4000.0, "tiktok"},
			{"清明踏青", 2500.0, "facebook"},
		}

		time.Sleep(1 * time.Second) // 错开启动时间

		for _, campaign := range campaigns {
			CreateCampaign(campaign.name, campaign.budget, campaign.platform)
			time.Sleep(3 * time.Second)
		}
	}()

	fmt.Println("💡 独立服务特性:")
	fmt.Println("   ✅ 通过go.mod引入即可使用")
	fmt.Println("   ✅ 业务代码零侵入")
	fmt.Println("   ✅ 方法调用自动拦截")
	fmt.Println("   ✅ 执行逻辑自动记录")
	fmt.Println("   ✅ SQLite数据持久化")
	fmt.Println("")
	fmt.Println("📊 访问 http://localhost:8090 查看实时执行记录")
	fmt.Println("⏰ 服务运行中，按 Ctrl+C 退出...")

	// 保持服务运行
	select {}
}
