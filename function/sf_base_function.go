package function

import (
	"context"
	"errors"
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
		return common.FunctionIdFirstVirtual
	}
	return base.P.GetId()
}

func (base *BaseFunction) GetNextId() string {
	if base.N == nil {
		//Function为尾结点
		return common.FunctionIdLastVirtual
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
	switch common.SfMode(config.FMode) {
	case common.V:
		f = new(SfFunctionV)
		break
	case common.S:
		f = new(SfFunctionS)
	case common.L:
		f = new(SfFunctionL)
	case common.C:
		f = new(SfFunctionC)
	case common.E:
		f = new(SfFunctionE)
	default:
		//LOG ERROR
		return nil
	}

	// 生成随机实例唯一ID
	f.CreateId()

	//设置基础信息属性
	if err := f.SetConfig(config); err != nil {
		panic(err)
	}

	if err := f.SetFlow(flow); err != nil {
		panic(err)
	}

	return f
}
