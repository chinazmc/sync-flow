package function

import (
	"context"
	"sync-flow/log"
	"sync-flow/sf"
)

type SfFunctionC struct {
	BaseFunction
}

func (f *SfFunctionC) Call(ctx context.Context, flow sf.Flow) error {
	log.GetLogger().InfoF("SfFunctionC, flow = %+v\n", flow)

	// 通过SfPool 路由到具体的执行计算Function中
	if err := sf.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.GetLogger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}

	return nil
}
func NewSfFunctionC() sf.Function {
	f := new(SfFunctionC)

	// 初始化metaData
	f.metaData = make(map[string]interface{})

	return f
}
