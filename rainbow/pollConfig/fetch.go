package pollConfig

import (
	"fmt"
	"git.code.oa.com/rainbow/golang-sdk/types"
	"git.code.oa.com/rainbow/golang-sdk/v3/confapi"
)

func FetchConfigKV() {
	appid := "2dea1ac9-5cd6-484c-a6f9-f0c54d464e77"
	env := "Test"
	group := "learn.kv"

	rainbow, err := confapi.New(
		// 北极星就近路由
		types.ConnectStr("polaris://65026305:65536"),
		// 开启本地内存缓存
		types.IsUsingLocalCache(true),
		// 开启文件缓存（必须同时开启本地内存缓存），七彩石服务端不可用时，可以支持程序正常重启
		types.IsUsingFileCache(true),
	)
	if err != nil {
		fmt.Printf("[confapi.New]%s\n", err.Error())
		return
	}

	getOpts := make([]types.AssignGetOption, 0)
	getOpts = append(getOpts, types.WithAppID(appid))
	getOpts = append(getOpts, types.WithGroup(group))
	getOpts = append(getOpts, types.WithEnvName(env))

	_, err = rainbow.GetGroup(getOpts...)
	if err != nil { // 必须处理ERR错误，用于感知七彩石服务异常
		fmt.Printf("[rainbow.Get]%s\n", err.Error())
		return
	}

	key1 := "confOne"
	val1, err := rainbow.Get(key1, getOpts...)
	if err != nil {
		// 必须处理ERR错误
		fmt.Printf("[rainbow.Get]%s\n", err.Error())
		return
	}

	fmt.Printf("val1: %s\n", val1)

	key2 := "confSecond"
	val2, err := rainbow.Get(key2, getOpts...)
	if err != nil {
		// 必须处理ERR错误
		fmt.Printf("[rainbow.Get]%s\n", err.Error())
		return
	}

	fmt.Printf("val2: %s\n", val2)

	key3 := "confThird"
	val3, err := rainbow.Get(key3, getOpts...)
	if err != nil {
		// 必须处理ERR错误
		fmt.Printf("[rainbow.Get]%s\n", err.Error())
		return
	}

	fmt.Printf("val3: %s\n", val3)
}

type Actor struct {
	actorId    int
	firstName  string
	lastName   string
	lastUpdate string
}

func FetchConfigTable() {
	appid := "2dea1ac9-5cd6-484c-a6f9-f0c54d464e77"
	env := "Test"
	group := "learn.table"

	rainbow, err := confapi.New(
		// 北极星就近路由
		types.ConnectStr("polaris://65026305:65536"),
		// 开启本地内存缓存
		types.IsUsingLocalCache(true),
		// 开启文件缓存（必须同时开启本地内存缓存），七彩石服务端不可用时，可以支持程序正常重启
		types.IsUsingFileCache(true),
	)
	if err != nil {
		fmt.Printf("[confapi.New]%s\n", err.Error())
		return
	}

	getOpts := make([]types.AssignGetOption, 0)
	getOpts = append(getOpts, types.WithAppID(appid))
	getOpts = append(getOpts, types.WithGroup(group))
	getOpts = append(getOpts, types.WithEnvName(env))

	_, err = rainbow.GetGroup(getOpts...)
	if err != nil { // 必须处理ERR错误，用于感知七彩石服务异常
		fmt.Printf("[rainbow.Get]%s\n", err.Error())
		return
	}

	// get table: actor
	actors, err := rainbow.GetTable(getOpts...)
	if err != nil {
		// 必须处理ERR错误
		fmt.Printf("[rainbow.GetTable] failed %s\n", err.Error())
		return
	}
	fmt.Println(actors)
}
