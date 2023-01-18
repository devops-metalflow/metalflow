package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"metalflow/pkg/global"
	"metalflow/pkg/grpc/tunepb"
	"time"
)

func SendTuneRequest(address string, port int, reqBody *tunepb.ServerRequest) (string, error) {
	conn, err := ConnectGrpc(address, port, context.Background())
	if err != nil {
		return "", err
	}
	global.Log.Infof("grpc连接%s:%d成功", address, port)
	defer func(conn *grpc.ClientConn) {
		e := conn.Close()
		if e != nil {
			fmt.Println("grpc关闭服务失败")
			return
		}
	}(conn)

	c := tunepb.NewServerProtoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	r, err := c.SendServer(ctx, reqBody)
	if err != nil {
		global.Log.Errorf("[%s]metaltune获取数据失败：%v", address, err)
		fmt.Printf("不能获取服务器%s的metaltune响应信息：%v", address, err)
		return "", err
	}
	if outputErr := r.GetError(); outputErr != "" {
		global.Log.Errorf("[%s]metaltune获取数据报错：%v", address, outputErr)
		fmt.Printf("获取服务器%s的metaltune信息报错：%s", address, outputErr)
		return "", fmt.Errorf("err: %s", outputErr)
	}
	return r.GetOutput(), nil
}

const (
	name = "metaltune"
)

func SetTune(address string, port int) (string, error) {
	tuneReq := &tunepb.ServerRequest{
		ApiVersion: apiVersion,
		Kind:       name,
		Metadata:   &tunepb.Metadata{Name: name},
		Spec:       &tunepb.Spec{Tuning: &tunepb.Tuning{Auto: true}},
	}
	ret, err := SendTuneRequest(address, port, tuneReq)
	if err != nil {
		return "", err
	}
	return ret, nil
}
