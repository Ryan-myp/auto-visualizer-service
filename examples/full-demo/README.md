# 完整演示程序

## 🎯 演示内容

这个演示程序模拟了一个完整的电商订单处理系统，展示 Auto-Visualizer 的所有功能：

### 业务流程

1. **创建订单**
   - 验证用户
   - 检查库存
   - 计算价格（含折扣）
   - 保存订单

2. **处理支付**
   - 调用支付网关
   - 更新订单状态

3. **订单发货**
   - 生成物流单号
   - 更新订单状态

### 演示场景

- **场景1**: 正常订单流程（完整流程）
- **场景2**: VIP用户订单（有折扣）
- **场景3**: 批量创建订单

## 🚀 运行演示

### 方式1: 使用脚本（推荐）

```bash
# 在项目根目录运行
./run-demo.sh
```

### 方式2: 手动运行

```bash
# 设置环境变量
export ENABLE_AUTO_VISUALIZER=true

# 进入演示目录
cd examples/full-demo

# 运行
go run main.go
```

## 📊 查看结果

演示运行后，访问以下地址：

### 1. 主页
```
http://localhost:8090
```

### 2. 追踪列表 API
```bash
curl http://localhost:8090/api/method-traces
```

### 3. 方法列表（动态生成）
```bash
curl http://localhost:8090/api/interceptors
```

返回示例：
```json
{
  "success": true,
  "methods": [
    {
      "name": "OrderService.CreateOrder",
      "package": "main.OrderService.CreateOrder",
      "call_count": 5,
      "success_count": 5,
      "error_count": 0,
      "avg_duration_ms": 150
    },
    {
      "name": "OrderService.ValidateUser",
      "call_count": 5,
      "success_count": 5,
      "error_count": 0,
      "avg_duration_ms": 30
    }
  ],
  "total": 10
}
```

### 4. 调用树
```bash
curl http://localhost:8090/api/method-traces/tree
```

### 5. 可视化详情页

```
http://localhost:8090/trace/{trace_id}
```

## 🎨 可视化展示

### 调用链示例

```
OrderService.CreateOrder (150ms) ✅
├── OrderService.ValidateUser (30ms) ✅
├── OrderService.CheckInventory (40ms) ✅
├── OrderService.CalculatePrice (35ms) ✅
│   └── OrderService.GetUserDiscount (15ms) ✅
└── OrderService.SaveOrder (25ms) ✅

OrderService.ProcessPayment (160ms) ✅
├── OrderService.CallPaymentGateway (80ms) ✅
└── OrderService.UpdateOrderStatus (20ms) ✅

OrderService.ShipOrder (80ms) ✅
├── OrderService.GenerateTrackingNumber (10ms) ✅
└── OrderService.UpdateOrderStatus (20ms) ✅
```

### 方法统计示例

```
📈 方法调用统计:
  • OrderService.CreateOrder: 5 次
  • OrderService.ValidateUser: 5 次
  • OrderService.CheckInventory: 5 次
  • OrderService.CalculatePrice: 5 次
  • OrderService.GetUserDiscount: 5 次
  • OrderService.SaveOrder: 5 次
  • OrderService.ProcessPayment: 2 次
  • OrderService.CallPaymentGateway: 2 次
  • OrderService.UpdateOrderStatus: 4 次
  • OrderService.ShipOrder: 2 次
  • OrderService.GenerateTrackingNumber: 2 次
```

## 💡 学习要点

### 1. 不同的追踪方式

```go
// 方式1: Begin + defer（记录入参出参）
func CreateOrder(userID int64, products []string, amount float64) (*Order, error) {
    end := autovisualizer.Begin("CreateOrder", userID, products, amount)
    var order *Order
    var err error
    defer func() { end(order, err) }()
    // 业务逻辑
    return order, err
}

// 方式2: Measure（只测量时间）
func ValidateUser(userID int64) bool {
    defer autovisualizer.Measure("ValidateUser")()
    // 业务逻辑
    return true
}
```

### 2. 嵌套调用追踪

```go
func CreateOrder() {
    // 父方法
    end := autovisualizer.Begin("CreateOrder")
    defer end(nil)
    
    // 子方法1
    ValidateUser()
    
    // 子方法2
    CheckInventory()
    
    // 自动形成调用树
}
```

### 3. 查看追踪数据

```go
// 获取所有追踪
traces := autovisualizer.GetAllTraces()

// 获取追踪器状态
tracer := autovisualizer.GetTracer()
status := tracer.GetCircuitStatus()
```

## 🔍 调试技巧

### 1. 查找性能瓶颈

1. 运行演示程序
2. 访问 `/api/method-traces/tree`
3. 查看每个方法的耗时
4. 找出最慢的方法

### 2. 查看调用关系

1. 访问可视化详情页
2. 查看"子调用链"部分
3. 了解方法之间的调用关系

### 3. 分析错误

1. 查看状态标签（红色表示失败）
2. 查看错误信息
3. 查看调用栈定位问题

## 📝 扩展练习

### 练习1: 添加新的业务方法

```go
func (s *OrderService) CancelOrder(orderID string) error {
    end := autovisualizer.Begin("OrderService.CancelOrder", orderID)
    defer end(nil, nil)
    
    // 实现取消订单逻辑
    
    return nil
}
```

### 练习2: 模拟错误场景

```go
func (s *OrderService) CreateOrder(...) (*Order, error) {
    end := autovisualizer.Begin("CreateOrder", ...)
    var err error
    defer func() { end(nil, err) }()
    
    // 模拟错误
    if amount > 10000 {
        err = fmt.Errorf("金额超限")
        return nil, err
    }
    
    return order, nil
}
```

### 练习3: 添加更多嵌套层级

```go
func Level1() {
    defer autovisualizer.Measure("Level1")()
    Level2()
}

func Level2() {
    defer autovisualizer.Measure("Level2")()
    Level3()
}

func Level3() {
    defer autovisualizer.Measure("Level3")()
    // 业务逻辑
}
```

## 🆘 常见问题

### Q1: 看不到追踪数据？

**A**: 确保设置了环境变量：
```bash
export ENABLE_AUTO_VISUALIZER=true
```

### Q2: 端口被占用？

**A**: 修改端口：
```bash
export AUTO_VISUALIZER_PORT=9090
```

### Q3: 追踪数据太多？

**A**: 降低采样率：
```bash
export AUTO_VISUALIZER_SAMPLE_RATE=0.1  # 只追踪10%
```

## 🔗 相关文档

- [快速开始](../../QUICKSTART.md)
- [使用指南](../../USAGE.md)
- [可视化说明](../../VISUALIZATION.md)
- [安全性说明](../../SAFETY.md)

---

**享受可视化带来的便利！** 🎉

