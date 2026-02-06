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

func (c *CustomerController) Disable(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	if err := c.customerService.UpdateStatus(ctx, uint(id), "suspended"); err != nil {
		c.Error(ctx, 500, "Failed to disable customer")
		return
	}
	c.Success(ctx, nil)
}
