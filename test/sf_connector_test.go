package test

import (
	"context"
	"sync-flow/common"
	"sync-flow/config"
	"sync-flow/flow"
	"sync-flow/sf"
	"sync-flow/test/caas"
	"sync-flow/test/faas"
	"testing"
)

func TestNewSfConnector(t *testing.T) {

	ctx := context.Background()

	// 0. 注册Function 回调业务
	sf.Pool().FaaS("funcName1", faas.FuncDemo1Handler)
	sf.Pool().FaaS("funcName2", faas.FuncDemo2Handler)
	sf.Pool().FaaS("funcName3", faas.FuncDemo3Handler)

	// 0. 注册ConnectorInit 和 Connector 回调业务
	sf.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	sf.Pool().CaaS("ConnName1", "funcName2", common.Save, caas.CaasDemoHanler1)

	// 1. 创建3个SfFunction配置实例, 其中myFuncConfig2 有Connector配置
	source1 := config.SfSource{
		Name: "订单数据",
		Must: []string{"order_id", "user_id"},
	}

	source2 := config.SfSource{
		Name: "用户订单错误率",
		Must: []string{"order_id", "user_id"},
	}

	myFuncConfig1 := config.NewFuncConfig("funcName1", common.Calculate, &source1, nil)
	if myFuncConfig1 == nil {
		panic("myFuncConfig1 is nil")
	}

	option := config.SfFuncOption{
		CName: "ConnName1",
	}

	myFuncConfig2 := config.NewFuncConfig("funcName2", common.Save, &source2, &option)
	if myFuncConfig2 == nil {
		panic("myFuncConfig2 is nil")
	}

	myFuncConfig3 := config.NewFuncConfig("funcName3", common.Expand, &source2, nil)
	if myFuncConfig3 == nil {
		panic("myFuncConfig3 is nil")
	}

	// 2. 创建一个SfConnector配置实例
	myConnConfig1 := config.NewConnConfig("ConnName1", "0.0.0.0:9998", common.REDIS, "redis-key", nil)
	if myConnConfig1 == nil {
		panic("myConnConfig1 is nil")
	}

	// 3. 将SfConnector配置实例绑定到SfFunction配置实例上
	_ = myFuncConfig2.AddConnConfig(myConnConfig1)

	// 4. 创建一个 SfFlow 配置实例
	myFlowConfig1 := config.NewFlowConfig("flowName1", common.FlowEnable)

	// 5. 创建一个SfFlow对象
	flow1 := flow.NewSfFlow(myFlowConfig1)

	// 6. 拼接Functioin 到 Flow 上
	if err := flow1.Link(myFuncConfig1, nil); err != nil {
		panic(err)
	}
	if err := flow1.Link(myFuncConfig2, nil); err != nil {
		panic(err)
	}
	if err := flow1.Link(myFuncConfig3, nil); err != nil {
		panic(err)
	}

	// 7. 提交原始数据
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 8. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}
