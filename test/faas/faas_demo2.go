package faas

import (
	"context"
	"fmt"
	"sync-flow/log"
	"sync-flow/sf"
)

// type FaaS func(context.Context, Flow) error

func FuncDemo2Handler(ctx context.Context, flow sf.Flow) error {
	fmt.Println("---> Call funcName2Handler ----")

	for index, row := range flow.Input() {
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FName, flow.GetThisFunction().GetId(), row)
		fmt.Println(str)

		conn, err := flow.GetConnector()
		if err != nil {
			log.GetLogger().ErrorFX(ctx, "FuncDemo2Handler(): GetConnector err = %s\n", err.Error())
			return err
		}

		if _, err := conn.Call(ctx, flow, row); err != nil {
			log.GetLogger().ErrorFX(ctx, "FuncDemo2Handler(): Call err = %s\n", err.Error())
			return err
		}

		// 计算结果数据
		resultStr := fmt.Sprintf("data from funcName[%s], index = %d", flow.GetThisFuncConf().FName, index)

		// 提交结果数据
		_ = flow.CommitRow(resultStr)
	}

	return nil
}
