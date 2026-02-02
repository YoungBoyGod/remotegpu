package response

import (
	"net/http"

	"github.com/YoungBoyGod/remotegpu/pkg/errors"
	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// HandleError 统一错误处理
// 自动识别 AppError 类型，记录日志并返回响应
func HandleError(c *gin.Context, err error) {
	if err == nil {
		Success(c, nil)
		return
	}

	// 检查是否为 AppError
	if appErr := errors.GetAppError(err); appErr != nil {
		// 记录错误日志
		logger.Error("业务错误",
			"code", appErr.Code,
			"message", appErr.Message,
			"error", appErr.Err,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
		)

		Error(c, appErr.Code, appErr.Message)
		return
	}

	// 普通错误，返回服务器错误
	logger.Error("服务器错误",
		"error", err,
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
	)

	Error(c, errors.ErrorServerError, err.Error())
}

// FailWithAppError 使用 AppError 返回错误响应
func FailWithAppError(c *gin.Context, appErr *errors.AppError) {
	if appErr == nil {
		Error(c, errors.ErrorServerError, "未知错误")
		return
	}

	// 记录错误日志
	logger.Error("业务错误",
		"code", appErr.Code,
		"message", appErr.Message,
		"error", appErr.Err,
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
	)

	Error(c, appErr.Code, appErr.Message)
}
