// Tencent is pleased to support the open source community by making Polaris available.
//
// Copyright (C) 2026 THL A29 Limited, a Tencent company. All rights reserved.
//
// Licensed under the BSD 3-Clause License (the "License");
// you may not use this file except in compliance with the License.
//
// 用途：
//   验证 polaris-server-api MR 170（feat: ScorerConfig 新增 AggregateMode）的 proto 改动
//   能否被 Go 代码正常引用与编解码。
//
//   MR: https://git.woa.com/polaris/polaris-server-api/-/merge_requests/170
//   依赖：GoLearn 顶层 go.mod 的 replace 指向本地 polaris-server-api 仓库（已切到 MR 分支）。

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"git.woa.com/polaris/polaris-server-api/api/schedule"
	"google.golang.org/protobuf/proto"
)

func main() {
	// 1) 直接引用 MR 新增的 enum 值，验证类型/常量是否存在且可用
	cases := []schedule.AggregateMode{
		schedule.AggregateMode_AGGREGATE_MODE_UNSPECIFIED,
		schedule.AggregateMode_AGGREGATE_MODE_NORMALIZE_THEN_WEIGHT,
		schedule.AggregateMode_AGGREGATE_MODE_WEIGHT_THEN_NORMALIZE,
	}
	fmt.Println("== AggregateMode 枚举值 ==")
	for _, m := range cases {
		fmt.Printf("  %-2d %s\n", int32(m), m.String())
	}

	// 2) 构造一个完整的 ScorerConfig，引用 MR 新增的 aggregate_mode 字段
	scorer := &schedule.ScorerConfig{
		Weight: 1.0,
		Metrics: []*schedule.MetricScorerItem{
			{
				Metric: "cpu_usage",
				Mode:   schedule.MetricMode_METRIC_MODE_LOWER_IS_BETTER,
				Weight: 0.6,
			},
			{
				Metric: "qps",
				Mode:   schedule.MetricMode_METRIC_MODE_HIGHER_IS_BETTER,
				Weight: 0.4,
			},
		},
		// 关键：使用 MR 新增的字段
		AggregateMode: schedule.AggregateMode_AGGREGATE_MODE_WEIGHT_THEN_NORMALIZE,
	}

	fmt.Println("\n== 原始 ScorerConfig ==")
	fmt.Printf("  weight=%v\n", scorer.GetWeight())
	fmt.Printf("  aggregate_mode=%s (%d)\n", scorer.GetAggregateMode().String(), scorer.GetAggregateMode())
	fmt.Printf("  metrics=%d\n", len(scorer.GetMetrics()))

	// 3) proto 二进制编解码：验证字段编号 3 的 wire 编/解
	raw, err := proto.Marshal(scorer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "proto.Marshal 失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nMarshal 字节数 = %d, hex = %x\n", len(raw), raw)

	got := &schedule.ScorerConfig{}
	if err := proto.Unmarshal(raw, got); err != nil {
		fmt.Fprintf(os.Stderr, "proto.Unmarshal 失败: %v\n", err)
		os.Exit(1)
	}

	if got.GetAggregateMode() != scorer.GetAggregateMode() {
		fmt.Fprintf(os.Stderr, "FAIL: aggregate_mode 不一致: want=%v got=%v\n",
			scorer.GetAggregateMode(), got.GetAggregateMode())
		os.Exit(1)
	}
	if got.GetWeight() != scorer.GetWeight() || len(got.GetMetrics()) != len(scorer.GetMetrics()) {
		fmt.Fprintf(os.Stderr, "FAIL: 其他字段不一致 got=%+v\n", got)
		os.Exit(1)
	}

	// 4) JSON 输出，肉眼确认字段名是否为 aggregateMode/aggregate_mode
	jsonBytes, _ := json.MarshalIndent(got, "", "  ")
	fmt.Println("\n== Unmarshal 后的 ScorerConfig (JSON) ==")
	fmt.Println(string(jsonBytes))

	// 5) 边界：未显式赋值时，默认值应为 AGGREGATE_MODE_UNSPECIFIED(=0)
	defaultScorer := &schedule.ScorerConfig{Weight: 0.5}
	if defaultScorer.GetAggregateMode() != schedule.AggregateMode_AGGREGATE_MODE_UNSPECIFIED {
		fmt.Fprintf(os.Stderr, "FAIL: 默认 aggregate_mode 应为 UNSPECIFIED, got=%v\n",
			defaultScorer.GetAggregateMode())
		os.Exit(1)
	}

	fmt.Println("\nPASS: MR 170 的 proto 改动可以被正常引用、编解码与读取默认值")
}
