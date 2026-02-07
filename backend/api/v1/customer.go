package v1

// CreateCustomerRequest 创建客户请求
type CreateCustomerRequest struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	CompanyCode string `json:"company_code" binding:"required"`
	Password    string `json:"password" binding:"omitempty,min=6"`
	Role        string `json:"role"`
	DisplayName string `json:"display_name"`
	FullName    string `json:"full_name"`
	Company     string `json:"company"`
	Phone       string `json:"phone"`
}

// UpdateCustomerRequest 更新客户请求
type UpdateCustomerRequest struct {
	Email       string `json:"email" binding:"omitempty,email"`
	DisplayName string `json:"display_name"`
	FullName    string `json:"full_name"`
	CompanyCode string `json:"company_code"`
	Company     string `json:"company"`
	Phone       string `json:"phone"`
	Role        string `json:"role"`
}
