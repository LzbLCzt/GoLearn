package test

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net"
	"testing"
	"time"
	"unsafe"
)

type Person struct {
	name string
	age  int
}

func TestPrac(t *testing.T) {
	x := Person{name: "aaa", age: 10}
	ptr := unsafe.Pointer(&x)
	off := uintptr(ptr) + unsafe.Offsetof(x.name)
	agePtr := (*string)(unsafe.Pointer(off))
	*agePtr = "ccc"
	fmt.Println(x)
}

func TestPrac2(t *testing.T) {
	g, ctx := errgroup.WithContext(context.Background())

	for i := 1; i <= 3; i++ {
		j := i
		g.Go(func() error {
			return worker(ctx, j)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("err: ", err)
	} else {
		fmt.Println("success")
	}
}

func worker(ctx context.Context, i int) error {
	fmt.Printf("worker %d, start\n", i)
	select {
	case <-time.After(time.Duration(i) * time.Second):
		fmt.Printf("worker %d finished \n", i)
	case <-ctx.Done():
		fmt.Printf("worker %d cancelled\n", i)
		return ctx.Err()
	}
	if i == 2 {
		return fmt.Errorf("err occur")
	}
	return nil
}

func TestGetIPByName(t *testing.T) {
	ips, err := GetIPByName("eth1")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ips)
	}
}

// 根据网卡名获取该网卡的ip列表
func GetIPByName(name string) ([]string, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(addrs))

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ret = append(ret, ipnet.IP.String())
			}
		}
	}

	return ret, nil
}
