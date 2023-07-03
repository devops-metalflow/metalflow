package service

import (
	"metalflow/models"
	"metalflow/pkg/response"
)

// GetCountData 获取首页的统计数据
func (s *MysqlService) GetCountData() (*response.CountDataResponseStruct, error) {
	var userCount int64
	err := s.TX.Model(new(models.SysUser)).Count(&userCount).Error
	if err != nil {
		return nil, err
	}

	var nodeCount int64
	err = s.TX.Model(new(models.SysNode)).Count(&nodeCount).Error
	if err != nil {
		return nil, err
	}

	var normalNodeCount int64
	err = s.TX.Model(new(models.SysNode)).Where("health = ?", 0).Count(&normalNodeCount).Error
	if err != nil {
		return nil, err
	}

	var highPerformanceNodeCount int64
	err = s.TX.Model(new(models.SysNode)).Where("performance = ?", 0).Count(&highPerformanceNodeCount).Error
	if err != nil {
		return nil, err
	}
	resp := &response.CountDataResponseStruct{
		UserCount:                userCount,
		NodeCount:                nodeCount,
		NormalNodeCount:          normalNodeCount,
		HighPerformanceNodeCount: highPerformanceNodeCount,
	}
	return resp, err
}

func (s *MysqlService) GetRegionNodeCount() (regionNodeData []response.RegionNodeItemResponseStruct, err error) {
	err = s.TX.Model(new(models.SysNode)).Select("region as region,count(*) as count").
		Group("region").Scan(&regionNodeData).Error
	return
}

func (s *MysqlService) GetManagerNodeCount() (managerNodeData []response.ManagerNodeItemResponseStruct, err error) {
	err = s.TX.Model(new(models.SysNode)).Select("manager as manager,count(*) as count").
		Group("manager").Scan(&managerNodeData).Error
	return
}

func (s *MysqlService) GetHealthNodeCount() (healthNodeData []response.HealthNodeItemResponseStruct, err error) {
	err = s.TX.Model(new(models.SysNode)).Select("health as health,count(*) as count").
		Group("health").Scan(&healthNodeData).Error
	return
}

func (s *MysqlService) GetPerformanceNodeCount() (performanceNodeData []response.PerformanceNodeItemResponseStruct, err error) {
	err = s.TX.Model(new(models.SysNode)).Select("performance as performance,count(*) as count").
		Group("performance").Scan(&performanceNodeData).Error
	return
}
