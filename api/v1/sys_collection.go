package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/utils"
)

// GetMyCollections 查询我的收藏机器节点信息
func GetMyCollections(c *gin.Context) {
	user := GetCurrentUser(c)
	// 获取工号收藏的所有机器节点信息
	collections := make([]models.SysCollection, 0)
	err := global.Mysql.Model(&models.SysCollection{}).Where("username = ?", user.Username).Find(&collections).Error
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	// 获取收藏的所有nodeIds
	responses := make([]response.CollectionNodeItem, 0)
	for _, c := range collections {
		var node models.SysNode
		err = global.Mysql.Model(&models.SysNode{}).Where("id = ?", c.NodeId).First(&node).Error
		if err != nil {
			global.Log.Errorf("数据库查找服务器%d失败：%v", c.NodeId, err)
			continue
		}
		item := response.CollectionNodeItem{
			CollectionId: c.Id,
			Description:  c.Description,
			SysNode:      node,
		}
		responses = append(responses, item)
	}

	response.SuccessWithData(responses)
}

func CreateCollection(c *gin.Context) {
	// 绑定参数
	var req request.CollectRequest
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}
	user := GetCurrentUser(c)
	collect := &models.SysCollection{
		Username: user.Username,
		NodeId:   req.NodeId,
	}
	// 检查是否重复收藏
	err = global.Mysql.Model(&models.SysCollection{}).
		Where(" username = ? AND node_id = ?", user.Username, req.NodeId).First(new(models.SysCollection)).Error
	if err == nil {
		response.FailWithMsg("无需重复收藏~")
		return
	}
	err = global.Mysql.Model(&models.SysCollection{}).Create(collect).Error
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

func UpdateDescriptionById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateDescRequest
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 获取path中的collectionId
	collectionId := utils.Str2Uint(c.Param("collectionId"))
	if collectionId == 0 {
		response.FailWithMsg("收藏Id不正确")
		return
	}
	err = global.Mysql.Model(&models.SysCollection{}).Where("id = ?", collectionId).
		Update("description", req.Description).Error
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("更新机器描述信息失败：%s", err.Error()))
	}
	response.Success()
}

func DeleteCollectByNodeId(c *gin.Context) {
	// 绑定参数
	var req request.CollectRequest
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	user := GetCurrentUser(c)
	err = global.Mysql.Model(&models.SysCollection{}).Where("node_id = ? AND username = ?", req.NodeId, user.Username).
		Delete(&models.SysCollection{}).Error
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
