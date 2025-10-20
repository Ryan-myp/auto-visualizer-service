package tracer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"sync"
	"time"
)

// MethodTrace 方法追踪信息
type MethodTrace struct {
	TraceID      string                 `json:"trace_id"`
	ParentID     string                 `json:"parent_id,omitempty"`
	MethodName   string                 `json:"method_name"`
	PackageName  string                 `json:"package_name"`
	FileName     string                 `json:"file_name"`
	LineNumber   int                    `json:"line_number"`
	Input        []interface{}          `json:"input"`
	Output       []interface{}          `json:"output"`
	Error        string                 `json:"error,omitempty"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      time.Time              `json:"end_time"`
	Duration     time.Duration          `json:"duration"`
	Goroutine    int                    `json:"goroutine"`
	CallStack    []string               `json:"call_stack"`
	Children     []*MethodTrace         `json:"children,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Status       string                 `json:"status"` // running, success, error
}

// Tracer 方法追踪器
type Tracer struct {
	enabled       bool
	traces        map[string]*MethodTrace
	activeTraces  map[int64]*MethodTrace // goroutine ID -> current trace
	mu            sync.RWMutex
	sampleRate    float64
	maxDepth      int
	captureStack  bool
	captureValues bool
	filters       []TraceFilter
	handlers      []TraceHandler
}

// TraceFilter 追踪过滤器
type TraceFilter func(methodName string) bool

// TraceHandler 追踪处理器
type TraceHandler func(trace *MethodTrace)

// globalTracer 全局追踪器实例
var globalTracer *Tracer
var once sync.Once

// InitTracer 初始化追踪器
func InitTracer(opts ...TracerOption) *Tracer {
	once.Do(func() {
		globalTracer = &Tracer{
			enabled:       true,
			traces:        make(map[string]*MethodTrace),
			activeTraces:  make(map[int64]*MethodTrace),
			sampleRate:    1.0,
			maxDepth:      10,
			captureStack:  true,
			captureValues: true,
			filters:       []TraceFilter{},
			handlers:      []TraceHandler{},
		}

		// 应用选项
		for _, opt := range opts {
			opt(globalTracer)
		}
	})
	return globalTracer
}

// GetTracer 获取全局追踪器
func GetTracer() *Tracer {
	if globalTracer == nil {
		return InitTracer()
	}
	return globalTracer
}

// TracerOption 追踪器选项
type TracerOption func(*Tracer)

// WithSampleRate 设置采样率
func WithSampleRate(rate float64) TracerOption {
	return func(t *Tracer) {
		t.sampleRate = rate
	}
}

// WithMaxDepth 设置最大追踪深度
func WithMaxDepth(depth int) TracerOption {
	return func(t *Tracer) {
		t.maxDepth = depth
	}
}

// WithCaptureStack 设置是否捕获调用栈
func WithCaptureStack(capture bool) TracerOption {
	return func(t *Tracer) {
		t.captureStack = capture
	}
}

// WithFilter 添加过滤器
func WithFilter(filter TraceFilter) TracerOption {
	return func(t *Tracer) {
		t.filters = append(t.filters, filter)
	}
}

// WithHandler 添加处理器
func WithHandler(handler TraceHandler) TracerOption {
	return func(t *Tracer) {
		t.handlers = append(t.handlers, handler)
	}
}

// Trace 追踪方法执行（装饰器模式）
func (t *Tracer) Trace(methodName string, fn interface{}) interface{} {
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		panic("Trace: fn must be a function")
	}

	fnType := fnValue.Type()

	// 创建包装函数
	wrapper := reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
		// 检查是否应该追踪
		if !t.shouldTrace(methodName) {
			return fnValue.Call(args)
		}

		// 开始追踪
		trace := t.StartTrace(methodName, reflectValuesToInterfaces(args))

		// 执行原函数
		var results []reflect.Value
		var panicErr interface{}

		func() {
			defer func() {
				if r := recover(); r != nil {
					panicErr = r
				}
			}()
			results = fnValue.Call(args)
		}()

		// 结束追踪
		if panicErr != nil {
			t.EndTrace(trace.TraceID, nil, fmt.Errorf("panic: %v", panicErr))
			panic(panicErr)
		} else {
			outputs := reflectValuesToInterfaces(results)
			var err error
			// 检查最后一个返回值是否是error
			if len(results) > 0 {
				lastResult := results[len(results)-1]
				if lastResult.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
					if !lastResult.IsNil() {
						err = lastResult.Interface().(error)
					}
				}
			}
			t.EndTrace(trace.TraceID, outputs, err)
		}

		return results
	})

	return wrapper.Interface()
}

// StartTrace 开始追踪
func (t *Tracer) StartTrace(methodName string, input []interface{}) *MethodTrace {
	if !t.enabled {
		return nil
	}

	// 获取调用信息
	pc, file, line, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	packageName := ""
	if fn != nil {
		packageName = fn.Name()
	}

	// 生成TraceID
	traceID := fmt.Sprintf("trace_%d_%d", time.Now().UnixNano(), getGoroutineID())

	// 获取父追踪
	goroutineID := getGoroutineID()
	t.mu.RLock()
	parent := t.activeTraces[goroutineID]
	t.mu.RUnlock()

	parentID := ""
	if parent != nil {
		parentID = parent.TraceID
	}

	// 创建追踪
	trace := &MethodTrace{
		TraceID:     traceID,
		ParentID:    parentID,
		MethodName:  methodName,
		PackageName: packageName,
		FileName:    file,
		LineNumber:  line,
		Input:       input,
		StartTime:   time.Now(),
		Goroutine:   int(goroutineID),
		Status:      "running",
		Metadata:    make(map[string]interface{}),
	}

	// 捕获调用栈
	if t.captureStack {
		trace.CallStack = captureCallStack(10)
	}

	// 保存追踪
	t.mu.Lock()
	t.traces[traceID] = trace
	t.activeTraces[goroutineID] = trace
	t.mu.Unlock()

	// 如果有父追踪，添加为子节点
	if parent != nil {
		parent.Children = append(parent.Children, trace)
	}

	log.Printf("🔍 [Tracer] 开始追踪: %s (ID: %s)", methodName, traceID)

	return trace
}

// EndTrace 结束追踪
func (t *Tracer) EndTrace(traceID string, output []interface{}, err error) {
	if !t.enabled || traceID == "" {
		return
	}

	t.mu.Lock()
	trace, exists := t.traces[traceID]
	if !exists {
		t.mu.Unlock()
		return
	}

	trace.EndTime = time.Now()
	trace.Duration = trace.EndTime.Sub(trace.StartTime)
	trace.Output = output

	if err != nil {
		trace.Error = err.Error()
		trace.Status = "error"
	} else {
		trace.Status = "success"
	}

	// 恢复父追踪为当前活动追踪
	goroutineID := getGoroutineID()
	if trace.ParentID != "" {
		if parent, ok := t.traces[trace.ParentID]; ok {
			t.activeTraces[goroutineID] = parent
		}
	} else {
		delete(t.activeTraces, goroutineID)
	}

	t.mu.Unlock()

	log.Printf("✅ [Tracer] 完成追踪: %s (耗时: %v)", trace.MethodName, trace.Duration)

	// 调用处理器
	for _, handler := range t.handlers {
		go handler(trace)
	}
}

// GetTrace 获取追踪信息
func (t *Tracer) GetTrace(traceID string) *MethodTrace {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.traces[traceID]
}

// GetAllTraces 获取所有追踪
func (t *Tracer) GetAllTraces() []*MethodTrace {
	t.mu.RLock()
	defer t.mu.RUnlock()

	traces := make([]*MethodTrace, 0, len(t.traces))
	for _, trace := range t.traces {
		// 只返回顶层追踪（没有父节点的）
		if trace.ParentID == "" {
			traces = append(traces, trace)
		}
	}
	return traces
}

// ClearTraces 清除所有追踪
func (t *Tracer) ClearTraces() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.traces = make(map[string]*MethodTrace)
	t.activeTraces = make(map[int64]*MethodTrace)
}

// shouldTrace 判断是否应该追踪
func (t *Tracer) shouldTrace(methodName string) bool {
	// 检查过滤器
	for _, filter := range t.filters {
		if !filter(methodName) {
			return false
		}
	}
	return true
}

// TraceMethod 追踪方法（手动方式）
func TraceMethod(methodName string) func() {
	tracer := GetTracer()
	trace := tracer.StartTrace(methodName, nil)

	return func() {
		tracer.EndTrace(trace.TraceID, nil, nil)
	}
}

// TraceMethodWithArgs 追踪方法（带参数）
func TraceMethodWithArgs(methodName string, args ...interface{}) (traceID string, endFunc func(results ...interface{})) {
	tracer := GetTracer()
	trace := tracer.StartTrace(methodName, args)

	endFunc = func(results ...interface{}) {
		var err error
		// 检查最后一个结果是否是error
		if len(results) > 0 {
			if e, ok := results[len(results)-1].(error); ok {
				err = e
			}
		}
		tracer.EndTrace(trace.TraceID, results, err)
	}

	return trace.TraceID, endFunc
}

// WithContext 在context中传递TraceID
func WithContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, "trace_id", traceID)
}

// FromContext 从context中获取TraceID
func FromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return ""
}

// 辅助函数

// getGoroutineID 获取goroutine ID
func getGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// 格式: "goroutine 123 [running]:"
	var id int64
	fmt.Sscanf(string(buf[:n]), "goroutine %d", &id)
	return id
}

// captureCallStack 捕获调用栈
func captureCallStack(maxDepth int) []string {
	stack := make([]string, 0, maxDepth)
	for i := 3; i < maxDepth+3; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
		}
	}
	return stack
}

// reflectValuesToInterfaces 将reflect.Value转换为interface{}
func reflectValuesToInterfaces(values []reflect.Value) []interface{} {
	result := make([]interface{}, len(values))
	for i, v := range values {
		if v.IsValid() && v.CanInterface() {
			result[i] = v.Interface()
		} else {
			result[i] = fmt.Sprintf("<invalid: %v>", v)
		}
	}
	return result
}

// SerializeValue 序列化值（处理复杂类型）
func SerializeValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}
		return SerializeValue(val.Elem().Interface())
	case reflect.Struct:
		// 尝试JSON序列化
		if data, err := json.Marshal(v); err == nil {
			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err == nil {
				return result
			}
		}
		return fmt.Sprintf("%+v", v)
	case reflect.Map, reflect.Slice, reflect.Array:
		if data, err := json.Marshal(v); err == nil {
			var result interface{}
			if err := json.Unmarshal(data, &result); err == nil {
				return result
			}
		}
		return fmt.Sprintf("%+v", v)
	default:
		return v
	}
}

// Enable 启用追踪器
func (t *Tracer) Enable() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.enabled = true
}

// Disable 禁用追踪器
func (t *Tracer) Disable() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.enabled = false
}

// IsEnabled 检查是否启用
func (t *Tracer) IsEnabled() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.enabled
}

