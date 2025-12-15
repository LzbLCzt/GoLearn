package Basic

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Config struct {
	Name        string         `yaml:"name"`
	ClusterInfo map[string]int `yaml:"clusterInfo"`
	OtherInfo   *int           `yaml:"otherInfo"`
}

var AppConf = &Config{
	ClusterInfo: map[string]int{},
}

// 初始化AppConf
func InitConf() {
	AppConf.Name = "test"
	AppConf.ClusterInfo = map[string]int{
		"aaa": 1,
		"bbb": 2,
	}
}

func InitConfV2() {
	AppConf.Name = "test"
	AppConf.ClusterInfo["aaa"] = 1
	AppConf.ClusterInfo["bbb"] = 2
}

type redisClientManager struct {
	clusterInfo map[string]int
	otherInfo   *int
}

var DefaultRedisClientManager = &redisClientManager{
	clusterInfo: AppConf.ClusterInfo,
	otherInfo:   AppConf.OtherInfo,
}

func TestMapWrong(t *testing.T) {
	m := DefaultRedisClientManager.clusterInfo
	fmt.Printf("before %+v", m)

	InitConf()
	/*这里DefaultRedisClientManager.clusterInfo依然是空，并没有因为InitConf()而改变
	虽然map是引用变量，但是InitConf是重新赋值了AppConf.ClusterInfo，对 Map 变量本身的赋值（=）操作会改变变量的引用指向，而不是在原有底层数据上进行修改
	*/
	assert.Equal(t, len(m), 0)
}

func TestMapRight(t *testing.T) {
	m := DefaultRedisClientManager.clusterInfo
	fmt.Printf("before %+v", m)

	InitConfV2()
	assert.Equal(t, len(m), 2)
}

/*
如果clusterInfo不是map，而是指针类型，也会这样吗：
会
*/

func InitConfV3() {
	AppConf.Name = "test"
	num := 1
	AppConf.OtherInfo = &num
}

func TestMapOtherTypeWrong(t *testing.T) {
	m := DefaultRedisClientManager.otherInfo
	fmt.Printf("before %+v", m)

	InitConfV3()
	assert.Nil(t, m)
}

func InitConfV4() {
	AppConf.Name = "test"
	num := 1
	DefaultRedisClientManager.otherInfo = &num
}
