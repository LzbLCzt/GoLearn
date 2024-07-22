package main

//
//import (
//	"fmt"
//	"os"
//)
//
//func main() {
//	fmt.Println("这是关于go流程控制和函数相关api")
//
//	//if else
//	//可以在if条件中声明变量
//	if x := 100; x > 1000 {
//		fmt.Println("")
//	} else if x > 100 && x < 1000 {
//		fmt.Println("")
//	} else {
//		fmt.Println("")
//	}
//
//	////跳转执行跳转
//	i := 0
//JumpToHere:
//	println(i)
//	i++
//	goto JumpToHere //JumpToHere这个并非关键词
//
//	//for loop
//	sum := 0
//	for index := 0; index < 10; index++ {
//		sum += index
//	}
//	fmt.Println("sum is equal to ", sum)
//
//	//for loop可以当成while loop使用
//	sum1 := 1
//	for sum1 < 1000 {
//		sum1 += sum1
//	}
//	//遍历map
//	maps := map[string]int{"a": 1, "b": 2}
//	for k, v := range maps {
//		fmt.Println("key: ", k, " value: ", v)
//	}
//	//遍历for中声明的数组
//	for i := range [10]int{1, 2, 3} {
//		fmt.Print(i, " ")
//	}
//	fmt.Println()
//	//遍历数组
//	arr := [10]int{1, 2, 3, 4, 5}
//	for ele := range arr {
//		fmt.Print(ele, " ")
//	}
//	fmt.Println()
//
//	//switch case
//	fmt.Println("Go里面switch默认相当于每个case最后带有break，匹配成功后不会自动向下执行其他case，而是跳出整个switch, 但是可以使用fallthrough强制执行后面的case代码")
//	integer := 6
//	switch integer {
//	case 4:
//		fmt.Println("The integer was <= 4")
//		fallthrough
//	case 5:
//		fmt.Println("The integer was <= 5")
//		fallthrough
//	case 6:
//		fmt.Println("The integer was <= 6")
//		fallthrough
//	case 7:
//		fmt.Println("The integer was <= 7")
//		fallthrough
//	case 8:
//		fmt.Println("The integer was <= 8")
//		fallthrough
//	default:
//		fmt.Println("default case")
//	}
//}
//
//// 函数
//// 函数返回值可以多个
//func SumAndProduct(a int, b int) (int, int) {
//	return a + b, a * b
//}
//
//// 返回值可以指定返回具体哪个变量
//func SumAndProduct2(a, b int) (add int, multiplied int) {
//	add = a + b
//	multiplied = a * b
//	return
//}
//
//// 不定参数
//func func1(arg ...int) {
//	fmt.Println("在函数体中，变量arg是一个int的slice")
//}
//
///*
//传值 or 传指针：
//默认情况下变量a作为形参传入函数时，会copy一个新的变量a_copy，函数对a_copy操作并不会影响变量a，这是传值
//如果想在函数中改变函数外声明的变量a，需要将变量a的地址值传给函数，这就是传指针
//*/
//func func2(a *int) int {
//	*a = *a + 100
//	return *a
//}
//
//// defer 关键词：你可以在函数中添加多个defer语句。当函数执行到最后时，这些defer语句会按照逆序执行，最后该函数返回
//func ReadWrite() bool {
//	file, _ := os.Open("file")
//	defer file.Close()
//	if 1 == 1 {
//		return false
//	}
//	if 2 == 2 {
//		return false
//	}
//	return true
//}
//
//// 如果有很多调用defer，那么defer是采用后进先出模式
//func func3() {
//	for i := 0; i < 5; i++ {
//		defer fmt.Printf("%d ", i)
//	}
//}
