# Auto-Visualizer 使用指南

## 🚀 快速开始（3步搞定）

### 步骤1：安装

```bash
go get github.com/Ryan-myp/auto-visualizer-service
```

### 步骤2：在代码中导入

```go
package main

import (
    _ "github.com/Ryan-myp/auto-visualizer-service"  // 只需这一行！
    autovisualizer "github.com/Ryan-myp/auto-visualizer-service"
)

func ProcessOrder(orderID string, amount float64) (string, error) {
    // 添加追踪（可选，推荐）
    end := autovisualizer.Begin("ProcessOrder", orderID, amount)
    var result string
    var err error
    defer func() { end(result, err) }()
    
    // 你的业务逻辑
    result = "success"
    return result, nil
}

func main() {
    ProcessOrder("ORD-001", 999.99)
    select {} // 保持运行
}
```

### 步骤3：启动并查看

```bash
# 设置环境变量启用追踪
export ENABLE_AUTO_VISUALIZER=true

# 运行程序
go run main.go

# 打开浏览器访问
open http://localhost:8090
```

## 📖 使用方式对比

### 方式1：最简单 - 零代码修改

适用场景：快速体验，不需要详细的参数追踪

```go
import _ "github.com/Ryan-myp/auto-visualizer-service"

func main() {
    // 业务代码，完全不需要修改
    DoSomething()
}
```

**优点**：零侵入
**缺点**：无法记录方法参数和返回值

---

### 方式2：Begin + defer（⭐推荐）

适用场景：需要追踪方法的入参、出参和执行时间

```go
func ProcessOrder(orderID string, amount float64) (string, error) {
    end := autovisualizer.Begin("ProcessOrder", orderID, amount)
    
    var result string
    var err error
    defer func() {
        end(result, err)  // 自动记录返回值
    }()
    
    // 业务逻辑
    result = "success"
    return result, nil
}
```

**优点**：
- ✅ 记录完整的入参和出参
- ✅ 自动计算执行时间
- ✅ 支持嵌套调用追踪
- ✅ 错误自动记录

**使用技巧**：
- 必须在 defer 中调用 end()
- 返回值必须是命名变量或在 defer 中赋值

---

### 方式3：Measure - 只测量时间

适用场景：只关心执行时间，不关心参数

```go
func ValidateOrder(orderID string) bool {
    defer autovisualizer.Measure("ValidateOrder")()
    
    // 业务逻辑
    return true
}
```

**优点**：
- ✅ 最简洁
- ✅ 自动计算执行时间

**缺点**：
- ❌ 不记录参数和返回值

---

### 方式4：TraceMethod - 简单追踪

适用场景：需要追踪但不需要参数细节

```go
func SaveOrder(orderID string) {
    defer autovisualizer.TraceMethod("SaveOrder")()
    
    // 业务逻辑
}
```

---

### 方式5：装饰器模式

适用场景：包装现有函数，不修改原函数

```go
// 原始函数
func calculateSum(a, b int) int {
    return a + b
}

// 包装后的函数（带追踪）
var tracedCalculateSum = autovisualizer.Trace(
    "calculateSum", 
    calculateSum,
).(func(int, int) int)

func main() {
    result := tracedCalculateSum(10, 20)  // 自动追踪
}
```

**优点**：
- ✅ 不修改原函数
- ✅ 可以选择性地使用追踪版本

**缺点**：
- ❌ 需要类型断言
- ❌ 语法稍复杂

---

## 🌲 嵌套调用追踪

Auto-Visualizer 会自动识别嵌套调用，构建调用树：

```go
func ProcessUserBatch(names []string) error {
    end := autovisualizer.Begin("ProcessUserBatch", names)
    defer end(nil)
    
    for _, name := range names {
        CreateUser(name)      // 子调用1
        ValidateUser(name)    // 子调用2
    }
    
    return nil
}

func CreateUser(name string) {
    defer autovisualizer.Measure("CreateUser")()
    // 业务逻辑
}

func ValidateUser(name string) {
    defer autovisualizer.Measure("ValidateUser")()
    // 业务逻辑
}
```

调用树结构：
```
ProcessUserBatch
├── CreateUser (Alice)
├── ValidateUser (Alice)
├── CreateUser (Bob)
└── ValidateUser (Bob)
```

---

## 🔧 配置选项

### 环境变量配置

```bash
# 必需：启用追踪
export ENABLE_AUTO_VISUALIZER=true

# 可选：配置端口
export AUTO_VISUALIZER_PORT=8090

# 可选：配置数据库路径
export AUTO_VISUALIZER_DB_PATH=./traces.db

# 可选：配置服务名称
export AUTO_VISUALIZER_SERVICE_NAME=my-service

# 可选：配置采样率（0.0-1.0）
export AUTO_VISUALIZER_SAMPLE_RATE=1.0
```

### 代码配置

```go
import "github.com/Ryan-myp/auto-visualizer-service/config"

func init() {
    config.SetServiceName("MyService")
    config.SetWebPort(9090)
    config.SetDBPath("./my_traces.db")
}
```

---

## 🌐 Web API

### 查看所有追踪

```bash
curl http://localhost:8090/api/method-traces
```

响应：
```json
{
  "success": true,
  "total": 10,
  "traces": [
    {
      "trace_id": "trace_123456",
      "method_name": "ProcessOrder",
      "input": ["ORD-001", 999.99],
      "output": ["success", null],
      "duration": "150ms",
      "status": "success"
    }
  ]
}
```

### 查看调用树

```bash
curl http://localhost:8090/api/method-traces/tree
```

### 清除追踪记录

```bash
curl -X DELETE http://localhost:8090/api/method-traces
```

---

## 💡 最佳实践

### 1. 选择合适的追踪方式

- **关键业务方法**：使用 `Begin + defer`，记录完整信息
- **性能敏感方法**：使用 `Measure`，只记录时间
- **快速调试**：使用 `TraceMethod`，快速添加追踪

### 2. 合理命名方法

```go
// ✅ 好的命名
autovisualizer.Begin("UserService.CreateUser", ...)
autovisualizer.Begin("OrderService.ProcessPayment", ...)

// ❌ 不好的命名
autovisualizer.Begin("func1", ...)
autovisualizer.Begin("process", ...)
```

### 3. 控制追踪粒度

```go
// ✅ 追踪业务方法
func ProcessOrder(orderID string) error {
    end := autovisualizer.Begin("ProcessOrder", orderID)
    defer end(nil)
    
    // ❌ 不要追踪内部小函数
    validateID(orderID)  // 不需要追踪
    saveToCache(orderID) // 不需要追踪
    
    return nil
}
```

### 4. 生产环境使用

```go
// 通过环境变量控制
// 开发环境：ENABLE_AUTO_VISUALIZER=true
// 生产环境：ENABLE_AUTO_VISUALIZER=false

// 或者使用采样率
// ENABLE_AUTO_VISUALIZER=true
// AUTO_VISUALIZER_SAMPLE_RATE=0.01  // 只追踪1%的请求
```

---

## 🎯 实际应用场景

### 场景1：调试复杂业务流程

```go
func CreateCampaign(req *CreateCampaignRequest) (*Campaign, error) {
    end := autovisualizer.Begin("CreateCampaign", req)
    var campaign *Campaign
    var err error
    defer func() { end(campaign, err) }()
    
    // 复杂的业务逻辑
    campaign, err = validateRequest(req)
    if err != nil {
        return nil, err
    }
    
    campaign, err = callExternalAPI(campaign)
    if err != nil {
        return nil, err
    }
    
    campaign, err = saveToDatabase(campaign)
    return campaign, err
}
```

通过 Web UI 可以看到：
- 每个步骤的执行时间
- 哪个步骤出错了
- 完整的调用链

### 场景2：性能分析

```go
func BatchProcess(items []Item) {
    defer autovisualizer.Measure("BatchProcess")()
    
    for _, item := range items {
        processItem(item)  // 也可以追踪
    }
}

func processItem(item Item) {
    defer autovisualizer.Measure("processItem")()
    // 处理逻辑
}
```

通过 API 可以获取：
- 总执行时间
- 单个 item 的平均处理时间
- 性能瓶颈在哪里

### 场景3：线上问题排查

当线上出现问题时，启用追踪：

```bash
# 临时启用追踪（只追踪1%的请求）
export ENABLE_AUTO_VISUALIZER=true
export AUTO_VISUALIZER_SAMPLE_RATE=0.01

# 重启服务
systemctl restart my-service

# 查看追踪数据
curl http://localhost:8090/api/method-traces/tree
```

---

## 📊 与 pprof 的对比

| 特性 | pprof | Auto-Visualizer |
|------|-------|-----------------|
| CPU分析 | ✅ | ❌ |
| 内存分析 | ✅ | ❌ |
| 方法追踪 | ⚠️ 需要手动埋点 | ✅ 自动追踪 |
| 参数记录 | ❌ | ✅ |
| 返回值记录 | ❌ | ✅ |
| 调用链可视化 | ⚠️ 有限 | ✅ 完整 |
| Web UI | ✅ | ✅ |
| 持久化 | ❌ | ✅ SQLite |
| 零侵入 | ✅ | ✅ |

**结论**：pprof 适合性能分析，Auto-Visualizer 适合业务逻辑追踪和调试。

---

## ❓ 常见问题

### Q1: 会影响性能吗？

A: 有轻微影响，但可控：
- 采样率设置为 0.1，只追踪 10% 的请求
- 只追踪关键方法，不追踪所有方法
- 生产环境可以完全关闭

### Q2: 如何在生产环境使用？

A: 建议：
- 默认关闭：`ENABLE_AUTO_VISUALIZER=false`
- 需要调试时临时开启
- 使用低采样率：`AUTO_VISUALIZER_SAMPLE_RATE=0.01`

### Q3: 追踪数据会占用多少空间？

A: 
- 默认保留 30 天
- 可配置：`AUTO_VISUALIZER_RETENTION_DAYS=7`
- 自动清理过期数据

### Q4: 支持分布式追踪吗？

A: 当前版本不支持，未来计划支持：
- 跨服务追踪
- TraceID 传递
- 与 OpenTelemetry 集成

---

## 🔗 相关链接

- GitHub: https://github.com/Ryan-myp/auto-visualizer-service
- 示例代码: `examples/quickstart/main.go`
- 完整示例: `examples/trace-demo/main.go`

---

## 📝 更新日志

### v1.1.0 (2025-10-20)
- ✨ 新增方法追踪功能
- ✨ 支持调用链可视化
- ✨ 自动记录入参出参
- 📝 完善文档和示例

### v1.0.0 (2025-10-19)
- 🎉 首次发布
- ✨ 基础流程可视化
- ✨ SQLite 持久化

