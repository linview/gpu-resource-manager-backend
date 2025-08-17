package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 日志接口
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

// LoggerWrapper zap.SugaredLogger的包装器
type LoggerWrapper struct {
	sugar *zap.SugaredLogger
}

// Info 实现Logger接口的Info方法
func (l *LoggerWrapper) Info(msg string, fields ...interface{}) {
	l.sugar.Infow(msg, fields...)
}

// Error 实现Logger接口的Error方法
func (l *LoggerWrapper) Error(msg string, fields ...interface{}) {
	l.sugar.Errorw(msg, fields...)
}

// Warn 实现Logger接口的Warn方法
func (l *LoggerWrapper) Warn(msg string, fields ...interface{}) {
	l.sugar.Warnw(msg, fields...)
}

// Debug 实现Logger接口的Debug方法
func (l *LoggerWrapper) Debug(msg string, fields ...interface{}) {
	l.sugar.Debugw(msg, fields...)
}

// Fatal 实现Logger接口的Fatal方法
func (l *LoggerWrapper) Fatal(msg string, fields ...interface{}) {
	l.sugar.Fatalw(msg, fields...)
}

// New 创建新的日志实例
func New(level string) Logger {
	config := zap.NewProductionConfig()
	
	// 设置日志级别
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &LoggerWrapper{
		sugar: logger.Sugar(),
	}
}
