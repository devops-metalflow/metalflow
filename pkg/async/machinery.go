package async

import (
	"context"
	"fmt"

	machinery "github.com/RichardKnop/machinery/v1"
	taskConfig "github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

// Machinery 初始化Machinery服务
type Machinery struct {
	// 获取当前Machinery服务端
	MachineryServer *machinery.Server
	// 当前会话
	Ctx context.Context
}

// GetMachinery 获取Machinery实例
func GetMachinery(redisIp, redisPasswd string, redisPort, redisDatabase int) *Machinery {
	ctx := context.Background()
	// 初始化.
	cnf := &taskConfig.Config{
		Broker:        fmt.Sprintf("redis://%s@%s:%d/%d", redisPasswd, redisIp, redisPort, redisDatabase),
		DefaultQueue:  "ServerTasksQueue",
		ResultBackend: fmt.Sprintf("redis://%s@%s:%d/%d", redisPasswd, redisIp, redisPort, redisDatabase),
	}
	server, err := machinery.NewServer(cnf)
	if err != nil {
		panic(fmt.Sprintf("获取machinery实例失败: %v", err))
	}
	return &Machinery{
		MachineryServer: server,
		Ctx:             ctx,
	}
}

// NewAsyncTaskWorker 创建任务的消费者worker
func (m *Machinery) NewAsyncTaskWorker(concurrency int) *machinery.Worker {
	consumerTag := "TaskWorker"
	worker := m.MachineryServer.NewWorker(consumerTag, concurrency)
	errorHandler := func(err error) {
		fmt.Println("执行失败: ", err)
	}
	preTaskHandler := func(signature *tasks.Signature) {
		fmt.Println("开始执行: ", signature.Name)
	}
	postTaskHandler := func(signature *tasks.Signature) {
		fmt.Println("执行结束: ", signature.Name)
	}
	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetErrorHandler(errorHandler)
	worker.SetPreTaskHandler(preTaskHandler)
	return worker
}

// RegisterTasks 将任务的参数进行签名并注册任务
func (m *Machinery) RegisterTasks(asyncTaskMap map[string]any) error {
	// 将任务进行注册
	return m.MachineryServer.RegisterTasks(asyncTaskMap)
}
