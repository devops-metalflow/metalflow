package cron

import (
	"fmt"
	"github.com/rfyiamcool/cronlib"
)

type DynamicJob struct {
	JobName string
	Job     *cronlib.JobModel
}

// InitJob 服务启动时需要执行的job结构体
type InitJob struct {
	Spec    string
	Handler func()
}

type Client struct {
	Cron     *cronlib.CronSchduler
	InitJobs map[string]*InitJob
	Start    chan *DynamicJob
	Stop     chan *DynamicJob
	Update   chan *DynamicJob
}

func NewCron() *Client {
	return &Client{
		Cron:     cronlib.New(),
		InitJobs: make(map[string]*InitJob),
		Start:    make(chan *DynamicJob),
		Stop:     make(chan *DynamicJob),
		Update:   make(chan *DynamicJob),
	}
}

func (c *Client) DoInitJobs() error {
	for name, initJob := range c.InitJobs {
		model, err := cronlib.NewJobModel(initJob.Spec, initJob.Handler)
		if err != nil {
			return err
		}
		err = c.Cron.Register(name, model)
		if err != nil {
			return err
		}
		fmt.Printf("添加初始化任务[%s]成功\n", name)
	}
	return nil
}

func (c *Client) Run() {
	c.Cron.Start()
	c.Cron.Wait()
}
