package main

import (
	"fmt"
	"reflect"
	"testing"
)

// 反射基础api
func TestReflection_test(tt *testing.T) {
	var x float64 = 3.4

	// 使用reflect.TypeOf()获取变量x的类型信息
	t := reflect.TypeOf(x)
	fmt.Println("Type:", t)

	// 使用reflect.ValueOf()获取变量x的值信息
	v := reflect.ValueOf(x)
	fmt.Println("Value:", v)

	// 反射可以用来检查变量的类型和值
	fmt.Println("Type is float64:", t.Kind() == reflect.Float64)
	fmt.Println("Value is 3.4:", v.Float() == 3.4)

	// 反射还可以用来修改变量的值
	// 注意：为了修改变量的值，需要传递变量的指针
	p := reflect.ValueOf(&x) // 注意：这里传递的是x的地址
	vp := p.Elem()           // Elem()用于获取指针指向的变量的Value
	vp.SetFloat(7.1)         // 修改变量的值
	fmt.Println("x after modification:", x)
}

func TestReflection2_test(_ *testing.T) {
	type Service struct {
		namespace string
		name      string
	}

	var i interface{}
	s := &Service{
		namespace: "test",
		name:      "lzb_test",
	}
	i = s

	t := reflect.TypeOf(i) // 获取类型信息：Service
	fmt.Println("type: ", t)

	// 获取一个切片类型：[]*Service
	sliceType := reflect.SliceOf(t)
	fmt.Println("type:", sliceType)

	/*
		作用：分配一个新的零值，类型为 []*Service，返回指向该值的指针的 reflect.Value
		内存分配：相当于 make([]*Service, 0) 但通过反射完成
		返回值：是一个 reflect.Value，但它实际上是一个指针（比如 *[]*dao.Service）
	*/
	slicePtr := reflect.New(sliceType)

	/*
		作用：将 reflect.Value 转换回普通的 Go 接口值（interface{}）
		原理：reflect.Value 是反射世界的表示，Interface() 将其转换回普通 Go 值
		用途：GORM 的 Find() 方法接受 interface{} 参数，这里传入切片指针让 GORM 填充数据
	*/
	interfacee := slicePtr.Interface()
	fmt.Println("interfacee", interfacee)

	/*
		查询数据库，将查询结果写入slicePtr，做后通过下面的slice对象(slice := slicePtr.Elem())查询数据
	*/
	//q := db.Model(queryConfig.Model)
	//q.Find(interfacee)

	/*
		作用：获取指针指向的实际值，对指针类型的 reflect.Value 进行解引用
		示例：如果 slicePtr 是 *[]*Service，slicePtr.Elem() 就是 []*Service
	*/
	slice := slicePtr.Elem()
	fmt.Println("slice: ", slice)

	/*
		将从db查询的结果(类型是reflect.Value)转换回普通的Go值( []interface{} )
	*/
	result := make([]interface{}, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		result[i] = slice.Index(i).Interface()
	}
}

// 通过反射判断一个对象是否包含字段名为"Name"的字段
func TestReflection3_test(_ *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	u := &User{Name: "张三", Age: 18}

	// 获取 reflect.Type
	t := reflect.TypeOf(u)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 通过反射判断结构体是否包含字段名为 "Name" 的字段
	field, ok := t.FieldByName("Name")
	if ok {
		fmt.Printf("找到字段: %s, 类型: %s\n", field.Name, field.Type)
	} else {
		fmt.Println("未找到字段 Name")
	}
}
