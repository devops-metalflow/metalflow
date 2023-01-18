package v1

import (
	"github.com/gin-gonic/gin"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
	"metalflow/pkg/utils"
)

// GetRiskCountById get the risk number of node.
func GetRiskCountById(c *gin.Context) {
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	count, err := s.GetRiskCountById(nodeId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData(count)
}

// GetNodeImagesById obtain all docker images of the corresponding node through nodeId.
func GetNodeImagesById(c *gin.Context) {
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	images, err := s.GetNodeImages(nodeId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	var resp response.SecureImages
	resp.Dockers = images
	response.SuccessWithData(resp)
}

const SBOM = "sbom"

// GetNodeImageSecureById obtain the security report of the docker container corresponding to node.
func GetNodeImageSecureById(c *gin.Context) {
	var req request.SecureImage
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}

	s := service.New(c)
	report, err := s.GetNodeDockerSecureInfo(nodeId, &req)
	if err != nil {
		response.FailWithMsg("failed to get docker security report for this node")
		return
	}
	var resp response.SecureImageReport
	if req.Category == SBOM {
		resp.Sbom = report
	} else {
		resp.Vul = report
	}
	response.SuccessWithData(resp)
}

// GetNodeBareSecureById obtain the bare metal security report of the node.
func GetNodeBareSecureById(c *gin.Context) {
	var req request.SecureBare
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}

	s := service.New(c)
	bareReport, err := s.GetNodeBareSecureInfo(nodeId, &req)
	if err != nil {
		response.FailWithMsg("failed to obtain the bare metal security report of this node")
		return
	}
	var resp response.SecureImageReport
	if req.Category == SBOM {
		resp.Sbom = bareReport
	} else {
		resp.Vul = bareReport
	}
	response.SuccessWithData(resp)
}

// GetNodeSecurityScoreById get the security score of the node.
func GetNodeSecurityScoreById(c *gin.Context) {
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}
	s := service.New(c)
	score, err := s.GetNodeSecureScore(nodeId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData(score)
}

// RunNodeBareSecurityById fixes bare risks for node.
func RunNodeBareSecurityById(c *gin.Context) {
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	err := s.RunBareSecure(nodeId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData("fix risk successfully")
}

// RunNodeContainerSecurityById fixes the security of the docker container first.
func RunNodeContainerSecurityById(c *gin.Context) {
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	err := s.RunContainerSecure(nodeId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData("fix risk successfully")
}

// FixNodeSecurityIssuesById fix node security issues.
func FixNodeSecurityIssuesById(c *gin.Context) {
	var req request.SecureFix
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	err = s.FixSecureRisk(nodeId, req.CveId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData("fix risk successfully")
}
