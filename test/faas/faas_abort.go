package faas

import (
	"context"
	"fmt"
	"sync-flow/sf"
)

// type FaaS func(context.Context, Flow) error

func AbortFuncHandler(ctx context.Context, flow sf.Flow) error {
	fmt.Println("---> Call AbortFuncHandler ----")

	for _, row := range flow.Input() {
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FuncName, flow.GetThisFunction().GetId(), row)
		fmt.Println(str)
	}

	return flow.Next(sf.ActionAbort)
}
