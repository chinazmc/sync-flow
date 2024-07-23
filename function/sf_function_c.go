package function

import (
	"context"
	"fmt"
	"sync-flow/log"
	"sync-flow/sf"
)

type SfFunctionC struct {
	BaseFunction
}

func (f *SfFunctionC) Call(ctx context.Context, flow sf.Flow) error {
	log.GetLogger().InfoF("SfFunctionC, flow = %+v\n", flow)

	//TODO 调用具体的Function执行方法
	//处理业务数据
	for i, row := range flow.Input() {
		fmt.Printf("In SfFunctionC, row = %+v\n", row)

		// 提交本层计算结果数据
		_ = flow.CommitRow("Data From SfFunctionC, index " + " " + fmt.Sprintf("%d", i))
	}

	return nil
}
