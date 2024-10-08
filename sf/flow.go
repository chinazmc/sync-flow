package sf

import (
	"context"
	"sync-flow/common"
	"sync-flow/config"
	"time"
)

type Flow interface {
	// Run 启动flow，获取数据进行func 链表的顺序处理
	Run(ctx context.Context) error
	// Link 按照func 配置来将func 加入到flow
	Link(fConf *config.SfFuncConfig, fParams config.FuncParam) error
	// AppendNewFunction 将Function追加到到Flow中
	AppendNewFunction(fConf *config.SfFuncConfig, fParams config.FuncParam) error
	//CommitRow 提交数据到即将执行的function层
	CommitRow(row interface{}) error
	// CommitRowBatch 提交Flow数据到即将执行的Function层(批量提交)
	CommitRowBatch(row interface{}) error
	// Input 得到flow当前执行Function的正在处理的数据
	Input() common.SfRowArr

	// GetName 获取Flow的名称
	GetName() string
	// GetThisFunction 得到当前正在执行的Function
	GetThisFunction() Function
	// GetThisFuncConf 得到当前正在执行的Function的配置
	GetThisFuncConf() *config.SfFuncConfig
	GetConnector() (Connector, error)
	// GetConnConf 得到当前正在执行的Function的Connector的配置
	GetConnConf() (*config.SfConnConfig, error)

	// GetConfig 得到当前Flow的配置
	GetConfig() *config.SfFlowConfig
	// GetFuncConfigByName 得到当前Flow的配置
	GetFuncConfigByName(funcName string) *config.SfFuncConfig
	// Next 当前Flow执行到的Function进入下一层Function所携带的Action动作
	Next(acts ...ActionFunc) error

	// GetCacheData 得到当前Flow的缓存数据
	GetCacheData(key string) interface{}
	// SetCacheData 设置当前Flow的缓存数据
	SetCacheData(key string, value interface{}, Exp time.Duration)

	// GetMetaData 得到当前Flow的临时数据
	GetMetaData(key string) interface{}
	// SetMetaData 设置当前Flow的临时数据
	SetMetaData(key string, value interface{})

	// GetFuncParam 得到Flow的当前正在执行的Function的配置默认参数，取出一对key-value
	GetFuncParam(key string) string
	// GetFuncParamAll 得到Flow的当前正在执行的Function的配置默认参数，取出全部Key-Value
	GetFuncParamAll() config.FuncParam
	// GetId 得到Flow的Id
	GetId() string
	// Fork 得到Flow的一个副本(深拷贝)
	Fork(ctx context.Context) Flow
	// GetFuncParamsAllFuncs 得到Flow中所有Function的FuncParams，取出全部Key-Value
	GetFuncParamsAllFuncs() map[string]config.FuncParam
}
