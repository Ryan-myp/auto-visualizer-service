package interceptor

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Ryan-myp/auto-visualizer-service/config"
	"github.com/Ryan-myp/auto-visualizer-service/storage"
)

// Manager 拦截器管理器
type Manager struct {
	config       *config.Config
	storage      *storage.SQLiteStorage
	interceptors map[string]*Interceptor
	execQueue    chan *ExecutionRequest
	mu           sync.RWMutex
}

// Interceptor 方法拦截器
type Interceptor struct {
	MethodName  string
	FlowName    string
	Description string
	Steps       []config.FlowStep
}

// ExecutionRequest 执行请求
type ExecutionRequest struct {
	ID        string                 `json:"id"`
	Method    string                 `json:"method"`
	FlowName  string                 `json:"flow_name"`
	Params    map[string]interface{} `json:"params"`
	StartTime time.Time              `json:"start_time"`
}

// NewManager 创建拦截器管理器
func NewManager(cfg *config.Config) (*Manager, error) {
	// 初始化SQLite存储
	storage, err := storage.NewSQLiteStorage(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("初始化存储失败: %v", err)
	}

	manager := &Manager{
		config:       cfg,
		storage:      storage,
		interceptors: make(map[string]*Interceptor),
		execQueue:    make(chan *ExecutionRequest, 1000),
	}

	// 注册默认拦截器
	manager.registerDefaultInterceptors()

	return manager, nil
}

// Start 启动拦截器管理器
func (m *Manager) Start() error {
	log.Printf("🎯 启动拦截器管理器...")

	// 启动执行引擎
	go m.startExecutionEngine()

	// 启动清理任务
	go m.startCleanupTask()

	log.Printf("✅ 拦截器管理器启动成功，已注册 %d 个拦截器", len(m.interceptors))
	return nil
}

// registerDefaultInterceptors 注册默认拦截器
func (m *Manager) registerDefaultInterceptors() {
	// 从配置中注册业务流程
	for methodName, flow := range m.config.BusinessFlows {
		m.interceptors[methodName] = &Interceptor{
			MethodName:  methodName,
			FlowName:    flow.Name,
			Description: flow.Description,
			Steps:       flow.Steps,
		}
	}

	// 不再预注册固定的方法列表
	// 方法将在实际调用时动态注册
}

// generateDefaultSteps 生成默认步骤（简化版）
func (m *Manager) generateDefaultSteps(method string) []config.FlowStep {
	return []config.FlowStep{
		{
			ID:          "execute",
			Name:        "执行业务逻辑",
			Description: fmt.Sprintf("执行 %s 方法", method),
			Method:      method,
			Timeout:     10 * time.Second,
			LogicFlow: []string{
				"1. 执行业务逻辑",
				"2. 处理返回结果",
			},
		},
	}
}

// RegisterMethod 注册方法拦截器
func (m *Manager) RegisterMethod(methodName, flowName string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.interceptors[methodName] = &Interceptor{
		MethodName:  methodName,
		FlowName:    flowName,
		Description: fmt.Sprintf("业务流程: %s", flowName),
		Steps:       m.generateDefaultSteps(methodName),
	}

	log.Printf("✅ 注册方法拦截器: %s -> %s", methodName, flowName)
}

// InterceptCall 拦截方法调用
func (m *Manager) InterceptCall(methodName string, params interface{}) string {
	m.mu.RLock()
	interceptor, exists := m.interceptors[methodName]
	m.mu.RUnlock()

	if !exists {
		// 自动注册未知方法
		m.RegisterMethod(methodName, fmt.Sprintf("自动拦截-%s", methodName))
		m.mu.RLock()
		interceptor = m.interceptors[methodName]
		m.mu.RUnlock()
	}

	// 创建执行请求
	req := &ExecutionRequest{
		ID:       fmt.Sprintf("trace_%d", time.Now().UnixNano()),
		Method:   methodName,
		FlowName: interceptor.FlowName,
		Params: map[string]interface{}{
			"method":    methodName,
			"params":    params,
			"timestamp": time.Now().Format(time.RFC3339),
		},
		StartTime: time.Now(),
	}

	// 异步处理
	select {
	case m.execQueue <- req:
		log.Printf("🔍 已拦截方法调用: %s (ID: %s)", methodName, req.ID)
	default:
		log.Printf("⚠️ 执行队列已满，跳过拦截: %s", methodName)
	}

	return req.ID
}

// RecordResult 记录方法执行结果
func (m *Manager) RecordResult(traceID string, result interface{}, err error) {
	// 这里可以更新已存在的trace记录
	// 为简化演示，暂时不实现
	log.Printf("📝 记录执行结果: %s", traceID)
}

// startExecutionEngine 启动执行引擎
func (m *Manager) startExecutionEngine() {
	log.Printf("🎯 执行引擎已启动，等待拦截...")

	for req := range m.execQueue {
		go m.processExecution(req)
	}
}

// processExecution 处理执行请求
func (m *Manager) processExecution(req *ExecutionRequest) {
	log.Printf("🚀 处理拦截调用: %s (ID: %s)", req.Method, req.ID)

	// 获取拦截器
	m.mu.RLock()
	interceptor, exists := m.interceptors[req.Method]
	m.mu.RUnlock()

	if !exists {
		log.Printf("❌ 未找到拦截器: %s", req.Method)
		return
	}

	// 创建执行追踪
	trace := &storage.ExecutionTrace{
		ID:          req.ID,
		ServiceName: m.config.ServiceName,
		FlowName:    req.FlowName,
		Method:      req.Method,
		Status:      "running",
		StartTime:   req.StartTime,
		Input:       req.Params,
		Steps:       []storage.ExecutionStep{},
		CreatedAt:   time.Now(),
	}

	// 保存初始状态
	m.storage.SaveTrace(trace)

	// 执行步骤
	for i, stepDef := range interceptor.Steps {
		stepStart := time.Now()

		step := storage.ExecutionStep{
			StepNumber:      i + 1,
			StepName:        stepDef.Name,
			Method:          stepDef.Method,
			Status:          "running",
			StartTime:       stepStart,
			LogicFlow:       stepDef.LogicFlow,
			BusinessContext: stepDef.Description,
			Input: map[string]interface{}{
				"step":      i + 1,
				"method":    req.Method,
				"timestamp": stepStart.Format(time.RFC3339),
			},
		}

		trace.Steps = append(trace.Steps, step)
		m.storage.SaveTrace(trace)

		// 模拟执行时间
		executionTime := time.Duration(100+i*50) * time.Millisecond
		if stepDef.Timeout > 0 && executionTime > stepDef.Timeout {
			executionTime = stepDef.Timeout / 2
		}
		time.Sleep(executionTime)

		// 完成步骤
		endTime := time.Now()
		step.EndTime = &endTime
		step.Duration = endTime.Sub(stepStart)
		step.Status = "completed"
		step.Output = map[string]interface{}{
			"success":  true,
			"duration": step.Duration.String(),
			"result":   fmt.Sprintf("步骤%d执行完成", i+1),
		}

		trace.Steps[i] = step
		m.storage.SaveTrace(trace)

		log.Printf("  ✅ 步骤%d: %s - 完成 (%v)", i+1, stepDef.Name, step.Duration)
	}

	// 完成追踪
	endTime := time.Now()
	trace.EndTime = &endTime
	trace.Duration = endTime.Sub(trace.StartTime)
	trace.Status = "completed"
	trace.Output = map[string]interface{}{
		"success":     true,
		"totalSteps":  len(interceptor.Steps),
		"totalTime":   trace.Duration.String(),
		"intercepted": true,
	}

	m.storage.SaveTrace(trace)

	log.Printf("🎉 拦截处理完成: %s (总耗时: %v)", req.Method, trace.Duration)
}

// startCleanupTask 启动清理任务
func (m *Manager) startCleanupTask() {
	ticker := time.NewTicker(24 * time.Hour) // 每天清理一次
	defer ticker.Stop()

	for range ticker.C {
		if err := m.storage.CleanupOldRecords(m.config.RetentionDays); err != nil {
			log.Printf("❌ 清理过期记录失败: %v", err)
		}
	}
}

// GetStorage 获取存储实例
func (m *Manager) GetStorage() *storage.SQLiteStorage {
	return m.storage
}

// GetInterceptors 获取所有拦截器
func (m *Manager) GetInterceptors() map[string]*Interceptor {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*Interceptor)
	for k, v := range m.interceptors {
		result[k] = v
	}
	return result
}
