package initialize

import (
	"bytes"
	"encoding/json"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"gorm.io/gorm"
	"io"
	"metalflow/models"
	"metalflow/pkg/async"
	"metalflow/pkg/consul"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/service"
	"net/http"
	"regexp"
)

// Consul start consul watch
// here doesn't use service discover is to ensure that the service status is monitored in real time
func Consul() {
	address := fmt.Sprintf("%s:%d", global.Conf.Consul.Address, global.Conf.Consul.Port)
	watch, err := consul.NewRegistry(address)
	if err != nil {
		panic(fmt.Sprintf("initialize consul watch failed: %v", err))
	}
	// 放到一个goroutine中处理节点状态变化，并更新数据库
	go func() {
		for {
			select {
			case statusService := <-watch.StatusChan:
				UpdateStateByConsul(statusService)
			case shutdownServiceName := <-watch.ShutdownChan:
				err = global.Mysql.Model(&models.SysNode{}).Where("address = ?", shutdownServiceName).
					Update("health", models.SysNodeHealthShutdown).Error
				if err != nil {
					global.Log.Errorf("服务%s掉线了，但更新数据库失败: %v", shutdownServiceName, err)
				}
				// 发送邮件通知对应节点负责人
				if global.Conf.Mail.Host != "" {
					sendMail(shutdownServiceName)
				} else {
					global.Log.Errorf("服务器%s已掉线，未配置邮件服务器，无法发送邮件通知", shutdownServiceName)
				}
			}
		}
	}()
	err = watch.StartWatch()
	if err != nil {
		panic(fmt.Sprintf("initialize consul watch failed: %v", err))
	}
}

func UpdateStateByConsul(svc *consul.Service) {
	statusMap := map[string]uint{
		consulapi.HealthCritical: models.SysNodeHealthAbnormal,
		consulapi.HealthPassing:  models.SysNodeHealthNormal,
	}
	var node models.SysNode
	query := global.Mysql.Model(&models.SysNode{}).Where("address = ?", svc.Address).First(&node)
	if query.Error == nil {
		err := query.Update("health", statusMap[svc.Status]).Error
		if err != nil {
			global.Log.Errorf("服务%s变化:%v，但更新数据库失败: %v", svc.Address, svc.Status, err)
		}
	} else if query.Error == gorm.ErrRecordNotFound {
		normalHealth := models.SysNodeHealthNormal
		newNode := &request.CreateNodeRequestStruct{
			Address:     svc.Address,
			SshPort:     22,
			ServicePort: svc.Port,
			Health:      &normalHealth,
			Creator:     "系统自动创建",
		}
		s := service.New(nil)
		if err := s.CreateNode(newNode); err != nil {
			fmt.Printf("create node：%s failed: %v", svc.Address, err)
			return
		}
		// let metalbeat execute the command that needs to initialize the deployment of workers
		err := DeployInitWorkers(svc.ServerOs, svc.Address, svc.Port)
		if err != nil {
			global.Log.Errorf("deploy workers failed for address: %s, error: %v", svc.Address, err)
			return
		}
		// query the metrics worker and communicate according to its configuration
		// TODO: 因进行过worker初始化，所以id为1的worker就是metalmetrics
		var worker models.SysWorker
		err = global.Mysql.Model(new(models.SysWorker)).Where("id = ?", 1).First(&worker).Error
		if err != nil {
			global.Log.Error("search worker from database failed")
			return
		}
		// get node info from metalmetrics using async task
		global.Machinery.SendGrpcTask(newNode.Address, worker.Port, worker.ServiceReq)
	}
	// update os info.
	if svc.ServerOs != "" {
		err := global.Mysql.Model(&models.SysNode{}).Where("address = ?", svc.Address).
			Update("os", svc.ServerOs).Error
		if err != nil {
			global.Log.Errorf("update server [%s] os info failed: %v", svc.Address, err)
		}
	}
}

func DeployInitWorkers(serverOs, address string, port int) error {
	workers := make([]models.SysWorker, 0)
	err := global.Mysql.Model(new(models.SysWorker)).Where("auto_deploy = ?", 1).Find(&workers).Error
	if err != nil {
		return err
	}
	for _, worker := range workers { //nolint:gocritic
		// send cmd for starting deploy workers.
		var cmdContent *models.CmdStruct
		cmdContent, err = worker.GetCmd(serverOs)
		if err != nil {
			global.Log.Errorf("Failed to get worker: [%s] deploy commands", worker.Name)
		}
		err = sendCmd2MetalBeat(address, cmdContent.Download, port)
		if err != nil {
			global.Log.Errorf("Failed to run worker: [%s download] command, err: %v", worker.Name, err)
		}
		// run stop command before start.
		err = sendCmd2MetalBeat(address, cmdContent.Stop, port)
		if err != nil {
			global.Log.Errorf("Failed to run worker: [%s stop] command, err: %v", worker.Name, err)
		}
		err = sendCmd2MetalBeat(address, cmdContent.Start, port)
		if err != nil {
			global.Log.Errorf("Failed to run worker: [%s start] command, err: %v", worker.Name, err)
		}
	}
	return nil
}

func sendCmd2MetalBeat(address, cmd string, port int) error {
	var statusOk = 201

	requestBody := map[string]string{
		"cmd": cmd,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal request body failed. address: %s, err: %v", address, err)
	}

	shellUrl := fmt.Sprintf("http://%s:%d/shell", address, port)
	req, err := http.NewRequest(http.MethodPost, shellUrl, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to request. address: %s, err: %v", address, err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req) // nolint:bodyclose
	if err != nil {
		return fmt.Errorf("failed to post. address: %s, err: %v", address, err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status to send metalbeat:%s cmd", address)
	}
	ret, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read metalbeat:%s resp", address)
	}
	metalBeatResp := MetalBeatResp{}
	err = json.Unmarshal(ret, &metalBeatResp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal metalbeat:%s resp", address)
	}
	if metalBeatResp.Code != statusOk {
		return fmt.Errorf("send cmd to metalbeat:%s failed, err: %s", address, metalBeatResp.Msg)
	}
	return nil
}

type MetalBeatResp struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func sendMail(address string) {
	var node models.SysNode
	err := global.Mysql.Model(&models.SysNode{}).Where("address = ?", address).First(&node).Error
	if err != nil {
		global.Log.Errorf("数据库查询节点：%s失败：%v", address, err)
		return
	}
	// 根据机器节点负责人获取对应邮件收件人
	// nolint:gocritic
	re, _ := regexp.Compile("[^0-9]")
	manager := re.ReplaceAllString(node.Manager, "")

	receivers := global.Conf.Mail.Cc
	if manager != "" {
		receivers = []string{manager}
	}
	// 发送邮件
	mail := &async.Mail{
		Title:     fmt.Sprintf("<节点告警>服务器[%s]异常，请知悉并处理！", address),
		Body:      fmt.Sprintf("%s状态异常，请确认该机器是否已关机。\n若未关机，请确认metalbeat服务是否正常，谢谢！", address),
		Receivers: getEmailAddr(receivers),
		CC:        getEmailAddr(global.Conf.Mail.Cc),
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
