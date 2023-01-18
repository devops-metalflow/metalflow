package async

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
)

const (
	JenkinsTaskName     string = "jenkins-Task"     // Jenkins任务的名字，可以自己取
	GrpcTaskName        string = "grpc-Task"        // grpc任务的名字
	MailTaskName        string = "mail-Task"        // mail任务名
	ImagesTaskName      string = "images-Task"      // 获取所有images的任务名
	ImageSecureTaskName string = "imageSecure-Task" // 镜像的安全任务名
	BareSecureTaskName  string = "bareSecure-Task"  // 裸金属的安全任务名
	SecureScoreTaskName string = "secureScore-Task" // 获取安全分数任务名
	SetTuneTaskName     string = "setTune-Task"     // 设置智能调优
)

// 任务发送的接口都放到这里，即每写一个任务便要在这里将对应任务的参数注册并暴露出去

// SendJenkinsTask Jenkins任务参数注册，提供发送Jenkins任务的接口
func (m *Machinery) SendJenkinsTask(jobName, jobParams string) {
	args := make([]tasks.Arg, 0)
	args = append(args,
		tasks.Arg{
			Name:  "jobName",
			Type:  "string",
			Value: jobName,
		},
		tasks.Arg{
			Name:  "jobParams",
			Type:  "string",
			Value: jobParams,
		},
	)
	// machinery最终是转化成json结构的，当前不支持map[string]string类型
	m.registerTask(JenkinsTaskName, args)
}

type Mail struct {
	Title, Body   string
	Receivers, CC []string
}

func (m *Machinery) SendMailTask(mail *Mail) {
	args := make([]tasks.Arg, 0)
	args = append(args,
		tasks.Arg{
			Name:  "title",
			Type:  "string",
			Value: mail.Title,
		},
		tasks.Arg{
			Name:  "body",
			Type:  "string",
			Value: mail.Body,
		},
		tasks.Arg{
			Name:  "receivers",
			Type:  "[]string",
			Value: mail.Receivers,
		},
		tasks.Arg{
			Name:  "cc",
			Type:  "[]string",
			Value: mail.CC,
		},
	)
	m.registerTask(MailTaskName, args)
}

func (m *Machinery) SendSecureImagesTask(address string, port int) {
	args := make([]tasks.Arg, 0)
	args = append(args,
		tasks.Arg{
			Name:  "address",
			Type:  "string",
			Value: address,
		},
		tasks.Arg{
			Name:  "port",
			Type:  "int",
			Value: port,
		},
	)
	m.registerTask(ImagesTaskName, args)
}

func (m *Machinery) SendSecureDockerTask(category, address string, port int, images []string) {
	args := make([]tasks.Arg, 0)
	args = append(args,
		tasks.Arg{
			Name:  "category",
			Type:  "string",
			Value: category,
		},
		tasks.Arg{
			Name:  "address",
			Type:  "string",
			Value: address,
		},
		tasks.Arg{
			Name:  "port",
			Type:  "int",
			Value: port,
		},
		tasks.Arg{
			Name:  "images",
			Type:  "[]string",
			Value: images,
		},
	)
	m.registerTask(ImageSecureTaskName, args)
}

func (m *Machinery) SendSecureBareTask(category, address string, port int, paths []string) {
	args := make([]tasks.Arg, 0)
	args = append(args,
		tasks.Arg{
			Name:  "category",
			Type:  "string",
			Value: category,
		},
		tasks.Arg{
			Name:  "address",
			Type:  "string",
			Value: address,
		},
		tasks.Arg{
			Name:  "port",
			Type:  "int",
			Value: port,
		},
		tasks.Arg{
			Name:  "paths",
			Type:  "[]string",
			Value: paths,
		},
	)
	m.registerTask(BareSecureTaskName, args)
}

func (m *Machinery) SendSecureScoreTask(address string, port int) {
	args := make([]tasks.Arg, 0)
	args = append(args,
		tasks.Arg{
			Name:  "address",
			Type:  "string",
			Value: address,
		},
		tasks.Arg{
			Name:  "port",
			Type:  "int",
			Value: port,
		},
	)
	m.registerTask(SecureScoreTaskName, args)
}

func (m *Machinery) SetTune(address string, port int, isSave bool) {
	args := make([]tasks.Arg, 0)
	args = append(args,
		tasks.Arg{
			Name:  "address",
			Type:  "string",
			Value: address,
		},
		tasks.Arg{
			Name:  "port",
			Type:  "int",
			Value: port,
		},
		tasks.Arg{
			Name:  "isSave",
			Type:  "bool",
			Value: isSave,
		},
	)
	m.registerTask(SetTuneTaskName, args)
}

// SendGrpcTask grpcTask
func (m *Machinery) SendGrpcTask(address string, port int, message string) {
	args := make([]tasks.Arg, 0)
	args = append(args,
		tasks.Arg{
			Name:  "address",
			Type:  "string",
			Value: address,
		},
		tasks.Arg{
			Name:  "port",
			Type:  "int",
			Value: port,
		},
		tasks.Arg{
			Name:  "message",
			Type:  "string",
			Value: message,
		},
	)
	m.registerTask(GrpcTaskName, args)
}

func (m *Machinery) registerTask(taskName string, args []tasks.Arg) {
	signature, _ := tasks.NewSignature(taskName, args)
	signature.RetryCount = 5
	_, err := m.MachineryServer.SendTaskWithContext(m.Ctx, signature)
	if err != nil {
		fmt.Println(fmt.Errorf("send %s err:%v", taskName, err))
	}
}
