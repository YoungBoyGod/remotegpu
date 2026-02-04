package errors

// 错误码定义
const (
	// 通用错误码 (0-999)
	Success            = 0
	ErrorInvalidParams = 400
	ErrorUnauthorized  = 401
	ErrorForbidden     = 403
	ErrorNotFound      = 404
	ErrorServerError   = 500

	// 基础设施错误 (1000-1999)
	ErrorDatabase = 1001
	ErrorRedis    = 1002
	ErrorCache    = 1003
	ErrorK8s      = 1004
	ErrorDocker   = 1005

	// 用户相关错误 (2000-2999)
	ErrorUserNotFound      = 2001
	ErrorUserExists        = 2002
	ErrorPasswordIncorrect = 2003
	ErrorTokenInvalid      = 2004
	ErrorTokenExpired      = 2005
	ErrorUserDisabled      = 2006

	// 工作空间相关错误 (3000-3999)
	ErrorWorkspaceNotFound     = 3001
	ErrorWorkspaceExists       = 3002
	ErrorWorkspacePermission   = 3003
	ErrorWorkspaceMemberExists = 3004
	ErrorWorkspaceMemberNotFound = 3005

	// 资源配额相关错误 (4000-4999)
	ErrorQuotaNotFound    = 4001
	ErrorQuotaExceeded    = 4002
	ErrorQuotaInvalid     = 4003
	ErrorInsufficientQuota = 4004

	// 环境相关错误 (5000-5999)
	ErrorEnvironmentNotFound   = 5001
	ErrorEnvironmentExists     = 5002
	ErrorEnvironmentCreating   = 5003
	ErrorEnvironmentStartFailed = 5004
	ErrorEnvironmentStopFailed  = 5005
	ErrorEnvironmentDeleteFailed = 5006
	ErrorEnvironmentInvalidStatus = 5007

	// 主机和GPU相关错误 (6000-6999)
	ErrorHostNotFound     = 6001
	ErrorHostOffline      = 6002
	ErrorGPUNotFound      = 6003
	ErrorGPUNotAvailable  = 6004
	ErrorGPUAllocationFailed = 6005

	// 存储相关错误 (7000-7999)
	ErrorStorageNotFound   = 7001
	ErrorStorageCreateFailed = 7002
	ErrorStorageMountFailed  = 7003
	ErrorStorageDeleteFailed = 7004

	// 网络相关错误 (8000-8999)
	ErrorPortNotAvailable = 8001
	ErrorPortAllocationFailed = 8002
	ErrorDNSConfigFailed  = 8003
	ErrorFirewallConfigFailed = 8004
)

// 错误信息映射
var ErrorMsg = map[int]string{
	// 通用错误
	Success:            "success",
	ErrorInvalidParams: "请求参数错误",
	ErrorUnauthorized:  "未授权",
	ErrorForbidden:     "禁止访问",
	ErrorNotFound:      "资源不存在",
	ErrorServerError:   "服务器内部错误",

	// 基础设施错误
	ErrorDatabase: "数据库错误",
	ErrorRedis:    "Redis错误",
	ErrorCache:    "缓存错误",
	ErrorK8s:      "Kubernetes错误",
	ErrorDocker:   "Docker错误",

	// 用户相关错误
	ErrorUserNotFound:      "用户不存在",
	ErrorUserExists:        "用户已存在",
	ErrorPasswordIncorrect: "密码错误",
	ErrorTokenInvalid:      "Token无效",
	ErrorTokenExpired:      "Token已过期",
	ErrorUserDisabled:      "账号已禁用",

	// 工作空间相关错误
	ErrorWorkspaceNotFound:       "工作空间不存在",
	ErrorWorkspaceExists:         "工作空间已存在",
	ErrorWorkspacePermission:     "无工作空间权限",
	ErrorWorkspaceMemberExists:   "成员已存在",
	ErrorWorkspaceMemberNotFound: "成员不存在",

	// 资源配额相关错误
	ErrorQuotaNotFound:     "配额不存在",
	ErrorQuotaExceeded:     "配额已超限",
	ErrorQuotaInvalid:      "配额参数无效",
	ErrorInsufficientQuota: "配额不足",

	// 环境相关错误
	ErrorEnvironmentNotFound:      "环境不存在",
	ErrorEnvironmentExists:        "环境已存在",
	ErrorEnvironmentCreating:      "环境正在创建中",
	ErrorEnvironmentStartFailed:   "环境启动失败",
	ErrorEnvironmentStopFailed:    "环境停止失败",
	ErrorEnvironmentDeleteFailed:  "环境删除失败",
	ErrorEnvironmentInvalidStatus: "环境状态无效",

	// 主机和GPU相关错误
	ErrorHostNotFound:        "主机不存在",
	ErrorHostOffline:         "主机离线",
	ErrorGPUNotFound:         "GPU不存在",
	ErrorGPUNotAvailable:     "GPU不可用",
	ErrorGPUAllocationFailed: "GPU分配失败",

	// 存储相关错误
	ErrorStorageNotFound:     "存储不存在",
	ErrorStorageCreateFailed: "存储创建失败",
	ErrorStorageMountFailed:  "存储挂载失败",
	ErrorStorageDeleteFailed: "存储删除失败",

	// 网络相关错误
	ErrorPortNotAvailable:     "端口不可用",
	ErrorPortAllocationFailed: "端口分配失败",
	ErrorDNSConfigFailed:      "DNS配置失败",
	ErrorFirewallConfigFailed: "防火墙配置失败",
}

// GetErrorMsg 获取错误信息
func GetErrorMsg(code int) string {
	if msg, ok := ErrorMsg[code]; ok {
		return msg
	}
	return "未知错误"
}

// AppError 应用错误类型
type AppError struct {
	Code    int    // 错误码
	Message string // 错误信息
	Err     error  // 原始错误
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap 实现 errors.Unwrap 接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建新的应用错误
func New(code int, message string) *AppError {
	if message == "" {
		message = GetErrorMsg(code)
	}
	return &AppError{
		Code:    code,
		Message: message,
		Err:     nil,
	}
}

// Wrap 包装错误
func Wrap(code int, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: GetErrorMsg(code),
		Err:     err,
	}
}

// WrapWithMessage 包装错误并自定义消息
func WrapWithMessage(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// IsAppError 判断是否为应用错误
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError 获取应用错误
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}
