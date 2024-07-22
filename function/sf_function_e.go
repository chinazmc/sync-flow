package function

import (
	"context"
	"fmt"
	"sync-flow/sf"
)

type SfFunctionE struct {
	BaseFunction
}

func (f *SfFunctionE) Call(ctx context.Context, flow sf.Flow) error {
	fmt.Printf("SfFunctionE, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法

	return nil
}
