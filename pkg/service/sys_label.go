package service

import (
	"fmt"
	"metalflow/models"
	"metalflow/pkg/request"
	"strings"
)

// GetLabels 获取label列表
func (s *MysqlService) GetLabels(req *request.LabelListRequestStruct) ([]models.SysLabel, error) {
	list := make([]models.SysLabel, 0)
	query := s.TX.Model(new(models.SysLabel)).Order("created_at DESC")

	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}

	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}

	// 查询列表
	err := s.Find(query, &req.PageInfo, &list)
	return list, err
}
