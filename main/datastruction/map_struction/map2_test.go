package map_struction

import (
	"fmt"
	"testing"
)

type OuterMaps struct {
	InnerMaps map[int]string
}

var outerMaps = OuterMaps{}

func TestMap2(t *testing.T) {
	maps := make(map[int]string)

	maps[1] = "1"
	outerMaps.InnerMaps = maps

	maps = make(map[int]string, len(maps)+1) // make 操作创建了一个全新的 map

	fmt.Printf("%+v", outerMaps)
}
