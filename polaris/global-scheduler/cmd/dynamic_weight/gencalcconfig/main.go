// Tencent is pleased to support the open source community by making Polaris available.
//
// Copyright (C) 2026 THL A29 Limited, a Tencent company. All rights reserved.
//
// Licensed under the BSD 3-Clause License (the "License");
// you may not use this file except in compliance with the License.

// gencalcconfig 用 pb 类型构造 dynamic-weight 的 calculate_config，
// 并以 JSON 字符串形式输出到 stdout，供 shell 脚本通过 $(...) 捕获。
//
// 设计说明:
//  1. 直接复用 polaris-server-api 中 schedule.CalculateConfig 等 pb 类型，
//     避免在 shell 里手写嵌套 JSON、字段名拼写错误、枚举 magic number 等问题。
//  2. 序列化使用 protojson.MarshalOptions{UseEnumNumbers: true}，让 mode/op
//     等枚举输出为数字（如 mode:2），与服务端历史协议保持一致。
package main

import (
	"fmt"
	"os"

	"git.woa.com/polaris/polaris-server-api/api/schedule"
	"google.golang.org/protobuf/encoding/protojson"
)

// defaultMetrics 默认参与评分的指标列表，权重均为 1，模式均为 HIGHER_IS_BETTER。
// 如需调整指标，直接改这里即可，无需关心 JSON 字段名/枚举数字。
var defaultMetrics = []*schedule.MetricScorerItem{
	//{
	//	Metric: "kv_cache_usage_perc",
	//	Weight: 0.2,
	//	Mode:   schedule.MetricMode_METRIC_MODE_LOWER_IS_BETTER,
	//},
	{
		Metric: "num_requests_running",
		Weight: 0.3,
		Mode:   schedule.MetricMode_METRIC_MODE_HIGHER_IS_BETTER,
	},
	{
		Metric: "num_requests_waiting",
		Weight: 0.7,
		Mode:   schedule.MetricMode_METRIC_MODE_LOWER_IS_BETTER,
	},
}

func main() {
	cfg := &schedule.CalculateConfig{
		Scorers: []*schedule.ScorerConfig{
			{
				Weight:        1,
				Metrics:       defaultMetrics,
				AggregateMode: schedule.AggregateMode_AGGREGATE_MODE_WEIGHT_THEN_NORMALIZE,
			},
		},
		Fallback: &schedule.FallbackStrategy{
			Type: schedule.FallbackType_FALLBACK_STATIC_WEIGHT,
		},
	}

	//cfg := &schedule.CalculateConfig{}

	// UseEnumNumbers: 让 mode 输出为数字 2，与服务端协议一致
	// EmitUnpopulated: false, 默认值字段不输出，保持 JSON 简洁
	marshaller := protojson.MarshalOptions{
		UseEnumNumbers:  true,
		EmitUnpopulated: false,
		UseProtoNames:   true, // 字段使用 proto 原名（snake_case），如 metric/weight/mode
	}
	data, err := marshaller.Marshal(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal calculate_config failed: %v\n", err)
		os.Exit(1)
	}

	// 直接打印到 stdout（不带换行），让 shell 用 $(...) 捕获完整 JSON
	fmt.Print(string(data))
}
