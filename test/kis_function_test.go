package test

import (
	"context"
	"sync-flow/common"
	"sync-flow/config"
	"sync-flow/flow"
	"sync-flow/function"
	"testing"
)

func TestNewSfFunction(t *testing.T) {
	ctx := context.Background()

	// 1. 创建一个SfFunction配置实例
	source := config.SfSource{
		Name: "订单数据",
		Must: []string{"order_id", "user_id"},
	}

	myFuncConfig1 := config.NewFuncConfig("funcName1", common.C, &source, nil)
	if myFuncConfig1 == nil {
		panic("myFuncConfig1 is nil")
	}

	// 2. 创建一个 SfFlow 配置实例
	myFlowConfig1 := config.NewFlowConfig("flowName1", common.FlowEnable)

	// 3. 创建一个SfFlow对象
	flow1 := flow.NewSfFlow(myFlowConfig1)

	// 4. 创建一个SfFunction对象
	func1 := function.NewSfFunction(flow1, myFuncConfig1)

	if err := func1.Call(ctx, flow1); err != nil {
		t.Errorf("func1.Call() error = %v", err)
	}
}
