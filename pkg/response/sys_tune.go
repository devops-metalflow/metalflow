package response

import "metalflow/models"

type TuneLogList struct {
	Id        uint             `json:"id"`
	CreatedAt models.LocalTime `json:"createdAt"`
}
