package sf

import (
	"context"
	"sync-flow/common"
	"sync-flow/config"
)

type Flow interface {
	// Run 调度Flow，依次调度Flow中的Function并且执行
	Run(ctx context.Context) error
	// Link 将Flow中的Function按照配置文件中的配置进行连接
	Link(fConf *config.SfFuncConfig, fParams config.FParam) error
	//CommitRow 提交flow数据到即将执行的function层
	CommitRow(row interface{}) error
	// Input 得到flow当前执行Function的输入源数据
	Input() common.SfRowArr

	// GetName 得到Flow的名称
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
}
