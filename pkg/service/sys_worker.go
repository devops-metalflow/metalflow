package service

import (
	"fmt"
	"metalflow/models"
	"metalflow/pkg/request"
	"strings"
)

func (s *MysqlService) GetWorkers(req *request.WorkerListRequestStruct) ([]models.SysWorker, error) {
	list := make([]models.SysWorker, 0)
	query := s.TX.Model(new(models.SysWorker)).Order("created_at DESC")

	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}

	if req.AutoDeploy != nil {
		query = query.Where("auto_deploy = ?", *req.AutoDeploy)
	}

	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}

	// 查询列表
	err := s.Find(query, &req.PageInfo, &list)
	return list, err
}
