package v1

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service"
	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	BaseController
	customerService *service.CustomerService
}

func NewCustomerController(cs *service.CustomerService) *CustomerController {
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

type CreateCustomerRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

func (c *CustomerController) Create(ctx *gin.Context) {
	var req CreateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	customer := &entity.Customer{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
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
