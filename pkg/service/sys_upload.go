package service

import (
	"errors"
	"fmt"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/grpc"
	"strings"
	"sync"
)

// BatchUploadByIds uses batch upload file to nodes by each database id.
func (s *MysqlService) BatchUploadByIds(m grpc.FileMetric, ids []uint) error {
	nodes := make([]*models.SysNode, 0)
	err := s.TX.Where("id in (?)", ids).Find(&nodes).Error
	if err != nil {
		return err
	}
	// get metaltask info from database
	// TODO 因进行过worker初始化，所以id为2的worker就是metaltask
	var metaltask models.SysWorker
	err = global.Mysql.Model(new(models.SysWorker)).Where("id = ?", 2).First(&metaltask).Error //nolint:gomnd
	if err != nil {
		return fmt.Errorf("search worker metaltask from database error:%v", err)
	}

	var (
		errAddrs = make([]string, 0)
		lock     = new(sync.RWMutex)
		wg       = sync.WaitGroup{}
	)

	// nolint:gomnd
	tokens := make(chan struct{}, 10)
	for _, node := range nodes {
		wg.Add(1)
		go func(addr string, port int, metric grpc.FileMetric) {
			defer wg.Done()
			// 限制grpc连接的并发数量
			tokens <- struct{}{}
			_, e := grpc.Upload(addr, port, metric)
			<-tokens
			if e != nil {
				lock.Lock()
				errAddrs = append(errAddrs, fmt.Sprintf("[%s]运行失败: %v。", addr, e))
				lock.Unlock()
			}
		}(node.Address, metaltask.Port, m)
	}
	wg.Wait()
	if len(errAddrs) > 0 {
		errStr := strings.Join(errAddrs, " ")
		return errors.New(errStr)
	}
	return nil
}
