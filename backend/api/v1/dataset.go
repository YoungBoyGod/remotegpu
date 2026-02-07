package v1

// InitMultipartRequest 初始化分片上传请求
type InitMultipartRequest struct {
	Filename string `json:"filename" binding:"required"`
	Size     int64  `json:"size" binding:"required"`
	MD5      string `json:"md5"`
}

// CompleteMultipartRequest 完成分片上传请求
type CompleteMultipartRequest struct {
	UploadID string `json:"upload_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Size     int64  `json:"size"`
}

// MountRequest 挂载数据集请求
type MountRequest struct {
	MachineID  string `json:"machine_id" binding:"required"`
	MountPoint string `json:"mount_point" binding:"required"`
	ReadOnly   bool   `json:"read_only"`
}
