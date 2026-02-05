package main

import (
	"net/http"

	agentErrors "github.com/YoungBoyGod/remotegpu-agent/internal/errors"
	"github.com/gin-gonic/gin"
)

type agentResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func respondSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, agentResponse{
		Success: true,
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func respondError(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, agentResponse{
		Success: false,
		Code:    code,
		Message: message,
	})
}

func respondErrorCode(c *gin.Context, httpStatus int, code int) {
	respondError(c, httpStatus, code, agentErrors.Message(code))
}

func respondWithData(c *gin.Context, success bool, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, agentResponse{
		Success: success,
		Code:    code,
		Message: message,
		Data:    data,
	})
}
