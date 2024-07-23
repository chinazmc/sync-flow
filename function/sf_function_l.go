package function

import (
	"context"
	"fmt"
	"sync-flow/sf"
)

type SfFunctionL struct {
	BaseFunction
}

func (f *SfFunctionL) Call(ctx context.Context, flow sf.Flow) error {
	fmt.Printf("SfFunctionL, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法

	return nil
}