package test

import (
	"sync-flow/common"
	"sync-flow/config"
	"sync-flow/log"
	"testing"
)

func TestNewFuncConfig(t *testing.T) {
	source := config.SfSource{
		Name: "订单数据",
		Must: []string{"order_id", "user_id"},
	}

	option := config.SfFuncOption{
		CName:        "connectorName1",
		RetryTimes:   3,
		RetryDuriton: 300,

		Params: config.FParam{
			"param1": "value1",
			"param2": "value2",
		},
	}

	myFunc1 := config.NewFuncConfig("funcName1", common.S, &source, &option)

	log.GetLogger().InfoF("funcName1: %+v\n", myFunc1)
}

func TestNewFlowConfig(t *testing.T) {

	flowFuncParams1 := config.SfFunctionParam{
		FuncName: "funcName1",
		Params: config.FParam{
			"flowSetFunParam1": "value1",
			"flowSetFunParam2": "value2",
		},
	}

	flowFuncParams2 := config.SfFunctionParam{
		FuncName: "funcName2",
		Params: config.FParam{
			"default": "value1",
		},
	}

	myFlow1 := config.NewFlowConfig("flowName1", common.FlowEnable)
	myFlow1.AppendFunctionConfig(flowFuncParams1)
	myFlow1.AppendFunctionConfig(flowFuncParams2)

	log.GetLogger().InfoF("myFlow1: %+v\n", myFlow1)
}
func TestNewConnConfig(t *testing.T) {

	source := config.SfSource{
		Name: "订单数据",
		Must: []string{"order_id", "user_id"},
	}

	option := config.SfFuncOption{
		CName:        "connectorName1",
		RetryTimes:   3,
		RetryDuriton: 300,

		Params: config.FParam{
			"param1": "value1",
			"param2": "value2",
		},
	}

	myFunc1 := config.NewFuncConfig("funcName1", common.S, &source, &option)

	connParams := config.FParam{
		"param1": "value1",
		"param2": "value2",
	}

	myConnector1 := config.NewConnConfig("connectorName1", "0.0.0.0:9987,0.0.0.0:9997", common.REDIS, "key", connParams)

	if err := myConnector1.WithFunc(myFunc1); err != nil {
		log.GetLogger().ErrorF("WithFunc err: %s\n", err.Error())
	}

	log.GetLogger().InfoF("myConnector1: %+v\n", myConnector1)
}
