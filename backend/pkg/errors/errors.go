package errors

// 错误码定义
const (
	Success            = 0
	ErrorInvalidParams = 400
	ErrorUnauthorized  = 401
	ErrorForbidden     = 403
	ErrorNotFound      = 404
	ErrorServerError   = 500
	ErrorDatabase      = 1001
	ErrorRedis         = 1002
)

// 错误信息映射
var ErrorMsg = map[int]string{
	Success:            "success",
	ErrorInvalidParams: "请求参数错误",
	ErrorUnauthorized:  "未授权",
	ErrorForbidden:     "禁止访问",
	ErrorNotFound:      "资源不存在",
	ErrorServerError:   "服务器内部错误",
	ErrorDatabase:      "数据库错误",
	ErrorRedis:         "Redis错误",
}

// GetErrorMsg 获取错误信息
func GetErrorMsg(code int) string {
	if msg, ok := ErrorMsg[code]; ok {
		return msg
	}
	return "未知错误"
}
