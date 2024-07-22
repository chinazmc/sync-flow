package config

import "sync-flow/common"

// SfFlowFunctionParam 一个Flow配置中Function的Id及携带固定配置参数
type SfFlowFunctionParam struct {
	FuncName string `yaml:"fname"`  //必须
	Params   FParam `yaml:"params"` //选填,在当前Flow中Function定制固定配置参数
}

// SfFlowConfig 用户贯穿整条流式计算上下文环境的对象
type SfFlowConfig struct {
	SfType   string                `yaml:"sfType"`
	Status   int                   `yaml:"status"`
	FlowName string                `yaml:"flow_name"`
	Flows    []SfFlowFunctionParam `yaml:"flows"`
}

// NewFlowConfig 创建一个Flow策略配置对象, 用于描述一个SfFlow信息
func NewFlowConfig(flowName string, enable common.SfOnOff) *SfFlowConfig {
	config := new(SfFlowConfig)
	config.FlowName = flowName
	config.Flows = make([]SfFlowFunctionParam, 0)

	config.Status = int(enable)

	return config
}

// AppendFunctionConfig 添加一个Function Config 到当前Flow中
func (fConfig *SfFlowConfig) AppendFunctionConfig(params SfFlowFunctionParam) {
	fConfig.Flows = append(fConfig.Flows, params)
}
