package function

import (
	"context"
	"fmt"
	"sync-flow/sf"
)

type SfFunctionS struct {
	BaseFunction
}

func (f *SfFunctionS) Call(ctx context.Context, flow sf.Flow) error {
	fmt.Printf("SfFunctionS, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法

	return nil
}
