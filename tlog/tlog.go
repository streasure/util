package tlog

import (
	"context"

	slog "github.com/streasure/treasure-slog"
)

// --- Re-export 全局函数 ---

// New 创建日志记录器
func New(configPath string) (slog.Logger, error) {
	return slog.New(configPath)
}

// Sync 同步并关闭 logger
func Sync() error {
	return slog.Sync()
}

// SetLevel 设置日志级别
func SetLevel(level string) {
	slog.SetLevel(level)
}

// GetLevel 获取日志级别
func GetLevel() string {
	return slog.GetLevel()
}

// Debug 调试日志
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// DebugContext 带上下文的调试日志
func DebugContext(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

// Info 信息日志
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// InfoContext 带上下文的信息日志
func InfoContext(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

// Warn 警告日志
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// WarnContext 带上下文的警告日志
func WarnContext(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

// Error 错误日志
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// ErrorContext 带上下文的错误日志
func ErrorContext(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

// Recover 恢复 panic
func Recover() {
	slog.Recover()
}
