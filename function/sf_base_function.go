package function

import (
	"context"
	"errors"
	"sync"
	"sync-flow/common"
	"sync-flow/config"
	"sync-flow/id"
	"sync-flow/sf"
)

type BaseFunction struct {
	// Id , SfFunction的实例ID，用于SfFlow内部区分不同的实例对象
	Id     string
	Config *config.SfFuncConfig

	// flow
	Flow sf.Flow //上下文环境SfFlow
	// connector
	connector sf.Connector

	// Function的自定义临时数据
	metaData map[string]interface{}
	// 管理metaData的读写锁
	mLock sync.RWMutex

	// link
	N sf.Function //下一个流计算Function
	P sf.Function //上一个流计算Function
}

// Call
// BaseFunction 为空实现，目的为了让其他具体类型的SfFunction，如SfFunction_V 来继承BaseFuncion来重写此方法
func (base *BaseFunction) Call(ctx context.Context, flow sf.Flow) error { return nil }

func (base *BaseFunction) Next() sf.Function {
	return base.N
}

func (base *BaseFunction) Prev() sf.Function {
	return base.P
}

func (base *BaseFunction) SetN(f sf.Function) {
	base.N = f
}

func (base *BaseFunction) SetP(f sf.Function) {
	base.P = f
}

func (base *BaseFunction) SetConfig(s *config.SfFuncConfig) error {
	if s == nil {
		return errors.New("SfFuncConfig is nil")
	}

	base.Config = s

	return nil
}

func (base *BaseFunction) GetId() string {
	return base.Id
}

func (base *BaseFunction) GetPrevId() string {
	if base.P == nil {
		//Function为首结点
		return common.FunctionLinkListFirstVirtualNode
	}
	return base.P.GetId()
}

func (base *BaseFunction) GetNextId() string {
	if base.N == nil {
		//Function为尾结点
		return common.FunctionLinkListLastVirtualNode
	}
	return base.N.GetId()
}

func (base *BaseFunction) GetConfig() *config.SfFuncConfig {
	return base.Config
}

func (base *BaseFunction) SetFlow(f sf.Flow) error {
	if f == nil {
		return errors.New("SfFlow is nil")
	}
	base.Flow = f
	return nil
}

func (base *BaseFunction) GetFlow() sf.Flow {
	return base.Flow
}

func (base *BaseFunction) CreateId() {
	base.Id = id.SfID(common.SfIdTypeFunction)
}

// NewSfFunction 创建一个NsFunction
// flow: 当前所属的flow实例
// s : 当前function的配置策略
func NewSfFunction(flow sf.Flow, config *config.SfFuncConfig) sf.Function {
	var f sf.Function

	//工厂生产泛化对象
	switch common.SfMode(config.FuncMode) {
	case common.Verify:
		f = NewSfFunctionV() // +++
	case common.Save:
		f = NewSfFunctionS() // +++
	case common.Load:
		f = NewSfFunctionL() // +++
	case common.Calculate:
		f = NewSfFunctionC() // +++
	case common.Expand:
		f = NewSfFunctionE() // +++
	default:
		//LOG ERROR
		return nil
	}

	// 生成随机实例唯一ID
	f.CreateId()

	// 设置基础信息属性
	if err := f.SetConfig(config); err != nil {
		panic(err)
	}

	// 设置Flow
	if err := f.SetFlow(flow); err != nil {
		panic(err)
	}

	return f
}

// AddConnector 给当前Function实例添加一个Connector
func (base *BaseFunction) AddConnector(conn sf.Connector) error {
	if conn == nil {
		return errors.New("conn is nil")
	}

	base.connector = conn

	return nil
}

// GetConnector 获取当前Function实例所关联的Connector
func (base *BaseFunction) GetConnector() sf.Connector {
	return base.connector
}

// GetMetaData 得到当前Function的临时数据
func (base *BaseFunction) GetMetaData(key string) interface{} {
	base.mLock.RLock()
	defer base.mLock.RUnlock()

	data, ok := base.metaData[key]
	if !ok {
		return nil
	}

	return data
}

// SetMetaData 设置当前Function的临时数据
func (base *BaseFunction) SetMetaData(key string, value interface{}) {
	base.mLock.Lock()
	defer base.mLock.Unlock()

	base.metaData[key] = value
}
