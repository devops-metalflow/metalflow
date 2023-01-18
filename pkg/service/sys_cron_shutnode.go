package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rfyiamcool/cronlib"
	"gorm.io/gorm"
	"io"
	"metalflow/models"
	"metalflow/pkg/async"
	"metalflow/pkg/cron"
	"metalflow/pkg/global"
	"metalflow/pkg/grpc"
	"metalflow/pkg/request"
	"strings"
	"sync"
	"time"
)

func (s *MysqlService) CreateCronShutNode(req *request.CreateCronShutNodeRequest) error {
	nodes := make([]*models.SysNode, 0)

	err := s.TX.Model(&models.SysNode{}).Where("id in (?)", req.NodeIds).Find(&nodes).Error
	if err != nil {
		return err
	}
	cronShutNode := &models.SysCronShutNode{
		Name:      req.Name,
		StartTime: req.StartTime,
		ShutTime:  req.ShutTime,
		Keyword:   req.Keyword,
		Status:    (*uint)(req.Status),
		Creator:   req.Creator,
		Nodes:     nodes,
	}
	// 如果定时任务状态为正常，则添加定时开关机任务
	if *cronShutNode.Status == models.SysCronShutNodeEnable {
		jobNodes := &JobNodes{Nodes: nodes}
		var startModel, shutModel *cronlib.JobModel
		startModel, err = cronlib.NewJobModel(req.StartTime, jobNodes.RunStartTask)
		if err != nil {
			global.Log.Errorf("获取定时开机任务%s的model错误：%v", req.Name, err)
		}
		startJob := &cron.DynamicJob{
			JobName: req.Keyword + ".start",
			Job:     startModel,
		}
		global.Cron.Start <- startJob
		shutModel, err = cronlib.NewJobModel(req.ShutTime, jobNodes.RunShutTask)
		if err != nil {
			global.Log.Errorf("获取定时关机任务%s的model错误：%v", req.Name, err)
		}
		shutJob := &cron.DynamicJob{
			JobName: req.Keyword + ".shut",
			Job:     shutModel,
		}
		global.Cron.Start <- shutJob
	}
	return s.TX.Create(cronShutNode).Error
}

func (s *MysqlService) GetCronShutNode(req *request.ListCronShutNodeRequest) ([]models.SysCronShutNode, error) {
	var err error
	list := make([]models.SysCronShutNode, 0)
	query := s.TX.
		Model(&models.SysCronShutNode{}).
		Preload("Nodes").
		Order("created_at DESC")
	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	// 查询列表
	err = s.Find(query, &req.PageInfo, &list)
	return list, err
}

func (s *MysqlService) UpdateCronShutNodeById(shutId uint, req *request.UpdateCronShutNodeRequest) (err error) {
	// 更新定时开关机任务
	var csn models.SysCronShutNode
	query := s.TX.Where("id = ?", shutId).First(&csn)
	if query.Error == gorm.ErrRecordNotFound {
		return fmt.Errorf("record does not exist, update failed")
	}
	nodes := make([]*models.SysNode, 0)
	err = s.TX.Where("id in (?)", req.NodeIds).Find(&nodes).Error
	if err != nil {
		return
	}

	// 根据机器状态判断是否需要停用
	if *csn.Status != uint(*req.Status) && uint(*req.Status) == models.SysCronShutNodeDisable {
		// 停用当前的定时任务
		stopShutJob := &cron.DynamicJob{JobName: csn.Keyword + ".shut"}
		stopStartJob := &cron.DynamicJob{JobName: csn.Keyword + ".start"}
		global.Cron.Stop <- stopShutJob
		global.Cron.Stop <- stopStartJob
	} else {
		// 更新并启动定时任务
		jobNodes := &JobNodes{Nodes: nodes}
		var startModel, shutModel *cronlib.JobModel
		startModel, err = cronlib.NewJobModel(req.StartTime, jobNodes.RunStartTask)
		if err != nil {
			global.Log.Errorf("获取定时开机任务%s的model错误：%v", req.Name, err)
		}
		startJob := &cron.DynamicJob{
			JobName: csn.Keyword + ".start",
			Job:     startModel,
		}

		shutModel, err = cronlib.NewJobModel(req.ShutTime, jobNodes.RunShutTask)
		if err != nil {
			global.Log.Errorf("获取定时关机任务%s的model错误：%v", req.Name, err)
		}
		shutJob := &cron.DynamicJob{
			JobName: csn.Keyword + ".shut",
			Job:     shutModel,
		}
		global.Cron.Update <- startJob
		global.Cron.Update <- shutJob
	}

	// 更新普通字段
	c := &models.SysCronShutNode{
		Name:      req.Name,
		StartTime: req.StartTime,
		ShutTime:  req.ShutTime,
		Status:    (*uint)(req.Status),
	}
	err = query.Updates(c).Error
	if err != nil {
		return err
	}
	// 更新机器节点
	if len(req.NodeIds) > 0 {
		// 更新机器节点对应的labels
		return s.TX.Model(&csn).Association("Nodes").Replace(nodes)
	}
	return err
}

func (s *MysqlService) DeleteCronShutTaskByIds(ids []uint) error {
	cronShutNodes := make([]*models.SysCronShutNode, 0)
	err := s.TX.Model(&models.SysCronShutNode{}).Where("id in (?)", ids).Find(&cronShutNodes).Error
	if err != nil {
		return err
	}
	// 停掉运行中的定时任务
	for _, shutNode := range cronShutNodes {
		if *shutNode.Status != models.SysCronShutNodeEnable {
			continue
		}
		stopShutJob := &cron.DynamicJob{JobName: shutNode.Keyword + ".shut"}
		stopStartJob := &cron.DynamicJob{JobName: shutNode.Keyword + ".start"}
		global.Cron.Stop <- stopShutJob
		global.Cron.Stop <- stopStartJob
	}
	err = s.DeleteByIds(ids, new(models.SysCronShutNode))
	return err
}

type JobNodes struct {
	Nodes []*models.SysNode
}

const (
	shutShellName  = "shut.sh"
	startShellName = "start.sh"
)

type addressInfo struct {
	address string
	id      uint
}

// RunStartTask 因路由及同网段限制，因此只能是相同网段去发送wol命令来执行远程唤醒
func (j *JobNodes) RunStartTask() {
	global.Log.Info("开始执行定时开机任务...")
	var metaltask models.SysWorker
	err := global.Mysql.Model(new(models.SysWorker)).Where("id = ?", 2).First(&metaltask).Error //nolint:gomnd
	if err != nil {
		global.Log.Errorf("search worker metaltask from database error:%v", err)
		return
	}
	grpcPort := metaltask.Port

	var wg sync.WaitGroup
	for _, sysNode := range j.Nodes {
		wg.Add(1)
		go func(node *models.SysNode, port int) {
			global.Log.Infof("远程唤醒:%s开始...", node.Address)
			defer wg.Done()
			// 解析该节点的 mac地址, 组成可执行远程开机唤醒命令的shell文件
			var metric grpc.Metric
			er := json.Unmarshal(node.Information, &metric)
			if er != nil {
				global.Log.Errorf("解析节点：%s的metric失败：%v", node.Address, er)
				return
			}
			startShellInfo := &cronShellInfo{
				content: fmt.Sprintf("#!/bin/bash\nwakeonlan %s", metric.Mac),
			}
			fileMetric := grpc.FileMetric{
				FilePath:   startShellName,
				RemoteDir:  remoteDir,
				IsRunnable: true,
				FileGetter: startShellInfo,
			}
			// 获取该ip地址相邻的同网段信息
			addressData := strings.Split(node.Address, ".")
			netSegment := fmt.Sprintf("%s.%s.%s.", addressData[0], addressData[1], addressData[2])
			// 查询该ip的其他同网段ip
			sameNetSegments := make([]*models.SysNode, 0)
			er = global.Mysql.Model(&models.SysNode{}).
				Where("address LIKE ? AND address != ?", fmt.Sprintf("%%%s%%", netSegment), node.Address).
				Find(&sameNetSegments).Error
			if er != nil {
				global.Log.Errorf("查询%s的同网段ip失败", node.Address)
				sendStartMail(node.Address, fmt.Sprintf("查询%s的同网段ip失败", node.Address))
				return
			}
			if len(sameNetSegments) == 0 {
				global.Log.Errorf("未找到%s的同网段ip，无法远程唤醒%s开机", node.Address, node.Address)
				sendStartMail(node.Address, fmt.Sprintf("未找到%s的同网段ip，无法远程唤醒开机", node.Address))
				return
			}
			// 获取最快连通metaltask的ip的id，然后发送远程唤醒命令的shell文件，使用带缓存channel，避免goroutine泄漏
			nodeChan := make(chan *addressInfo, len(sameNetSegments))
			for _, n := range sameNetSegments {
				go func(in chan<- *addressInfo, address string, nodeId uint, port int, ctx context.Context) {
					_, e := grpc.ConnectGrpc(address, port, ctx)
					if e != nil {
						return
					}
					in <- &addressInfo{
						address: address,
						id:      nodeId,
					}
				}(nodeChan, n.Address, n.Id, port, context.Background())
			}
			select {
			case ai := <-nodeChan:
				s := New(nil)
				ids := []uint{ai.id}
				global.Log.Infof("使用同网段[%s]执行远程唤醒:%s...", ai.address, node.Address)
				e := s.BatchUploadByIds(fileMetric, ids)
				if e != nil {
					global.Log.Errorf("使用同网段[%s]远程唤醒:%s失败:%v", ai.address, node.Address, e)
					sendStartMail(node.Address, fmt.Sprintf("使用同网段[%s]远程唤醒:%s失败:%v", ai.address, node.Address, e))
				}
				global.Log.Infof("使用同网段[%s]远程唤醒:%s完成...", ai.address, node.Address)
			case <-time.After(time.Duration(global.Conf.System.ConnectTimeout) * time.Second):
				global.Log.Errorf("%s的所有相邻网段连接grpc超时，无法执行远程唤醒...", node.Address)
				sendStartMail(node.Address, fmt.Sprintf("%s的所有相邻网段连接grpc超时，无法执行远程唤醒...", node.Address))
			}
		}(sysNode, grpcPort)
	}
	wg.Wait()
	global.Log.Info("定时开机任务执行结束...")
}

func sendStartMail(address, content string) {
	// 发送邮件
	mail := &async.Mail{
		Title:     fmt.Sprintf("<开机失败>服务器[%s]远程唤醒失败，请知悉并处理！", address),
		Body:      content,
		Receivers: getEmailAddr(global.Conf.Mail.Cc),
		CC:        nil,
	}
	global.Machinery.SendMailTask(mail)
}

func getEmailAddr(usernames []string) []string {
	addrs := make([]string, 0)
	for _, u := range usernames {
		addrs = append(addrs, u+global.Conf.Mail.Suffix)
	}
	return addrs
}

func (j *JobNodes) RunShutTask() {
	global.Log.Info("开始执行定时关机任务...")
	ids := make([]uint, 0)
	for _, node := range j.Nodes {
		ids = append(ids, node.Id)
	}
	shutShellInfo := &cronShellInfo{
		content: "#!/bin/bash\npoweroff",
	}
	fileMetric := grpc.FileMetric{
		FilePath:   shutShellName,
		RemoteDir:  remoteDir,
		IsRunnable: true,
		FileGetter: shutShellInfo,
	}
	s := New(nil)
	err := s.BatchUploadByIds(fileMetric, ids)
	if err != nil {
		global.Log.Errorf("定时关机任务执行失败：%v", err)
	}
	global.Log.Info("定时关机任务执行结束...")
}

type cronShellInfo struct {
	content string
}

func (c *cronShellInfo) GetFile() (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(c.content)), nil
}
