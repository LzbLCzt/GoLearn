package util

import (
	"fmt"
)

func MockIPs(num int) []string {
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
