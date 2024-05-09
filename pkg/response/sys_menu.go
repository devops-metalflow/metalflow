package response

// MenuTreeResponseStruct 菜单树信息响应, 字段含义见models.SysMenu
type MenuTreeResponseStruct struct {
	BaseData
	ParentId    uint                     `json:"parentId"`
	Name        string                   `json:"name"`
	Title       string                   `json:"title"`
	Icon        string                   `json:"icon"`
	Path        string                   `json:"path"`
	Redirect    string                   `json:"redirect"`
	Component   string                   `json:"component"`
	Permission  string                   `json:"permission"`
	Creator     string                   `json:"creator"`
	Sort        int                      `json:"sort"`
	Status      uint                     `json:"status"`
	OnlyContent uint                     `json:"onlyContent"`
	NewTab      uint                     `json:"newTab"`
	Visible     uint                     `json:"visible"`
	Breadcrumb  uint                     `json:"breadcrumb"`
	Affix       uint                     `json:"affix"`
	Children    []MenuTreeResponseStruct `json:"children"`
}

// MenuTreeWithAccessResponseStruct 菜单树信息响应, 包含有权限访问的id列表
type MenuTreeWithAccessResponseStruct struct {
	List      []MenuTreeResponseStruct `json:"list"`
	AccessIds []uint                   `json:"accessIds"`
}
