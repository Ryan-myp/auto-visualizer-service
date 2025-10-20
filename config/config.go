package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config 配置结构
type Config struct {
	// 基础配置
	ServiceName string `json:"service_name"`
	WebPort     int    `json:"web_port"`
	DBPath      string `json:"db_path"`
	LogLevel    string `json:"log_level"`

	// 存储配置
	RetentionDays int  `json:"retention_days"`
	MaxRecords    int  `json:"max_records"`
	EnableIndex   bool `json:"enable_index"`

	// 拦截配置
	EnableTracing   bool          `json:"enable_tracing"`
	EnableRecording bool          `json:"enable_recording"`
	SampleRate      float64       `json:"sample_rate"`
	MaxTraceSize    int           `json:"max_trace_size"`
	TraceTimeout    time.Duration `json:"trace_timeout"`

	// 业务流程配置
	BusinessFlows map[string]*BusinessFlow `json:"business_flows"`
}

// BusinessFlow 业务流程定义
type BusinessFlow struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Steps       []FlowStep `json:"steps"`
	Tags        []string   `json:"tags"`
}

// FlowStep 流程步骤
type FlowStep struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Method      string        `json:"method"`
	Timeout     time.Duration `json:"timeout"`
	LogicFlow   []string      `json:"logic_flow"`
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	config := DefaultConfig()

	// 从环境变量加载配置
	if serviceName := os.Getenv("AUTO_VISUALIZER_SERVICE_NAME"); serviceName != "" {
		config.ServiceName = serviceName
	}

	if port := os.Getenv("AUTO_VISUALIZER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.WebPort = p
		}
	}

	if dbPath := os.Getenv("AUTO_VISUALIZER_DB_PATH"); dbPath != "" {
		config.DBPath = dbPath
	}

	if logLevel := os.Getenv("AUTO_VISUALIZER_LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	if retentionDays := os.Getenv("AUTO_VISUALIZER_RETENTION_DAYS"); retentionDays != "" {
		if days, err := strconv.Atoi(retentionDays); err == nil {
			config.RetentionDays = days
		}
	}

	if sampleRate := os.Getenv("AUTO_VISUALIZER_SAMPLE_RATE"); sampleRate != "" {
		if rate, err := strconv.ParseFloat(sampleRate, 64); err == nil {
			config.SampleRate = rate
		}
	}

	// 确保数据库路径包含服务名
	if config.DBPath == "./auto_visualizer.db" {
		config.DBPath = fmt.Sprintf("./%s_visualizer.db", config.ServiceName)
	}

	return config
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		ServiceName:     "unknown-service",
		WebPort:         8090,
		DBPath:          "./auto_visualizer.db",
		LogLevel:        "info",
		RetentionDays:   30,
		MaxRecords:      100000,
		EnableIndex:     true,
		EnableTracing:   true,
		EnableRecording: true,
		SampleRate:      1.0,
		MaxTraceSize:    10000,
		TraceTimeout:    30 * time.Second,
		BusinessFlows:   make(map[string]*BusinessFlow),
	}
}

// 全局配置实例
var globalConfig *Config

// SetServiceName 设置服务名称
func SetServiceName(name string) {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
	}
	globalConfig.ServiceName = name
}

// SetWebPort 设置Web端口
func SetWebPort(port int) {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
	}
	globalConfig.WebPort = port
}

// SetDBPath 设置数据库路径
func SetDBPath(path string) {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
	}
	globalConfig.DBPath = path
}

// RegisterInterceptor 注册方法拦截器
func RegisterInterceptor(methodName, flowName string) {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
	}

	if globalConfig.BusinessFlows == nil {
		globalConfig.BusinessFlows = make(map[string]*BusinessFlow)
	}

	globalConfig.BusinessFlows[methodName] = &BusinessFlow{
		ID:          methodName,
		Name:        flowName,
		Description: fmt.Sprintf("自动拦截的业务流程: %s", flowName),
		Steps:       []FlowStep{},
		Tags:        []string{"auto-intercepted"},
	}
}

// SetFlowSteps 设置流程步骤
func SetFlowSteps(methodName string, steps []FlowStep) {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
	}

	if flow, exists := globalConfig.BusinessFlows[methodName]; exists {
		flow.Steps = steps
	}
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if globalConfig == nil {
		globalConfig = LoadConfig()
	}
	return globalConfig
}
