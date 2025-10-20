package autovisualizer

import (
	"log"
	"os"

	"github.com/Ryan-myp/auto-visualizer-service/config"
	"github.com/Ryan-myp/auto-visualizer-service/interceptor"
	"github.com/Ryan-myp/auto-visualizer-service/web"
)

// 全局可视化器实例
var globalVisualizer *AutoVisualizer

// init 函数 - 包导入时自动执行
func init() {
	// 检查环境变量是否启用
	if os.Getenv("ENABLE_AUTO_VISUALIZER") != "true" {
		return
	}

	// 加载配置
	cfg := config.LoadConfig()

	// 创建可视化器实例
	visualizer, err := NewAutoVisualizer(cfg)
	if err != nil {
		log.Printf("❌ Auto-Visualizer初始化失败: %v", err)
		return
	}

	// 启动服务
	if err := visualizer.Start(); err != nil {
		log.Printf("❌ Auto-Visualizer启动失败: %v", err)
		return
	}

	globalVisualizer = visualizer

	log.Printf("🎉 Auto-Visualizer插件已自动启动!")
	log.Printf("🌐 访问地址: http://localhost:%d", cfg.WebPort)
	log.Printf("💾 数据库: %s", cfg.DBPath)
}

// AutoVisualizer 自动可视化器
type AutoVisualizer struct {
	config      *config.Config
	interceptor *interceptor.Manager
	webServer   *web.Server
}

// NewAutoVisualizer 创建自动可视化器
func NewAutoVisualizer(cfg *config.Config) (*AutoVisualizer, error) {
	// 创建拦截器管理器
	interceptorMgr, err := interceptor.NewManager(cfg)
	if err != nil {
		return nil, err
	}

	// 创建Web服务器
	webSrv, err := web.NewServer(cfg, interceptorMgr)
	if err != nil {
		return nil, err
	}

	return &AutoVisualizer{
		config:      cfg,
		interceptor: interceptorMgr,
		webServer:   webSrv,
	}, nil
}

// Start 启动可视化器
func (av *AutoVisualizer) Start() error {
	log.Printf("🚀 启动Auto-Visualizer - 服务: %s", av.config.ServiceName)

	// 启动拦截器
	if err := av.interceptor.Start(); err != nil {
		return err
	}

	// 启动Web服务器
	if err := av.webServer.Start(); err != nil {
		return err
	}

	return nil
}

// GetVisualizer 获取全局可视化器实例
func GetVisualizer() *AutoVisualizer {
	return globalVisualizer
}

// IsEnabled 检查插件是否启用
func IsEnabled() bool {
	return globalVisualizer != nil
}

// RegisterMethod 注册方法拦截（供业务代码调用）
func RegisterMethod(methodName, flowName string) {
	if globalVisualizer != nil {
		globalVisualizer.interceptor.RegisterMethod(methodName, flowName)
	}
}

// InterceptCall 拦截方法调用（供业务代码调用）
func InterceptCall(methodName string, params interface{}) string {
	if globalVisualizer != nil {
		return globalVisualizer.interceptor.InterceptCall(methodName, params)
	}
	return ""
}

// RecordResult 记录方法执行结果（供业务代码调用）
func RecordResult(traceID string, result interface{}, err error) {
	if globalVisualizer != nil {
		globalVisualizer.interceptor.RecordResult(traceID, result, err)
	}
}
