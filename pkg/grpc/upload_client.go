package grpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"metalflow/pkg/global"
	proto "metalflow/pkg/grpc/uploadpb"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const ChunkSize = 64 * 1024 // 64kib

type FileGetter interface {
	GetFile() (io.ReadCloser, error)
}

type Uploader struct {
	client      proto.TaskProtoClient
	ctx         context.Context
	doneRequest chan string
	failRequest chan string
}

type FileMetric struct {
	FilePath   string
	RemoteDir  string
	IsRunnable bool
	FileGetter FileGetter
}

func NewUploader(ctx context.Context, client proto.TaskProtoClient) *Uploader {
	u := &Uploader{
		ctx:         ctx,
		client:      client,
		doneRequest: make(chan string),
		failRequest: make(chan string),
	}
	return u
}

// Upload 传输文件 tmpFilepath为本地服务器临时的文件所在目录，remoteDir为期望放置目的机器节点文件的目录，isRun为是否传输完成后执行
func Upload(address string, port int, metric FileMetric) (string, error) {
	ctx := context.Background()
	conn, err := ConnectGrpc(address, port, ctx)
	if err != nil {
		return "", err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("grpc关闭服务失败")
		}
	}(conn)
	return uploadFiles(ctx, proto.NewTaskProtoClient(conn), metric)
}

func uploadFiles(ctx context.Context, client proto.TaskProtoClient, metric FileMetric) (output string, err error) {
	ul := NewUploader(ctx, client)
	// 添加传输执行任务超时时间
	c, cancel := context.WithTimeout(ctx, time.Duration(global.Conf.System.ExecuteTimeout)*time.Second)
	defer cancel()

	go ul.upload(metric)

	select {
	// 如果传输成功，则接收成功后的输出
	case output = <-ul.doneRequest:
	case failedStr := <-ul.failRequest:
		err = fmt.Errorf("文件传输失败： %s", failedStr)
	case <-c.Done():
		err = errors.New("文件传输或执行超时")
	}
	return
}

func (u *Uploader) upload(fileMetric FileMetric) {
	var buf []byte
	// start upload
	streamUploader, err := u.client.SendTask(u.ctx)
	if err != nil {
		_ = fmt.Errorf("failed to create upload stream for file %s", fileMetric.FilePath) //nolint:ineffassign
	}
	defer func(streamUploader proto.TaskProto_SendTaskClient) {
		_ = streamUploader.CloseSend()
	}(streamUploader)

	// split filename
	_, filename := filepath.Split(fileMetric.FilePath)
	remoteFilePath := filepath.Join(fileMetric.RemoteDir, filename)
	// convert path Separator to linux
	remoteFilePath = strings.ReplaceAll(remoteFilePath, string(os.PathSeparator), "/")

	// create a buffer of chunkSize to be streamed
	buf = make([]byte, ChunkSize)
	file, err := fileMetric.FileGetter.GetFile()
	defer func(file io.ReadCloser) {
		_ = file.Close()
	}(file)
	if err != nil {
		u.failRequest <- filename
		return
	}
	for {
		var n int
		// 不能使用并发直接file.Read(),因file类型是一个指针，因此多个goroutine读取的是同一个file。
		// 而file.Read()会不断地消耗file，故成了并发对同一个file写了。因此file的open及读取关闭要当成原子操作。
		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {
				err = nil //nolint:ineffassign
				// finish the file
				break
			}
			return
		}
		err = streamUploader.Send(&proto.TaskRequest{
			Data:     buf[:n],
			Path:     remoteFilePath,
			Runnable: fileMetric.IsRunnable,
		})
		if err != nil {
			fmt.Println("send chunk failed: ", err)
			break
		}
	}
	status, err := streamUploader.CloseAndRecv()
	if err != nil {
		u.failRequest <- filename
		return
	}

	if errOutput := status.Error; errOutput != "" {
		u.failRequest <- errOutput
		return
	}
	// 如果成功，则将远程metaltask的消息返回
	u.doneRequest <- status.Output
}
