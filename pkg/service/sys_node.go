package service

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/grpc"
	"metalflow/pkg/request"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GetNodes 获取机器列表
func (s *MysqlService) GetNodes(req *request.NodeListRequestStruct) ([]models.SysNode, error) {
	var err error
	list := make([]models.SysNode, 0)
	query := s.TX.
		Model(&models.SysNode{}).
		Preload("Labels").
		Order("created_at DESC")
	// Eliminate machines that need to be hidden
	hide := strings.TrimSpace(global.Conf.NodeConf.Hide)
	if hide != "" {
		hideNodes := strings.Split(hide, ",")
		query = query.Where("address NOT IN (?)", hideNodes)
	}

	address := strings.TrimSpace(req.Address)
	if address != "" {
		query = query.Where("address LIKE ?", fmt.Sprintf("%%%s%%", address))
	}
	manager := strings.TrimSpace(req.Manager)
	if manager != "" {
		query = query.Where("manager LIKE ?", fmt.Sprintf("%%%s%%", manager))
	}
	region := strings.TrimSpace(req.Region)
	if region != "" {
		query = query.Where("region LIKE ?", fmt.Sprintf("%%%s%%", region))
	}
	if req.Performance != nil {
		query = query.Where("performance = ?", *req.Performance)
	}

	asset := strings.TrimSpace(req.Asset)
	if asset != "" {
		query = query.Where("asset LIKE ?", fmt.Sprintf("%%%s%%", asset))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	if req.Health != nil {
		query = query.Where("health = ?", *req.Health)
	}
	// 查询列表
	err = s.Find(query, &req.PageInfo, &list)
	return list, err
}

func (s *MysqlService) CreateNode(req *request.CreateNodeRequestStruct) error {
	// Eliminate machines that need to be hidden
	hide := strings.TrimSpace(global.Conf.NodeConf.Hide)
	if hide != "" {
		hideNodes := strings.Split(hide, ",")
		for _, node := range hideNodes {
			if req.Address == node {
				return fmt.Errorf("for security reasons, address [%s] is not allowed", req.Address)
			}
		}
	}

	query := s.TX.Where("address = ?", req.Address).First(&models.SysNode{})
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		labels := make([]models.SysLabel, 0)
		if len(req.LabelIds) > 0 {
			err := s.TX.Where("id in (?)", req.LabelIds).Find(&labels).Error
			if err != nil {
				return err
			}
		}
		node := models.SysNode{
			Address:     req.Address,
			SshPort:     req.SshPort,
			Os:          req.Os,
			Asset:       req.Asset,
			Health:      req.Health,
			Performance: req.Performance,
			Region:      GetAddrByIp(req.Address),
			ServicePort: req.ServicePort,
			Remark:      req.Remark,
			Creator:     req.Creator,
			Labels:      labels,
		}
		return s.TX.Create(&node).Error
	} else {
		return fmt.Errorf("the machine node already exists, please do not repeat the creation")
	}
}

// UpdateNodeById 更新机器节点
func (s *MysqlService) UpdateNodeById(nodeId uint, req *request.UpdateNodeRequestStruct) (err error) {
	// 更新机器节点
	var node models.SysNode
	query := s.TX.Where("id = ?", nodeId).First(&node)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("record does not exist, update failed")
	}
	// 更新责任人
	if req.Manager != "" {
		err = query.Update("manager", req.Manager).Error
		if err != nil {
			return
		}
	}
	// 更新机器标签
	if len(req.LabelIds) > 0 {
		labels := make([]models.SysLabel, 0)
		err = s.TX.Where("id in (?)", req.LabelIds).Find(&labels).Error
		if err != nil {
			return
		}
		// 更新机器节点对应的labels
		return s.TX.Model(&node).Association("Labels").Replace(labels)
	}
	return
}

func (s *MysqlService) DeleteNodeByIds(ids []uint) error {
	// 为了解决软删除与索引唯一的冲突，删除前先将unique键address进行更新重写
	for _, id := range ids {
		var node models.SysNode
		err := s.TX.Where("id = ?", id).First(&node).Error
		if err != nil {
			return err
		}
		err = s.TX.Model(new(models.SysNode)).Where("id = ?", id).Update("address",
			fmt.Sprintf("%s|D|%s", time.Now().Format(global.MsecLocalTimeFormat), node.Address)).Error
		if err != nil {
			return err
		}
		err = s.TX.Model(new(models.SysNode)).Where("id = ?", id).Delete(new(models.SysNode)).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// RefreshNodeInfoById 刷新机器配置信息
func (s *MysqlService) RefreshNodeInfoById(id uint) error {
	var node models.SysNode
	query := s.TX.Where("id = ?", id).First(&node)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("机器节点不存在")
	}
	// 判断本次刷新时间与上次刷新时间的间隔，如果小于5分钟，则不进行刷新
	if time.Since(node.RefreshLastTime.Time).Minutes() < 5 { //nolint:gomnd
		return fmt.Errorf("五分钟内已有其他人刷新该节点信息，无需重复刷新")
	}
	// 启动异步任务刷新机器节点信息
	if global.Machinery != nil {
		var worker models.SysWorker
		err := global.Mysql.Model(new(models.SysWorker)).Where("id = ?", 1).First(&worker).Error
		if err != nil {
			global.Log.Error("search worker from database failed")
			return err
		}
		global.Machinery.SendGrpcTask(node.Address, worker.Port, worker.ServiceReq)
		// 计入刷新时间与刷新次数
		err = query.Updates(map[string]any{
			"refresh_count":     *node.RefreshCount + 1,
			"refresh_last_time": time.Now(),
		}).Error
		if err != nil {
			return fmt.Errorf("更新机器节点信息失败")
		}
	}
	return nil
}

//go:embed script/reboot.sh
var f embed.FS

const (
	rebootScriptPath = "script/reboot.sh"
	remoteDir        = "/tmp/"
)

// BatchRebootNodesByIds used to batch reboot nodes
func (s *MysqlService) BatchRebootNodesByIds(ids []uint) error {
	fsInfo := &FSInfo{path: rebootScriptPath}
	metric := grpc.FileMetric{
		FilePath:   rebootScriptPath,
		RemoteDir:  remoteDir,
		FileGetter: fsInfo,
		IsRunnable: true,
	}
	return s.BatchUploadByIds(metric, ids)
}

type FSInfo struct {
	path string
}

func (f2 *FSInfo) GetFile() (io.ReadCloser, error) {
	file, err := f.Open(f2.path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
