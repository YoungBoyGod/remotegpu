package errors

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// HandleNotFoundError 处理记录不存在错误
// 如果是 gorm.ErrRecordNotFound，返回 AppError
func HandleNotFoundError(err error, entityName string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return New(ErrorNotFound, fmt.Sprintf("%s不存在", entityName))
	}
	return err
}

// HandleDuplicateError 处理重复记录错误
// 检查是否是唯一约束冲突，返回 AppError
func HandleDuplicateError(err error, fieldName string) error {
	if err == nil {
		return nil
	}
	errMsg := err.Error()
	if strings.Contains(errMsg, "duplicate") ||
	   strings.Contains(errMsg, "UNIQUE constraint failed") ||
	   strings.Contains(errMsg, "Duplicate entry") {
		return New(ErrorInvalidParams, fmt.Sprintf("%s已存在", fieldName))
	}
	return err
}

// HandleForeignKeyError 处理外键约束错误，返回 AppError
func HandleForeignKeyError(err error, message string) error {
	if err == nil {
		return nil
	}
	errMsg := err.Error()
	if strings.Contains(errMsg, "foreign key constraint") ||
	   strings.Contains(errMsg, "FOREIGN KEY constraint failed") {
		return New(ErrorInvalidParams, message)
	}
	return err
}

// HandleDatabaseError 处理数据库错误，返回 AppError
func HandleDatabaseError(err error, operation string) error {
	if err == nil {
		return nil
	}
	return WrapWithMessage(ErrorDatabase, fmt.Sprintf("数据库%s失败", operation), err)
}
