// Package initialize is used to initialize some configurations and global variables.
// such as mysql database, routing, log configuration, etc.
package initialize

import (
	"fmt"
	"metalflow/pkg/async"
	"metalflow/pkg/global"
	"metalflow/pkg/service"
)

// Async Task 初始化任务队列
func Async() {
	// 初始化machinery
	Machinery()
}

// Machinery 初始化任务队列Machinery
func Machinery() {
	if global.Redis == nil {
		global.Log.Info("未初始化redis, 将无法执行后台异步任务")
		return
	}
	// 获取machinery的server实例
	machinery := async.GetMachinery(
		global.Conf.Redis.Host,
		global.Conf.Redis.Password,
		global.Conf.Redis.Port,
		global.Conf.Redis.Database,
	)
	// 获取worker,启动任务消费者
	worker := machinery.NewAsyncTaskWorker(10) //nolint:gomnd
	// 这里一定要使用go协程启动异步任务在后台，不然会直接卡住，无法初始化
	go func() {
		err := worker.Launch()
		if err != nil {
			panic(fmt.Sprintf("启动machinery worker失败，%v", err))
		}
	}()

	// 将任务映射集注册到machinery中
	err := machinery.RegisterTasks(service.InitAsyncTaskMap())
	if err != nil {
		panic(fmt.Sprintf("注册异步任务失败：%v", err))
	}
	global.Machinery = machinery
	global.Log.Info("初始化任务队列: machinery完成")
}
