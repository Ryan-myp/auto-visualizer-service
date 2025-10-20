# 演示程序输出示例

## 🎬 运行演示

```bash
./run-demo.sh
```

## 📺 预期输出

```
🚀 Auto-Visualizer 完整演示
==========================================

✅ 环境变量已设置:
   ENABLE_AUTO_VISUALIZER=true
   AUTO_VISUALIZER_PORT=8090
   AUTO_VISUALIZER_SERVICE_NAME=demo-service

📦 安装依赖...

🎬 启动演示程序...

==========================================

🚀 Auto-Visualizer 完整演示
==================================================

✅ Auto-Visualizer 已启用

⏳ 等待 Web 服务器启动...

🚀 启动Auto-Visualizer - 服务: demo-service
🎯 启动拦截器管理器...
✅ 拦截器管理器启动成功，已注册 0 个拦截器
🌐 Web服务器启动中: http://localhost:8090
✅ Web服务器已启动: http://localhost:8090
🎉 Auto-Visualizer插件已自动启动!
🌐 访问地址: http://localhost:8090
💾 数据库: ./demo-service_visualizer.db
🔍 方法追踪器已启用 (采样率: 1.00)

🎬 开始演示订单处理流程
==================================================

📋 场景1: 正常订单流程
--------------------------------------------------

📦 创建订单...
  🔍 验证用户: 1001
  📦 检查库存: [商品A 商品B]
  💰 计算价格: 999.99
    🎁 获取用户折扣: 1001
  💾 保存订单: ORD-1729425045
✅ 订单创建成功: ORD-1729425045 (金额: 999.99)

💳 处理支付...
  🌐 调用支付网关: 999.99
  📝 更新订单状态: ORD-1729425045 -> paid
✅ 支付成功: ORD-1729425045

📮 发货处理...
  🔢 生成物流单号
  📝 更新订单状态: ORD-1729425045 -> shipped
✅ 发货成功: ORD-1729425045 (物流单号: TRK-1729425045)

📋 场景2: VIP用户订单（有折扣）
--------------------------------------------------

📦 创建订单...
  🔍 验证用户: 1002
  📦 检查库存: [商品C]
  💰 计算价格: 1299.50
    🎁 获取用户折扣: 1002
  💾 保存订单: ORD-1729425046
✅ 订单创建成功: ORD-1729425046 (金额: 1169.55)

💳 处理支付...
  🌐 调用支付网关: 1169.55
  📝 更新订单状态: ORD-1729425046 -> paid
✅ 支付成功: ORD-1729425046

📮 发货处理...
  🔢 生成物流单号
  📝 更新订单状态: ORD-1729425046 -> shipped
✅ 发货成功: ORD-1729425046 (物流单号: TRK-1729425046)

📋 场景3: 批量创建订单
--------------------------------------------------

📦 创建订单...
  🔍 验证用户: 2000
  📦 检查库存: [商品-0]
  💰 计算价格: 987.00
    🎁 获取用户折扣: 2000
  💾 保存订单: ORD-1729425047
✅ 订单创建成功: ORD-1729425047 (金额: 888.30)
✅ 订单创建成功: ORD-1729425047

📦 创建订单...
  🔍 验证用户: 2001
  📦 检查库存: [商品-1]
  💰 计算价格: 756.00
    🎁 获取用户折扣: 2001
  💾 保存订单: ORD-1729425048
✅ 订单创建成功: ORD-1729425048 (金额: 756.00)
✅ 订单创建成功: ORD-1729425048

📦 创建订单...
  🔍 验证用户: 2002
  📦 检查库存: [商品-2]
  💰 计算价格: 1234.00
    🎁 获取用户折扣: 2002
  💾 保存订单: ORD-1729425049
✅ 订单创建成功: ORD-1729425049 (金额: 1110.60)
✅ 订单创建成功: ORD-1729425049

==================================================
📊 演示完成统计
==================================================
✅ 总追踪记录数: 35

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

🔧 追踪器状态:
  • 熔断器: false
  • 错误计数: 0
  • 启用状态: true

==================================================
🌐 访问可视化界面
==================================================

📍 主页:
   http://localhost:8090

📍 追踪列表 API:
   http://localhost:8090/api/method-traces

📍 方法列表 API:
   http://localhost:8090/api/interceptors

📍 调用树 API:
   http://localhost:8090/api/method-traces/tree

📍 查看第一个追踪的可视化详情:
   http://localhost:8090/trace/trace_1729425045123456789

💡 提示:
   1. 打开浏览器访问上面的链接
   2. 查看美观的可视化界面
   3. 点击追踪记录查看详细信息
   4. 查看方法调用树和性能指标

按 Ctrl+C 退出...
```

## 🌐 访问 Web 界面

### 1. 主页
打开浏览器访问: http://localhost:8090

### 2. API 示例

#### 获取方法列表
```bash
curl http://localhost:8090/api/interceptors
```

响应：
```json
{
  "success": true,
  "methods": [
    {
      "name": "OrderService.CreateOrder",
      "package": "main.(*OrderService).CreateOrder",
      "call_count": 5,
      "success_count": 5,
      "error_count": 0,
      "avg_duration_ms": 150
    },
    {
      "name": "OrderService.ValidateUser",
      "package": "main.(*OrderService).ValidateUser",
      "call_count": 5,
      "success_count": 5,
      "error_count": 0,
      "avg_duration_ms": 30
    },
    {
      "name": "OrderService.CheckInventory",
      "package": "main.(*OrderService).CheckInventory",
      "call_count": 5,
      "success_count": 5,
      "error_count": 0,
      "avg_duration_ms": 40
    }
  ],
  "total": 11
}
```

#### 获取调用树
```bash
curl http://localhost:8090/api/method-traces/tree
```

响应：
```json
{
  "success": true,
  "tree": [
    {
      "id": "trace_1729425045123456789",
      "method": "OrderService.CreateOrder",
      "status": "success",
      "duration": "150ms",
      "children": [
        {
          "id": "trace_1729425045123456790",
          "method": "OrderService.ValidateUser",
          "status": "success",
          "duration": "30ms"
        },
        {
          "id": "trace_1729425045123456791",
          "method": "OrderService.CheckInventory",
          "status": "success",
          "duration": "40ms"
        },
        {
          "id": "trace_1729425045123456792",
          "method": "OrderService.CalculatePrice",
          "status": "success",
          "duration": "35ms",
          "children": [
            {
              "id": "trace_1729425045123456793",
              "method": "OrderService.GetUserDiscount",
              "status": "success",
              "duration": "15ms"
            }
          ]
        }
      ]
    }
  ]
}
```

## 🎨 可视化详情页

访问: http://localhost:8090/trace/trace_1729425045123456789

页面展示：
- ✅ 美观的渐变色背景
- ✅ 时间线展示
- ✅ 性能指标卡片
- ✅ 输入参数格式化
- ✅ 返回值格式化
- ✅ 调用栈展示
- ✅ 子调用树形展示

## 📊 调用关系图

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

## 🎯 功能亮点

1. ✅ **自动方法追踪** - 无需修改业务代码
2. ✅ **嵌套调用链** - 完整展示方法调用关系
3. ✅ **性能统计** - 每个方法的调用次数和平均耗时
4. ✅ **动态方法列表** - 自动发现被追踪的方法
5. ✅ **可视化展示** - 美观的 Web 界面
6. ✅ **参数记录** - 完整记录输入输出
7. ✅ **调用栈** - 完整的调用栈信息
8. ✅ **零性能影响** - 异步处理，不阻塞业务

---

**让业务流程可视化变得简单！** 🚀

