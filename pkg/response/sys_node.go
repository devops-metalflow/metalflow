package response

import (
	"gorm.io/datatypes"
	"metalflow/models"
)

type NodeListResponseStruct struct {
	Id          uint              `json:"id"`
	Address     string            `json:"address"`
	CreatedAt   string            `json:"createdAt"`
	Manager     string            `json:"manager"`
	Metrics     string            `json:"metrics"`
	SshPort     int               `json:"sshPort"`
	Asset       string            `json:"asset"`
	Health      *uint             `json:"health"`
	PingStat    *uint             `json:"pingStat"`
	Performance *uint             `json:"performance"`
	Region      string            `json:"region"`
	Remark      string            `json:"remark"`
	Creator     string            `json:"creator"`
	Labels      []models.SysLabel `json:"labels"`
	Information datatypes.JSON    `json:"information"`
}

type ShellWsFilesResponseStruct struct {
	Files      []map[string]string `json:"files"`
	FileCount  uint                `json:"fileCount"`
	Paths      []map[string]string `json:"paths"`
	DirCount   uint                `json:"dirCount"`
	CurrentDir string              `json:"currentDir"`
}
