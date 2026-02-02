package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger 初始化日志
func InitLogger(mode string) error {
	var config zap.Config

	if mode == "release" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.OutputPaths = []string{"stdout", "./logs/app.log"}
	config.ErrorOutputPaths = []string{"stderr", "./logs/error.log"}

	// 创建日志目录
	if err := os.MkdirAll("./logs", 0755); err != nil {
		return err
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	if Logger == nil {
		// 如果未初始化，使用默认配置
		Logger, _ = zap.NewDevelopment()
	}
	return Logger
}

// Info 记录信息日志
func Info(msg string, fields ...interface{}) {
	GetLogger().Sugar().Infow(msg, fields...)
}

// Error 记录错误日志
func Error(msg string, fields ...interface{}) {
	GetLogger().Sugar().Errorw(msg, fields...)
}

// Warn 记录警告日志
func Warn(msg string, fields ...interface{}) {
	GetLogger().Sugar().Warnw(msg, fields...)
}

// Debug 记录调试日志
func Debug(msg string, fields ...interface{}) {
	GetLogger().Sugar().Debugw(msg, fields...)
}
