package v1

import (
	"github.com/gin-gonic/gin"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
)

// GetMyCollections queries my favorite machines.
func GetMyCollections(c *gin.Context) {
	user := GetCurrentUser(c)
	// obtain all machine node information collected by current user.
	collections := make([]models.SysCollection, 0)
	err := global.Mysql.Model(&models.SysCollection{}).Where("username = ?", user.Username).Find(&collections).Error
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	// get all nodeIds of collection.
	nodeIds := make([]uint, 0)
	for _, c := range collections {
		nodeIds = append(nodeIds, c.NodeId)
	}

	collectNodes := make([]models.SysNode, 0)
	err = global.Mysql.Model(&models.SysNode{}).Where("id in (?)", nodeIds).Find(&collectNodes).Error
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData(collectNodes)
}

// CreateCollection used to favorite a node.
func CreateCollection(c *gin.Context) {
	// bind request body to struct.
	var req request.CollectRequest
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}
	user := GetCurrentUser(c)
	collect := &models.SysCollection{
		Username: user.Username,
		NodeId:   req.NodeId,
	}
	// check for duplicate favorites.
	err = global.Mysql.Model(&models.SysCollection{}).
		Where(" username = ? and node_id = ?", user.Username, req.NodeId).First(new(models.SysCollection)).Error
	if err == nil {
		response.FailWithMsg("you don't need to repeat favorites~")
		return
	}
	err = global.Mysql.Model(&models.SysCollection{}).Create(collect).Error
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// DeleteCollectByNodeId for canceling favorites.
func DeleteCollectByNodeId(c *gin.Context) {
	// bind request body to struct.
	var req request.CollectRequest
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}
	err = global.Mysql.Model(&models.SysCollection{}).Where("node_id = ?", req.NodeId).Delete(&models.SysCollection{}).Error
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
