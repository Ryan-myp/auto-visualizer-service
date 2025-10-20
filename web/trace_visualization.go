package web

import (
	"encoding/json"
	"fmt"
	"html"
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
	
	// 将追踪数据转换为 JSON
	traceJSON, _ := json.Marshal(trace)
	
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
        
        .container { 
            max-width: 1800px; 
            margin: 0 auto; 
        }
        
        /* 顶部导航 */
        .top-nav {
            background: rgba(255,255,255,0.98);
            padding: 20px 30px;
            border-radius: 16px;
            margin-bottom: 20px;
            box-shadow: 0 8px 32px rgba(0,0,0,0.12);
            display: flex;
            justify-content: space-between;
            align-items: center;
            backdrop-filter: blur(10px);
        }
        
        .btn-back {
            padding: 12px 24px;
            background: linear-gradient(135deg, #667eea, #764ba2);
            color: white;
            text-decoration: none;
            border-radius: 10px;
            font-weight: 600;
            transition: all 0.3s;
            box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
        }
        
        .btn-back:hover { 
            transform: translateY(-2px); 
            box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
        }
        
        /* 头部卡片 */
        .header-card {
            background: white;
            padding: 35px;
            border-radius: 16px;
            margin-bottom: 20px;
            box-shadow: 0 8px 32px rgba(0,0,0,0.12);
        }
        
        .method-title {
            font-size: 36px;
            font-weight: 800;
            background: linear-gradient(135deg, #667eea, #764ba2);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 20px;
            display: flex;
            align-items: center;
            gap: 15px;
        }
        
        .status-badge {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            padding: 10px 20px;
            border-radius: 25px;
            font-size: 14px;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        
        .status-success { 
            background: linear-gradient(135deg, #d4edda, #c3e6cb);
            color: #155724; 
        }
        
        .status-error { 
            background: linear-gradient(135deg, #f8d7da, #f5c6cb);
            color: #721c24; 
        }
        
        .meta-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-top: 25px;
        }
        
        .meta-item {
            padding: 20px;
            background: linear-gradient(135deg, #f8f9fa, #e9ecef);
            border-radius: 12px;
            border-left: 5px solid #667eea;
            transition: all 0.3s;
        }
        
        .meta-item:hover {
            transform: translateY(-3px);
            box-shadow: 0 6px 20px rgba(0,0,0,0.1);
        }
        
        .meta-label {
            font-size: 13px;
            color: #666;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 8px;
            font-weight: 600;
        }
        
        .meta-value {
            font-size: 20px;
            font-weight: 700;
            color: #333;
        }
        
        /* 调用树 */
        .call-tree-container {
            background: white;
            padding: 30px;
            border-radius: 16px;
            box-shadow: 0 8px 32px rgba(0,0,0,0.12);
        }
        
        .tree-title {
            font-size: 24px;
            font-weight: 700;
            color: #333;
            margin-bottom: 25px;
            padding-bottom: 15px;
            border-bottom: 3px solid #667eea;
            display: flex;
            align-items: center;
            gap: 12px;
        }
        
        .tree-node {
            margin-bottom: 12px;
            position: relative;
        }
        
        .tree-node-header {
            display: flex;
            align-items: flex-start;
            padding: 18px;
            background: linear-gradient(135deg, #f8f9fa, #ffffff);
            border-radius: 12px;
            border-left: 5px solid #667eea;
            cursor: pointer;
            transition: all 0.3s;
            box-shadow: 0 2px 8px rgba(0,0,0,0.05);
        }
        
        .tree-node-header:hover {
            background: linear-gradient(135deg, #e9ecef, #f8f9fa);
            transform: translateX(5px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
        }
        
        .tree-node-icon {
            font-size: 24px;
            margin-right: 15px;
            flex-shrink: 0;
        }
        
        .tree-node-content {
            flex: 1;
            min-width: 0;
        }
        
        .tree-node-method {
            font-weight: 700;
            color: #333;
            font-size: 16px;
            margin-bottom: 10px;
            word-break: break-word;
        }
        
        .tree-node-meta {
            display: flex;
            gap: 20px;
            flex-wrap: wrap;
            margin-bottom: 10px;
        }
        
        .tree-node-meta-item {
            display: flex;
            align-items: center;
            gap: 6px;
            font-size: 13px;
            color: #666;
        }
        
        .tree-node-meta-item strong {
            color: #667eea;
            font-weight: 600;
        }
        
        .tree-node-params {
            margin-top: 12px;
            padding: 12px;
            background: #f8f9fa;
            border-radius: 8px;
            font-size: 13px;
            display: none;
        }
        
        .tree-node-params.show {
            display: block;
        }
        
        .param-section {
            margin-bottom: 10px;
        }
        
        .param-label {
            font-weight: 600;
            color: #667eea;
            margin-bottom: 5px;
            display: flex;
            align-items: center;
            gap: 6px;
        }
        
        .param-value {
            background: #2d2d2d;
            color: #f8f8f2;
            padding: 10px;
            border-radius: 6px;
            overflow-x: auto;
            font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
            font-size: 12px;
            line-height: 1.5;
            max-height: 200px;
            overflow-y: auto;
        }
        
        .tree-node-toggle {
            padding: 6px 12px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 12px;
            font-weight: 600;
            transition: all 0.3s;
            margin-left: auto;
            flex-shrink: 0;
        }
        
        .tree-node-toggle:hover {
            background: #5568d3;
        }
        
        .tree-node-children {
            margin-left: 40px;
            margin-top: 12px;
            padding-left: 20px;
            border-left: 3px dashed #ddd;
            display: none;
        }
        
        .tree-node-children.show {
            display: block;
        }
        
        .collapse-btn {
            padding: 4px 10px;
            background: rgba(102, 126, 234, 0.1);
            color: #667eea;
            border: 1px solid #667eea;
            border-radius: 6px;
            cursor: pointer;
            font-size: 11px;
            font-weight: 600;
            transition: all 0.3s;
            margin-left: 10px;
        }
        
        .collapse-btn:hover {
            background: #667eea;
            color: white;
        }
        
        /* 性能指标 */
        .metrics-container {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 20px;
        }
        
        .metric-card {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 30px;
            border-radius: 16px;
            text-align: center;
            box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
            transition: all 0.3s;
        }
        
        .metric-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 12px 32px rgba(102, 126, 234, 0.6);
        }
        
        .metric-value {
            font-size: 42px;
            font-weight: 800;
            margin-bottom: 10px;
        }
        
        .metric-label {
            font-size: 14px;
            opacity: 0.95;
            text-transform: uppercase;
            letter-spacing: 1.5px;
            font-weight: 600;
        }
        
        /* 工具栏 */
        .toolbar {
            background: white;
            padding: 20px;
            border-radius: 12px;
            margin-bottom: 20px;
            box-shadow: 0 4px 16px rgba(0,0,0,0.08);
            display: flex;
            gap: 15px;
            flex-wrap: wrap;
        }
        
        .toolbar-btn {
            padding: 10px 20px;
            background: linear-gradient(135deg, #667eea, #764ba2);
            color: white;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-weight: 600;
            transition: all 0.3s;
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
        }
        
        .toolbar-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 16px rgba(102, 126, 234, 0.5);
        }
        
        /* 响应式 */
        @media (max-width: 768px) {
            .meta-grid {
                grid-template-columns: 1fr;
            }
            .tree-node-children {
                margin-left: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <!-- 顶部导航 -->
        <div class="top-nav">
            <a href="/" class="btn-back">← 返回首页</a>
            <div style="color: #667eea; font-weight: 600; font-size: 14px;">追踪ID: %s</div>
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
                    <div class="meta-label">📦 包路径</div>
                    <div class="meta-value" style="font-size: 14px; word-break: break-all;">%s</div>
                </div>
                <div class="meta-item">
                    <div class="meta-label">⏱️ 执行耗时</div>
                    <div class="meta-value">%s</div>
                </div>
                <div class="meta-item">
                    <div class="meta-label">🔢 Goroutine</div>
                    <div class="meta-value">#%d</div>
                </div>
                <div class="meta-item">
                    <div class="meta-label">🕐 开始时间</div>
                    <div class="meta-value" style="font-size: 14px;">%s</div>
                </div>
                <div class="meta-item">
                    <div class="meta-label">🕑 结束时间</div>
                    <div class="meta-value" style="font-size: 14px;">%s</div>
                </div>
                <div class="meta-item">
                    <div class="meta-label">🌲 子调用数</div>
                    <div class="meta-value">%d</div>
                </div>
            </div>
        </div>
        
        <!-- 性能指标 -->
        <div class="metrics-container">
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
        
        <!-- 工具栏 -->
        <div class="toolbar">
            <button class="toolbar-btn" onclick="expandAll()">📂 展开全部</button>
            <button class="toolbar-btn" onclick="collapseAll()">📁 折叠全部</button>
            <button class="toolbar-btn" onclick="showAllParams()">👁️ 显示所有参数</button>
            <button class="toolbar-btn" onclick="hideAllParams()">🙈 隐藏所有参数</button>
        </div>
        
        <!-- 调用树 -->
        <div class="call-tree-container">
            <div class="tree-title">
                <span>🌲</span>
                <span>完整调用链路</span>
            </div>
            <div id="callTree"></div>
        </div>
    </div>
    
    <script>
        // 追踪数据
        const traceData = %s;
        
        // 渲染调用树
        function renderTree() {
            const container = document.getElementById('callTree');
            container.innerHTML = renderNode(traceData, 0);
        }
        
        // 渲染单个节点
        function renderNode(node, level) {
            const hasChildren = node.Children && node.Children.length > 0;
            const nodeId = 'node-' + node.TraceID;
            
            let html = '<div class="tree-node" style="margin-left: ' + (level * 30) + 'px;">';
            
            // 节点头部
            html += '<div class="tree-node-header">';
            html += '<span class="tree-node-icon">' + getStatusIcon(node.Status) + '</span>';
            html += '<div class="tree-node-content">';
            html += '<div class="tree-node-method">' + escapeHtml(node.MethodName) + '</div>';
            
            // 元数据
            html += '<div class="tree-node-meta">';
            html += '<div class="tree-node-meta-item"><strong>⏱️</strong> ' + formatDuration(node.Duration) + '</div>';
            html += '<div class="tree-node-meta-item"><strong>🔢</strong> Goroutine #' + node.Goroutine + '</div>';
            if (hasChildren) {
                html += '<div class="tree-node-meta-item"><strong>🌲</strong> ' + node.Children.length + ' 个子调用</div>';
            }
            html += '</div>';
            
            // 参数区域（默认隐藏）
            html += '<div class="tree-node-params" id="params-' + nodeId + '">';
            
            // 输入参数
            html += '<div class="param-section">';
            html += '<div class="param-label">📥 输入参数</div>';
            html += '<div class="param-value">' + formatJSON(node.Input) + '</div>';
            html += '</div>';
            
            // 返回值
            html += '<div class="param-section">';
            html += '<div class="param-label">📤 返回值</div>';
            html += '<div class="param-value">' + formatJSON(node.Output) + '</div>';
            html += '</div>';
            
            // 错误信息
            if (node.Error) {
                html += '<div class="param-section">';
                html += '<div class="param-label">❌ 错误信息</div>';
                html += '<div class="param-value" style="background: #f8d7da; color: #721c24;">' + escapeHtml(node.Error) + '</div>';
                html += '</div>';
            }
            
            html += '</div>';
            
            html += '</div>';
            
            // 切换按钮
            html += '<button class="tree-node-toggle" onclick="toggleParams(\'' + nodeId + '\')">查看参数</button>';
            if (hasChildren) {
                html += '<button class="collapse-btn" onclick="toggleChildren(\'' + nodeId + '\')">展开 (' + node.Children.length + ')</button>';
            }
            
            html += '</div>';
            
            // 子节点
            if (hasChildren) {
                html += '<div class="tree-node-children" id="children-' + nodeId + '">';
                node.Children.forEach(child => {
                    html += renderNode(child, level + 1);
                });
                html += '</div>';
            }
            
            html += '</div>';
            
            return html;
        }
        
        // 切换参数显示
        function toggleParams(nodeId) {
            const paramsEl = document.getElementById('params-' + nodeId);
            const btn = event.target;
            
            if (paramsEl.classList.contains('show')) {
                paramsEl.classList.remove('show');
                btn.textContent = '查看参数';
            } else {
                paramsEl.classList.add('show');
                btn.textContent = '隐藏参数';
            }
        }
        
        // 切换子节点显示
        function toggleChildren(nodeId) {
            const childrenEl = document.getElementById('children-' + nodeId);
            const btn = event.target;
            
            if (childrenEl.classList.contains('show')) {
                childrenEl.classList.remove('show');
                btn.textContent = '展开 (' + childrenEl.querySelectorAll(':scope > .tree-node').length + ')';
            } else {
                childrenEl.classList.add('show');
                btn.textContent = '折叠 (' + childrenEl.querySelectorAll(':scope > .tree-node').length + ')';
            }
        }
        
        // 展开全部
        function expandAll() {
            document.querySelectorAll('.tree-node-children').forEach(el => {
                el.classList.add('show');
            });
            document.querySelectorAll('.collapse-btn').forEach(btn => {
                const count = btn.closest('.tree-node').querySelector('.tree-node-children').querySelectorAll(':scope > .tree-node').length;
                btn.textContent = '折叠 (' + count + ')';
            });
        }
        
        // 折叠全部
        function collapseAll() {
            document.querySelectorAll('.tree-node-children').forEach(el => {
                el.classList.remove('show');
            });
            document.querySelectorAll('.collapse-btn').forEach(btn => {
                const count = btn.closest('.tree-node').querySelector('.tree-node-children').querySelectorAll(':scope > .tree-node').length;
                btn.textContent = '展开 (' + count + ')';
            });
        }
        
        // 显示所有参数
        function showAllParams() {
            document.querySelectorAll('.tree-node-params').forEach(el => {
                el.classList.add('show');
            });
            document.querySelectorAll('.tree-node-toggle').forEach(btn => {
                btn.textContent = '隐藏参数';
            });
        }
        
        // 隐藏所有参数
        function hideAllParams() {
            document.querySelectorAll('.tree-node-params').forEach(el => {
                el.classList.remove('show');
            });
            document.querySelectorAll('.tree-node-toggle').forEach(btn => {
                btn.textContent = '查看参数';
            });
        }
        
        // 格式化 JSON
        function formatJSON(data) {
            if (data === null || data === undefined) {
                return '<span style="color: #c678dd;">null</span>';
            }
            try {
                return escapeHtml(JSON.stringify(data, null, 2));
            } catch (e) {
                return escapeHtml(String(data));
            }
        }
        
        // 格式化时长
        function formatDuration(ns) {
            if (ns < 1000) {
                return ns + 'ns';
            } else if (ns < 1000000) {
                return (ns / 1000).toFixed(2) + 'µs';
            } else if (ns < 1000000000) {
                return (ns / 1000000).toFixed(2) + 'ms';
            } else {
                return (ns / 1000000000).toFixed(2) + 's';
            }
        }
        
        // 获取状态图标
        function getStatusIcon(status) {
            switch (status) {
                case 'success': return '✅';
                case 'error': return '❌';
                case 'running': return '🔄';
                default: return '⏸️';
            }
        }
        
        // HTML 转义
        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
        
        // 页面加载时渲染
        window.onload = function() {
            renderTree();
        };
    </script>
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
		string(traceJSON),
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
