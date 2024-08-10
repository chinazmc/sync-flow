package log

import "context"

type SfLogger interface {
	InfoFX(ctx context.Context, str string, v ...interface{})
	ErrorFX(ctx context.Context, str string, v ...interface{})
	DebugFX(ctx context.Context, str string, v ...interface{})
	InfoF(str string, v ...interface{})
	ErrorF(str string, v ...interface{})
	DebugF(str string, v ...interface{})
	SetDebugMode(enable bool)
}

// sfLog 默认的Log 对象
var sfLog SfLogger

// SetLogger 设置Log对象, 可以是用户自定义的Logger对象
func SetLogger(newlog SfLogger) {
	sfLog = newlog
}

// Logger 获取到Log对象
func GetLogger() SfLogger {
	return sfLog
}
