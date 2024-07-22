package pollConfig

import (
	"fmt"
	"git.code.oa.com/rainbow/golang-sdk/types"
	"git.code.oa.com/rainbow/golang-sdk/v3/confapi"
	"git.code.oa.com/rainbow/golang-sdk/v3/watch"
	v3 "git.code.oa.com/rainbow/proto/api/configv3"
	"time"
)

// watchCallBackV3 监听回调函数
func watchCallBackV3(oldVal watch.Result, newVal []*v3.Item) error {
	fmt.Printf("\n---------------------\n")
	fmt.Printf("old value:%+v\n", oldVal)
	fmt.Printf("new value:")
	for i := 0; i < len(newVal); i++ {
		fmt.Printf("%+v", *newVal[i])
	}
	fmt.Printf("\n---------------------\n")

	return nil
}

func FetchConfigKVAndListen() {
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

	//拉取配置
	key := "confOne"
	val, err := rainbow.Get(key, getOpts...)
	if err != nil {
		// 必须处理ERR错误
		fmt.Printf("[rainbow.Get]%s\n", err.Error())
		return
	}
	// 需要进一步判断val是否为空，是否满足业务需要
	fmt.Printf("val: %s\n", val)

	var watch = watch.Watcher{
		GetOptions: types.GetOptions{
			AppID:   appid,
			EnvName: env,
			Group:   group,
		},
		CB: watchCallBackV3, //指定Call Back 函数
	}
	//添加监听
	rainbow.AddWatcher(watch, getOpts...)
	for {
		time.Sleep(5 * time.Second)
		fmt.Printf("[rainbow.watch] watching config...\n")
	}
}
