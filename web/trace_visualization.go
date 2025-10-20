package web

import (
	"fmt"
	"net/http"

	"github.com/Ryan-myp/auto-visualizer-service/tracer"
	"github.com/gin-gonic/gin"
)

// handleTraceVisualization 可视化展示追踪详情
func (s *Server) handleTraceVisualization(c *gin.Context) {
	traceID := c.Param("id")
	
	t := tracer.GetTracer()
	trace := t.GetTrace(traceID)
	
	if trace == nil {
		c.HTML(http.StatusNotFound, "", gin.H{})
		c.String(http.StatusNotFound, "追踪记录不存在")
		return
	}
	
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>追踪详情 - %s</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            min-height: 100vh;
            padding: 20px;
        }
        .container { max-width: 1400px; margin: 0 auto; }
        
        /* 头部 */
        .header {
            background: white;
            padding: 30px;
            border-radius: 15px;
            margin-bottom: 20px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
        }
        .header h1 { 
            font-size: 28px;
            color: #333;
            margin-bottom: 10px;
        }
        .header .meta {
            display: flex;
            gap: 20px;
            flex-wrap: wrap;
            margin-top: 15px;
        }
        .meta-item {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 8px 15px;
            background: #f8f9fa;
            border-radius: 8px;
            font-size: 14px;
        }
        .meta-item .label { color: #666; }
        .meta-item .value { 
            font-weight: 600;
            color: #333;
        }
        
        /* 状态标签 */
        .status-badge {
            display: inline-block;
            padding: 6px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 600;
            text-transform: uppercase;
        }
        .status-success { background: #d4edda; color: #155724; }
        .status-error { background: #f8d7da; color: #721c24; }
        .status-running { background: #fff3cd; color: #856404; }
        
        /* 主内容区 */
        .content {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
            margin-bottom: 20px;
        }
        
        .card {
            background: white;
            padding: 25px;
            border-radius: 15px;
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        .card h2 {
            font-size: 20px;
            color: #333;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid #667eea;
        }
        
        /* 时间线 */
        .timeline {
            position: relative;
            padding-left: 30px;
        }
        .timeline::before {
            content: '';
            position: absolute;
            left: 10px;
            top: 0;
            bottom: 0;
            width: 2px;
            background: linear-gradient(to bottom, #667eea, #764ba2);
        }
        .timeline-item {
            position: relative;
            margin-bottom: 20px;
            padding-left: 20px;
        }
        .timeline-item::before {
            content: '';
            position: absolute;
            left: -24px;
            top: 5px;
            width: 12px;
            height: 12px;
            border-radius: 50%%;
            background: #667eea;
            border: 3px solid white;
            box-shadow: 0 0 0 2px #667eea;
        }
        .timeline-item .time {
            font-size: 12px;
            color: #999;
            margin-bottom: 5px;
        }
        .timeline-item .event {
            font-weight: 600;
            color: #333;
            margin-bottom: 5px;
        }
        .timeline-item .detail {
            font-size: 14px;
            color: #666;
        }
        
        /* 参数展示 */
        .param-box {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 15px;
            border-left: 4px solid #667eea;
        }
        .param-box .title {
            font-weight: 600;
            color: #333;
            margin-bottom: 10px;
            font-size: 14px;
        }
        .param-box pre {
            background: #2d2d2d;
            color: #f8f8f2;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
            font-size: 13px;
            line-height: 1.5;
        }
        
        /* 调用栈 */
        .call-stack {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 8px;
        }
        .stack-item {
            padding: 10px;
            margin-bottom: 8px;
            background: white;
            border-radius: 5px;
            font-size: 13px;
            font-family: 'Courier New', monospace;
            border-left: 3px solid #667eea;
        }
        
        /* 性能指标 */
        .metrics {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
        }
        .metric-card {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 20px;
            border-radius: 10px;
            text-align: center;
        }
        .metric-card .value {
            font-size: 32px;
            font-weight: bold;
            margin-bottom: 5px;
        }
        .metric-card .label {
            font-size: 14px;
            opacity: 0.9;
        }
        
        /* 调用树 */
        .call-tree {
            margin-top: 20px;
        }
        .tree-node {
            margin-left: 20px;
            padding: 10px;
            border-left: 2px solid #ddd;
            margin-bottom: 10px;
        }
        .tree-node .node-header {
            display: flex;
            align-items: center;
            gap: 10px;
            padding: 10px;
            background: #f8f9fa;
            border-radius: 5px;
            cursor: pointer;
        }
        .tree-node .node-header:hover {
            background: #e9ecef;
        }
        .tree-node .node-method {
            font-weight: 600;
            color: #667eea;
        }
        .tree-node .node-duration {
            color: #666;
            font-size: 13px;
        }
        .tree-node .node-status {
            margin-left: auto;
        }
        
        /* 按钮 */
        .btn {
            display: inline-block;
            padding: 10px 20px;
            background: linear-gradient(45deg, #667eea, #764ba2);
            color: white;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 600;
            transition: transform 0.2s;
        }
        .btn:hover {
            transform: translateY(-2px);
        }
        
        /* 响应式 */
        @media (max-width: 768px) {
            .content {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <!-- 头部 -->
        <div class="header">
            <a href="/" class="btn">← 返回首页</a>
            <h1>🔍 追踪详情</h1>
            <div class="meta">
                <div class="meta-item">
                    <span class="label">方法:</span>
                    <span class="value">%s</span>
                </div>
                <div class="meta-item">
                    <span class="label">包:</span>
                    <span class="value">%s</span>
                </div>
                <div class="meta-item">
                    <span class="label">状态:</span>
                    <span class="status-badge status-%s">%s</span>
                </div>
                <div class="meta-item">
                    <span class="label">耗时:</span>
                    <span class="value">%s</span>
                </div>
                <div class="meta-item">
                    <span class="label">Goroutine:</span>
                    <span class="value">#%d</span>
                </div>
            </div>
        </div>
        
        <!-- 主内容 -->
        <div class="content">
            <!-- 左侧：时间线 -->
            <div class="card">
                <h2>📅 执行时间线</h2>
                <div class="timeline">
                    <div class="timeline-item">
                        <div class="time">开始时间</div>
                        <div class="event">方法调用开始</div>
                        <div class="detail">%s</div>
                    </div>
                    <div class="timeline-item">
                        <div class="time">结束时间</div>
                        <div class="event">方法执行完成</div>
                        <div class="detail">%s</div>
                    </div>
                    <div class="timeline-item">
                        <div class="time">总耗时</div>
                        <div class="event">%s</div>
                        <div class="detail">状态: %s</div>
                    </div>
                </div>
            </div>
            
            <!-- 右侧：性能指标 -->
            <div class="card">
                <h2>📊 性能指标</h2>
                <div class="metrics">
                    <div class="metric-card">
                        <div class="value">%s</div>
                        <div class="label">执行时间</div>
                    </div>
                    <div class="metric-card">
                        <div class="value">%s</div>
                        <div class="label">状态</div>
                    </div>
                    <div class="metric-card">
                        <div class="value">%d</div>
                        <div class="label">子调用数</div>
                    </div>
                </div>
            </div>
        </div>
        
        <!-- 参数和返回值 -->
        <div class="content">
            <div class="card">
                <h2>📥 输入参数</h2>
                <div class="param-box">
                    <div class="title">参数列表</div>
                    <pre>%s</pre>
                </div>
            </div>
            
            <div class="card">
                <h2>📤 返回值</h2>
                <div class="param-box">
                    <div class="title">返回结果</div>
                    <pre>%s</pre>
                </div>
                %s
            </div>
        </div>
        
        <!-- 调用栈 -->
        <div class="card">
            <h2>📚 调用栈</h2>
            <div class="call-stack">
                %s
            </div>
        </div>
        
        <!-- 子调用 -->
        %s
    </div>
</body>
</html>
`, 
		trace.MethodName,
		trace.MethodName,
		trace.PackageName,
		trace.Status,
		trace.Status,
		trace.Duration.String(),
		trace.Goroutine,
		trace.StartTime.Format("2006-01-02 15:04:05.000"),
		trace.EndTime.Format("2006-01-02 15:04:05.000"),
		trace.Duration.String(),
		trace.Status,
		trace.Duration.String(),
		trace.Status,
		len(trace.Children),
		formatJSON(trace.Input),
		formatJSON(trace.Output),
		formatError(trace.Error),
		formatCallStack(trace.CallStack),
		formatChildren(trace.Children),
	)
	
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// formatJSON 格式化 JSON
func formatJSON(data interface{}) string {
	if data == nil {
		return "null"
	}
	
	// 简单格式化
	return fmt.Sprintf("%+v", data)
}

// formatError 格式化错误信息
func formatError(err string) string {
	if err == "" {
		return ""
	}
	
	return fmt.Sprintf(`
		<div class="param-box" style="border-left-color: #dc3545;">
			<div class="title" style="color: #dc3545;">❌ 错误信息</div>
			<pre style="background: #f8d7da; color: #721c24;">%s</pre>
		</div>
	`, err)
}

// formatCallStack 格式化调用栈
func formatCallStack(stack []string) string {
	if len(stack) == 0 {
		return "<p style='color: #999;'>无调用栈信息</p>"
	}
	
	html := ""
	for i, item := range stack {
		html += fmt.Sprintf(`
			<div class="stack-item">
				<strong>#%d</strong> %s
			</div>
		`, i+1, item)
	}
	
	return html
}

// formatChildren 格式化子调用
func formatChildren(children []*tracer.MethodTrace) string {
	if len(children) == 0 {
		return ""
	}
	
	html := `
		<div class="card">
			<h2>🌲 子调用链</h2>
			<div class="call-tree">
	`
	
	for _, child := range children {
		html += formatTreeNode(child, 0)
	}
	
	html += `
			</div>
		</div>
	`
	
	return html
}

// formatTreeNode 格式化树节点
func formatTreeNode(node *tracer.MethodTrace, level int) string {
	statusClass := "status-" + node.Status
	
	html := fmt.Sprintf(`
		<div class="tree-node" style="margin-left: %dpx;">
			<div class="node-header">
				<span class="node-method">%s</span>
				<span class="node-duration">%s</span>
				<span class="node-status">
					<span class="status-badge %s">%s</span>
				</span>
			</div>
	`, level*20, node.MethodName, node.Duration.String(), statusClass, node.Status)
	
	// 递归处理子节点
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			html += formatTreeNode(child, level+1)
		}
	}
	
	html += `</div>`
	
	return html
}

