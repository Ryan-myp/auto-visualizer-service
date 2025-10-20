# Auto-Visualizer Service

🔌 **非侵入式业务流程可视化插件**

## 🎯 功能特性

- **🔌 非侵入式设计**: 业务代码零修改，只需引入包即可
- **💾 SQLite持久化**: 执行记录自动保存到数据库
- **📊 分页查询**: 支持大量数据的高效查询和筛选
- **🎨 可视化界面**: 美观的Web界面展示执行逻辑
- **📈 统计分析**: 自动计算性能指标和执行趋势
- **🔍 方法拦截**: 自动拦截和记录业务方法调用
- **🎮 环境控制**: 通过环境变量灵活启用/禁用

## 🚀 快速开始

### 1. 引入依赖

在你的服务的 `go.mod` 中添加：

```go
require (
    github.com/Ryan-myp/auto-visualizer-service v1.0.0
)
```

### 2. 导入包

在你的 `main.go` 或任意入口文件中导入：

```go
import (
    _ "github.com/Ryan-myp/auto-visualizer-service" // 自动启动插件
)

func main() {
    // 你的业务代码，无需任何修改
    // 插件会自动拦截和记录方法调用
}
```

### 3. 启用插件

设置环境变量启用插件：

```bash
export ENABLE_AUTO_VISUALIZER=true
export AUTO_VISUALIZER_PORT=8090  # 可选，默认8090
export AUTO_VISUALIZER_DB_PATH=./visualizer.db  # 可选

# 运行你的服务
go run main.go
```

### 4. 访问可视化界面

打开浏览器访问: `http://localhost:8090`

## 📖 使用示例

### 基础使用

```go
package main

import (
    _ "github.com/Ryan-myp/auto-visualizer-service" // 导入插件
)

// 你的业务方法 - 会被自动拦截
func ProcessOrder(orderID string, userID int64) error {
    // 业务逻辑...
    return nil
}

func main() {
    // 业务代码，插件自动工作
    ProcessOrder("order_123", 888888)
    
    // 服务继续运行...
    select {}
}
```

### 高级配置

```go
package main

import (
    "github.com/Ryan-myp/auto-visualizer-service/config"
    _ "github.com/Ryan-myp/auto-visualizer-service" // 自动启动
)

func init() {
    // 可选：自定义配置
    config.SetServiceName("MyService")
    config.SetWebPort(9090)
    config.SetDBPath("./my_service_traces.db")
    
    // 注册自定义方法拦截器
    config.RegisterInterceptor("ProcessOrder", "订单处理流程")
    config.RegisterInterceptor("CreateCampaign", "广告创建流程")
}
```

## 🔧 配置选项

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `ENABLE_AUTO_VISUALIZER` | `false` | 是否启用插件 |
| `AUTO_VISUALIZER_PORT` | `8090` | Web界面端口 |
| `AUTO_VISUALIZER_DB_PATH` | `./auto_visualizer.db` | SQLite数据库路径 |
| `AUTO_VISUALIZER_SERVICE_NAME` | `unknown-service` | 服务名称 |
| `AUTO_VISUALIZER_LOG_LEVEL` | `info` | 日志级别 |

### 代码配置

```go
import "github.com/Ryan-myp/auto-visualizer-service/config"

// 设置服务名称
config.SetServiceName("AdMgmt")

// 设置Web端口
config.SetWebPort(8090)

// 设置数据库路径
config.SetDBPath("./traces.db")

// 注册方法拦截器
config.RegisterInterceptor("MethodName", "业务流程描述")

// 设置业务流程步骤
config.SetFlowSteps("MethodName", []config.FlowStep{
    {Name: "步骤1", Description: "业务含义1"},
    {Name: "步骤2", Description: "业务含义2"},
})
```

## 📊 API接口

### 获取执行记录

```bash
# 分页查询
GET /api/traces?page=1&page_size=10&status=completed

# 查询特定方法
GET /api/traces?method=ProcessOrder

# 时间范围查询
GET /api/traces?start_time=2024-01-01&end_time=2024-01-31
```

### 获取详细信息

```bash
# 获取执行详情
GET /api/trace/{trace_id}

# 获取统计信息
GET /api/stats?service_name=MyService&days=30

# 健康检查
GET /health
```

## 🎨 Web界面功能

- **📊 执行历史**: 分页展示所有拦截记录
- **🔍 详细追踪**: 每次执行的完整步骤和参数
- **📈 统计分析**: 成功率、平均耗时、执行趋势
- **🎯 实时监控**: 正在执行的方法实时状态
- **🔧 参数查看**: 输入输出参数的JSON格式展示
- **📋 业务上下文**: 每个步骤的业务含义解释

## 🏗️ 架构设计

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   业务服务       │    │  Auto-Visualizer │    │   SQLite DB     │
│                │    │     插件          │    │                │
│ ┌─────────────┐ │    │ ┌──────────────┐ │    │ ┌─────────────┐ │
│ │业务方法调用  │ ├────┤ │方法拦截器     │ ├────┤ │执行记录表    │ │
│ └─────────────┘ │    │ └──────────────┘ │    │ └─────────────┘ │
│                │    │ ┌──────────────┐ │    │ ┌─────────────┐ │
│ ┌─────────────┐ │    │ │执行引擎      │ │    │ │统计信息表    │ │
│ │参数和返回值  │ ├────┤ └──────────────┘ ├────┤ └─────────────┘ │
│ └─────────────┘ │    │ ┌──────────────┐ │    │                │
│                │    │ │Web服务器     │ │    │                │
└─────────────────┘    │ └──────────────┘ │    └─────────────────┘
                       └──────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │   Web界面        │
                       │  (localhost:8090) │
                       └──────────────────┘
```

## 🔍 工作原理

1. **自动启动**: 通过 `init()` 函数在包导入时自动启动
2. **方法拦截**: 使用反射和代理模式拦截业务方法调用
3. **执行记录**: 记录方法的输入参数、输出结果、执行时间
4. **步骤分析**: 自动分析方法内部的执行步骤和业务逻辑
5. **数据持久化**: 将所有记录保存到SQLite数据库
6. **Web展示**: 提供美观的界面展示执行历史和统计信息

## 🛠️ 开发和贡献

### 本地开发

```bash
# 克隆项目
git clone https://github.com/Ryan-myp/auto-visualizer-service.git
cd auto-visualizer

# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 构建
go build -o auto-visualizer ./cmd/
```

### 项目结构

```
auto-visualizer/
├── README.md
├── go.mod
├── go.sum
├── init.go              # 自动启动入口
├── config/              # 配置管理
│   ├── config.go
│   └── env.go
├── interceptor/         # 方法拦截器
│   ├── interceptor.go
│   └── registry.go
├── storage/             # 数据存储
│   ├── sqlite.go
│   └── models.go
├── web/                 # Web服务
│   ├── server.go
│   ├── handlers.go
│   └── templates/
├── analyzer/            # 代码分析
│   └── flow_analyzer.go
└── examples/            # 使用示例
    ├── basic/
    └── advanced/
```

## 📝 更新日志

### v1.0.0 (2024-10-20)
- 🎉 首次发布
- ✅ 非侵入式方法拦截
- ✅ SQLite持久化存储
- ✅ Web可视化界面
- ✅ 分页查询和统计
- ✅ 环境变量配置

## 📄 许可证

MIT License

## 🤝 支持

如有问题或建议，请提交 Issue 或 Pull Request。

---

**让业务流程可视化变得简单！** 🚀
