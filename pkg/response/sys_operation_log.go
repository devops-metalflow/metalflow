package response

import (
	"gorm.io/datatypes"
	"metalflow/models"
)

// OperationLogListResponseStruct is operation log response struct.
type OperationLogListResponseStruct struct {
	Id        uint             `json:"id"`
	CreatedAt models.LocalTime `json:"createdAt"`
	ApiDesc   string           `json:"apiDesc"`
	Path      string           `json:"path"`
	Method    string           `json:"method"`
	Header    datatypes.JSON   `json:"header"`
	Body      datatypes.JSON   `json:"body"`
	Data      datatypes.JSON   `json:"data"`
	Status    int              `json:"status"`
	Username  string           `json:"username"`
	RoleName  string           `json:"roleName"`
	Latency   int64            `json:"latency"`
	UserAgent string           `json:"userAgent"`
}
