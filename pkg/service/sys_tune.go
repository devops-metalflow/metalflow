package service

import (
	"gorm.io/gorm"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/grpc"
	"metalflow/pkg/grpc/tunepb"
	"metalflow/pkg/request"
	"time"
)

const (
	apiVersion = "1.0.0"
	name       = "metaltune"
)

func (s *MysqlService) GetTuneScoreByNodeId(nodeId uint) (uint, error) {
	var (
		score, cleanupScore, turboScore uint
		cleanupLog, turboLog            models.SysNodeTuneLog
	)
	err := s.TX.Model(&models.SysNodeTuneLog{}).Where("node_id = ? AND tune_type = ?", nodeId, models.CleanupType).
		Order("created_at DESC").First(&cleanupLog).Error
	if err == gorm.ErrRecordNotFound {
		cleanupScore = 50
	} else if err != nil {
		return 0, err
	}
	err = s.TX.Model(&models.SysNodeTuneLog{}).Where("node_id = ? AND tune_type = ?", nodeId, models.TurboType).
		Order("created_at DESC").First(&turboLog).Error
	if err == gorm.ErrRecordNotFound {
		turboScore = 50
	} else if err != nil {
		return 0, err
	}
	// 根据记录计算分数
	nowDate := time.Now()
	if cleanupLog.NodeId > 0 {
		cleanupSubDays := uint(nowDate.Sub(cleanupLog.CreatedAt.Time).Hours() / 24) // nolint:gomnd
		cleanupScore = calculateScore(cleanupSubDays)
	}
	if turboLog.NodeId > 0 {
		turboSubDays := uint(nowDate.Sub(turboLog.CreatedAt.Time).Hours() / 24) // nolint:gomnd
		turboScore = calculateScore(turboSubDays)
	}
	score = 100 - (cleanupScore + turboScore) // nolint:gomnd
	return score, nil
}

// nolint:gomnd
func calculateScore(day uint) uint {
	switch {
	case day < 7:
		return 0
	case day < 14:
		return 10
	case day < 21:
		return 20
	case day < 28:
		return 30
	case day < 35:
		return 40
	default:
		return 50
	}
}

func (s *MysqlService) GetTuneLogsByNodeId(nodeId uint, req *request.TuneLogsRequest) ([]*models.SysNodeTuneLog, error) {
	list := make([]*models.SysNodeTuneLog, 0)
	query := s.TX.Model(&models.SysNodeTuneLog{}).Where("tune_type = ? AND node_id = ?", models.AutoTuneType, nodeId).
		Order("created_at DESC")
	// 查询列表
	err := s.Find(query, &req.PageInfo, &list)
	return list, err
}

func (s *MysqlService) Cleanup(nodeId uint) error {
	var metaltune models.SysWorker
	err := s.TX.Model(new(models.SysWorker)).Where("id = ?", 4).First(&metaltune).Error //nolint:gomnd
	if err != nil {
		return err
	}
	var node models.SysNode
	err = s.TX.Model(models.SysNode{}).Where("id = ?", nodeId).First(&node).Error
	if err != nil {
		return err
	}

	tuneRequest := &tunepb.ServerRequest{
		ApiVersion: apiVersion,
		Kind:       name,
		Metadata:   &tunepb.Metadata{Name: name},
		Spec:       &tunepb.Spec{Cleanup: true},
	}
	_, err = grpc.SendTuneRequest(node.Address, metaltune.Port, tuneRequest)
	if err != nil {
		return err
	}
	// save cleanup log in database
	tuneLog := &models.SysNodeTuneLog{
		NodeId:      nodeId,
		TuneType:    models.CleanupType,
		RespProfile: "",
	}
	return s.TX.Model(&models.SysNodeTuneLog{}).Create(tuneLog).Error
}

func (s *MysqlService) Rollback(nodeId, logId uint) error {
	var metaltune models.SysWorker
	err := s.TX.Model(new(models.SysWorker)).Where("id = ?", 4).First(&metaltune).Error //nolint:gomnd
	if err != nil {
		return err
	}
	var node models.SysNode
	err = s.TX.Model(models.SysNode{}).Where("id = ?", nodeId).First(&node).Error
	if err != nil {
		return err
	}

	// 根据logId获取临近上一次的resp profile
	var tuneLog models.SysNodeTuneLog
	err = s.TX.Model(&models.SysNodeTuneLog{}).Where("id = ?", logId).First(&tuneLog).Error
	if err != nil {
		return err
	}
	tuneLogs := make([]*models.SysNodeTuneLog, 0)
	err = s.TX.Model(&models.SysNodeTuneLog{}).
		Where("node_id = ? AND tune_type = ? AND created_at < ?", nodeId, models.AutoTuneType, tuneLog.CreatedAt).
		Order("created_at DESC").Find(&tuneLogs).Error
	if err != nil {
		return err
	}

	var preProfile string
	// 如果有数据，则其为排序后的第一条数据。否则，则为原始的profile
	if len(tuneLogs) > 0 {
		preProfile = tuneLogs[0].RespProfile
	} else {
		var tuneScene models.SysNodeTuneScene
		err = s.TX.Model(&models.SysNodeTuneScene{}).Where("scene = ?", "origin").First(&tuneScene).Error
		if err != nil {
			return err
		}
		preProfile = tuneScene.Profile
	}
	tuneReq := &tunepb.ServerRequest{
		ApiVersion: apiVersion,
		Kind:       name,
		Metadata:   &tunepb.Metadata{Name: name},
		Spec: &tunepb.Spec{
			Tuning: &tunepb.Tuning{
				Profile: preProfile,
			},
		},
	}
	_, err = grpc.SendTuneRequest(node.Address, metaltune.Port, tuneReq)
	if err != nil {
		return err
	}
	return nil
}

func (s *MysqlService) SetTune(nodeId uint, req *request.TuneAutoSetRequest) error {
	var metaltune models.SysWorker
	err := s.TX.Model(new(models.SysWorker)).Where("id = ?", 4).First(&metaltune).Error //nolint:gomnd
	if err != nil {
		return err
	}
	var node models.SysNode
	err = s.TX.Model(models.SysNode{}).Where("id = ?", nodeId).First(&node).Error
	if err != nil {
		return err
	}
	// 异步请求同步进行智能调优
	global.Machinery.SetTune(node.Address, metaltune.Port, req.IsSave)
	return nil
}

func (s *MysqlService) SetScene(nodeId uint, scene string) error {
	var metaltune models.SysWorker
	err := s.TX.Model(new(models.SysWorker)).Where("id = ?", 4).First(&metaltune).Error //nolint:gomnd
	if err != nil {
		return err
	}
	var node models.SysNode
	err = s.TX.Model(models.SysNode{}).Where("id = ?", nodeId).First(&node).Error
	if err != nil {
		return err
	}
	// 获取对应场景的profile
	var sceneProfile models.SysNodeTuneScene
	err = s.TX.Model(&models.SysNodeTuneScene{}).Where("scene = ?", scene).First(&sceneProfile).Error
	if err != nil {
		return err
	}
	tuneReq := &tunepb.ServerRequest{
		ApiVersion: apiVersion,
		Kind:       name,
		Metadata:   &tunepb.Metadata{Name: name},
		Spec: &tunepb.Spec{
			Tuning: &tunepb.Tuning{
				Profile: sceneProfile.Profile,
			},
		},
	}
	_, err = grpc.SendTuneRequest(node.Address, metaltune.Port, tuneReq)
	if err != nil {
		return err
	}
	return nil
}

func (s *MysqlService) Turbo(nodeId uint) error {
	var metaltune models.SysWorker
	err := s.TX.Model(new(models.SysWorker)).Where("id = ?", 4).First(&metaltune).Error //nolint:gomnd
	if err != nil {
		return err
	}
	var node models.SysNode
	err = s.TX.Model(models.SysNode{}).Where("id = ?", nodeId).First(&node).Error
	if err != nil {
		return err
	}

	tuneReq := &tunepb.ServerRequest{
		ApiVersion: apiVersion,
		Kind:       name,
		Metadata:   &tunepb.Metadata{Name: name},
		Spec: &tunepb.Spec{
			Turbo: true,
		},
	}
	_, err = grpc.SendTuneRequest(node.Address, metaltune.Port, tuneReq)
	if err != nil {
		return err
	}
	// save turbo log in database.
	tuneLog := &models.SysNodeTuneLog{NodeId: nodeId, TuneType: models.TurboType, RespProfile: ""}
	return s.TX.Model(&models.SysNodeTuneLog{}).Create(tuneLog).Error
}
