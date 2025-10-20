package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ryan-myp/auto-visualizer-service/config"
	"github.com/Ryan-myp/auto-visualizer-service/interceptor"
	"github.com/Ryan-myp/auto-visualizer-service/storage"
	"github.com/Ryan-myp/auto-visualizer-service/tracer"
	"github.com/gin-gonic/gin"
)

// Server Web服务器
type Server struct {
	config      *config.Config
	interceptor *interceptor.Manager
	engine      *gin.Engine
}

// NewServer 创建Web服务器
func NewServer(cfg *config.Config, interceptorMgr *interceptor.Manager) (*Server, error) {
	// 设置Gin模式
	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	server := &Server{
		config:      cfg,
		interceptor: interceptorMgr,
		engine:      engine,
	}

	server.setupRoutes()
	return server, nil
}

// Start 启动Web服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.WebPort)

	go func() {
		log.Printf("🌐 Web服务器启动中: http://localhost:%d", s.config.WebPort)
		if err := s.engine.Run(addr); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ Web服务器启动失败: %v", err)
		}
	}()

	// 等待服务器启动（简单的延迟）
	time.Sleep(100 * time.Millisecond)
	log.Printf("✅ Web服务器已启动: http://localhost:%d", s.config.WebPort)

	return nil
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 主页
	s.engine.GET("/", s.handleIndex)

	// API路由
	api := s.engine.Group("/api")
	{
		api.GET("/traces", s.handleGetTraces)
		api.GET("/traces/:id", s.handleGetTraceDetail)
		api.POST("/simulate", s.handleSimulate)
		api.GET("/stats", s.handleGetStats)
		api.GET("/interceptors", s.handleGetInterceptors)
		api.DELETE("/traces/cleanup", s.handleCleanupTraces)

		// 新增：方法追踪API
		api.GET("/method-traces", s.handleGetMethodTraces)
		api.GET("/method-traces/:id", s.handleGetMethodTraceDetail)
		api.DELETE("/method-traces", s.handleClearMethodTraces)
		api.GET("/method-traces/tree", s.handleGetMethodTraceTree)
	}

	// 健康检查
	s.engine.GET("/health", s.handleHealth)
}

// handleIndex 主页
func (s *Server) handleIndex(c *gin.Context) {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>%s - Auto-Visualizer</title>
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
        .pagination { text-align: center; margin: 20px 0; }
        .pagination button { padding: 8px 16px; margin: 0 5px; border: none; border-radius: 5px; cursor: pointer; }
        .pagination .active { background: #007bff; color: white; }
        .filter-bar { background: #f8f9fa; padding: 15px; border-radius: 10px; margin-bottom: 20px; }
        .filter-bar select, .filter-bar input { padding: 8px 12px; margin: 0 10px; border: 1px solid #ddd; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔌 %s</h1>
            <p>Auto-Visualizer 独立服务 - 非侵入式业务流程可视化</p>
        </div>
        
        <div class="highlight">
            <h2>🎯 独立服务特性</h2>
            <p>通过go.mod引入即可自动启用，业务代码零侵入，执行逻辑自动记录！</p>
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
                <button class="btn" onclick="loadTraces(1)">📊 查看拦截记录</button>
                <button class="btn" onclick="loadStats()">📈 查看统计信息</button>
                <button class="btn" onclick="loadInterceptors()">🔧 查看拦截器</button>
            </div>
        </div>
        
        <div id="traces">
            <h3>📊 拦截记录 (SQLite存储)</h3>
            <div class="filter-bar">
                <label>状态筛选:</label>
                <select id="statusFilter" onchange="loadTraces(1)">
                    <option value="">全部</option>
                    <option value="running">执行中</option>
                    <option value="completed">已完成</option>
                    <option value="failed">失败</option>
                </select>
                <label>方法筛选:</label>
                <select id="methodFilter" onchange="loadTraces(1)">
                    <option value="">全部方法</option>
                </select>
                <label>每页显示:</label>
                <select id="pageSizeSelect" onchange="loadTraces(1)">
                    <option value="5">5条</option>
                    <option value="10" selected>10条</option>
                    <option value="20">20条</option>
                </select>
            </div>
            <div id="tracesContent"></div>
            <div id="pagination" class="pagination"></div>
        </div>
        
        <div class="card">
            <h2>💡 使用说明</h2>
            <div style="background: #e8f5e8; padding: 20px; border-radius: 10px;">
                <h3>🚀 引入服务</h3>
                <pre style="background: #f8f9fa; padding: 15px; border-radius: 5px;">
# 在你的服务的go.mod中添加
require github.com/Ryan-myp/auto-visualizer-service v1.0.0

# 在main.go中导入
import _ "github.com/Ryan-myp/auto-visualizer-service"

# 设置环境变量启用
export ENABLE_AUTO_VISUALIZER=true

# 运行你的服务，插件自动工作
go run main.go
                </pre>
            </div>
        </div>
    </div>
    
    <script>
        let currentPage = 1;
        
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
                    setTimeout(() => loadTraces(1), 1000);
                }
            } catch (error) {
                alert('❌ 模拟失败: ' + error.message);
            }
        }
        
        async function loadTraces(page) {
            const tracesDiv = document.getElementById('traces');
            const tracesContent = document.getElementById('tracesContent');
            
            tracesDiv.style.display = 'block';
            tracesContent.innerHTML = '🔍 正在从SQLite查询拦截记录...';
            
            const status = document.getElementById('statusFilter').value;
            const method = document.getElementById('methodFilter').value;
            const pageSize = document.getElementById('pageSizeSelect').value;
            
            try {
                let url = '/api/traces?page=' + page + '&page_size=' + pageSize;
                if (status) url += '&status=' + status;
                if (method) url += '&method=' + method;
                
                const response = await fetch(url);
                const result = await response.json();
                
                if (result.success) {
                    currentPage = result.page;
                    
                    let html = '<div style="background: #e8f5e8; padding: 15px; border-radius: 8px; margin-bottom: 20px;">';
                    html += '<h4>🔌 独立服务拦截记录</h4>';
                    html += '<p>总记录数: ' + result.total + ' | 当前页: ' + result.page + '/' + result.total_pages + '</p>';
                    html += '</div>';
                    
                    if (result.traces && result.traces.length > 0) {
                        result.traces.forEach(trace => {
                            html += '<div class="trace-item">';
                            html += '<h4>🎯 ' + trace.flow_name + '</h4>';
                            html += '<p><strong>方法:</strong> ' + trace.method + '</p>';
                            html += '<p><strong>状态:</strong> ' + getStatusIcon(trace.status) + ' ' + trace.status + '</p>';
                            html += '<p><strong>拦截时间:</strong> ' + new Date(trace.start_time).toLocaleString() + '</p>';
                            if (trace.end_time) {
                                html += '<p><strong>执行耗时:</strong> ' + Math.round(trace.duration / 1000000) + 'ms</p>';
                            }
                            html += '<p><strong>步骤数:</strong> ' + trace.steps.length + '</p>';
                            html += '<button onclick="viewTraceDetail(\'' + trace.id + '\')" style="padding: 8px 16px; background: #007bff; color: white; border: none; border-radius: 15px; cursor: pointer;">查看详情</button>';
                            html += '</div>';
                        });
                        
                        // 分页控件
                        let paginationHtml = '';
                        if (result.page > 1) {
                            paginationHtml += '<button onclick="loadTraces(' + (result.page - 1) + ')">上一页</button>';
                        }
                        for (let i = Math.max(1, result.page - 2); i <= Math.min(result.total_pages, result.page + 2); i++) {
                            const className = i === result.page ? 'active' : '';
                            paginationHtml += '<button class="' + className + '" onclick="loadTraces(' + i + ')">' + i + '</button>';
                        }
                        if (result.page < result.total_pages) {
                            paginationHtml += '<button onclick="loadTraces(' + (result.page + 1) + ')">下一页</button>';
                        }
                        
                        document.getElementById('pagination').innerHTML = paginationHtml;
                    } else {
                        html += '<p>📭 暂无拦截记录</p>';
                    }
                    
                    tracesContent.innerHTML = html;
                }
            } catch (error) {
                tracesContent.innerHTML = '<p>❌ 查询失败: ' + error.message + '</p>';
            }
        }
        
        async function loadInterceptors() {
            try {
                const response = await fetch('/api/interceptors');
                const result = await response.json();
                
                if (result.success) {
                    const methodFilter = document.getElementById('methodFilter');
                    methodFilter.innerHTML = '<option value="">全部方法</option>';
                    
                    Object.keys(result.interceptors).forEach(method => {
                        const option = document.createElement('option');
                        option.value = method;
                        option.textContent = method;
                        methodFilter.appendChild(option);
                    });
                    
                    alert('✅ 已加载 ' + Object.keys(result.interceptors).length + ' 个拦截器');
                }
            } catch (error) {
                alert('❌ 加载拦截器失败: ' + error.message);
            }
        }
        
        async function viewTraceDetail(traceId) {
            try {
                const response = await fetch('/api/traces/' + traceId);
                const result = await response.json();
                
                if (result.success && result.trace) {
                    const trace = result.trace;
                    let detailHtml = '<h2>🔍 拦截详情 (独立服务)</h2>';
                    detailHtml += '<p><strong>ID:</strong> ' + trace.id + '</p>';
                    detailHtml += '<p><strong>服务:</strong> ' + trace.service_name + '</p>';
                    detailHtml += '<p><strong>方法:</strong> ' + trace.method + '</p>';
                    detailHtml += '<p><strong>流程:</strong> ' + trace.flow_name + '</p>';
                    
                    trace.steps.forEach(step => {
                        detailHtml += '<div style="margin: 20px 0; padding: 20px; background: #f5f5f5; border-radius: 10px;">';
                        detailHtml += '<h4>步骤' + step.step_number + ': ' + step.step_name + '</h4>';
                        detailHtml += '<p><strong>业务含义:</strong> ' + step.business_context + '</p>';
                        detailHtml += '<h5>🔄 执行逻辑:</h5><ul>';
                        step.logic_flow.forEach(logic => {
                            detailHtml += '<li>' + logic + '</li>';
                        });
                        detailHtml += '</ul></div>';
                    });
                    
                    const newWindow = window.open('', '_blank', 'width=900,height=700,scrollbars=yes');
                    newWindow.document.write('<html><head><title>独立服务拦截详情</title><style>body{font-family:Arial,sans-serif;margin:20px;}</style></head><body>' + detailHtml + '</body></html>');
                }
            } catch (error) {
                alert('查看详情失败: ' + error.message);
            }
        }
        
        async function loadStats() {
            try {
                const response = await fetch('/api/stats');
                const result = await response.json();
                
                if (result.success) {
                    let statsHtml = '📈 统计信息:\\n';
                    statsHtml += '总执行次数: ' + result.total + '\\n';
                    statsHtml += '成功次数: ' + result.completed + '\\n';
                    statsHtml += '失败次数: ' + result.failed + '\\n';
                    statsHtml += '成功率: ' + result.success_rate.toFixed(2) + '%\\n';
                    statsHtml += '平均耗时: ' + result.avg_duration_ms + 'ms';
                    
                    alert(statsHtml);
                }
            } catch (error) {
                alert('❌ 加载统计失败: ' + error.message);
            }
        }
        
        function getStatusIcon(status) {
            switch (status) {
                case 'running': return '🔄';
                case 'completed': return '✅';
                case 'failed': return '❌';
                default: return '⏸️';
            }
        }
        
        // 页面加载时自动加载数据
        window.onload = function() {
            loadInterceptors();
            setTimeout(() => loadTraces(1), 1000);
        };
    </script>
</body>
</html>`, s.config.ServiceName, s.config.ServiceName)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// handleGetTraces 获取拦截记录
func (s *Server) handleGetTraces(c *gin.Context) {
	// 解析查询参数
	options := &storage.QueryOptions{
		Page:        1,
		PageSize:    10,
		Status:      c.Query("status"),
		Method:      c.Query("method"),
		ServiceName: c.Query("service_name"),
		UserID:      c.Query("user_id"),
		OrderBy:     c.Query("order_by"),
		OrderDesc:   c.Query("order_desc") == "true",
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			options.Page = p
		}
	}
	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			options.PageSize = ps
		}
	}

	// 从存储查询
	storage := s.interceptor.GetStorage()
	result, err := storage.GetTraces(options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"traces":      result.Traces,
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// handleGetTraceDetail 获取拦截详情
func (s *Server) handleGetTraceDetail(c *gin.Context) {
	traceID := c.Param("id")

	storage := s.interceptor.GetStorage()
	trace, err := storage.GetTraceByID(traceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"trace":   trace,
	})
}

// handleSimulate 模拟方法拦截
func (s *Server) handleSimulate(c *gin.Context) {
	var req struct {
		Method string      `json:"method"`
		Params interface{} `json:"params"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 触发方法拦截
	traceID := s.interceptor.InterceptCall(req.Method, req.Params)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"interceptId": traceID,
		"message":     "方法拦截模拟成功",
	})
}

// handleGetStats 获取统计信息
func (s *Server) handleGetStats(c *gin.Context) {
	serviceName := c.Query("service_name")
	method := c.Query("method")
	days := 30
	if d := c.Query("days"); d != "" {
		if parsedDays, err := strconv.Atoi(d); err == nil {
			days = parsedDays
		}
	}

	storage := s.interceptor.GetStorage()
	stats, err := storage.GetStatistics(serviceName, method, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// handleGetInterceptors 获取拦截器列表（动态从追踪数据中获取）
func (s *Server) handleGetInterceptors(c *gin.Context) {
	// 从追踪器中获取实际被追踪的方法
	t := tracer.GetTracer()
	traces := t.GetAllTraces()
	
	// 统计每个方法的调用次数和状态
	methodStats := make(map[string]map[string]interface{})
	for _, trace := range traces {
		if _, exists := methodStats[trace.MethodName]; !exists {
			methodStats[trace.MethodName] = map[string]interface{}{
				"name":          trace.MethodName,
				"package":       trace.PackageName,
				"call_count":    0,
				"success_count": 0,
				"error_count":   0,
				"total_duration": int64(0),
			}
		}
		
		stats := methodStats[trace.MethodName]
		stats["call_count"] = stats["call_count"].(int) + 1
		
		if trace.Status == "success" {
			stats["success_count"] = stats["success_count"].(int) + 1
		} else if trace.Status == "error" {
			stats["error_count"] = stats["error_count"].(int) + 1
		}
		
		stats["total_duration"] = stats["total_duration"].(int64) + trace.Duration.Milliseconds()
	}
	
	// 转换为数组并计算平均耗时
	result := make([]map[string]interface{}, 0, len(methodStats))
	for _, stats := range methodStats {
		callCount := stats["call_count"].(int)
		if callCount > 0 {
			avgDuration := stats["total_duration"].(int64) / int64(callCount)
			stats["avg_duration_ms"] = avgDuration
		}
		delete(stats, "total_duration") // 移除临时字段
		result = append(result, stats)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"methods": result,
		"total":   len(result),
	})
}

// handleCleanupTraces 清理过期记录
func (s *Server) handleCleanupTraces(c *gin.Context) {
	storage := s.interceptor.GetStorage()
	if err := storage.CleanupOldRecords(s.config.RetentionDays); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "清理完成",
	})
}

// handleHealth 健康检查
func (s *Server) handleHealth(c *gin.Context) {
	interceptors := s.interceptor.GetInterceptors()

	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"service":      s.config.ServiceName,
		"timestamp":    time.Now().Format(time.RFC3339),
		"database":     s.config.DBPath,
		"web_port":     s.config.WebPort,
		"interceptors": len(interceptors),
		"independent":  true,
		"version":      "1.0.0",
	})
}

// ============ 新增：方法追踪处理函数 ============

// handleGetMethodTraces 获取所有方法追踪
func (s *Server) handleGetMethodTraces(c *gin.Context) {
	t := tracer.GetTracer()
	traces := t.GetAllTraces()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"total":   len(traces),
		"traces":  traces,
	})
}

// handleGetMethodTraceDetail 获取方法追踪详情
func (s *Server) handleGetMethodTraceDetail(c *gin.Context) {
	traceID := c.Param("id")
	t := tracer.GetTracer()
	trace := t.GetTrace(traceID)

	if trace == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "追踪记录不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"trace":   trace,
	})
}

// handleClearMethodTraces 清除所有方法追踪
func (s *Server) handleClearMethodTraces(c *gin.Context) {
	t := tracer.GetTracer()
	t.ClearTraces()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "已清除所有追踪记录",
	})
}

// handleGetMethodTraceTree 获取方法调用树
func (s *Server) handleGetMethodTraceTree(c *gin.Context) {
	t := tracer.GetTracer()
	traces := t.GetAllTraces()

	// 构建树形结构
	tree := buildTraceTree(traces)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tree":    tree,
	})
}

// buildTraceTree 构建追踪树
func buildTraceTree(traces []*tracer.MethodTrace) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(traces))

	for _, trace := range traces {
		node := map[string]interface{}{
			"id":         trace.TraceID,
			"method":     trace.MethodName,
			"package":    trace.PackageName,
			"status":     trace.Status,
			"duration":   trace.Duration.String(),
			"start_time": trace.StartTime.Format(time.RFC3339),
			"input":      trace.Input,
			"output":     trace.Output,
			"error":      trace.Error,
			"goroutine":  trace.Goroutine,
			"call_stack": trace.CallStack,
			"children":   buildTraceChildren(trace.Children),
		}
		result = append(result, node)
	}

	return result
}

// buildTraceChildren 构建子追踪
func buildTraceChildren(children []*tracer.MethodTrace) []map[string]interface{} {
	if len(children) == 0 {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0, len(children))
	for _, child := range children {
		node := map[string]interface{}{
			"id":         child.TraceID,
			"method":     child.MethodName,
			"package":    child.PackageName,
			"status":     child.Status,
			"duration":   child.Duration.String(),
			"start_time": child.StartTime.Format(time.RFC3339),
			"input":      child.Input,
			"output":     child.Output,
			"error":      child.Error,
			"goroutine":  child.Goroutine,
			"children":   buildTraceChildren(child.Children),
		}
		result = append(result, node)
	}

	return result
}
