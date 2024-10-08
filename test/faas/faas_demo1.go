package faas

import (
	"context"
	"fmt"
	"sync-flow/sf"
)

// type FaaS func(context.Context, Flow) error

func FuncDemo1Handler(ctx context.Context, flow sf.Flow) error {
	fmt.Println("---> Call funcName1Handler ----")

	for index, row := range flow.Input() {
		// 打印数据
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FuncName, flow.GetThisFunction().GetId(), row)
		fmt.Println(str)

		// 计算结果数据
		resultStr := fmt.Sprintf("data from funcName[%s], index = %d", flow.GetThisFuncConf().FuncName, index)

		// 提交结果数据
		_ = flow.CommitRow(resultStr)
	}

	return nil
}
