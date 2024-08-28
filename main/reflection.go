package main

//import (
//	"fmt"
//	"reflect"
//)
//
//func main() {
//	var x float64 = 3.4
//
//	// 使用reflect.TypeOf()获取变量x的类型信息
//	t := reflect.TypeOf(x)
//	fmt.Println("Type:", t)
//
//	// 使用reflect.ValueOf()获取变量x的值信息
//	v := reflect.ValueOf(x)
//	fmt.Println("Value:", v)
//
//	// 反射可以用来检查变量的类型和值
//	fmt.Println("Type is float64:", t.Kind() == reflect.Float64)
//	fmt.Println("Value is 3.4:", v.Float() == 3.4)
//
//	// 反射还可以用来修改变量的值
//	// 注意：为了修改变量的值，需要传递变量的指针
//	p := reflect.ValueOf(&x) // 注意：这里传递的是x的地址
//	vp := p.Elem()           // Elem()用于获取指针指向的变量的Value
//	vp.SetFloat(7.1)         // 修改变量的值
//	fmt.Println("x after modification:", x)
//}
