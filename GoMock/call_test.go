package GoMock

import (
	"errors"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
)

// TestUserProvider_CallMethods 展示gomock.Call主要方法的使用
func TestUserProvider_CallMethods(t *testing.T) {
	// 1. 创建控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 确保验证所有期望

	// 2. 创建Mock对象
	mockUserProvider := NewMockUserProvider(ctrl)

	// 3. 设置各种调用期望
	// 3.1 基本返回值设置 (.Return())
	mockUserProvider.EXPECT().
		GetUser(1).
		Return(&User{ID: 1, Name: "Alice"}, nil)

	// 3.2 指定调用次数 (.Times())
	mockUserProvider.EXPECT().
		GetUser(2).
		Return(&User{ID: 2, Name: "Bob"}, nil).
		Times(2) // 期望被调用2次

	// 3.3 任意调用次数 (.AnyTimes())
	mockUserProvider.EXPECT().
		GetUser(gomock.Any()). // 匹配任何参数
		Return(&User{ID: 999, Name: "Default"}, nil).
		AnyTimes() // 可以调用任意次（包括0次）

	// 4. 执行测试调用
	fmt.Println("=== 测试调用开始 ===")

	// 调用1: 基本返回
	user1, err1 := mockUserProvider.GetUser(1)
	if err1 != nil {
		t.Errorf("GetUser(1)失败: %v", err1)
	} else {
		fmt.Printf("调用1成功: %+v\n", user1)
	}

	// 调用2: 指定次数（调用2次）
	user2, _ := mockUserProvider.GetUser(2)
	fmt.Printf("调用2-1成功: %+v\n", user2)

	user2, _ = mockUserProvider.GetUser(2) // 第二次调用
	fmt.Printf("调用2-2成功: %+v\n", user2)

	// 调用3: 任意参数匹配
	userAny, _ := mockUserProvider.GetUser(999) // 匹配gomock.Any()
	fmt.Printf("调用任意参数成功: %+v\n", userAny)

	fmt.Println("=== 测试调用结束 ===")

	// 注意：ctrl.Finish()会在defer中自动验证所有期望
}

// TestUserProvider_CallOrder 测试调用顺序
func TestUserProvider_CallOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProvider := NewMockUserProvider(ctrl)

	// 设置调用顺序期望
	firstCall := mockUserProvider.EXPECT().
		GetUser(1).
		Return(&User{ID: 1, Name: "First"}, nil)

	secondCall := mockUserProvider.EXPECT().
		GetUser(2).
		Return(&User{ID: 2, Name: "Second"}, nil).
		After(firstCall) // 必须在firstCall之后调用

	thirdCall := mockUserProvider.EXPECT().
		GetUser(3).
		Return(&User{ID: 3, Name: "Third"}, nil).
		After(secondCall) // 必须在secondCall之后调用

	// 按正确顺序调用
	fmt.Println("=== 顺序调用测试 ===")
	mockUserProvider.GetUser(1)
	mockUserProvider.GetUser(2)
	mockUserProvider.GetUser(3)
	fmt.Println("顺序调用测试通过")

	// 验证调用对象
	_ = firstCall
	_ = secondCall
	_ = thirdCall
}

// TestUserProvider_DoMethods 测试Do和DoAndReturn方法
func TestUserProvider_DoMethods(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProvider := NewMockUserProvider(ctrl)
	callCount := 0

	// 使用.Do()执行额外操作
	mockUserProvider.EXPECT().
		GetUser(gomock.Any()).
		Do(func(id int) {
			callCount++
			fmt.Printf("Do回调执行，调用次数: %d, 参数: %d\n", callCount, id)
		}).
		Return(&User{ID: 100, Name: "From Do"}, nil).
		Times(2)

	// 使用.DoAndReturn()动态返回
	mockUserProvider.EXPECT().
		GetUser(gomock.Any()).
		DoAndReturn(func(id int) (*User, error) {
			if id > 0 {
				return &User{ID: id, Name: fmt.Sprintf("User%d", id)}, nil
			}
			return nil, errors.New("invalid id")
		}).
		AnyTimes()

	fmt.Println("=== Do方法测试 ===")

	// 测试.Do()
	user1, _ := mockUserProvider.GetUser(10)
	fmt.Printf("Do测试1: %+v\n", user1)

	user2, _ := mockUserProvider.GetUser(20)
	fmt.Printf("Do测试2: %+v\n", user2)

	// 测试.DoAndReturn()
	userDynamic, _ := mockUserProvider.GetUser(30)
	fmt.Printf("DoAndReturn测试: %+v\n", userDynamic)

	_, err := mockUserProvider.GetUser(-1)
	if err != nil {
		fmt.Printf("错误处理测试: %v\n", err)
	}
}
