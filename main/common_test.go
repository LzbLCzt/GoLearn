package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Echo(s string) {
	for i := 0; i < 3; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func Test_1(t *testing.T) {
	go Echo("协程执行")
	Echo("主线程执行")
}

func Test_2(t *testing.T) {
	m := make(map[int]int, 2000)
	//for循环100000次，统计每个数字出现频率
	for i := 0; i <= 100000; i++ {
		num := rand.Intn(2000)
		m[num]++
	}
	fmt.Println(m)
}

type RouterInfo struct {
	Ip   string
	Name string
}

func Test_3(t *testing.T) {
	ips := make([]string, 1, 2)

	routers := make(map[string]RouterInfo, 2)
	fmt.Println(len(ips))
	routers["1"] = RouterInfo{Name: "1"}
	routers["2"] = RouterInfo{Name: "2"}

	for _, router := range routers {
		ips = append(ips, router.Ip)
	}

	//fmt.Println(len(ips))
	//fmt.Println(ips)
}

type Person2 struct {
	Name            string
	Age             int
	lastProduceTime time.Time
}

func Test_4(t *testing.T) {
	p := Person2{}
	fmt.Println(p.lastProduceTime)
	fmt.Println(time.Since(p.lastProduceTime) > 6*time.Hour)
}
