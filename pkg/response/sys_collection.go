package response

import "metalflow/models"

type CollectionNodeItem struct {
	CollectionId uint   `json:"collectionId"`
	Description  string `json:"description"`
	models.SysNode
}
