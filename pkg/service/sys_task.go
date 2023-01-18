package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"metalflow/models"
	"metalflow/pkg/async"
	"metalflow/pkg/global"
	"metalflow/pkg/grpc"
	"metalflow/pkg/grpc/securepb"
	"regexp"
	"strconv"
	"strings"

	"github.com/bndr/gojenkins"
	"github.com/go-mail/mail"
)

// 创建任务,所有的任务都可以写到这里
const (
	JenkinsUrl      string = "http://10.89.34.100:8080"           // jenkins的url地址
	JenkinsUser     string = "devops-app"                         // 用户名
	JenkinsPassword string = "1106a776bc163a8d8884ff7bd7befd416b" //nolint:gosec // 密码
)

// JenkinsTask 创建Jenkins任务，使用go来构建Jenkins的job
func JenkinsTask(jobName, jobParams string, ctx context.Context) (err error) {
	jenkins := gojenkins.CreateJenkins(nil, JenkinsUrl, JenkinsUser, JenkinsPassword)
	_, err = jenkins.Init(ctx)
	if err != nil {
		fmt.Println("连接Jenkins失败:", err)
		return
	}

	var mapParams map[string]string // 将对应的jobParams转化为map结构

	err = json.Unmarshal([]byte(jobParams), &mapParams)
	if err != nil {
		fmt.Println("转化job参数失败:", err)
		return
	}

	// 开始构建一次对应的Job
	var buildId int64
	buildId, err = jenkins.BuildJob(ctx, jobName, mapParams)
	if err != nil {
		fmt.Println("执行Jenkins任务失败:", err)
		return
	}
	fmt.Printf("执行了任务:%s的第%d次构建", jobName, buildId)
	return
}

// MailTask 邮件发送任务
func MailTask(title, body string, receivers, cc []string) (err error) {
	msg := mail.NewMessage()
	// Set E-Mail sender
	msg.SetAddressHeader("From", global.Conf.Mail.From+global.Conf.Mail.Suffix, global.Conf.Mail.Header)
	// Set E-Mail subject
	msg.SetHeader("Subject", title)
	// Set E-Mail receivers
	msg.SetHeader("To", receivers...)
	// Set E-Mail cc
	msg.SetHeader("Cc", cc...)
	// Set E-Mail body. You can set plain text or html with text/html
	msg.SetBody("text/plain", body)

	// Settings for SMTP server
	dialer := mail.NewDialer(
		global.Conf.Mail.Host,
		global.Conf.Mail.Port,
		global.Conf.Mail.Username,
		global.Conf.Mail.Password,
	)
	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	// nolint:gosec
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Now send E-Mail
	if err = dialer.DialAndSend(msg); err != nil {
		global.Log.Errorf("邮件发送任务失败：%v", err)
	}
	return err
}

// GrpcTask grpc通信的任务
func GrpcTask(address string, port int, message string) error {
	// 获取grpc通信后的响应
	response, err := grpc.SendFlow(address, port, message)
	if err != nil {
		return err
	}
	// 提取数据库需要的字段并更新数据库
	metric, err := grpc.ParseMetrics(response)
	if err != nil {
		return err
	}
	information, err := grpc.GetNodeInformation(&metric)
	if err != nil {
		return err
	}
	var metricStr string
	metricStr, err = grpc.GetModelNodeMetric(&metric)
	if err != nil {
		return err
	}
	// 保存数据库
	var node models.SysNode
	if err := global.Mysql.Model(&models.SysNode{}).Where("address = ?", address).First(&node).Error; err != nil {
		return err
	}
	// 根据规则计算性能
	var performance uint
	cupData := GetPerformanceByMetrics(metric.Cpu)
	ramData := GetPerformanceByMetrics(metric.Ram)
	score := (cupData*40 + ramData*10) / 100
	// nolint:gomnd
	if score < 19 {
		performance = 2
	} else if score >= 19 && score < 38 {
		performance = 1
	} else {
		performance = 0
	}

	node.Metrics = metricStr
	node.Information = []byte(information)
	node.Asset = metric.Assets
	node.Region = GetAddrByIp(address)
	node.Performance = &performance
	return global.Mysql.Save(&node).Error
}

func SecureImagesTask(address string, port int) error {
	modelCol := "images"
	images, err := grpc.GetImages(address, port)
	if err != nil {
		return err
	}
	ns := &models.SysNodeSecure{
		Images: images,
	}
	return saveSecure2DB(address, modelCol, images, ns)
}

const (
	SBOM     = "sbom"
	imageLen = 2
)

func SecureDockerTask(category, address string, port int, images []string) error {
	// 因task不支持结构体切片，故使用字符串切片来转化取代
	grpcImages := make([]*securepb.ServerRequest_Spec_Docker_Image, 0)
	for _, image := range images {
		imageInfos := strings.Split(image, ":")
		if len(imageInfos) < imageLen {
			continue
		}
		grpcImages = append(grpcImages, &securepb.ServerRequest_Spec_Docker_Image{
			Repo: imageInfos[0],
			Tag:  imageInfos[1],
		})
	}
	var (
		dockerSecure datatypes.JSON
		modelColumn  string
		nodeSecure   *models.SysNodeSecure
		err          error
	)
	if category == SBOM {
		dockerSecure, err = grpc.GetDockerSBOM(address, port, grpcImages)
		modelColumn = "docker_sbom"
		nodeSecure = &models.SysNodeSecure{
			DockerSbom: dockerSecure,
		}
	} else {
		dockerSecure, err = grpc.GetDockerVul(address, port, grpcImages)
		modelColumn = "docker_vul"
		nodeSecure = &models.SysNodeSecure{
			DockerVul: dockerSecure,
		}
	}
	if err != nil {
		global.Log.Errorf("[%s]获取metalsecure信息失败%v", address, err)
		return err
	}
	return saveSecure2DB(address, modelColumn, dockerSecure, nodeSecure)
}

func SecureBareTask(category, address string, port int, paths []string) error {
	var (
		bareSecure  datatypes.JSON
		modelColumn string
		nodeSecure  *models.SysNodeSecure
		err         error
	)
	if category == SBOM {
		bareSecure, err = grpc.GetBareSBOM(address, port, paths)
		modelColumn = "bare_sbom"
		nodeSecure = &models.SysNodeSecure{
			BareSbom: bareSecure,
		}
	} else {
		bareSecure, err = grpc.GetBareVul(address, port, paths)
		modelColumn = "bare_vul"
		nodeSecure = &models.SysNodeSecure{
			DockerVul: bareSecure,
		}
	}
	if err != nil {
		return err
	}
	return saveSecure2DB(address, modelColumn, bareSecure, nodeSecure)
}

func saveSecure2DB(address, modelCol string, ret datatypes.JSON, ns *models.SysNodeSecure) error {
	var node models.SysNode
	err := global.Mysql.Model(&models.SysNode{}).Where("address = ?", address).First(&node).Error
	if err != nil {
		return err
	}
	var nodeSecure models.SysNodeSecure
	query := global.Mysql.Model(&models.SysNodeSecure{}).Where("node_id = ?", node.Id).First(&nodeSecure)
	if query.Error == nil {
		if err != nil {
			return err
		}
		err = query.Update(modelCol, ret).Error
	} else if query.Error == gorm.ErrRecordNotFound {
		ns.NodeId = node.Id
		err = global.Mysql.Model(&models.SysNodeSecure{}).Create(ns).Error
	}
	global.Log.Errorf("[%s]保存metalsecure信息结果：%v", address, err)
	return err
}

func SetTuneTask(address string, port int, isSave bool) error {
	var node models.SysNode
	err := global.Mysql.Model(&models.SysNode{}).Where("address = ?", address).First(&node).Error
	if err != nil {
		return err
	}
	profile, err := grpc.SetTune(address, port)
	if err != nil {
		return err
	}
	if isSave {
		// 保存调优后的profile到数据库, 调优类型为tune
		tuneLog := &models.SysNodeTuneLog{
			TuneType:    models.AutoTuneType,
			NodeId:      node.Id,
			RespProfile: profile,
		}
		return global.Mysql.Model(&models.SysNodeTuneLog{}).Create(tuneLog).Error
	}
	return nil
}

// InitAsyncTaskMap 将Jenkins任务添加到map里，其他创建的任务也依次添加到map里。(以待后续注册到machinery中)
func InitAsyncTaskMap() map[string]any {
	asyncTaskMap := make(map[string]any)
	asyncTaskMap[async.JenkinsTaskName] = JenkinsTask
	asyncTaskMap[async.GrpcTaskName] = GrpcTask
	asyncTaskMap[async.MailTaskName] = MailTask
	asyncTaskMap[async.ImagesTaskName] = SecureImagesTask
	asyncTaskMap[async.ImageSecureTaskName] = SecureDockerTask
	asyncTaskMap[async.BareSecureTaskName] = SecureBareTask
	asyncTaskMap[async.SetTuneTaskName] = SetTuneTask
	return asyncTaskMap
}

func GetAddrByIp(ip string) string {
	addrInfos := global.Conf.NodeConf.AddrBind
	otherAddr := "Other"
	if len(addrInfos) == 0 {
		return otherAddr
	}
	ip = strings.TrimSpace(ip)
	for _, addrInfo := range addrInfos {
		for _, s := range addrInfo.Ips {
			if strings.HasPrefix(ip, s) {
				return addrInfo.Addr
			}
		}
	}
	return otherAddr
}

func GetPerformanceByMetrics(metric string) int {
	var data int
	reg := regexp.MustCompile(`(\d+).*\(.*`)
	strs := reg.FindStringSubmatch(metric)
	if len(strs) > 0 {
		data, _ = strconv.Atoi(strs[1])
	}
	return data
}
