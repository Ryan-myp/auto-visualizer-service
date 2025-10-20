package main

import (
	"fmt"
	"math/rand"
	"time"

	_ "github.com/Ryan-myp/auto-visualizer-service"
	autovisualizer "github.com/Ryan-myp/auto-visualizer-service"
)

// ============ 模拟电商订单处理系统 ============

// Order 订单
type Order struct {
	OrderID  string
	UserID   int64
	Amount   float64
	Status   string
	Products []string
}

// OrderService 订单服务
type OrderService struct {
	orders map[string]*Order
}

func NewOrderService() *OrderService {
	return &OrderService{
		orders: make(map[string]*Order),
	}
}

// CreateOrder 创建订单（主流程）
func (s *OrderService) CreateOrder(userID int64, products []string, amount float64) (*Order, error) {
	end := autovisualizer.Begin("OrderService.CreateOrder", userID, products, amount)

	var order *Order
	var err error

	defer func() {
		end(order, err)
	}()

	fmt.Printf("\n📦 创建订单...\n")
	time.Sleep(50 * time.Millisecond)

	// 步骤1: 验证用户
	if !s.ValidateUser(userID) {
		err = fmt.Errorf("用户验证失败: %d", userID)
		return nil, err
	}

	// 步骤2: 检查库存
	if !s.CheckInventory(products) {
		err = fmt.Errorf("库存不足")
		return nil, err
	}

	// 步骤3: 计算价格
	finalAmount := s.CalculatePrice(amount, userID)

	// 步骤4: 创建订单
	orderID := fmt.Sprintf("ORD-%d", time.Now().Unix())
	order = &Order{
		OrderID:  orderID,
		UserID:   userID,
		Amount:   finalAmount,
		Status:   "created",
		Products: products,
	}

	// 步骤5: 保存订单
	s.SaveOrder(order)

	fmt.Printf("✅ 订单创建成功: %s (金额: %.2f)\n", orderID, finalAmount)
	return order, nil
}

// ValidateUser 验证用户
func (s *OrderService) ValidateUser(userID int64) bool {
	defer autovisualizer.Measure("OrderService.ValidateUser")()

	fmt.Printf("  🔍 验证用户: %d\n", userID)
	time.Sleep(30 * time.Millisecond)

	// 模拟验证逻辑
	return userID > 0
}

// CheckInventory 检查库存
func (s *OrderService) CheckInventory(products []string) bool {
	defer autovisualizer.Measure("OrderService.CheckInventory")()

	fmt.Printf("  📦 检查库存: %v\n", products)
	time.Sleep(40 * time.Millisecond)

	// 模拟库存检查
	return len(products) > 0
}

// CalculatePrice 计算价格
func (s *OrderService) CalculatePrice(amount float64, userID int64) float64 {
	end := autovisualizer.Begin("OrderService.CalculatePrice", amount, userID)
	defer end(amount*0.9, nil)

	fmt.Printf("  💰 计算价格: %.2f\n", amount)
	time.Sleep(20 * time.Millisecond)

	// 应用折扣
	discount := s.GetUserDiscount(userID)
	finalAmount := amount * (1 - discount)

	return finalAmount
}

// GetUserDiscount 获取用户折扣
func (s *OrderService) GetUserDiscount(userID int64) float64 {
	defer autovisualizer.Measure("OrderService.GetUserDiscount")()

	fmt.Printf("    🎁 获取用户折扣: %d\n", userID)
	time.Sleep(15 * time.Millisecond)

	// VIP 用户 10% 折扣
	if userID%2 == 0 {
		return 0.1
	}
	return 0
}

// SaveOrder 保存订单
func (s *OrderService) SaveOrder(order *Order) {
	defer autovisualizer.Measure("OrderService.SaveOrder")()

	fmt.Printf("  💾 保存订单: %s\n", order.OrderID)
	time.Sleep(25 * time.Millisecond)

	// 保存到内存
	s.orders[order.OrderID] = order
}

// ProcessPayment 处理支付
func (s *OrderService) ProcessPayment(orderID string) error {
	end := autovisualizer.Begin("OrderService.ProcessPayment", orderID)

	var err error
	defer func() {
		end(nil, err)
	}()

	fmt.Printf("\n💳 处理支付...\n")
	time.Sleep(60 * time.Millisecond)

	order, exists := s.orders[orderID]
	if !exists {
		err = fmt.Errorf("订单不存在: %s", orderID)
		return err
	}

	// 调用支付网关
	if !s.CallPaymentGateway(order.Amount) {
		err = fmt.Errorf("支付失败")
		return err
	}

	// 更新订单状态
	s.UpdateOrderStatus(orderID, "paid")

	fmt.Printf("✅ 支付成功: %s\n", orderID)
	return nil
}

// CallPaymentGateway 调用支付网关
func (s *OrderService) CallPaymentGateway(amount float64) bool {
	defer autovisualizer.Measure("OrderService.CallPaymentGateway")()

	fmt.Printf("  🌐 调用支付网关: %.2f\n", amount)
	time.Sleep(80 * time.Millisecond)

	// 模拟支付成功
	return true
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(orderID string, status string) {
	defer autovisualizer.Measure("OrderService.UpdateOrderStatus")()

	fmt.Printf("  📝 更新订单状态: %s -> %s\n", orderID, status)
	time.Sleep(20 * time.Millisecond)

	if order, exists := s.orders[orderID]; exists {
		order.Status = status
	}
}

// ShipOrder 发货
func (s *OrderService) ShipOrder(orderID string) error {
	end := autovisualizer.Begin("OrderService.ShipOrder", orderID)

	var err error
	defer func() {
		end(nil, err)
	}()

	fmt.Printf("\n📮 发货处理...\n")
	time.Sleep(50 * time.Millisecond)

	order, exists := s.orders[orderID]
	if !exists {
		err = fmt.Errorf("订单不存在: %s", orderID)
		return err
	}

	if order.Status != "paid" {
		err = fmt.Errorf("订单未支付: %s", orderID)
		return err
	}

	// 生成物流单号
	trackingNumber := s.GenerateTrackingNumber()

	// 更新订单状态
	s.UpdateOrderStatus(orderID, "shipped")

	fmt.Printf("✅ 发货成功: %s (物流单号: %s)\n", orderID, trackingNumber)
	return nil
}

// GenerateTrackingNumber 生成物流单号
func (s *OrderService) GenerateTrackingNumber() string {
	defer autovisualizer.Measure("OrderService.GenerateTrackingNumber")()

	fmt.Printf("  🔢 生成物流单号\n")
	time.Sleep(10 * time.Millisecond)

	return fmt.Sprintf("TRK-%d", time.Now().Unix())
}

// ============ 主程序 ============

func main() {
	fmt.Println("🚀 Auto-Visualizer 完整演示")
	fmt.Println("=" + "==================================================")
	fmt.Println()

	// 检查是否启用
	if !autovisualizer.IsEnabled() {
		fmt.Println("❌ Auto-Visualizer 未启用")
		fmt.Println("请运行: export ENABLE_AUTO_VISUALIZER=true")
		return
	}

	fmt.Println("✅ Auto-Visualizer 已启用")
	fmt.Println()

	// 等待 Web 服务器启动
	fmt.Println("⏳ 等待 Web 服务器启动...")
	time.Sleep(1 * time.Second)
	fmt.Println()

	// 创建服务
	orderService := NewOrderService()

	// 模拟多个订单处理流程
	fmt.Println("🎬 开始演示订单处理流程")
	fmt.Println("=" + "==================================================")

	// 场景1: 正常订单流程
	fmt.Println("\n📋 场景1: 正常订单流程")
	fmt.Println("-" + "--------------------------------------------------")
	order1, err := orderService.CreateOrder(1001, []string{"商品A", "商品B"}, 999.99)
	if err != nil {
		fmt.Printf("❌ 创建订单失败: %v\n", err)
	} else {
		// 支付
		orderService.ProcessPayment(order1.OrderID)

		// 发货
		orderService.ShipOrder(order1.OrderID)
	}

	// 场景2: VIP用户订单（有折扣）
	fmt.Println("\n📋 场景2: VIP用户订单（有折扣）")
	fmt.Println("-" + "--------------------------------------------------")
	order2, err := orderService.CreateOrder(1002, []string{"商品C"}, 1299.50)
	if err != nil {
		fmt.Printf("❌ 创建订单失败: %v\n", err)
	} else {
		orderService.ProcessPayment(order2.OrderID)
		orderService.ShipOrder(order2.OrderID)
	}

	// 场景3: 批量创建订单
	fmt.Println("\n📋 场景3: 批量创建订单")
	fmt.Println("-" + "--------------------------------------------------")
	for i := 0; i < 3; i++ {
		userID := int64(2000 + i)
		products := []string{fmt.Sprintf("商品-%d", i)}
		amount := float64(rand.Intn(1000) + 500)

		order, err := orderService.CreateOrder(userID, products, amount)
		if err != nil {
			fmt.Printf("❌ 订单创建失败: %v\n", err)
		} else {
			fmt.Printf("✅ 订单创建成功: %s\n", order.OrderID)
		}

		time.Sleep(100 * time.Millisecond)
	}

	// 统计信息
	fmt.Println()
	fmt.Println("=" + "==================================================")
	fmt.Println("📊 演示完成统计")
	fmt.Println("=" + "==================================================")

	traces := autovisualizer.GetAllTraces()
	fmt.Printf("✅ 总追踪记录数: %d\n", len(traces))

	// 按方法统计
	methodStats := make(map[string]int)
	for _, trace := range traces {
		methodStats[trace.MethodName]++
	}

	fmt.Println("\n📈 方法调用统计:")
	for method, count := range methodStats {
		fmt.Printf("  • %s: %d 次\n", method, count)
	}

	// 获取追踪器状态
	tracer := autovisualizer.GetTracer()
	status := tracer.GetCircuitStatus()

	fmt.Println("\n🔧 追踪器状态:")
	fmt.Printf("  • 熔断器: %v\n", status["circuit_open"])
	fmt.Printf("  • 错误计数: %v\n", status["error_count"])
	fmt.Printf("  • 启用状态: %v\n", status["enabled"])

	// 访问链接
	fmt.Println()
	fmt.Println("=" + "==================================================")
	fmt.Println("🌐 访问可视化界面")
	fmt.Println("=" + "==================================================")
	fmt.Println()
	fmt.Println("📍 主页:")
	fmt.Println("   http://localhost:8090")
	fmt.Println()
	fmt.Println("📍 追踪列表 API:")
	fmt.Println("   http://localhost:8090/api/method-traces")
	fmt.Println()
	fmt.Println("📍 方法列表 API:")
	fmt.Println("   http://localhost:8090/api/interceptors")
	fmt.Println()
	fmt.Println("📍 调用树 API:")
	fmt.Println("   http://localhost:8090/api/method-traces/tree")
	fmt.Println()

	if len(traces) > 0 {
		// 显示第一个追踪的可视化链接
		firstTrace := traces[0]
		fmt.Println("📍 查看第一个追踪的可视化详情:")
		fmt.Printf("   http://localhost:8090/trace/%s\n", firstTrace.TraceID)
		fmt.Println()
	}

	fmt.Println("💡 提示:")
	fmt.Println("   1. 打开浏览器访问上面的链接")
	fmt.Println("   2. 查看美观的可视化界面")
	fmt.Println("   3. 点击追踪记录查看详细信息")
	fmt.Println("   4. 查看方法调用树和性能指标")
	fmt.Println()
	fmt.Println("按 Ctrl+C 退出...")

	// 保持运行
	select {}
}
