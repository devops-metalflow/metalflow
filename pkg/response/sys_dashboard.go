package response

type CountDataResponseStruct struct {
	UserCount                int64 `json:"userCount,omitempty"`
	NodeCount                int64 `json:"nodeCount,omitempty"`
	NormalNodeCount          int64 `json:"normalNodeCount,omitempty"`          // 运行中的机器数量
	HighPerformanceNodeCount int64 `json:"highPerformanceNodeCount,omitempty"` // 高性能机器数量
}

type RegionNodeItemResponseStruct struct {
	Region string `json:"region,omitempty"`
	Count  int64  `json:"count,omitempty"`
}

type ManagerNodeItemResponseStruct struct {
	Manager string `json:"manager,omitempty"`
	Count   int64  `json:"count,omitempty"`
}

type HealthNodeItemResponseStruct struct {
	Health int64 `json:"health"`
	Count  int64 `json:"count"`
}

type PerformanceNodeItemResponseStruct struct {
	Performance int64 `json:"performance"`
	Count       int64 `json:"count"`
}
