package web

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"

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
	
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>追踪详情 - %s</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            min-height: 100vh;
            padding: 20px;
            color: #333;
        }
        .container { max-width: 1600px; margin: 0 auto; }
        
        /* 顶部导航 */
        .top-nav {
            background: rgba(255,255,255,0.95);
            padding: 15px 30px;
            border-radius: 15px;
            margin-bottom: 20px;
            box-shadow: 0 5px 20px rgba(0,0,0,0.1);
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .btn-back {
            padding: 10px 20px;
            background: linear-gradient(45deg, #667eea, #764ba2);
            color: white;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 600;
            transition: transform 0.2s;
        }
        .btn-back:hover { transform: translateY(-2px); }
        
        /* 头部卡片 */
        .header-card {
            background: white;
            padding: 30px;
            border-radius: 15px;
            margin-bottom: 20px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.15);
        }
        .method-title {
            font-size: 32px;
            font-weight: 700;
            color: #667eea;
            margin-bottom: 15px;
            display: flex;
            align-items: center;
            gap: 15px;
        }
        .status-badge {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            padding: 8px 16px;
            border-radius: 20px;
            font-size: 14px;
            font-weight: 600;
        }
        .status-success { background: #d4edda; color: #155724; }
        .status-error { background: #f8d7da; color: #721c24; }
        
        .meta-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-top: 20px;
        }
        .meta-item {
            padding: 15px;
            background: #f8f9fa;
            border-radius: 10px;
            border-left: 4px solid #667eea;
        }
        .meta-label {
            font-size: 12px;
            color: #666;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 5px;
        }
        .meta-value {
            font-size: 18px;
            font-weight: 600;
            color: #333;
        }
        
        /* 主内容区域 */
        .content-grid {
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
        .card-title {
            font-size: 20px;
            font-weight: 700;
            color: #333;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 3px solid #667eea;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        /* 参数展示 */
        .param-container {
            background: #f8f9fa;
            border-radius: 10px;
            overflow: hidden;
            border: 1px solid #e9ecef;
        }
        .param-header {
            background: linear-gradient(135deg, #667eea, #764ba2);
            color: white;
            padding: 12px 20px;
            font-weight: 600;
            font-size: 14px;
        }
        .param-content {
            padding: 20px;
        }
        .param-content pre {
            background: #2d2d2d;
            color: #f8f8f2;
            padding: 20px;
            border-radius: 8px;
            overflow-x: auto;
            font-size: 14px;
            line-height: 1.6;
            font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
        }
        .param-empty {
            color: #999;
            font-style: italic;
            text-align: center;
            padding: 20px;
        }
        
        /* JSON 语法高亮 */
        .json-key { color: #e06c75; }
        .json-string { color: #98c379; }
        .json-number { color: #d19a66; }
        .json-boolean { color: #56b6c2; }
        .json-null { color: #c678dd; }
        
        /* 调用树 */
        .call-tree {
            margin-top: 10px;
        }
        .tree-node {
            margin-bottom: 10px;
            position: relative;
        }
        .tree-node-content {
            display: flex;
            align-items: center;
            padding: 15px;
            background: #f8f9fa;
            border-radius: 10px;
            border-left: 4px solid #667eea;
            transition: all 0.3s;
            cursor: pointer;
        }
        .tree-node-content:hover {
            background: #e9ecef;
            transform: translateX(5px);
        }
        .tree-node-icon {
            font-size: 20px;
            margin-right: 12px;
        }
        .tree-node-method {
            flex: 1;
            font-weight: 600;
            color: #333;
            font-size: 15px;
        }
        .tree-node-duration {
            padding: 4px 12px;
            background: white;
            border-radius: 15px;
            font-size: 13px;
            color: #667eea;
            font-weight: 600;
            margin-right: 10px;
        }
        .tree-node-children {
            margin-left: 30px;
            margin-top: 10px;
            padding-left: 20px;
            border-left: 2px dashed #ddd;
        }
        
        /* 时间线 */
        .timeline {
            position: relative;
            padding: 20px 0;
        }
        .timeline-item {
            display: flex;
            gap: 20px;
            margin-bottom: 20px;
            position: relative;
        }
        .timeline-time {
            min-width: 150px;
            font-weight: 600;
            color: #667eea;
            text-align: right;
        }
        .timeline-marker {
            width: 16px;
            height: 16px;
            background: #667eea;
            border-radius: 50%%;
            border: 4px solid white;
            box-shadow: 0 0 0 2px #667eea;
            position: relative;
            z-index: 1;
        }
        .timeline-content {
            flex: 1;
            padding: 15px;
            background: #f8f9fa;
            border-radius: 10px;
        }
        .timeline-title {
            font-weight: 600;
            color: #333;
            margin-bottom: 5px;
        }
        .timeline-desc {
            color: #666;
            font-size: 14px;
        }
        
        /* 性能指标 */
        .metrics-grid {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 15px;
        }
        .metric-card {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 25px;
            border-radius: 12px;
            text-align: center;
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.3);
        }
        .metric-value {
            font-size: 36px;
            font-weight: 700;
            margin-bottom: 8px;
        }
        .metric-label {
            font-size: 14px;
            opacity: 0.9;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        
        /* 全宽卡片 */
        .card-full {
            grid-column: 1 / -1;
        }
        
        /* 响应式 */
        @media (max-width: 1024px) {
            .content-grid {
                grid-template-columns: 1fr;
            }
            .metrics-grid {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <!-- 顶部导航 -->
        <div class="top-nav">
            <a href="/" class="btn-back">← 返回首页</a>
            <div style="color: #667eea; font-weight: 600;">追踪ID: %s</div>
        </div>
        
        <!-- 头部信息 -->
        <div class="header-card">
            <div class="method-title">
                <span>🔍</span>
                <span>%s</span>
                <span class="status-badge status-%s">
                    %s %s
                </span>
            </div>
            
            <div class="meta-grid">
                <div class="meta-item">
                    <div class="meta-label">包路径</div>
                    <div class="meta-value" style="font-size: 14px;">%s</div>
                </div>
                <div class="meta-item">
                    <div class="meta-label">执行耗时</div>
                    <div class="meta-value">%s</div>
                </div>
                <div class="meta-item">
                    <div class="meta-label">Goroutine</div>
                    <div class="meta-value">#%d</div>
                </div>
                <div class="meta-item">
                    <div class="meta-label">开始时间</div>
                    <div class="meta-value" style="font-size: 14px;">%s</div>
                </div>
                <div class="meta-item">
                    <div class="meta-label">结束时间</div>
                    <div class="meta-value" style="font-size: 14px;">%s</div>
                </div>
                <div class="meta-item">
                    <div class="meta-label">子调用数</div>
                    <div class="meta-value">%d</div>
                </div>
            </div>
        </div>
        
        <!-- 性能指标 -->
        <div class="card" style="margin-bottom: 20px;">
            <div class="card-title">📊 性能指标</div>
            <div class="metrics-grid">
                <div class="metric-card">
                    <div class="metric-value">%s</div>
                    <div class="metric-label">执行时间</div>
                </div>
                <div class="metric-card">
                    <div class="metric-value">%s</div>
                    <div class="metric-label">状态</div>
                </div>
                <div class="metric-card">
                    <div class="metric-value">%d</div>
                    <div class="metric-label">子调用</div>
                </div>
            </div>
        </div>
        
        <!-- 参数和返回值 -->
        <div class="content-grid">
            <div class="card">
                <div class="card-title">📥 输入参数</div>
                <div class="param-container">
                    <div class="param-header">Input Parameters</div>
                    <div class="param-content">
                        %s
                    </div>
                </div>
            </div>
            
            <div class="card">
                <div class="card-title">📤 返回值</div>
                <div class="param-container">
                    <div class="param-header">Output Values</div>
                    <div class="param-content">
                        %s
                    </div>
                </div>
                %s
            </div>
        </div>
        
        <!-- 调用树 -->
        %s
    </div>
</body>
</html>
`, 
		trace.MethodName,
		traceID,
		trace.MethodName,
		trace.Status,
		getStatusIcon(trace.Status),
		trace.Status,
		trace.PackageName,
		formatDuration(trace.Duration.Nanoseconds()),
		trace.Goroutine,
		trace.StartTime.Format("2006-01-02 15:04:05.000"),
		trace.EndTime.Format("2006-01-02 15:04:05.000"),
		len(trace.Children),
		formatDuration(trace.Duration.Nanoseconds()),
		trace.Status,
		len(trace.Children),
		formatParamHTML(trace.Input),
		formatParamHTML(trace.Output),
		formatErrorHTML(trace.Error),
		formatCallTreeHTML(trace.Children),
	)
	
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
}

// getStatusIcon 获取状态图标
func getStatusIcon(status string) string {
	switch status {
	case "success":
		return "✅"
	case "error":
		return "❌"
	case "running":
		return "🔄"
	default:
		return "⏸️"
	}
}

// formatDuration 格式化时长
func formatDuration(ns int64) string {
	if ns < 1000 {
		return fmt.Sprintf("%dns", ns)
	} else if ns < 1000000 {
		return fmt.Sprintf("%.2fµs", float64(ns)/1000)
	} else if ns < 1000000000 {
		return fmt.Sprintf("%.2fms", float64(ns)/1000000)
	} else {
		return fmt.Sprintf("%.2fs", float64(ns)/1000000000)
	}
}

// formatParamHTML 格式化参数为 HTML
func formatParamHTML(data interface{}) string {
	if data == nil {
		return `<div class="param-empty">无数据</div>`
	}
	
	// 尝试转换为 JSON
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		// 如果无法转换为 JSON，使用字符串表示
		return fmt.Sprintf(`<pre>%s</pre>`, html.EscapeString(fmt.Sprintf("%+v", data)))
	}
	
	// 对 JSON 进行语法高亮
	jsonStr := string(jsonBytes)
	highlighted := highlightJSON(jsonStr)
	
	return fmt.Sprintf(`<pre>%s</pre>`, highlighted)
}

// highlightJSON 简单的 JSON 语法高亮
func highlightJSON(jsonStr string) string {
	// 转义 HTML
	jsonStr = html.EscapeString(jsonStr)
	
	// 简单的语法高亮（可以进一步优化）
	jsonStr = strings.ReplaceAll(jsonStr, `"`, `<span class="json-string">"</span>`)
	
	return jsonStr
}

// formatErrorHTML 格式化错误信息
func formatErrorHTML(errMsg string) string {
	if errMsg == "" {
		return ""
	}
	
	return fmt.Sprintf(`
		<div class="param-container" style="margin-top: 15px; border-left-color: #dc3545;">
			<div class="param-header" style="background: linear-gradient(135deg, #dc3545, #c82333);">
				❌ 错误信息
			</div>
			<div class="param-content">
				<pre style="background: #f8d7da; color: #721c24;">%s</pre>
			</div>
		</div>
	`, html.EscapeString(errMsg))
}

// formatCallTreeHTML 格式化调用树
func formatCallTreeHTML(children []*tracer.MethodTrace) string {
	if len(children) == 0 {
		return ""
	}
	
	html := `
		<div class="card card-full">
			<div class="card-title">🌲 调用链路图</div>
			<div class="call-tree">
	`
	
	for _, child := range children {
		html += formatTreeNodeHTML(child, 0)
	}
	
	html += `
			</div>
		</div>
	`
	
	return html
}

// formatTreeNodeHTML 格式化树节点
func formatTreeNodeHTML(node *tracer.MethodTrace, level int) string {
	statusClass := "status-" + node.Status
	icon := getStatusIcon(node.Status)
	
	// 构建节点的详细信息
	inputSummary := "无参数"
	if node.Input != nil {
		inputJSON, _ := json.Marshal(node.Input)
		inputStr := string(inputJSON)
		if len(inputStr) > 50 {
			inputSummary = inputStr[:50] + "..."
		} else {
			inputSummary = inputStr
		}
	}
	
	outputSummary := "无返回值"
	if node.Output != nil {
		outputJSON, _ := json.Marshal(node.Output)
		outputStr := string(outputJSON)
		if len(outputStr) > 50 {
			outputSummary = outputStr[:50] + "..."
		} else {
			outputSummary = outputStr
		}
	}
	
	nodeHTML := fmt.Sprintf(`
		<div class="tree-node" style="margin-left: %dpx;">
			<div class="tree-node-content" onclick="window.location.href='/trace/%s'">
				<span class="tree-node-icon">%s</span>
				<div class="tree-node-method">
					<div>%s</div>
					<div style="font-size: 12px; color: #666; margin-top: 5px;">
						<span style="margin-right: 15px;">📥 %s</span>
						<span>📤 %s</span>
					</div>
				</div>
				<span class="tree-node-duration">%s</span>
				<span class="status-badge %s" style="font-size: 11px;">%s</span>
			</div>
	`, 
		level*30,
		node.TraceID,
		icon,
		node.MethodName,
		html.EscapeString(inputSummary),
		html.EscapeString(outputSummary),
		formatDuration(node.Duration.Nanoseconds()),
		statusClass,
		node.Status,
	)
	
	// 递归处理子节点
	if len(node.Children) > 0 {
		nodeHTML += `<div class="tree-node-children">`
		for _, child := range node.Children {
			nodeHTML += formatTreeNodeHTML(child, level+1)
		}
		nodeHTML += `</div>`
	}
	
	nodeHTML += `</div>`
	
	return nodeHTML
}
