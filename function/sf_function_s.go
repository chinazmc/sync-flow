package function

import (
	"context"
	"sync-flow/log"
	"sync-flow/sf"
)

type SfFunctionS struct {
	BaseFunction
}

func (f *SfFunctionS) Call(ctx context.Context, flow sf.Flow) error {
	log.GetLogger().InfoF("SfFunctionS, flow = %+v\n", flow)

	// 通过SfPool 路由到具体的执行计算Function中
	if err := sf.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.GetLogger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}

	return nil
}
