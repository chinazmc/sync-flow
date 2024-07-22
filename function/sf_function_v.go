package function

import (
	"context"
	"fmt"
	"sync-flow/sf"
)

type SfFunctionV struct {
	BaseFunction
}

func (f *SfFunctionV) Call(ctx context.Context, flow sf.Flow) error {
	fmt.Printf("SfFunctionV, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法

	return nil
}
