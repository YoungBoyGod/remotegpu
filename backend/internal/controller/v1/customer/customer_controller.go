package customer

import (
	"strconv"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceCustomer "github.com/YoungBoyGod/remotegpu/internal/service/customer"
	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	common.BaseController
	customerService *serviceCustomer.CustomerService
}

func NewCustomerController(cs *serviceCustomer.CustomerService) *CustomerController {
	return &CustomerController{
		customerService: cs,
	}
}

func (c *CustomerController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	customers, total, err := c.customerService.ListCustomers(ctx, page, pageSize)
	if err != nil {
		c.Error(ctx, 500, "Failed to list customers")
		return
	}

	c.Success(ctx, gin.H{
		"list":      customers,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (c *CustomerController) Create(ctx *gin.Context) {
	var req apiV1.CreateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}
	mustChangePassword := false
	if req.Password == "" {
		req.Password = "ChangeME_123"
		mustChangePassword = true
	} else if req.Password == "ChangeME_123" {
		mustChangePassword = true
	}

	customer := &entity.Customer{
		Username:           req.Username,
		Email:              req.Email,
		Role:               req.Role,
		DisplayName:        req.DisplayName,
		FullName:           req.FullName,
		CompanyCode:        req.CompanyCode,
		Company:            req.Company,
		Phone:              req.Phone,
		MustChangePassword: mustChangePassword,
	}
	if customer.Role == "" {
		customer.Role = "customer_owner"
	}

	if err := c.customerService.CreateCustomer(ctx, customer, req.Password); err != nil {
		c.Error(ctx, 500, "Failed to create customer")
		return
	}

	c.Success(ctx, customer)
}

// Detail 获取客户详情（包含机器分配信息）
func (c *CustomerController) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "Invalid customer ID")
		return
	}

	detail, err := c.customerService.GetCustomerDetail(ctx, uint(id))
	if err != nil {
		c.Error(ctx, 404, "Customer not found")
		return
	}
	c.Success(ctx, detail)
}

// Update 更新客户信息
func (c *CustomerController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "Invalid customer ID")
		return
	}

	var req apiV1.UpdateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	fields := make(map[string]interface{})
	if req.Email != "" {
		fields["email"] = req.Email
	}
	if req.DisplayName != "" {
		fields["display_name"] = req.DisplayName
	}
	if req.FullName != "" {
		fields["full_name"] = req.FullName
	}
	if req.CompanyCode != "" {
		fields["company_code"] = req.CompanyCode
	}
	if req.Company != "" {
		fields["company"] = req.Company
	}
	if req.Phone != "" {
		fields["phone"] = req.Phone
	}
	if req.Role != "" {
		fields["role"] = req.Role
	}

	if len(fields) == 0 {
		c.Error(ctx, 400, "No fields to update")
		return
	}

	if err := c.customerService.UpdateCustomer(ctx, uint(id), fields); err != nil {
		c.Error(ctx, 500, "Failed to update customer")
		return
	}
	c.Success(ctx, nil)
}

// Disable 禁用客户
func (c *CustomerController) Disable(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	if err := c.customerService.UpdateStatus(ctx, uint(id), "suspended"); err != nil {
		c.Error(ctx, 500, "Failed to disable customer")
		return
	}
	c.Success(ctx, nil)
}

// Enable 启用客户
func (c *CustomerController) Enable(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "Invalid customer ID")
		return
	}

	if err := c.customerService.UpdateStatus(ctx, uint(id), "active"); err != nil {
		c.Error(ctx, 500, "Failed to enable customer")
		return
	}
	c.Success(ctx, nil)
}

// UpdateQuota 更新客户配额
func (c *CustomerController) UpdateQuota(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的客户 ID")
		return
	}

	var req apiV1.UpdateQuotaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.customerService.UpdateQuota(ctx, uint(id), req.QuotaGPU, req.QuotaStorage); err != nil {
		c.Error(ctx, 500, "更新配额失败")
		return
	}
	c.Success(ctx, nil)
}

// ResourceUsage 获取客户资源使用统计
func (c *CustomerController) ResourceUsage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的客户 ID")
		return
	}

	usage, err := c.customerService.GetResourceUsage(ctx, uint(id))
	if err != nil {
		c.Error(ctx, 500, "获取资源使用统计失败")
		return
	}
	c.Success(ctx, usage)
}
