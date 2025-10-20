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
            max-width: 1400px; 
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
            padding: 12px 24px;
            background: #667eea;
            color: white;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 600;
            font-size: 15px;
            transition: all 0.3s;
        }
        
        .btn-back:hover { 
            background: #5568d3;
        }
        
        /* 头部信息 */
        .header-card {
            background: white;
            padding: 30px;
            border-radius: 12px;
            margin-bottom: 20px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        
        .method-title {
            font-size: 32px;
            font-weight: 700;
            color: #2c3e50;
            margin-bottom: 20px;
        }
        
        .status-badge {
            display: inline-block;
            padding: 8px 20px;
            border-radius: 20px;
            font-size: 14px;
            font-weight: 600;
            margin-left: 15px;
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
            gap: 40px;
            flex-wrap: wrap;
            font-size: 15px;
            color: #666;
        }
        
        .meta-item strong {
            color: #333;
            margin-right: 8px;
        }
        
        /* 工具栏 */
        .toolbar {
            background: white;
            padding: 20px;
            border-radius: 12px;
            margin-bottom: 20px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            display: flex;
            gap: 15px;
            flex-wrap: wrap;
        }
        
        .toolbar-btn {
            padding: 10px 20px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-weight: 600;
            font-size: 14px;
            transition: all 0.3s;
        }
        
        .toolbar-btn:hover {
            background: #5568d3;
            transform: translateY(-2px);
        }
        
        /* 调用树容器 */
        .tree-container {
            background: white;
            padding: 30px;
            border-radius: 12px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        
        .tree-title {
            font-size: 24px;
            font-weight: 700;
            color: #2c3e50;
            margin-bottom: 25px;
            padding-bottom: 15px;
            border-bottom: 2px solid #e9ecef;
        }
        
        /* 树节点 - 列表式布局 */
        .tree-node {
            margin-bottom: 15px;
        }
        
        .node-card {
            background: #ffffff;
            border: 2px solid #e9ecef;
            border-radius: 10px;
            padding: 20px;
            transition: all 0.3s;
        }
        
        .node-card:hover {
            border-color: #667eea;
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
        }
        
        .node-header {
            display: flex;
            align-items: center;
            gap: 15px;
            margin-bottom: 12px;
        }
        
        .node-icon {
            font-size: 28px;
            flex-shrink: 0;
        }
        
        .node-method {
            font-size: 20px;
            font-weight: 700;
            color: #2c3e50;
            flex: 1;
        }
        
        .node-duration {
            font-size: 18px;
            font-weight: 700;
            color: #667eea;
            padding: 6px 16px;
            background: #f0f3ff;
            border-radius: 20px;
        }
        
        .node-meta {
            display: flex;
            gap: 25px;
            flex-wrap: wrap;
            font-size: 14px;
            color: #666;
            margin-bottom: 12px;
            padding-left: 43px;
        }
        
        .node-meta span {
            display: flex;
            align-items: center;
            gap: 6px;
        }
        
        .node-actions {
            display: flex;
            gap: 10px;
            padding-left: 43px;
        }
        
        .node-btn {
            padding: 8px 18px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 13px;
            font-weight: 600;
            transition: all 0.3s;
        }
        
        .node-btn:hover {
            background: #5568d3;
            transform: translateY(-1px);
        }
        
        .node-btn.secondary {
            background: #6c757d;
        }
        
        .node-btn.secondary:hover {
            background: #5a6268;
        }
        
        /* 参数面板 */
        .params-panel {
            margin-top: 15px;
            padding: 20px;
            background: #f8f9fa;
            border-radius: 8px;
            display: none;
            margin-left: 43px;
        }
        
        .params-panel.show {
            display: block;
        }
        
        .param-section {
            margin-bottom: 15px;
        }
        
        .param-section:last-child {
            margin-bottom: 0;
        }
        
        .param-label {
            font-weight: 700;
            color: #667eea;
            margin-bottom: 8px;
            font-size: 14px;
        }
        
        .param-content {
            background: #2d2d2d;
            color: #f8f8f2;
            padding: 15px;
            border-radius: 6px;
            overflow-x: auto;
            font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
            font-size: 14px;
            line-height: 1.6;
            max-height: 400px;
            overflow-y: auto;
        }
        
        .param-empty {
            color: #999;
            font-style: italic;
        }
        
        /* 子节点容器 - 缩进显示 */
        .node-children {
            margin-top: 15px;
            margin-left: 40px;
            padding-left: 20px;
            border-left: 3px solid #dee2e6;
            display: none;
        }
        
        .node-children.show {
            display: block;
        }
        
        /* 层级指示器 */
        .level-indicator {
            display: inline-block;
            width: 30px;
            height: 3px;
            background: #dee2e6;
            margin-right: 10px;
            vertical-align: middle;
        }
        
        /* 响应式 */
        @media (max-width: 768px) {
            .node-children {
                margin-left: 20px;
                padding-left: 10px;
            }
            
            .node-meta, .node-actions, .params-panel {
                padding-left: 0;
                margin-left: 0;
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
                🔍 %s
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
            <button class="toolbar-btn" onclick="expandAll()">📂 展开全部调用</button>
            <button class="toolbar-btn" onclick="collapseAll()">📁 折叠全部调用</button>
            <button class="toolbar-btn" onclick="showAllParams()">👁️ 显示全部参数</button>
            <button class="toolbar-btn" onclick="hideAllParams()">🙈 隐藏全部参数</button>
        </div>
        
        <!-- 调用树 -->
        <div class="tree-container">
            <div class="tree-title">🌲 完整调用链路（列表视图）</div>
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
                container.innerHTML = '<p style="color: #999; font-size: 16px;">无追踪数据</p>';
                return;
            }
            container.innerHTML = renderNode(traceData, 0);
        }
        
        // 渲染单个节点
        function renderNode(node, level) {
            if (!node) return '';
            
            const hasChildren = node.Children && node.Children.length > 0;
            const nodeId = 'node-' + (node.TraceID || Math.random().toString(36).substr(2, 9));
            
            let html = '<div class="tree-node">';
            html += '<div class="node-card">';
            
            // 节点头部
            html += '<div class="node-header">';
            
            // 层级指示器
            for (let i = 0; i < level; i++) {
                html += '<span class="level-indicator"></span>';
            }
            
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
            if (node.StartTime) {
                html += '<span><strong>开始:</strong> ' + formatTime(node.StartTime) + '</span>';
            }
            html += '</div>';
            
            // 操作按钮
            html += '<div class="node-actions">';
            html += '<button class="node-btn" onclick="toggleParams(\'' + nodeId + '\')">查看参数</button>';
            if (hasChildren) {
                html += '<button class="node-btn secondary" onclick="toggleChildren(\'' + nodeId + '\')">展开子调用 (' + node.Children.length + ')</button>';
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
            const buttons = document.querySelectorAll('button[onclick*="' + nodeId + '"]');
            
            if (panel && panel.classList.contains('show')) {
                panel.classList.remove('show');
                buttons.forEach(btn => {
                    if (btn.textContent.includes('隐藏')) {
                        btn.textContent = '查看参数';
                    }
                });
            } else if (panel) {
                panel.classList.add('show');
                buttons.forEach(btn => {
                    if (btn.textContent.includes('查看')) {
                        btn.textContent = '隐藏参数';
                    }
                });
            }
        }
        
        // 切换子节点显示
        function toggleChildren(nodeId) {
            const children = document.getElementById('children-' + nodeId);
            const buttons = document.querySelectorAll('button[onclick*="' + nodeId + '"]');
            
            if (children && children.classList.contains('show')) {
                children.classList.remove('show');
                buttons.forEach(btn => {
                    if (btn.textContent.includes('折叠')) {
                        const count = children.querySelectorAll(':scope > .tree-node').length;
                        btn.textContent = '展开子调用 (' + count + ')';
                    }
                });
            } else if (children) {
                children.classList.add('show');
                buttons.forEach(btn => {
                    if (btn.textContent.includes('展开')) {
                        const count = children.querySelectorAll(':scope > .tree-node').length;
                        btn.textContent = '折叠子调用 (' + count + ')';
                    }
                });
            }
        }
        
        // 展开全部
        function expandAll() {
            document.querySelectorAll('.node-children').forEach(el => {
                el.classList.add('show');
            });
            document.querySelectorAll('.node-btn.secondary').forEach(btn => {
                if (btn.textContent.includes('展开')) {
                    const match = btn.textContent.match(/\\d+/);
                    if (match) {
                        btn.textContent = '折叠子调用 (' + match[0] + ')';
                    }
                }
            });
        }
        
        // 折叠全部
        function collapseAll() {
            document.querySelectorAll('.node-children').forEach(el => {
                el.classList.remove('show');
            });
            document.querySelectorAll('.node-btn.secondary').forEach(btn => {
                if (btn.textContent.includes('折叠')) {
                    const match = btn.textContent.match(/\\d+/);
                    if (match) {
                        btn.textContent = '展开子调用 (' + match[0] + ')';
                    }
                }
            });
        }
        
        // 显示所有参数
        function showAllParams() {
            document.querySelectorAll('.params-panel').forEach(el => {
                el.classList.add('show');
            });
            document.querySelectorAll('.node-btn:not(.secondary)').forEach(btn => {
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
            document.querySelectorAll('.node-btn:not(.secondary)').forEach(btn => {
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
        
        // 格式化时间
        function formatTime(timeStr) {
            if (!timeStr) return '';
            try {
                const date = new Date(timeStr);
                return date.toLocaleTimeString('zh-CN', { hour12: false });
            } catch (e) {
                return timeStr;
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
                document.getElementById('callTree').innerHTML = '<p style="color: #dc3545; font-size: 16px;">渲染失败: ' + e.message + '</p>';
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
