package function

import (
	"context"
	"fmt"
	"sync-flow/sf"
)

type SfFunctionC struct {
	BaseFunction
}

func (f *SfFunctionC) Call(ctx context.Context, flow sf.Flow) error {
	fmt.Printf("SfFunction_C, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法

	return nil
}
