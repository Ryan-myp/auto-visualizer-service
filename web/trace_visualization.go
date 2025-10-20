package web

import (
	"encoding/json"
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
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f5f7fa;
            min-height: 100vh;
            padding: 20px;
        }
        
        .container { 
            max-width: 1600px; 
            margin: 0 auto; 
        }
        
        /* 顶部导航 */
        .top-bar {
            background: white;
            padding: 20px 30px;
            border-radius: 12px;
            margin-bottom: 20px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .btn-back {
            padding: 10px 20px;
            background: #667eea;
            color: white;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 600;
            transition: all 0.3s;
        }
        
        .btn-back:hover { 
            background: #5568d3;
            transform: translateY(-2px);
        }
        
        /* 头部信息卡片 */
        .header-card {
            background: white;
            padding: 30px;
            border-radius: 12px;
            margin-bottom: 20px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        
        .method-title {
            font-size: 28px;
            font-weight: 700;
            color: #2c3e50;
            margin-bottom: 20px;
            display: flex;
            align-items: center;
            gap: 12px;
        }
        
        .status-badge {
            padding: 6px 16px;
            border-radius: 20px;
            font-size: 13px;
            font-weight: 600;
        }
        
        .status-success { 
            background: #d4edda;
            color: #155724; 
        }
        
        .status-error { 
            background: #f8d7da;
            color: #721c24; 
        }
        
        .meta-row {
            display: flex;
            gap: 30px;
            flex-wrap: wrap;
            font-size: 14px;
            color: #666;
        }
        
        .meta-item strong {
            color: #333;
            margin-right: 8px;
        }
        
        /* 工具栏 */
        .toolbar {
            background: white;
            padding: 15px 20px;
            border-radius: 12px;
            margin-bottom: 20px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            display: flex;
            gap: 12px;
            flex-wrap: wrap;
        }
        
        .toolbar-btn {
            padding: 8px 16px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-weight: 600;
            font-size: 13px;
            transition: all 0.3s;
        }
        
        .toolbar-btn:hover {
            background: #5568d3;
        }
        
        .toolbar-btn.secondary {
            background: #6c757d;
        }
        
        .toolbar-btn.secondary:hover {
            background: #5a6268;
        }
        
        /* 调用树容器 */
        .tree-container {
            background: white;
            padding: 25px;
            border-radius: 12px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        
        .tree-title {
            font-size: 20px;
            font-weight: 700;
            color: #2c3e50;
            margin-bottom: 20px;
            padding-bottom: 12px;
            border-bottom: 2px solid #e9ecef;
        }
        
        /* 调用树节点 */
        .tree-node {
            margin-bottom: 8px;
        }
        
        .node-card {
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            border-radius: 8px;
            padding: 16px;
            transition: all 0.3s;
        }
        
        .node-card:hover {
            background: #e9ecef;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        
        .node-header {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-bottom: 8px;
        }
        
        .node-icon {
            font-size: 20px;
            flex-shrink: 0;
        }
        
        .node-method {
            font-size: 16px;
            font-weight: 700;
            color: #2c3e50;
            flex: 1;
        }
        
        .node-duration {
            font-size: 14px;
            font-weight: 600;
            color: #667eea;
            padding: 4px 12px;
            background: white;
            border-radius: 12px;
        }
        
        .node-meta {
            display: flex;
            gap: 20px;
            flex-wrap: wrap;
            font-size: 13px;
            color: #666;
            margin-bottom: 10px;
        }
        
        .node-actions {
            display: flex;
            gap: 8px;
            margin-top: 10px;
        }
        
        .node-btn {
            padding: 6px 14px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 12px;
            font-weight: 600;
            transition: all 0.3s;
        }
        
        .node-btn:hover {
            background: #5568d3;
        }
        
        .node-btn.collapse {
            background: #6c757d;
        }
        
        .node-btn.collapse:hover {
            background: #5a6268;
        }
        
        /* 参数展示 */
        .params-panel {
            margin-top: 12px;
            display: none;
            background: white;
            border-radius: 8px;
            padding: 15px;
            border: 1px solid #e9ecef;
        }
        
        .params-panel.show {
            display: block;
        }
        
        .param-section {
            margin-bottom: 12px;
        }
        
        .param-section:last-child {
            margin-bottom: 0;
        }
        
        .param-label {
            font-weight: 700;
            color: #667eea;
            margin-bottom: 6px;
            font-size: 13px;
        }
        
        .param-content {
            background: #2d2d2d;
            color: #f8f8f2;
            padding: 12px;
            border-radius: 6px;
            overflow-x: auto;
            font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
            font-size: 13px;
            line-height: 1.6;
            max-height: 300px;
            overflow-y: auto;
        }
        
        .param-empty {
            color: #999;
            font-style: italic;
        }
        
        /* 子节点容器 */
        .node-children {
            margin-left: 30px;
            margin-top: 8px;
            padding-left: 20px;
            border-left: 2px solid #dee2e6;
            display: none;
        }
        
        .node-children.show {
            display: block;
        }
        
        /* 响应式 */
        @media (max-width: 768px) {
            .node-children {
                margin-left: 15px;
                padding-left: 10px;
            }
            
            .meta-row {
                flex-direction: column;
                gap: 10px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <!-- 顶部导航 -->
        <div class="top-bar">
            <a href="/" class="btn-back">← 返回首页</a>
            <div style="color: #666; font-size: 14px;">追踪ID: %s</div>
        </div>
        
        <!-- 头部信息 -->
        <div class="header-card">
            <div class="method-title">
                <span>🔍</span>
                <span>%s</span>
                <span class="status-badge status-%s">%s %s</span>
            </div>
            
            <div class="meta-row">
                <div class="meta-item">
                    <strong>📦 包路径:</strong> %s
                </div>
                <div class="meta-item">
                    <strong>⏱️ 执行耗时:</strong> %s
                </div>
                <div class="meta-item">
                    <strong>🔢 Goroutine:</strong> #%d
                </div>
                <div class="meta-item">
                    <strong>🌲 子调用:</strong> %d 个
                </div>
            </div>
        </div>
        
        <!-- 工具栏 -->
        <div class="toolbar">
            <button class="toolbar-btn" onclick="expandAll()">📂 展开全部</button>
            <button class="toolbar-btn" onclick="collapseAll()">📁 折叠全部</button>
            <button class="toolbar-btn" onclick="showAllParams()">👁️ 显示所有参数</button>
            <button class="toolbar-btn secondary" onclick="hideAllParams()">🙈 隐藏所有参数</button>
        </div>
        
        <!-- 调用树 -->
        <div class="tree-container">
            <div class="tree-title">🌲 完整调用链路</div>
            <div id="callTree"></div>
        </div>
    </div>
    
    <script>
        // 追踪数据
        const traceData = %s;
        
        // 渲染调用树
        function renderTree() {
            const container = document.getElementById('callTree');
            if (!traceData) {
                container.innerHTML = '<p style="color: #999;">无追踪数据</p>';
                return;
            }
            container.innerHTML = renderNode(traceData, 0);
        }
        
        // 渲染单个节点
        function renderNode(node, level) {
            if (!node) return '';
            
            const hasChildren = node.Children && node.Children.length > 0;
            const nodeId = 'node-' + (node.TraceID || Math.random().toString(36).substr(2, 9));
            
            let html = '<div class="tree-node" style="margin-left: ' + (level * 0) + 'px;">';
            html += '<div class="node-card">';
            
            // 节点头部
            html += '<div class="node-header">';
            html += '<span class="node-icon">' + getStatusIcon(node.Status) + '</span>';
            html += '<span class="node-method">' + escapeHtml(node.MethodName || '未知方法') + '</span>';
            html += '<span class="node-duration">' + formatDuration(node.Duration || 0) + '</span>';
            html += '</div>';
            
            // 元数据
            html += '<div class="node-meta">';
            html += '<span><strong>Goroutine:</strong> #' + (node.Goroutine || 0) + '</span>';
            if (hasChildren) {
                html += '<span><strong>子调用:</strong> ' + node.Children.length + ' 个</span>';
            }
            if (node.PackageName) {
                html += '<span><strong>包:</strong> ' + escapeHtml(node.PackageName) + '</span>';
            }
            html += '</div>';
            
            // 操作按钮
            html += '<div class="node-actions">';
            html += '<button class="node-btn" onclick="toggleParams(\'' + nodeId + '\')">查看参数</button>';
            if (hasChildren) {
                html += '<button class="node-btn collapse" onclick="toggleChildren(\'' + nodeId + '\')">展开 (' + node.Children.length + ')</button>';
            }
            html += '</div>';
            
            // 参数面板
            html += '<div class="params-panel" id="params-' + nodeId + '">';
            
            // 输入参数
            html += '<div class="param-section">';
            html += '<div class="param-label">📥 输入参数</div>';
            html += '<div class="param-content">' + formatParam(node.Input) + '</div>';
            html += '</div>';
            
            // 返回值
            html += '<div class="param-section">';
            html += '<div class="param-label">📤 返回值</div>';
            html += '<div class="param-content">' + formatParam(node.Output) + '</div>';
            html += '</div>';
            
            // 错误信息
            if (node.Error) {
                html += '<div class="param-section">';
                html += '<div class="param-label">❌ 错误信息</div>';
                html += '<div class="param-content" style="background: #f8d7da; color: #721c24;">' + escapeHtml(node.Error) + '</div>';
                html += '</div>';
            }
            
            html += '</div>';
            
            html += '</div>';
            
            // 子节点
            if (hasChildren) {
                html += '<div class="node-children" id="children-' + nodeId + '">';
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
            const panel = document.getElementById('params-' + nodeId);
            const btn = event.target;
            
            if (panel && panel.classList.contains('show')) {
                panel.classList.remove('show');
                btn.textContent = '查看参数';
            } else if (panel) {
                panel.classList.add('show');
                btn.textContent = '隐藏参数';
            }
        }
        
        // 切换子节点显示
        function toggleChildren(nodeId) {
            const children = document.getElementById('children-' + nodeId);
            const btn = event.target;
            
            if (children && children.classList.contains('show')) {
                children.classList.remove('show');
                const count = children.querySelectorAll(':scope > .tree-node').length;
                btn.textContent = '展开 (' + count + ')';
            } else if (children) {
                children.classList.add('show');
                const count = children.querySelectorAll(':scope > .tree-node').length;
                btn.textContent = '折叠 (' + count + ')';
            }
        }
        
        // 展开全部
        function expandAll() {
            document.querySelectorAll('.node-children').forEach(el => {
                el.classList.add('show');
            });
            document.querySelectorAll('.node-btn.collapse').forEach(btn => {
                const nodeId = btn.onclick.toString().match(/'([^']+)'/)[1];
                const children = document.getElementById('children-' + nodeId);
                if (children) {
                    const count = children.querySelectorAll(':scope > .tree-node').length;
                    btn.textContent = '折叠 (' + count + ')';
                }
            });
        }
        
        // 折叠全部
        function collapseAll() {
            document.querySelectorAll('.node-children').forEach(el => {
                el.classList.remove('show');
            });
            document.querySelectorAll('.node-btn.collapse').forEach(btn => {
                const nodeId = btn.onclick.toString().match(/'([^']+)'/)[1];
                const children = document.getElementById('children-' + nodeId);
                if (children) {
                    const count = children.querySelectorAll(':scope > .tree-node').length;
                    btn.textContent = '展开 (' + count + ')';
                }
            });
        }
        
        // 显示所有参数
        function showAllParams() {
            document.querySelectorAll('.params-panel').forEach(el => {
                el.classList.add('show');
            });
            document.querySelectorAll('.node-btn:not(.collapse)').forEach(btn => {
                if (btn.textContent.includes('查看')) {
                    btn.textContent = '隐藏参数';
                }
            });
        }
        
        // 隐藏所有参数
        function hideAllParams() {
            document.querySelectorAll('.params-panel').forEach(el => {
                el.classList.remove('show');
            });
            document.querySelectorAll('.node-btn:not(.collapse)').forEach(btn => {
                if (btn.textContent.includes('隐藏')) {
                    btn.textContent = '查看参数';
                }
            });
        }
        
        // 格式化参数
        function formatParam(data) {
            if (data === null || data === undefined) {
                return '<span class="param-empty">null</span>';
            }
            try {
                return escapeHtml(JSON.stringify(data, null, 2));
            } catch (e) {
                return escapeHtml(String(data));
            }
        }
        
        // 格式化时长
        function formatDuration(ns) {
            if (typeof ns !== 'number') ns = 0;
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
            if (typeof text !== 'string') text = String(text);
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
        
        // 页面加载时渲染
        window.onload = function() {
            try {
                renderTree();
            } catch (e) {
                console.error('渲染失败:', e);
                document.getElementById('callTree').innerHTML = '<p style="color: #dc3545;">渲染失败: ' + e.message + '</p>';
            }
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
