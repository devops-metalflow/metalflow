package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"metalflow/pkg/global"
	"metalflow/pkg/grpc/pb"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func SendFlow(address string, port int, message string) (string, error) {
	conn, err := ConnectGrpc(address, port, context.Background())
	if err != nil {
		return "", err
	}
	global.Log.Infof("grpc连接%s:%d成功", address, port)
	defer func(conn *grpc.ClientConn) {
		e := conn.Close()
		if e != nil {
			fmt.Println("grpc关闭服务失败")
		}
	}(conn)
	c := pb.NewMetricsProtoClient(conn)
	global.Log.Infof("设置grpc与%s metrics执行时长为 %d 秒", address, global.Conf.System.ExecuteTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Conf.System.ExecuteTimeout)*time.Second)
	defer cancel()
	r, err := c.SendMetrics(ctx, &pb.MetricsRequest{Message: message})
	if err != nil {
		fmt.Printf("不能获取服务器%s的响应信息：%v", address, err)
		return "", err
	}
	if outputErr := r.GetError(); outputErr != "" {
		fmt.Printf("获取服务器响应信息报错：%s", outputErr)
		return "", fmt.Errorf("err: %s", outputErr)
	}
	return r.GetOutput(), err
}

// ModelMetric 解析json中的部分字符串以保存到数据库
type ModelMetric struct {
	Cpu  string `json:"cpu,omitempty"`
	Disk string `json:"disk,omitempty"`
	Ram  string `json:"ram,omitempty"`
}

// Metric grpc通信返回的所有字段解析
type Metric struct {
	ModelMetric
	Io      string `json:"io,omitempty"`
	Ip      string `json:"ip,omitempty"`
	Kernel  string `json:"kernel,omitempty"`
	Mac     string `json:"mac,omitempty"`
	Network string `json:"network,omitempty"`
	Os      string `json:"os,omitempty"`
	Assets  string `json:"assets,omitempty"`
	Users   string `json:"users,omitempty"`
	Eth     string `json:"eth,omitempty"`
	Wake    string `json:"wake,omitempty"`
}

type Response struct {
	Metrics []Metric `json:"metrics,omitempty"`
}

// ParseMetrics 解析metalmetrics返回的数据
func ParseMetrics(message string) (metric Metric, err error) {
	var m Response
	// grpc返回的数据最外层有双引号，内部可能有单引号
	message = strings.Trim(message, "\"")
	message = strings.Replace(message, "'", "\"", -1)
	err = json.Unmarshal([]byte(message), &m)
	if err != nil {
		return
	}
	if len(m.Metrics) == 0 {
		err = fmt.Errorf("未知的数据：%s", message)
		return
	}
	metric = m.Metrics[0]
	return
}

// GetModelNodeMetric 获取数据库标node需要的字段metrics
func GetModelNodeMetric(metric *Metric) (string, error) {
	modelMetric := metric.ModelMetric
	modelStr := []string{
		modelMetric.Cpu,
		modelMetric.Disk,
		modelMetric.Ram,
	}
	return strings.Join(modelStr, ","), nil
}

// GetNodeInformation 获取数据库表node需要的字段information
func GetNodeInformation(metric *Metric) (string, error) {
	information, err := json.Marshal(metric)
	if err != nil {
		return "", err
	}
	return string(information), nil
}

func ConnectGrpc(address string, port int, ctx context.Context) (conn *grpc.ClientConn, err error) {
	// 根据发现的服务的ip和port建立连接
	cox, cancel := context.WithTimeout(ctx, time.Duration(global.Conf.System.ConnectTimeout)*time.Second)
	defer cancel()
	// 设置可接受的字节最大为100m
	conn, err = grpc.DialContext(cox, fmt.Sprintf("%s:%d", address, port), grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*100))) //nolint:gomnd
	if err != nil {
		global.Log.Errorf("grpc连接%s:%d失败", address, port)
		return nil, errors.New("grpc连接失败")
	}
	return conn, nil
}
