package request

type CollectRequest struct {
	NodeId uint `json:"nodeId" form:"nodeId"`
}

type UpdateDescRequest struct {
	Username    string `json:"username" form:"username" validate:"required"`
	Address     string `json:"address" form:"address" validate:"required"`
	Description string `json:"description" form:"description"`
}
