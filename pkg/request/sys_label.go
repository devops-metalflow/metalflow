package request

import "metalflow/pkg/response"

type LabelListRequestStruct struct {
	Name              string `json:"name,omitempty" form:"name"`
	Creator           string `json:"creator,omitempty" form:"creator"`
	response.PageInfo        // 分页参数
}

type CreateLabelRequestStruct struct {
	Name    string `json:"name,omitempty" form:"name" validate:"required"`
	Creator string `json:"creator,omitempty" form:"creator"`
}

type UpdateLabelRequestStruct struct {
	Name string `json:"name,omitempty" form:"name"`
}

func (s *CreateLabelRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "标签名称"
	return m
}
