package main

import (
	poll "GoLearn/rainbow/pollConfig"
)

func main() {
	//获取七彩石配置-kv类型
	//poll.FetchConfigKV()

	//获取七彩石配置-table类型
	//poll.FetchConfigTable()

	//获取七彩石配置-kv类型 && 监听配置变更
	poll.FetchConfigKVAndListen()
}
