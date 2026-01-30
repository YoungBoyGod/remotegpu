package k8s

import (
	"errors"
	"fmt"
)

// K8s相关错误定义
var (
	// ErrPodNotFound Pod不存在错误
	ErrPodNotFound = errors.New("pod not found")

	// ErrPodCreationFailed Pod创建失败错误
	ErrPodCreationFailed = errors.New("pod creation failed")

	// ErrPodDeletionFailed Pod删除失败错误
	ErrPodDeletionFailed = errors.New("pod deletion failed")

	// ErrConnectionFailed K8s连接失败错误
	ErrConnectionFailed = errors.New("kubernetes connection failed")

	// ErrInvalidConfig 无效配置错误
	ErrInvalidConfig = errors.New("invalid kubernetes config")

	// ErrClientNotInitialized 客户端未初始化错误
	ErrClientNotInitialized = errors.New("kubernetes client not initialized")

	// ErrPodStatusTimeout Pod状态超时错误
	ErrPodStatusTimeout = errors.New("pod status timeout")

	// ErrLogsFetchFailed 日志获取失败错误
	ErrLogsFetchFailed = errors.New("failed to fetch pod logs")
)

// WrapError 包装错误，添加上下文信息
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// WrapErrorf 包装错误，支持格式化消息
func WrapErrorf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	message := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %w", message, err)
}

// IsPodNotFound 判断是否为Pod不存在错误
func IsPodNotFound(err error) bool {
	return errors.Is(err, ErrPodNotFound)
}

// IsConnectionFailed 判断是否为连接失败错误
func IsConnectionFailed(err error) bool {
	return errors.Is(err, ErrConnectionFailed)
}
