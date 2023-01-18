package response

import "metalflow/models"

type CronShutNodeResponse struct {
	Id        uint              `json:"id"`
	Name      string            `json:"name"`
	Keyword   string            `json:"keyword"`
	StartTime string            `json:"startTime"`
	ShutTime  string            `json:"shutTime"`
	Creator   string            `json:"creator"`
	Status    *uint             `json:"status"`
	Nodes     []*models.SysNode `json:"nodes"`
	CreatedAt models.LocalTime  `json:"createdAt"`
}
