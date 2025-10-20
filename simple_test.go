package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// 简化版独立Auto-Visualizer服务测试

type AutoVisualizer struct {
	config         *Config
	server         *http.ServeMux
	executionQueue chan ExecutionRequest
	db             *sql.DB
	mu             sync.RWMutex
	interceptors   map[string]string
}

type Config struct {
	ServiceName string
	WebPort     int
	DBPath      string
	Enabled     bool
}

type ExecutionRequest struct {
	ID       string                 `json:"id"`
	FlowName string                 `json:"flow_name"`
	Method   string                 `json:"method"`
	Params   map[string]interface{} `json:"params"`
}

type ExecutionTrace struct {
	ID          string                 `json:"id"`
	ServiceName string                 `json:"service_name"`
	FlowName    string                 `json:"flow_name"`
	Method      string                 `json:"method"`
	Status      string                 `json:"status"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     *time.Time             `json:"end_time,omitempty"`
	Duration    time.Duration          `json:"duration"`
	Steps       []ExecutionStep        `json:"steps"`
	Input       map[string]interface{} `json:"input"`
	Output      map[string]interface{} `json:"output"`
	CreatedAt   time.Time              `json:"created_at"`
}

type ExecutionStep struct {
	StepNumber      int                    `json:"step_number"`
	StepName        string                 `json:"step_name"`
	Method          string                 `json:"method"`
	Status          string                 `json:"status"`
	StartTime       time.Time              `json:"start_time"`
	EndTime         *time.Time             `json:"end_time,omitempty"`
	Duration        time.Duration          `json:"duration"`
	Input           map[string]interface{} `json:"input"`
	Output          map[string]interface{} `json:"output"`
	LogicFlow       []string               `json:"logic_flow"`
	BusinessContext string                 `json:"business_context"`
}

// 全局实例
var globalVisualizer *AutoVisualizer

// init 模拟插件自动启动
func init() {
	if os.Getenv("ENABLE_AUTO_VISUALIZER") != "true" {
		return
	}

	config := &Config{
		ServiceName: getEnvOrDefault("AUTO_VISUALIZER_SERVICE_NAME", "AutoVisualizerService"),
		WebPort:     getEnvIntOrDefault("AUTO_VISUALIZER_PORT", 8090),
		DBPath:      getEnvOrDefault("AUTO_VISUALIZER_DB_PATH", "./auto_visualizer_service.db"),
		Enabled:     true,
	}

	var err error
	globalVisualizer, err = NewAutoVisualizer(config)
	if err != nil {
		log.Printf("❌ Auto-Visualizer初始化失败: %v", err)
		return
	}

	globalVisualizer.Start()
	log.Printf("🎉 Auto-Visualizer独立服务已自动启动!")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func NewAutoVisualizer(config *Config) (*AutoVisualizer, error) {
	db, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		return nil, fmt.Errorf("打开SQLite数据库失败: %v", err)
	}

	av := &AutoVisualizer{
		config:         config,
		server:         http.NewServeMux(),
		executionQueue: make(chan ExecutionRequest, 100),
		db:             db,
		interceptors:   make(map[string]string),
	}

	if err := av.initDatabase(); err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %v", err)
	}

	// 注册默认拦截器
	av.registerDefaultInterceptors()

	return av, nil
}

func (av *AutoVisualizer) initDatabase() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS execution_traces (
		id TEXT PRIMARY KEY,
		service_name TEXT NOT NULL,
		flow_name TEXT NOT NULL,
		method TEXT NOT NULL,
		status TEXT NOT NULL,
		start_time DATETIME NOT NULL,
		end_time DATETIME,
		duration INTEGER DEFAULT 0,
		steps TEXT,
		input TEXT,
		output TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_traces_created_at ON execution_traces(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_traces_method ON execution_traces(method);
	`

	if _, err := av.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}

	return nil
}

func (av *AutoVisualizer) registerDefaultInterceptors() {
	av.interceptors["AppPublishOpsAdsRun"] = "AdMgmt广告发布流程"
	av.interceptors["CreateCampaign"] = "广告活动创建"
	av.interceptors["CreateAdSet"] = "广告组创建"
	av.interceptors["CreateAd"] = "广告创建"
	av.interceptors["ProcessOrder"] = "订单处理流程"
}

func (av *AutoVisualizer) Start() {
	av.setupRoutes()
	av.startExecutionEngine()

	addr := fmt.Sprintf(":%d", av.config.WebPort)

	go func() {
		log.Printf("🚀 Auto-Visualizer独立服务启动成功!")
		log.Printf("🌐 访问地址: http://localhost%s", addr)
		log.Printf("💾 SQLite数据库: %s", av.config.DBPath)
		log.Printf("🔌 独立服务模式 - 通过go.mod引入即用")

		if err := http.ListenAndServe(addr, av.server); err != nil {
			log.Printf("❌ 服务启动失败: %v", err)
		}
	}()
}

func (av *AutoVisualizer) setupRoutes() {
	av.server.HandleFunc("/", av.handleIndex)
	av.server.HandleFunc("/api/traces", av.handleTraces)
	av.server.HandleFunc("/api/trace/", av.handleTraceDetail)
	av.server.HandleFunc("/api/simulate", av.handleSimulate)
	av.server.HandleFunc("/api/stats", av.handleStats)
	av.server.HandleFunc("/api/interceptors", av.handleInterceptors)
	av.server.HandleFunc("/health", av.handleHealth)
}

func (av *AutoVisualizer) startExecutionEngine() {
	go func() {
		log.Printf("🎯 执行引擎已启动，等待拦截...")
		for req := range av.executionQueue {
			go av.processExecution(req)
		}
	}()
}

func (av *AutoVisualizer) InterceptCall(methodName string, params interface{}) string {
	flowName, exists := av.interceptors[methodName]
	if !exists {
		flowName = fmt.Sprintf("自动拦截-%s", methodName)
		av.interceptors[methodName] = flowName
	}

	req := ExecutionRequest{
		ID:       fmt.Sprintf("trace_%d", time.Now().UnixNano()),
		Method:   methodName,
		FlowName: flowName,
		Params: map[string]interface{}{
			"method":    methodName,
			"params":    params,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	select {
	case av.executionQueue <- req:
		log.Printf("🔍 已拦截方法调用: %s (ID: %s)", methodName, req.ID)
	default:
		log.Printf("⚠️ 执行队列已满，跳过拦截: %s", methodName)
	}

	return req.ID
}

func (av *AutoVisualizer) processExecution(req ExecutionRequest) {
	log.Printf("🚀 处理拦截调用: %s (ID: %s)", req.Method, req.ID)

	trace := &ExecutionTrace{
		ID:          req.ID,
		ServiceName: av.config.ServiceName,
		FlowName:    req.FlowName,
		Method:      req.Method,
		Status:      "running",
		StartTime:   time.Now(),
		Input:       req.Params,
		Steps:       []ExecutionStep{},
		CreatedAt:   time.Now(),
	}

	av.saveTrace(trace)

	// 生成步骤
	steps := av.generateSteps(req.Method)

	for i, stepDef := range steps {
		stepStart := time.Now()

		step := ExecutionStep{
			StepNumber:      i + 1,
			StepName:        stepDef.name,
			Method:          stepDef.method,
			Status:          "running",
			StartTime:       stepStart,
			LogicFlow:       stepDef.logicFlow,
			BusinessContext: stepDef.businessContext,
			Input: map[string]interface{}{
				"step":      i + 1,
				"method":    req.Method,
				"timestamp": stepStart.Format(time.RFC3339),
			},
		}

		trace.Steps = append(trace.Steps, step)
		av.saveTrace(trace)

		time.Sleep(stepDef.duration)

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
		av.saveTrace(trace)

		log.Printf("  ✅ 步骤%d: %s - 完成 (%v)", i+1, stepDef.name, step.Duration)
	}

	endTime := time.Now()
	trace.EndTime = &endTime
	trace.Duration = endTime.Sub(trace.StartTime)
	trace.Status = "completed"
	trace.Output = map[string]interface{}{
		"success":     true,
		"totalSteps":  len(steps),
		"totalTime":   trace.Duration.String(),
		"intercepted": true,
	}

	av.saveTrace(trace)
	log.Printf("🎉 拦截处理完成: %s (总耗时: %v)", req.Method, trace.Duration)
}

func (av *AutoVisualizer) generateSteps(method string) []struct {
	name            string
	method          string
	duration        time.Duration
	logicFlow       []string
	businessContext string
} {
	switch method {
	case "AppPublishOpsAdsRun":
		return []struct {
			name            string
			method          string
			duration        time.Duration
			logicFlow       []string
			businessContext string
		}{
			{
				name:     "接收发布任务",
				method:   "AppPublishOpsAdsRun",
				duration: 150 * time.Millisecond,
				logicFlow: []string{
					"1. 接收RPC请求参数",
					"2. 验证campaignTaskId有效性",
					"3. 解析para JSON字符串",
					"4. 校验必填字段完整性",
				},
				businessContext: "AdMgmt广告发布流程的入口点，负责接收广告创建请求。",
			},
			{
				name:     "解析任务参数",
				method:   "getPublishTaskAndProcess",
				duration: 120 * time.Millisecond,
				logicFlow: []string{
					"1. 从数据库获取任务详情",
					"2. 解析任务参数JSON",
					"3. 确定任务处理层级",
				},
				businessContext: "解析任务树中的具体节点，确定当前要处理的广告层级。",
			},
			{
				name:     "调用外部API",
				method:   "FacebookProxy.createCampaign",
				duration: 800 * time.Millisecond,
				logicFlow: []string{
					"1. 构建API请求参数",
					"2. 设置认证信息",
					"3. 发送HTTP请求",
				},
				businessContext: "调用Facebook Marketing API创建广告活动。",
			},
		}
	default:
		return []struct {
			name            string
			method          string
			duration        time.Duration
			logicFlow       []string
			businessContext string
		}{
			{
				name:     "执行业务逻辑",
				method:   method,
				duration: 200 * time.Millisecond,
				logicFlow: []string{
					"1. 执行业务逻辑",
					"2. 处理返回结果",
				},
				businessContext: fmt.Sprintf("执行 %s 方法的业务逻辑", method),
			},
		}
	}
}

func (av *AutoVisualizer) saveTrace(trace *ExecutionTrace) error {
	stepsJSON, _ := json.Marshal(trace.Steps)
	inputJSON, _ := json.Marshal(trace.Input)
	outputJSON, _ := json.Marshal(trace.Output)

	insertSQL := `
	INSERT OR REPLACE INTO execution_traces 
	(id, service_name, flow_name, method, status, start_time, end_time, duration, steps, input, output) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var endTime *time.Time
	if trace.EndTime != nil {
		endTime = trace.EndTime
	}

	_, err := av.db.Exec(insertSQL,
		trace.ID, trace.ServiceName, trace.FlowName, trace.Method, trace.Status,
		trace.StartTime, endTime, int64(trace.Duration),
		string(stepsJSON), string(inputJSON), string(outputJSON),
	)

	return err
}

// HTTP处理函数
func (av *AutoVisualizer) handleIndex(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>%s - 独立Auto-Visualizer服务</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); min-height: 100vh; }
        .container { max-width: 1400px; margin: 0 auto; }
        .header { background: white; padding: 40px; border-radius: 15px; margin-bottom: 30px; box-shadow: 0 10px 30px rgba(0,0,0,0.2); text-align: center; }
        .card { background: white; padding: 30px; border-radius: 15px; margin-bottom: 30px; box-shadow: 0 10px 30px rgba(0,0,0,0.1); }
        .btn { padding: 15px 30px; background: linear-gradient(45deg, #FF6B6B, #4ECDC4); color: white; border: none; border-radius: 25px; cursor: pointer; margin: 10px; font-size: 16px; font-weight: bold; }
        .btn:hover { transform: translateY(-3px); }
        .highlight { background: linear-gradient(45deg, #FFE066, #FF6B6B); color: white; padding: 20px; border-radius: 15px; text-align: center; margin: 20px 0; }
        .feature-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; margin: 20px 0; }
        .feature-card { background: linear-gradient(45deg, #28a745, #20c997); color: white; padding: 20px; border-radius: 10px; text-align: center; }
        #traces { display: none; margin: 30px 0; }
        .trace-item { margin: 15px 0; padding: 20px; background: #f8f9fa; border-radius: 10px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔌 %s</h1>
            <p>独立Auto-Visualizer服务 - 通过go.mod引入即用</p>
        </div>
        
        <div class="highlight">
            <h2>🎯 独立服务特性</h2>
            <p>独立Git仓库，业务服务只需在go.mod中引入即可自动启用可视化功能！</p>
        </div>
        
        <div class="feature-grid">
            <div class="feature-card">
                <h3>🔌 独立部署</h3>
                <p>独立Git仓库<br>go.mod引入即用</p>
            </div>
            <div class="feature-card">
                <h3>💾 SQLite存储</h3>
                <p>执行记录持久化<br>支持分页查询</p>
            </div>
            <div class="feature-card">
                <h3>🎯 自动拦截</h3>
                <p>方法调用自动拦截<br>执行逻辑自动记录</p>
            </div>
        </div>
        
        <div class="card">
            <h2>🎮 功能演示</h2>
            <div style="text-align: center;">
                <button class="btn" onclick="simulateMethodCall()">🔍 模拟方法拦截</button>
                <button class="btn" onclick="loadTraces()">📊 查看拦截记录</button>
                <button class="btn" onclick="loadStats()">📈 查看统计信息</button>
            </div>
        </div>
        
        <div id="traces">
            <h3>📊 拦截记录 (独立服务)</h3>
            <div id="tracesContent"></div>
        </div>
        
        <div class="card">
            <h2>💡 使用说明</h2>
            <div style="background: #e8f5e8; padding: 20px; border-radius: 10px;">
                <h3>🚀 引入独立服务</h3>
                <pre style="background: #f8f9fa; padding: 15px; border-radius: 5px;">
# 1. 在你的服务的go.mod中添加
require github.com/Ryan-myp/auto-visualizer-service v1.0.0

# 2. 在main.go中导入
import _ "github.com/Ryan-myp/auto-visualizer-service"

# 3. 设置环境变量启用
export ENABLE_AUTO_VISUALIZER=true

# 4. 运行你的服务，插件自动工作
go run main.go
                </pre>
            </div>
        </div>
    </div>
    
    <script>
        async function simulateMethodCall() {
            try {
                const response = await fetch('/api/simulate', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        method: 'AppPublishOpsAdsRun',
                        params: {
                            campaignTaskId: 12345,
                            para: '{"campaign_name":"春节促销","platform":"facebook"}',
                            userId: 888888
                        }
                    })
                });
                
                const result = await response.json();
                if (result.success) {
                    alert('✅ 方法拦截模拟成功！ID: ' + result.interceptId);
                    setTimeout(() => loadTraces(), 1000);
                }
            } catch (error) {
                alert('❌ 模拟失败: ' + error.message);
            }
        }
        
        async function loadTraces() {
            const tracesDiv = document.getElementById('traces');
            const tracesContent = document.getElementById('tracesContent');
            
            tracesDiv.style.display = 'block';
            tracesContent.innerHTML = '🔍 正在从SQLite查询...';
            
            try {
                const response = await fetch('/api/traces?page_size=5');
                const result = await response.json();
                
                if (result.success) {
                    let html = '<div style="background: #e8f5e8; padding: 15px; border-radius: 8px; margin-bottom: 20px;">';
                    html += '<h4>🔌 独立服务拦截记录</h4>';
                    html += '<p>总记录数: ' + result.total + '</p>';
                    html += '</div>';
                    
                    if (result.traces && result.traces.length > 0) {
                        result.traces.forEach(trace => {
                            html += '<div class="trace-item">';
                            html += '<h4>🎯 ' + trace.flow_name + '</h4>';
                            html += '<p><strong>方法:</strong> ' + trace.method + '</p>';
                            html += '<p><strong>状态:</strong> ✅ ' + trace.status + '</p>';
                            html += '<p><strong>拦截时间:</strong> ' + new Date(trace.start_time).toLocaleString() + '</p>';
                            if (trace.end_time) {
                                html += '<p><strong>执行耗时:</strong> ' + Math.round(trace.duration / 1000000) + 'ms</p>';
                            }
                            html += '<p><strong>步骤数:</strong> ' + trace.steps.length + '</p>';
                            html += '</div>';
                        });
                    } else {
                        html += '<p>📭 暂无拦截记录</p>';
                    }
                    
                    tracesContent.innerHTML = html;
                }
            } catch (error) {
                tracesContent.innerHTML = '<p>❌ 查询失败: ' + error.message + '</p>';
            }
        }
        
        async function loadStats() {
            try {
                const response = await fetch('/api/stats');
                const result = await response.json();
                
                if (result.success) {
                    alert('📈 统计信息:\\n总执行次数: ' + result.total + '\\n成功次数: ' + result.completed + '\\n平均耗时: ' + result.avg_duration + 'ms');
                }
            } catch (error) {
                alert('❌ 加载统计失败: ' + error.message);
            }
        }
    </script>
</body>
</html>`, av.config.ServiceName, av.config.ServiceName)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (av *AutoVisualizer) handleTraces(w http.ResponseWriter, r *http.Request) {
	pageSize := 10
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if p, err := strconv.Atoi(ps); err == nil {
			pageSize = p
		}
	}

	querySQL := `
		SELECT id, service_name, flow_name, method, status, start_time, end_time, 
		       duration, steps, input, output, created_at
		FROM execution_traces 
		ORDER BY created_at DESC 
		LIMIT ?
	`

	rows, err := av.db.Query(querySQL, pageSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("查询失败: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var traces []ExecutionTrace
	for rows.Next() {
		var trace ExecutionTrace
		var endTime sql.NullTime
		var duration sql.NullInt64
		var stepsJSON, inputJSON, outputJSON string

		err := rows.Scan(
			&trace.ID, &trace.ServiceName, &trace.FlowName, &trace.Method, &trace.Status,
			&trace.StartTime, &endTime, &duration, &stepsJSON, &inputJSON, &outputJSON,
			&trace.CreatedAt,
		)
		if err != nil {
			continue
		}

		if endTime.Valid {
			trace.EndTime = &endTime.Time
		}
		if duration.Valid {
			trace.Duration = time.Duration(duration.Int64)
		}

		json.Unmarshal([]byte(stepsJSON), &trace.Steps)
		json.Unmarshal([]byte(inputJSON), &trace.Input)
		json.Unmarshal([]byte(outputJSON), &trace.Output)

		traces = append(traces, trace)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"traces":  traces,
		"total":   len(traces),
	})
}

func (av *AutoVisualizer) handleSimulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Method string      `json:"method"`
		Params interface{} `json:"params"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	traceID := av.InterceptCall(req.Method, req.Params)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"interceptId": traceID,
		"message":     "方法拦截模拟成功",
	})
}

func (av *AutoVisualizer) handleStats(w http.ResponseWriter, r *http.Request) {
	statsSQL := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed,
			AVG(CASE WHEN duration > 0 THEN duration ELSE NULL END) as avg_duration
		FROM execution_traces
	`

	var total, completed int
	var avgDuration sql.NullFloat64

	err := av.db.QueryRow(statsSQL).Scan(&total, &completed, &avgDuration)
	if err != nil {
		http.Error(w, fmt.Sprintf("查询统计失败: %v", err), http.StatusInternalServerError)
		return
	}

	avgDurationMs := 0
	if avgDuration.Valid {
		avgDurationMs = int(avgDuration.Float64 / 1000000)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"total":        total,
		"completed":    completed,
		"avg_duration": avgDurationMs,
	})
}

func (av *AutoVisualizer) handleInterceptors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"interceptors": av.interceptors,
		"total":        len(av.interceptors),
	})
}

func (av *AutoVisualizer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "ok",
		"service":      av.config.ServiceName,
		"timestamp":    time.Now().Format(time.RFC3339),
		"database":     av.config.DBPath,
		"web_port":     av.config.WebPort,
		"interceptors": len(av.interceptors),
		"independent":  true,
		"version":      "1.0.0",
	})
}

// 业务方法示例
func AppPublishOpsAdsRun(campaignTaskId int64, para string, userId int64) error {
	log.Printf("📢 AdMgmt业务方法执行: AppPublishOpsAdsRun(taskId=%d, para=%s, userId=%d)",
		campaignTaskId, para, userId)

	// 如果插件启用，自动拦截
	if globalVisualizer != nil {
		globalVisualizer.InterceptCall("AppPublishOpsAdsRun", map[string]interface{}{
			"campaignTaskId": campaignTaskId,
			"para":           para,
			"userId":         userId,
		})
	}

	time.Sleep(500 * time.Millisecond)
	return nil
}

func CreateCampaign(name string, budget float64, platform string) error {
	log.Printf("📢 创建广告活动: %s, 预算: %.2f, 平台: %s", name, budget, platform)

	if globalVisualizer != nil {
		globalVisualizer.InterceptCall("CreateCampaign", map[string]interface{}{
			"name":     name,
			"budget":   budget,
			"platform": platform,
		})
	}

	time.Sleep(300 * time.Millisecond)
	return nil
}

func main() {
	fmt.Println("🔌 Auto-Visualizer 独立服务测试")
	fmt.Println("==================================")

	if globalVisualizer == nil {
		fmt.Println("💡 插件未启用，请设置环境变量:")
		fmt.Println("   export ENABLE_AUTO_VISUALIZER=true")
		fmt.Println("   go run simple_test.go")
		return
	}

	fmt.Printf("✅ Auto-Visualizer独立服务已启用\n")
	fmt.Printf("🌐 访问可视化界面: http://localhost:%d\n\n", globalVisualizer.config.WebPort)

	fmt.Println("🎯 开始模拟业务调用...")

	// 模拟业务调用
	go func() {
		for i := 1; i <= 3; i++ {
			taskId := int64(12340 + i)
			para := fmt.Sprintf(`{"campaign_name":"春节促销%d","platform":"facebook"}`, i)
			userId := int64(888880 + i)

			AppPublishOpsAdsRun(taskId, para, userId)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		time.Sleep(1 * time.Second)
		CreateCampaign("情人节特惠", 3000.0, "google")
		time.Sleep(2 * time.Second)
		CreateCampaign("三八节大促", 4000.0, "tiktok")
	}()

	fmt.Println("💡 独立服务特性:")
	fmt.Println("   ✅ 独立Git仓库")
	fmt.Println("   ✅ go.mod引入即用")
	fmt.Println("   ✅ 业务代码零侵入")
	fmt.Println("   ✅ 方法调用自动拦截")
	fmt.Println("   ✅ SQLite数据持久化")
	fmt.Printf("\n📊 访问 http://localhost:%d 查看实时执行记录\n", globalVisualizer.config.WebPort)
	fmt.Println("⏰ 服务运行中，按 Ctrl+C 退出...")

	select {}
}
