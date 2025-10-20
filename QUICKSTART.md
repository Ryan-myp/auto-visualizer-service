# 快速开始指南

## ⚡ 5分钟快速集成

### 步骤 1: 安装依赖

```bash
go get github.com/Ryan-myp/auto-visualizer-service@latest
```

### 步骤 2: 设置环境变量

**必须设置这个环境变量，否则插件不会启动！**

```bash
export ENABLE_AUTO_VISUALIZER=true
```

可选配置：
```bash
export AUTO_VISUALIZER_PORT=8090              # Web UI 端口，默认 8090
export AUTO_VISUALIZER_SERVICE_NAME=my-app    # 服务名称
export AUTO_VISUALIZER_SAMPLE_RATE=1.0        # 采样率，1.0 = 100%
```

### 步骤 3: 在代码中导入

**方式 A: 只导入（零代码修改）**

```go
package main

import (
    _ "github.com/Ryan-myp/auto-visualizer-service"  // 只需这一行！
)

func main() {
    // 你的业务代码
    ProcessOrder("ORD-001")
    
    // 保持程序运行以便查看 Web UI
    select {}
}
```

**方式 B: 手动追踪（推荐）**

```go
package main

import (
    _ "github.com/Ryan-myp/auto-visualizer-service"
    autovisualizer "github.com/Ryan-myp/auto-visualizer-service"
)

func ProcessOrder(orderID string) (string, error) {
    // 添加追踪
    end := autovisualizer.Begin("ProcessOrder", orderID)
    var result string
    var err error
    defer func() {
        end(result, err)  // 自动记录返回值
    }()
    
    // 你的业务逻辑
    result = "success"
    return result, nil
}

func main() {
    ProcessOrder("ORD-001")
    select {}
}
```

### 步骤 4: 运行并查看

```bash
# 运行程序
go run main.go
```

你应该看到类似的输出：
```
🚀 启动Auto-Visualizer - 服务: unknown-service
🌐 Web服务器启动中: http://localhost:8090
✅ Web服务器已启动: http://localhost:8090
🎉 Auto-Visualizer插件已自动启动!
🌐 访问地址: http://localhost:8090
💾 数据库: ./unknown-service_visualizer.db
🔍 方法追踪器已启用 (采样率: 1.00)
```

### 步骤 5: 访问 Web UI

打开浏览器访问：
- **主页**: http://localhost:8090
- **追踪列表**: http://localhost:8090/api/method-traces
- **调用树**: http://localhost:8090/api/method-traces/tree
- **健康检查**: http://localhost:8090/health

---

## ⚠️ 常见问题

### 问题 1: 端口没有被监听

**症状**: 无法访问 http://localhost:8090

**原因**: 
1. 没有设置环境变量 `ENABLE_AUTO_VISUALIZER=true`
2. 端口被占用
3. 程序启动后立即退出

**解决方案**:

```bash
# 1. 确保设置环境变量
export ENABLE_AUTO_VISUALIZER=true

# 2. 检查端口是否被占用
lsof -i :8090

# 3. 确保程序保持运行
# 在 main 函数末尾添加:
select {}  // 或 time.Sleep(time.Hour)
```

### 问题 2: 看不到追踪数据

**原因**: 
1. 没有调用追踪方法
2. 采样率设置过低
3. 熔断器打开

**解决方案**:

```go
// 检查是否启用
if autovisualizer.IsEnabled() {
    fmt.Println("✅ Auto-Visualizer 已启用")
} else {
    fmt.Println("❌ Auto-Visualizer 未启用")
}

// 检查追踪数据
traces := autovisualizer.GetAllTraces()
fmt.Printf("已记录 %d 条追踪\n", len(traces))

// 检查熔断器状态
status := autovisualizer.GetTracer().GetCircuitStatus()
fmt.Printf("熔断器状态: %+v\n", status)
```

### 问题 3: 性能影响

**如果担心性能影响**:

```bash
# 方案 1: 降低采样率（只追踪 10% 的请求）
export AUTO_VISUALIZER_SAMPLE_RATE=0.1

# 方案 2: 生产环境关闭
export ENABLE_AUTO_VISUALIZER=false

# 方案 3: 只在需要时开启
# 默认不设置环境变量，需要调试时再设置
```

---

## 📝 完整示例

```go
package main

import (
    "fmt"
    "time"
    
    _ "github.com/Ryan-myp/auto-visualizer-service"
    autovisualizer "github.com/Ryan-myp/auto-visualizer-service"
)

// 示例 1: 使用 Begin (推荐)
func CreateOrder(orderID string, amount float64) (string, error) {
    end := autovisualizer.Begin("CreateOrder", orderID, amount)
    var result string
    var err error
    defer func() { end(result, err) }()
    
    // 业务逻辑
    time.Sleep(50 * time.Millisecond)
    
    // 调用其他方法（会自动形成调用链）
    ValidateOrder(orderID)
    SaveOrder(orderID)
    
    result = "success"
    return result, nil
}

// 示例 2: 使用 Measure (只测量时间)
func ValidateOrder(orderID string) bool {
    defer autovisualizer.Measure("ValidateOrder")()
    
    time.Sleep(20 * time.Millisecond)
    return true
}

// 示例 3: 使用 TraceMethod (简单追踪)
func SaveOrder(orderID string) {
    defer autovisualizer.TraceMethod("SaveOrder")()
    
    time.Sleep(30 * time.Millisecond)
}

func main() {
    fmt.Println("🚀 Auto-Visualizer 示例程序")
    fmt.Println()
    
    // 检查是否启用
    if !autovisualizer.IsEnabled() {
        fmt.Println("❌ Auto-Visualizer 未启用")
        fmt.Println("请运行: export ENABLE_AUTO_VISUALIZER=true")
        return
    }
    
    fmt.Println("✅ Auto-Visualizer 已启用")
    fmt.Println()
    
    // 等待 Web 服务器完全启动
    time.Sleep(500 * time.Millisecond)
    
    // 执行业务逻辑
    fmt.Println("📝 执行业务逻辑...")
    CreateOrder("ORD-001", 999.99)
    CreateOrder("ORD-002", 1299.50)
    
    // 查看追踪统计
    traces := autovisualizer.GetAllTraces()
    fmt.Printf("\n📊 已记录 %d 条追踪\n", len(traces))
    
    fmt.Println()
    fmt.Println("🌐 访问 Web UI: http://localhost:8090")
    fmt.Println("🌐 查看追踪: http://localhost:8090/api/method-traces")
    fmt.Println()
    fmt.Println("按 Ctrl+C 退出...")
    
    // 保持运行
    select {}
}
```

运行：
```bash
export ENABLE_AUTO_VISUALIZER=true
go run main.go
```

---

## 🎯 验证清单

在集成后，请验证以下几点：

- [ ] 设置了环境变量 `ENABLE_AUTO_VISUALIZER=true`
- [ ] 程序启动时看到 "Auto-Visualizer插件已自动启动!" 日志
- [ ] 可以访问 http://localhost:8090/health
- [ ] 执行业务方法后可以在 http://localhost:8090/api/method-traces 看到数据
- [ ] Web UI 正常显示

---

## 🆘 还有问题？

1. 查看详细文档: [README.md](README.md)
2. 查看使用指南: [USAGE.md](USAGE.md)
3. 查看安全说明: [SAFETY.md](SAFETY.md)
4. 运行测试示例: `cd examples/test-integration && go run main.go`

---

## 💡 提示

- **开发环境**: 设置 `ENABLE_AUTO_VISUALIZER=true`
- **生产环境**: 不设置环境变量（默认关闭）或设置 `ENABLE_AUTO_VISUALIZER=false`
- **调试时**: 临时设置环境变量并重启服务
- **性能测试**: 设置 `AUTO_VISUALIZER_SAMPLE_RATE=0.01` (1% 采样)

