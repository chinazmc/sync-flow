package function

import (
	"context"
	"sync-flow/log"
	"sync-flow/sf"
)

type SfFunctionV struct {
	BaseFunction
}

func (f *SfFunctionV) Call(ctx context.Context, flow sf.Flow) error {
	// 通过SfPool 路由到具体的执行计算Function中
	if err := sf.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.GetLogger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}

	return nil
}
func NewSfFunctionV() sf.Function {
	f := new(SfFunctionV)

	// 初始化metaData
	f.metaData = make(map[string]interface{})

	return f
}
