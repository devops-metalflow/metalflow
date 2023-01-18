package v1

import (
	"github.com/gin-gonic/gin"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
)

// GetCountData gets homepage statistics.
func GetCountData(c *gin.Context) {
	s := service.New(c)
	resp, err := s.GetCountData()
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData(resp)
}

// GetRegionNodeData count the number of each machine node according to the region.
func GetRegionNodeData(c *gin.Context) {
	s := service.New(c)
	resp, err := s.GetRegionNodeCount()
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData(resp)
}

// GetManagerNodeData count the number of each machine node according to the person in charge.
func GetManagerNodeData(c *gin.Context) {
	s := service.New(c)
	resp, err := s.GetManagerNodeCount()
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData(resp)
}

// GetHealthNodeData count the number of each machine node according to the health.
func GetHealthNodeData(c *gin.Context) {
	s := service.New(c)
	resp, err := s.GetHealthNodeCount()
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData(resp)
}

// GetPerformanceNodeData count the number of machine nodes according to machine performance.
func GetPerformanceNodeData(c *gin.Context) {
	s := service.New(c)
	resp, err := s.GetPerformanceNodeCount()
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData(resp)
}
