package function

import (
	"context"
	"fmt"
	"sync-flow/log"
	"sync-flow/sf"
)

type SfFunctionE struct {
	BaseFunction
}

func (f *SfFunctionE) Call(ctx context.Context, flow sf.Flow) error {
	log.GetLogger().InfoF("SfFunctionE, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法
	//处理业务数据
	for _, row := range flow.Input() {
		fmt.Printf("In SfFunctionE, row = %+v\n", row)
	}

	return nil
}
