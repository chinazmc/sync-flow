package function

import (
	"context"
	"sync-flow/log"
	"sync-flow/sf"
)

type SfFunctionE struct {
	BaseFunction
}

func (f *SfFunctionE) Call(ctx context.Context, flow sf.Flow) error {
	// 通过SfPool 路由到具体的执行计算Function中
	if err := sf.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.GetLogger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}

	return nil
}
func NewSfFunctionE() sf.Function {
	f := new(SfFunctionE)

	// 初始化metaData
	f.metaData = make(map[string]interface{})

	return f
}
