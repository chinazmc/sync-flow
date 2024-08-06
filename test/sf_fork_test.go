package test

import (
	"context"
	"sync-flow/common"
	"sync-flow/file"
	"sync-flow/sf"
	"sync-flow/test/caas"
	"sync-flow/test/faas"
	"testing"
)

func TestForkFlow(t *testing.T) {
	ctx := context.Background()

	// 0. 注册Function 回调业务
	sf.Pool().FaaS("funcName1", faas.FuncDemo1Handler)
	sf.Pool().FaaS("funcName2", faas.FuncDemo2Handler)
	sf.Pool().FaaS("funcName3", faas.FuncDemo3Handler)

	// 0. 注册ConnectorInit 和 Connector 回调业务
	sf.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	sf.Pool().CaaS("ConnName1", "funcName2", common.Save, caas.CaasDemoHanler1)

	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("./load_conf/"); err != nil {
		panic(err)
	}

	// 2. 获取Flow
	flow1 := sf.Pool().GetFlow("flowName1")

	flow1Clone1 := flow1.Fork(ctx)

	// 3. 提交原始数据
	_ = flow1Clone1.CommitRow("This is Data1 from Test")
	_ = flow1Clone1.CommitRow("This is Data2 from Test")
	_ = flow1Clone1.CommitRow("This is Data3 from Test")

	// 4. 执行flow1
	if err := flow1Clone1.Run(ctx); err != nil {
		panic(err)
	}
}
