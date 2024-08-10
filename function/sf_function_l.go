package function

import (
	"context"
	"sync-flow/log"
	"sync-flow/sf"
)

type SfFunctionL struct {
	BaseFunction
}

func (f *SfFunctionL) Call(ctx context.Context, flow sf.Flow) error {
	// 通过SfPool 路由到具体的执行计算Function中
	if err := sf.Pool().CallFunction(ctx, f.Config.FuncName, flow); err != nil {
		log.GetLogger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}

	return nil
}
func NewSfFunctionL() sf.Function {
	f := new(SfFunctionL)

	// 初始化metaData
	f.metaData = make(map[string]interface{})

	return f
}
