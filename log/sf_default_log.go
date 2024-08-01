package log

import (
	"context"
	"fmt"
	"sync"
)

// defaultLog 默认提供的日志对象
type defaultLog struct {
	debugMode bool
	mu        sync.Mutex
}

func (log *defaultLog) SetDebugMode(enable bool) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.debugMode = enable
}
func (log *defaultLog) InfoF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
}

func (log *defaultLog) ErrorF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
}

func (log *defaultLog) DebugF(str string, v ...interface{}) {
	log.mu.Lock()
	defer log.mu.Unlock()
	if log.debugMode {
		fmt.Printf(str, v...)
		fmt.Printf("\n")
	}
}

func (log *defaultLog) InfoFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
}

func (log *defaultLog) ErrorFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
}

func (log *defaultLog) DebugFX(ctx context.Context, str string, v ...interface{}) {
	log.mu.Lock()
	defer log.mu.Unlock()
	if log.debugMode {
		fmt.Println(ctx)
		fmt.Printf(str, v...)
		fmt.Printf("\n")
	}
}

func init() {
	// 如果没有设置Logger, 则启动时使用默认的defaultLog对象
	if GetLogger() == nil {
		SetLogger(&defaultLog{})
	}
}
