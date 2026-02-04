package v1

// CreateSSHKeyRequest 创建 SSH 密钥请求
type CreateSSHKeyRequest struct {
	Name      string `json:"name" binding:"required,min=1,max=64"`
	PublicKey string `json:"public_key" binding:"required"`
}

// SSHKeyResponse SSH 密钥响应
type SSHKeyResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	CreatedAt   string `json:"created_at"`
}
