package common

import (
	"net/http"

	"github.com/YoungBoyGod/remotegpu/pkg/response"
	"github.com/gin-gonic/gin"
)

type BaseController struct{}

func (c *BaseController) Success(ctx *gin.Context, data interface{}) {
	response.Success(ctx, data)
}

func (c *BaseController) Error(ctx *gin.Context, code int, msg string) {
	// Assuming pkg/response has a standardized error response
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}