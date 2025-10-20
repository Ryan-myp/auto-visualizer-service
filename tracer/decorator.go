package tracer

import (
	"fmt"
	"reflect"
)

// Wrap 包装函数以支持自动追踪
// 使用示例:
//   tracedFunc := tracer.Wrap("MyMethod", myFunc).(func(int) (string, error))
//   result, err := tracedFunc(123)
func Wrap(methodName string, fn interface{}) interface{} {
	tracer := GetTracer()
	return tracer.Trace(methodName, fn)
}

// WrapStruct 包装结构体的所有方法
func WrapStruct(obj interface{}) interface{} {
	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
		objValue = objValue.Elem()
	}

	if objType.Kind() != reflect.Struct {
		panic("WrapStruct: obj must be a struct or pointer to struct")
	}

	// 创建新的结构体类型（带追踪的方法）
	// 注意：由于Go的限制，这里我们返回一个代理对象
	return createProxy(obj)
}

// createProxy 创建代理对象
func createProxy(obj interface{}) interface{} {
	// 这是一个简化版本，实际使用中可能需要更复杂的代理逻辑
	return obj
}

// Auto 自动追踪装饰器（使用defer）
// 使用示例:
//   func MyMethod(a int) (string, error) {
//       defer tracer.Auto("MyMethod", a)()
//       // 业务逻辑
//       return "result", nil
//   }
func Auto(methodName string, args ...interface{}) func(...interface{}) {
	// 添加 panic 保护
	defer func() {
		if r := recover(); r != nil {
			// 静默处理，不影响业务
		}
	}()

	tracer := GetTracer()
	if tracer == nil {
		return func(...interface{}) {} // 返回空函数
	}

	trace := tracer.StartTrace(methodName, args)
	if trace == nil {
		return func(...interface{}) {} // 返回空函数
	}

	return func(results ...interface{}) {
		defer func() {
			if r := recover(); r != nil {
				// 静默处理
			}
		}()

		var err error
		if len(results) > 0 {
			if e, ok := results[len(results)-1].(error); ok {
				err = e
			}
		}
		tracer.EndTrace(trace.TraceID, results, err)
	}
}

// Track 简化的追踪装饰器（自动检测返回值）
// 使用示例:
//   func MyMethod(a int) (result string, err error) {
//       defer Track("MyMethod", a)(&result, &err)
//       // 业务逻辑
//       result = "success"
//       return
//   }
func Track(methodName string, args ...interface{}) func(...interface{}) {
	tracer := GetTracer()
	trace := tracer.StartTrace(methodName, args)

	return func(resultPtrs ...interface{}) {
		results := make([]interface{}, len(resultPtrs))
		var err error

		for i, ptr := range resultPtrs {
			if ptr == nil {
				continue
			}

			ptrValue := reflect.ValueOf(ptr)
			if ptrValue.Kind() == reflect.Ptr && !ptrValue.IsNil() {
				value := ptrValue.Elem()
				if value.IsValid() && value.CanInterface() {
					results[i] = value.Interface()

					// 检查是否是error类型
					if e, ok := results[i].(error); ok && e != nil {
						err = e
					}
				}
			}
		}

		tracer.EndTrace(trace.TraceID, results, err)
	}
}

// Begin 开始追踪（返回结束函数）
// 使用示例:
//   end := tracer.Begin("MyMethod", arg1, arg2)
//   defer end(result1, result2, err)
func Begin(methodName string, args ...interface{}) func(...interface{}) {
	return Auto(methodName, args...)
}

// Measure 测量函数执行时间
// 使用示例:
//   defer tracer.Measure("MyMethod")()
func Measure(methodName string) func() {
	// 添加 panic 保护
	defer func() {
		if r := recover(); r != nil {
			// 静默处理
		}
	}()

	tracer := GetTracer()
	if tracer == nil {
		return func() {} // 返回空函数
	}

	trace := tracer.StartTrace(methodName, nil)
	if trace == nil {
		return func() {} // 返回空函数
	}

	return func() {
		defer func() {
			if r := recover(); r != nil {
				// 静默处理
			}
		}()
		tracer.EndTrace(trace.TraceID, nil, nil)
	}
}

// TraceFunc 追踪函数（泛型辅助）
type TracedFunc0 func()
type TracedFunc1[R any] func() R
type TracedFunc2[R any] func() (R, error)
type TracedFunc1In1Out[T, R any] func(T) R
type TracedFunc1In2Out[T, R any] func(T) (R, error)

// WrapFunc0 包装无参数无返回值函数
func WrapFunc0(methodName string, fn func()) TracedFunc0 {
	return func() {
		defer Measure(methodName)()
		fn()
	}
}

// WrapFunc1 包装无参数单返回值函数
func WrapFunc1[R any](methodName string, fn func() R) TracedFunc1[R] {
	return func() R {
		end := Begin(methodName)
		result := fn()
		end(result)
		return result
	}
}

// WrapFunc2 包装无参数双返回值函数（含error）
func WrapFunc2[R any](methodName string, fn func() (R, error)) TracedFunc2[R] {
	return func() (R, error) {
		end := Begin(methodName)
		result, err := fn()
		end(result, err)
		return result, err
	}
}

// WrapFunc1In1Out 包装单参数单返回值函数
func WrapFunc1In1Out[T, R any](methodName string, fn func(T) R) TracedFunc1In1Out[T, R] {
	return func(arg T) R {
		end := Begin(methodName, arg)
		result := fn(arg)
		end(result)
		return result
	}
}

// WrapFunc1In2Out 包装单参数双返回值函数（含error）
func WrapFunc1In2Out[T, R any](methodName string, fn func(T) (R, error)) TracedFunc1In2Out[T, R] {
	return func(arg T) (R, error) {
		end := Begin(methodName, arg)
		result, err := fn(arg)
		end(result, err)
		return result, err
	}
}

// MethodInterceptor 方法拦截器接口
type MethodInterceptor interface {
	Before(methodName string, args []interface{})
	After(methodName string, results []interface{}, err error)
}

// InterceptorChain 拦截器链
type InterceptorChain struct {
	interceptors []MethodInterceptor
}

// NewInterceptorChain 创建拦截器链
func NewInterceptorChain() *InterceptorChain {
	return &InterceptorChain{
		interceptors: []MethodInterceptor{},
	}
}

// Add 添加拦截器
func (c *InterceptorChain) Add(interceptor MethodInterceptor) {
	c.interceptors = append(c.interceptors, interceptor)
}

// ExecuteBefore 执行前置拦截
func (c *InterceptorChain) ExecuteBefore(methodName string, args []interface{}) {
	for _, interceptor := range c.interceptors {
		interceptor.Before(methodName, args)
	}
}

// ExecuteAfter 执行后置拦截
func (c *InterceptorChain) ExecuteAfter(methodName string, results []interface{}, err error) {
	for i := len(c.interceptors) - 1; i >= 0; i-- {
		c.interceptors[i].After(methodName, results, err)
	}
}

// LoggingInterceptor 日志拦截器
type LoggingInterceptor struct{}

func (l *LoggingInterceptor) Before(methodName string, args []interface{}) {
	fmt.Printf("→ Calling %s with args: %v\n", methodName, args)
}

func (l *LoggingInterceptor) After(methodName string, results []interface{}, err error) {
	if err != nil {
		fmt.Printf("← %s returned error: %v\n", methodName, err)
	} else {
		fmt.Printf("← %s returned: %v\n", methodName, results)
	}
}

// TimingInterceptor 计时拦截器
type TimingInterceptor struct {
	startTimes map[string]int64
}

func NewTimingInterceptor() *TimingInterceptor {
	return &TimingInterceptor{
		startTimes: make(map[string]int64),
	}
}

func (t *TimingInterceptor) Before(methodName string, args []interface{}) {
	t.startTimes[methodName] = GetTracer().StartTrace(methodName, args).StartTime.UnixNano()
}

func (t *TimingInterceptor) After(methodName string, results []interface{}, err error) {
	// 计时在tracer中已经处理
}

