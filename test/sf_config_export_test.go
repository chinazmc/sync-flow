package test

import (
	"sync-flow/common"
	"sync-flow/file"
	"sync-flow/sf"
	"sync-flow/test/caas"
	"sync-flow/test/faas"
	"testing"
)

func TestConfigExportYmal(t *testing.T) {
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

	// 2. 讲构建的内存SfFlow结构配置导出的文件当中
	flows := sf.Pool().GetFlows()
	for _, flow := range flows {
		if err := file.ConfigExportYaml(flow, "./export_conf/"); err != nil {
			panic(err)
		}
	}
}
