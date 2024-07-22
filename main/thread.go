package main

//func say(s string) {
//	for i := 0; i < 5; i++ {
//		runtime.Gosched()
//		fmt.Println(s)
//	}
//}
//
//func main() {
//	runtime.GOMAXPROCS(10) //告诉调度器同时使用多个线程
//	go say("world")        //开一个新的Goroutines执行
//	say("hello")           //当前Goroutines执行
//}

//func main() {
//	c := make(chan int, 2) //修改2为1就报错，修改2为3可以正常运行
//	c <- 1
//	c <- 2
//	fmt.Println(<-c)
//	fmt.Println(<-c)
//}

//修改为1报如下的错误:
//fatal error: all goroutines are asleep - deadlock!
