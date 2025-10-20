package main

import (
	"fmt"
	"time"

	// 导入auto-visualizer，自动启动追踪功能
	_ "github.com/Ryan-myp/auto-visualizer-service"
	autovisualizer "github.com/Ryan-myp/auto-visualizer-service"
)

// User 用户结构
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// UserService 用户服务
type UserService struct {
	users map[int]*User
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{
		users: make(map[int]*User),
	}
}

// CreateUser 创建用户（方式1：使用 defer + Begin）
func (s *UserService) CreateUser(name string, age int) (*User, error) {
	// 开始追踪，自动记录入参和出参
	end := autovisualizer.Begin("UserService.CreateUser", name, age)
	
	var user *User
	var err error
	
	defer func() {
		end(user, err)
	}()
	
	// 模拟业务逻辑
	time.Sleep(50 * time.Millisecond)
	
	id := len(s.users) + 1
	user = &User{
		ID:   id,
		Name: name,
		Age:  age,
	}
	
	s.users[id] = user
	
	fmt.Printf("✅ 创建用户成功: %+v\n", user)
	return user, nil
}

// GetUser 获取用户（方式2：使用 Measure 只测量时间）
func (s *UserService) GetUser(id int) (*User, error) {
	defer autovisualizer.Measure("UserService.GetUser")()
	
	time.Sleep(20 * time.Millisecond)
	
	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("用户不存在: %d", id)
	}
	
	fmt.Printf("✅ 获取用户成功: %+v\n", user)
	return user, nil
}

// UpdateUser 更新用户（方式3：使用 TraceMethod 简单追踪）
func (s *UserService) UpdateUser(id int, name string, age int) error {
	defer autovisualizer.TraceMethod("UserService.UpdateUser")()
	
	time.Sleep(30 * time.Millisecond)
	
	user, exists := s.users[id]
	if !exists {
		return fmt.Errorf("用户不存在: %d", id)
	}
	
	user.Name = name
	user.Age = age
	
	fmt.Printf("✅ 更新用户成功: %+v\n", user)
	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int) error {
	end := autovisualizer.Begin("UserService.DeleteUser", id)
	defer end(nil)
	
	time.Sleep(15 * time.Millisecond)
	
	if _, exists := s.users[id]; !exists {
		return fmt.Errorf("用户不存在: %d", id)
	}
	
	delete(s.users, id)
	fmt.Printf("✅ 删除用户成功: ID=%d\n", id)
	return nil
}

// ListUsers 列出所有用户
func (s *UserService) ListUsers() []*User {
	defer autovisualizer.Measure("UserService.ListUsers")()
	
	time.Sleep(10 * time.Millisecond)
	
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	
	fmt.Printf("✅ 列出用户: 共 %d 个\n", len(users))
	return users
}

// ProcessUserBatch 批量处理用户（演示嵌套调用追踪）
func (s *UserService) ProcessUserBatch(names []string) error {
	end := autovisualizer.Begin("UserService.ProcessUserBatch", names)
	defer end(nil)
	
	fmt.Printf("\n🔄 开始批量处理 %d 个用户...\n", len(names))
	
	for i, name := range names {
		// 这里会产生嵌套的追踪调用
		user, err := s.CreateUser(name, 20+i)
		if err != nil {
			return err
		}
		
		// 再次调用，形成调用链
		_, err = s.GetUser(user.ID)
		if err != nil {
			return err
		}
	}
	
	// 列出所有用户
	s.ListUsers()
	
	fmt.Printf("✅ 批量处理完成\n\n")
	return nil
}

// 高阶函数示例：使用装饰器模式
func calculateSum(a, b int) int {
	time.Sleep(10 * time.Millisecond)
	return a + b
}

// 包装后的函数
var tracedCalculateSum = autovisualizer.Trace("calculateSum", calculateSum).(func(int, int) int)

func main() {
	fmt.Println("🚀 Auto-Visualizer 方法追踪演示")
	fmt.Println("=" + string(make([]byte, 50)) + "=")
	fmt.Println()
	
	// 创建服务
	userService := NewUserService()
	
	// 示例1：单个方法调用
	fmt.Println("📝 示例1：单个方法调用")
	user1, _ := userService.CreateUser("张三", 25)
	userService.GetUser(user1.ID)
	userService.UpdateUser(user1.ID, "张三丰", 26)
	
	fmt.Println()
	
	// 示例2：批量处理（嵌套调用）
	fmt.Println("📝 示例2：批量处理（嵌套调用）")
	userService.ProcessUserBatch([]string{"李四", "王五", "赵六"})
	
	// 示例3：使用装饰器包装的函数
	fmt.Println("📝 示例3：装饰器模式")
	result := tracedCalculateSum(10, 20)
	fmt.Printf("✅ 计算结果: %d\n\n", result)
	
	// 示例4：删除用户
	fmt.Println("📝 示例4：删除操作")
	userService.DeleteUser(1)
	
	fmt.Println()
	
	// 打印追踪统计
	printTraceStats()
	
	// 保持程序运行，以便查看Web UI
	fmt.Println("\n🌐 访问 http://localhost:8090 查看可视化界面")
	fmt.Println("🌐 访问 http://localhost:8090/api/method-traces 查看追踪数据")
	fmt.Println("🌐 访问 http://localhost:8090/api/method-traces/tree 查看调用树")
	fmt.Println("\n按 Ctrl+C 退出...")
	
	select {}
}

// printTraceStats 打印追踪统计
func printTraceStats() {
	traces := autovisualizer.GetAllTraces()
	
	fmt.Println("📊 追踪统计:")
	fmt.Printf("   总追踪数: %d\n", len(traces))
	
	for _, trace := range traces {
		printTrace(trace, 0)
	}
}

// printTrace 打印追踪信息
func printTrace(trace interface{}, level int) {
	// 由于类型限制，这里简化处理
	indent := ""
	for i := 0; i < level; i++ {
		indent += "  "
	}
	fmt.Printf("%s└─ 追踪记录\n", indent)
}

