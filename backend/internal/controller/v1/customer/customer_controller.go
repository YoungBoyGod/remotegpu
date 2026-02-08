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

// List 获取客户列表
// @Summary 获取客户列表
// @Description 分页获取所有客户信息
// @Tags Admin - Customers
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/customers [get]
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

// Create 创建客户
// @Summary 创建客户
// @Description 创建新客户账号，如未指定密码则使用默认密码
// @Tags Admin - Customers
// @Accept json
// @Produce json
// @Param request body v1.CreateCustomerRequest true "创建客户请求"
// @Security Bearer
// @Success 200 {object} entity.Customer
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/customers [post]
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
// @Summary 获取客户详情
// @Description 根据客户 ID 获取详细信息，包含机器分配信息
// @Tags Admin - Customers
// @Accept json
// @Produce json
// @Param id path int true "客户 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /admin/customers/{id} [get]
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
// @Summary 更新客户信息
// @Description 根据客户 ID 更新客户字段，仅更新传入的非空字段
// @Tags Admin - Customers
// @Accept json
// @Produce json
// @Param id path int true "客户 ID"
// @Param request body v1.UpdateCustomerRequest true "更新客户请求"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/customers/{id} [put]
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
// @Summary 禁用客户
// @Description 将客户状态设置为 suspended
// @Tags Admin - Customers
// @Accept json
// @Produce json
// @Param id path int true "客户 ID"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/customers/{id}/disable [post]
func (c *CustomerController) Disable(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的客户 ID")
		return
	}

	if err := c.customerService.UpdateStatus(ctx, uint(id), "suspended"); err != nil {
		c.Error(ctx, 500, "Failed to disable customer")
		return
	}
	c.Success(ctx, nil)
}

// Enable 启用客户
// @Summary 启用客户
// @Description 将客户状态设置为 active
// @Tags Admin - Customers
// @Accept json
// @Produce json
// @Param id path int true "客户 ID"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/customers/{id}/enable [post]
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
// @Summary 更新客户配额
// @Description 更新指定客户的 GPU 和存储配额
// @Tags Admin - Customers
// @Accept json
// @Produce json
// @Param id path int true "客户 ID"
// @Param request body v1.UpdateQuotaRequest true "更新配额请求"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/customers/{id}/quota [put]
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
// @Summary 获取客户资源使用统计
// @Description 获取指定客户的 GPU、存储等资源使用情况
// @Tags Admin - Customers
// @Accept json
// @Produce json
// @Param id path int true "客户 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/customers/{id}/usage [get]
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
