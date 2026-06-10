package test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
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

func test() string {
	defer func() {
		fmt.Println("defering")
	}()

	return "test"
}
func TestPrac3(t *testing.T) {
	test()
}

type A struct {
	Name *string
	Age  int
}

var DefaultA = A{}

func TestPrac4(t *testing.T) {
	fmt.Printf("DefaultA: %v\n", DefaultA)
	DefaultACopy := DefaultA
	n := "test"
	DefaultA.Name = &n
	fmt.Printf("DefaultA: %v\n", *DefaultA.Name)
	fmt.Printf("DefaultACopy: %v\n", DefaultACopy)
}

type Cluster struct {
	Name string `yaml:"name"`
}

type Config struct {
	Name        string             `yaml:"name"`
	ClusterInfo map[string]Cluster `yaml:"clusterInfo"`
}

var DefaultConfig = &Config{}

func TestPrac5(t *testing.T) {
	clusterInfo := DefaultConfig.ClusterInfo
	err := LoadConfigFromYAML("./config.yaml", DefaultConfig)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Printf("DefaultConfig: %v\n", *DefaultConfig)
	fmt.Printf("clusterInfo from source: %v\n", DefaultConfig.ClusterInfo)
	fmt.Printf("cfg: %v\n", clusterInfo)
}

func LoadConfigFromYAML(yamlFile string, config interface{}) error {
	if yamlFile == "" {
		return fmt.Errorf("empty file name")
	}

	fmt.Printf("[INFO] load config from yaml file %s\n", yamlFile)

	file, err := os.Open(yamlFile)
	if err != nil {
		return fmt.Errorf("[ERROR] open yaml file fail: %v\n", err)
	}
	defer file.Close()

	return yaml.NewDecoder(file).Decode(config)
}

type KVcache struct {
	Name  string
	Cache *sync.Map
}

type Data struct {
	ID   int
	Data string
}

func TestPrac6(t *testing.T) {
	cache := &KVcache{
		Name:  "abc",
		Cache: &sync.Map{},
	}

	data := &Data{
		ID:   1,
		Data: "aaa",
	}

	cache.Cache.Store("abc", data)

	value, ok := cache.Cache.Load("abc")
	if ok {
		if x, y := value.(*Data); y {
			copyData := *x
			copyData.ID = 2323

			// 查看copyData的类型
			fmt.Printf("copyData: %v\n", reflect.TypeOf(copyData))
			fmt.Printf("value: %v\n", reflect.TypeOf(value))
		}
	}
}

func TestPrac7(t *testing.T) {
	m := map[string]string{
		"k": "v",
	}

	fmt.Printf("m: %v\n", m)

	f := func(source map[string]string) {
		source["k1"] = "v1"
	}
	f(m)
	fmt.Printf("m: %v\n", m)
}

type CallStatOne struct {
	CallCount *bool
	IpCount   *bool
}

func TestPrac8(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	m := make(map[int]float64, len(arr))
	for _, idx := range arr {
		m[idx] = 0
	}
	fmt.Printf("m: %v\n", m)

	fmt.Println(float64(5) / float64(10))
}

// CallStat 调用统计
type CallStat struct {
	CallCount int
	IpCount   int
}

func (c *CallStat) String() string {
	return fmt.Sprintf("call count: %d, ip count: %d", c.CallCount, c.IpCount)
}

func CallCount(refreshStats map[string]*CallStat) {
	if refreshStats == nil {
		refreshStats = make(map[string]*CallStat)
	}
	refreshStats["test"] = &CallStat{
		CallCount: 1,
		IpCount:   1,
	}
	return
}

func CallCount2(x *CallStat) {
	if x == nil {
		x = &CallStat{
			CallCount: 1,
			IpCount:   1,
		}
		return
	}
	x.CallCount += 1
	x.IpCount += 1
}
