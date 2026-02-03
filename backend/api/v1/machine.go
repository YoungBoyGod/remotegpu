package v1

// AllocateRequest 分配机器请求
type AllocateRequest struct {
	CustomerID     uint   `json:"customer_id" binding:"required"`
	HostID         string `json:"host_id" binding:"required"`
	DurationMonths int    `json:"duration_months" binding:"required,min=1"`
	Remark         string `json:"remark"`
}

// ReclaimRequest 回收机器请求
type ReclaimRequest struct {
	Reason string `json:"reason"`
	Force  bool   `json:"force"`
}
