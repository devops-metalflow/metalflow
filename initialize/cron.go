package initialize

import (
	probing "github.com/prometheus-community/pro-bing"
	"metalflow/models"
	"metalflow/pkg/cron"
	"metalflow/pkg/global"
	"metalflow/pkg/service"
	"sync"
	"time"
)

// Cron 初始化定时任务
func Cron() {
	c := cron.NewCron()
	go func(c *cron.Client) {
		for {
			select {
			case startJob := <-c.Start:
				err := c.Cron.DynamicRegister(startJob.JobName, startJob.Job)
				if err != nil {
					global.Log.Errorf("动态添加定时任务[%s]失败：%v", startJob.JobName, err)
				}
				global.Log.Infof("动态添加定时任务[%s]成功", startJob.JobName)
			case stopJob := <-c.Stop:
				c.Cron.StopService(stopJob.JobName)
				global.Log.Infof("移除定时任务[%s]成功", stopJob.JobName)
			case updateJob := <-c.Update:
				err := c.Cron.UpdateJobModel(updateJob.JobName, updateJob.Job)
				if err != nil {
					global.Log.Errorf("更新定时任务[%s]失败：%v", updateJob.JobName, err)
				}
				global.Log.Infof("更新定时任务[%s]成功", updateJob.JobName)
			}
		}
	}(c)
	// 添加初始启动时的定时任务并运行
	go func(c *cron.Client) {
		addRefreshNodeMetricsTask(c)
		addRefreshNodePingStatsTask(c)
		addShutStartNodeTask(c)
		err := c.DoInitJobs()
		if err != nil {
			panic("执行初始化定时任务失败")
		}
		c.Run()
	}(c)
	global.Cron = c
	global.Log.Debug("初始化定时任务完成")
}

const refreshNodeMetricsTask = "refresh.node.metrics.10m"

func addRefreshNodeMetricsTask(c *cron.Client) {
	if global.Conf.System.NodeMetricsCronTask != "" {
		c.InitJobs[refreshNodeMetricsTask] = &cron.InitJob{
			Spec:    global.Conf.System.NodeMetricsCronTask,
			Handler: runRefreshNodeMetrics,
		}
	}
}

func runRefreshNodeMetrics() {
	global.Log.Info("[定时任务][机器节点信息刷新]准备开始...")
	// 获取所有状态正常的机器节点
	nodes := make([]models.SysNode, 0)
	err := global.Mysql.Model(new(models.SysNode)).Find(&nodes).Error
	if err != nil {
		global.Log.Error("查询数据库机器节点失败：", err)
		return
	}
	for _, node := range nodes { //nolint:gocritic
		// 判断本次刷新时间与上次刷新时间的间隔，如果小于5分钟，则不进行刷新
		if time.Since(node.RefreshLastTime.Time).Minutes() < 5 { //nolint:gomnd
			global.Log.Debugf("五分钟内已有其他人刷新节点[%s]，跳过本次刷新", node.Address)
			continue
		}
		// 启动异步任务刷新机器节点信息
		var worker models.SysWorker
		err = global.Mysql.Model(new(models.SysWorker)).Where("id = ?", 1).First(&worker).Error
		if err != nil {
			global.Log.Error("search worker from database failed")
			return
		}
		global.Machinery.SendGrpcTask(node.Address, worker.Port, worker.ServiceReq)
		// 计入刷新时间与刷新次数
		err = global.Mysql.Model(new(models.SysNode)).Where("address = ?", node.Address).
			Updates(map[string]any{
				"refresh_count":     *node.RefreshCount + 1,
				"refresh_last_time": time.Now(),
			}).Error
		if err != nil {
			global.Log.Errorf("更新机器节点[%s]信息失败", node.Address)
		}
	}
	global.Log.Info("[定时任务][机器节点信息刷新]任务结束")
}

func addShutStartNodeTask(c *cron.Client) {
	cronShutNodeTasks := make([]*models.SysCronShutNode, 0)
	err := global.Mysql.Model(&models.SysCronShutNode{}).Preload("Nodes").
		Where("status = ?", 1).Find(&cronShutNodeTasks).Error
	if err != nil {
		global.Log.Errorf("查询定时开关机任务失败：%v", err)
		return
	}
	if len(cronShutNodeTasks) == 0 {
		return
	}
	for _, task := range cronShutNodeTasks {
		if len(task.Nodes) == 0 {
			continue
		}
		jobNodes := &service.JobNodes{
			Nodes: task.Nodes,
		}
		// 添加定时开机任务
		c.InitJobs[task.Keyword+".start"] = &cron.InitJob{
			Spec:    task.StartTime,
			Handler: jobNodes.RunStartTask,
		}
		// 添加定时关机任务
		c.InitJobs[task.Keyword+".shut"] = &cron.InitJob{
			Spec:    task.ShutTime,
			Handler: jobNodes.RunShutTask,
		}
	}
}

// Add cron ping servers task
const refreshNodePingStatsName = "refresh.node.ping.1m"

func addRefreshNodePingStatsTask(c *cron.Client) {
	if global.Conf.System.NodePingCronTask != "" {
		c.InitJobs[refreshNodePingStatsName] = &cron.InitJob{
			Spec:    global.Conf.System.NodePingCronTask,
			Handler: DoPingIps,
		}
		global.Log.Infof("Enable refresh ping status scheduled task [%s] successfully", global.Conf.System.NodePingCronTask)
	}
}

type ServerStats struct {
	Ip     string
	Status bool
}

func DoPingIps() {
	nodes := make([]*models.SysNode, 0)
	err := global.Mysql.Model(&models.SysNode{}).Find(&nodes).Error
	if err != nil {
		global.Log.Errorf("search nodes failed: %v", err)
		return
	}

	wg := sync.WaitGroup{}
	serverStatsChan := make(chan ServerStats)
	for _, node := range nodes {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			stat, msg := ping(addr)
			global.Log.Infof("Ping %s %s(%v)\n", addr, msg, stat)
			serverStat := ServerStats{
				Ip:     addr,
				Status: stat,
			}
			serverStatsChan <- serverStat
		}(node.Address)
	}
	go func() {
		wg.Wait()
		close(serverStatsChan)
	}()

	for s := range serverStatsChan {
		// update server ping stat in database
		var p uint
		if s.Status {
			p = 1
		}
		err = global.Mysql.Model(&models.SysNode{}).Where("address = ?", s.Ip).Update("ping_stat", p).Error
		if err != nil {
			global.Log.Errorf("Update %s ping stat failed: %v", s.Ip, err)
		}
	}
}

func ping(ip string) (result bool, message string) {
	pinger, err := probing.NewPinger(ip)
	if err != nil {
		return false, err.Error()
	}
	// ping timeout
	pinger.Timeout = 10 * time.Second
	pinger.Count = 3
	pinger.SetPrivileged(true)
	err = pinger.Run()
	if err != nil {
		return false, err.Error()
	}
	stats := pinger.Statistics()

	result = true
	message = "Succeeded!"
	// check packets loss rate, if it is greater than 20%, return false
	if stats.PacketLoss > 20 { // nolint:gomnd
		result = false
		message = "Failed!"
	}
	return result, message
}
