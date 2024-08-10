package caas

import (
	"context"
	"fmt"
	"sync-flow/sf"
)

// type CaaS func(context.Context, Connector, Function, Flow, interface{}) error

func CaasDemoHanler1(ctx context.Context, conn sf.Connector, fn sf.Function, flow sf.Flow, args interface{}) (interface{}, error) {
	fmt.Printf("===> In CaasDemoHanler1: flowName: %s, cName:%s, fnName:%s, mode:%s\n",
		flow.GetName(), conn.GetName(), fn.GetConfig().FuncName, fn.GetConfig().FuncMode)

	fmt.Printf("===> Call Connector CaasDemoHanler1, args from funciton: %s\n", args)

	return nil, nil
}
