package response

// ApiListResponseStruct 接口信息响应, 字段含义见models
type ApiListResponseStruct struct {
	BaseData
	Method   string `json:"method"`
	Path     string `json:"path"`
	Category string `json:"category"`
	Creator  string `json:"creator"`
	Desc     string `json:"desc"`
	Title    string `json:"title"`
}

// ApiGroupByCategoryResponseStruct 权限接口信息响应, 字段含义见models
type ApiGroupByCategoryResponseStruct struct {
	Title    string                  `json:"title"`    // 标题
	Category string                  `json:"category"` // 分组名称
	Children []ApiListResponseStruct `json:"children"` // 前端以树形图结构展示, 这里用children表示
}

// ApiTreeWithAccessResponseStruct 接口树信息响应, 包含有权限访问的id列表
type ApiTreeWithAccessResponseStruct struct {
	List      []ApiGroupByCategoryResponseStruct `json:"list"`
	AccessIds []uint                             `json:"accessIds"`
}
