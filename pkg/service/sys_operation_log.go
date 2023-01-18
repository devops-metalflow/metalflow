package service

import (
	"fmt"
	"metalflow/models"
	"metalflow/pkg/request"
	"strings"
)

// GetOperationLogs get operation logs from database.
func (s *MysqlService) GetOperationLogs(req *request.OperationLogRequestStruct) ([]models.SysOperationLog, error) {
	var err error
	list := make([]models.SysOperationLog, 0)
	query := s.TX.
		Model(&models.SysOperationLog{}).
		Order("created_at DESC")
	method := strings.TrimSpace(req.Method)
	if method != "" {
		query = query.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	username := strings.TrimSpace(req.Username)
	if username != "" {
		query = query.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		query = query.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	status := strings.TrimSpace(req.Status)
	if status != "" {
		query = query.Where("status LIKE ?", fmt.Sprintf("%%%s%%", status))
	}
	// 查询列表
	err = s.Find(query, &req.PageInfo, &list)
	return list, err
}
