package main

//
//import (
//	"fmt"
//	"math"
//)
//
///*
//method：指带有接收者的函数，method归属于接收者对应的类型
//method语法：func (r ReceiverType) funcName(parameters) (returnType) {}
//*/
//type Rectangle struct {
//	width, height float64
//}
//type Circle struct {
//	radius float64
//}
//
//// method命名可以重复，只要接收者不同即可
//func (r Rectangle) area() float64 { //r是接收者，这个method只归属于Rectangle
//	return r.width * r.height
//}
//func (c Circle) area() float64 { //r Circle是接收者
//	return c.radius * c.radius * math.Pi
//}
//func main() {
//	r1 := Rectangle{12, 2}
//	r2 := Rectangle{9, 4}
//	c1 := Circle{10}
//	c2 := Circle{25}
//
//	fmt.Println("Area of r1 is: ", r1.area())
//	fmt.Println("Area of r2 is: ", r2.area())
//	fmt.Println("Area of c1 is: ", c1.area())
//	fmt.Println("Area of c2 is: ", c2.area())
//}
