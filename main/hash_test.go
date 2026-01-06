package main

import (
	"fmt"
	"sort"
	"testing"
)

/*
case_name: 测试 HashString
case_path: /common/utils
case_description: 测试 HashString
*/

var length = 500

func TestHashString(t *testing.T) {
	ips := mockIPs(10000 * length)
	for _, ip := range ips {
		fmt.Println(ip)
	}
	//fmt.Printf("ips: %v\n", ips)
	statistics := make([]int, length)

	for _, ip := range ips {
		h := HashString(ip)
		selectedIndex := h % uint32(length)
		statistics[selectedIndex] += 1
	}

	sort.Slice(statistics, func(i, j int) bool {
		return statistics[i] > statistics[j]
	})
}

// mock 常用的IPV4 地址，要求IP地址灵活，不能是连续的或前缀相同或后缀相同
func mockIPs(num int) []string {
	ips := make([]string, num)
	for i := 0; i < num; i++ {
		// 使用伪随机算法生成更灵活的IP地址
		// 避免连续和模式化的IP，确保分布均匀
		segment1 := (i*17 + 123) % 255
		segment2 := (i*23 + 456) % 255
		segment3 := (i*31 + 789) % 255
		segment4 := (i*41+101)%254 + 1 // 确保第四段不为0

		ips[i] = fmt.Sprintf("%d.%d.%d.%d", segment1, segment2, segment3, segment4)
	}
	return ips
}

// HashString 采用FNV-1a算法
func HashString(s string) uint32 {
	const (
		offset32 uint32 = 2166136261
		prime32  uint32 = 16777619
	)

	hash := offset32
	for i := 0; i < len(s); i++ {
		hash ^= uint32(s[i])
		hash *= prime32
	}
	return hash
}
