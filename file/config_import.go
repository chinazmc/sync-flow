package file

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync-flow/common"
	"sync-flow/config"
	"sync-flow/flow"
	"sync-flow/metrics"
	"sync-flow/sf"
)

type allConfig struct {
	FlowConfigMap map[string]*config.SfFlowConfig
	FuncConfigMap map[string]*config.SfFuncConfig
	ConnConfigMap map[string]*config.SfConnConfig
}

// sfTypeFlowConfigure 解析Flow配置文件，yaml格式
func sfTypeFlowConfigure(all *allConfig, confData []byte, fileName string, sfType interface{}) error {
	flow := new(config.SfFlowConfig)
	if ok := yaml.Unmarshal(confData, flow); ok != nil {
		return errors.New(fmt.Sprintf("%s has wrong format sfType = %s", fileName, sfType))
	}

	// 如果FLow状态为关闭，则不做配置加载
	if common.SfOnOff(flow.Status) == common.FlowDisable {
		return nil
	}

	if _, ok := all.FlowConfigMap[flow.FlowName]; ok {
		return errors.New(fmt.Sprintf("%s set repeat flow_id:%s", fileName, flow.FlowName))
	}

	// 加入配置集合中
	all.FlowConfigMap[flow.FlowName] = flow

	return nil
}

// sfTypeFuncConfigure 解析Function配置文件，yaml格式
func sfTypeFuncConfigure(all *allConfig, confData []byte, fileName string, sfType interface{}) error {
	function := new(config.SfFuncConfig)
	if ok := yaml.Unmarshal(confData, function); ok != nil {
		return errors.New(fmt.Sprintf("%s has wrong format sfType = %s", fileName, sfType))
	}
	if _, ok := all.FuncConfigMap[function.FuncName]; ok {
		return errors.New(fmt.Sprintf("%s set repeat function_id:%s", fileName, function.FuncName))
	}

	// 加入配置集合中
	all.FuncConfigMap[function.FuncName] = function

	return nil
}

// sfTypeConnConfigure 解析Connector配置文件，yaml格式
func sfTypeConnConfigure(all *allConfig, confData []byte, fileName string, sfType interface{}) error {
	conn := new(config.SfConnConfig)
	if ok := yaml.Unmarshal(confData, conn); ok != nil {
		return errors.New(fmt.Sprintf("%s is wrong format nsType = %s", fileName, sfType))
	}

	if _, ok := all.ConnConfigMap[conn.ConnName]; ok {
		return errors.New(fmt.Sprintf("%s set repeat conn_id:%s", fileName, conn.ConnName))
	}

	// 加入配置集合中
	all.ConnConfigMap[conn.ConnName] = conn

	return nil
}

// parseConfigTotalYaml 全盘解析配置文件，yaml格式, 讲配置信息解析到allConfig中
func parseConfigTotalYaml(loadPath string) (*allConfig, error) {

	all := new(allConfig)

	all.FlowConfigMap = make(map[string]*config.SfFlowConfig)
	all.FuncConfigMap = make(map[string]*config.SfFuncConfig)
	all.ConnConfigMap = make(map[string]*config.SfConnConfig)
	if !filepath.IsAbs(loadPath) {
		var err error
		loadPath, err = filepath.Abs(loadPath)
		if err != nil {
			return nil, err
		}
	}
	err := filepath.Walk(loadPath, func(filePath string, info os.FileInfo, err error) error {
		// 校验文件后缀是否合法
		if suffix := path.Ext(filePath); suffix != ".yml" && suffix != ".yaml" {
			return nil
		}

		// 读取文件内容
		confData, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		confMap := make(map[string]interface{})

		// 校验yaml合法性
		if err := yaml.Unmarshal(confData, confMap); err != nil {
			return err
		}

		// 判断sfType是否存在
		if sfType, ok := confMap["sftype"]; !ok {
			return errors.New(fmt.Sprintf("yaml file %s has no file [sftype]!", filePath))
		} else {
			switch sfType {
			case common.SfIdTypeFlow:
				return sfTypeFlowConfigure(all, confData, filePath, sfType)

			case common.SfIdTypeFunction:
				return sfTypeFuncConfigure(all, confData, filePath, sfType)

			case common.SfIdTypeConnnector:
				return sfTypeConnConfigure(all, confData, filePath, sfType)
			case common.SfIdTypeGlobal:
				return sfTypeGlobalConfigure(confData, filePath, sfType)
			default:
				return errors.New(fmt.Sprintf("%s set wrong sftype %s", filePath, sfType))
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return all, nil
}

// ConfigImportYaml 全盘解析配置文件，yaml格式
func ConfigImportYaml(loadPath string) error {

	all, err := parseConfigTotalYaml(loadPath)
	if err != nil {
		return err
	}

	for flowName, flowConfig := range all.FlowConfigMap {

		// 构建一个Flow
		newFlow := flow.NewSfFlow(flowConfig)

		for _, fp := range flowConfig.Funcs {
			if err := buildFlow(all, fp, newFlow, flowName); err != nil {
				return err
			}
		}

		//将flow添加到FlowPool中
		sf.Pool().AddFlow(flowName, newFlow)
	}

	return nil
}
func buildFlow(all *allConfig, fp config.SfFunctionParam, newFlow sf.Flow, flowName string) error {
	//加载当前Flow依赖的Function
	if funcConfig, ok := all.FuncConfigMap[fp.FuncName]; !ok {
		return errors.New(fmt.Sprintf("FlowName [%s] need FuncName [%s], But has No This FuncName Config", flowName, fp.FuncName))
	} else {
		//flow add connector
		if funcConfig.Option.ConnName != "" {
			// 加载当前Function依赖的Connector
			if connConf, ok := all.ConnConfigMap[funcConfig.Option.ConnName]; !ok {
				return errors.New(fmt.Sprintf("FuncName [%s] need ConnName [%s], But has No This ConnName Config", fp.FuncName, funcConfig.Option.ConnName))
			} else {
				// Function Config 关联 Connector Config
				_ = funcConfig.AddConnConfig(connConf)
			}
		}

		//flow add function
		if err := newFlow.AppendNewFunction(funcConfig, fp.Params); err != nil {
			return err
		}
	}

	return nil
}

// sfTypeGlobalConfigure 解析Global配置文件，yaml格式
func sfTypeGlobalConfigure(confData []byte, fileName string, sfType interface{}) error {
	// 全局配置
	if ok := yaml.Unmarshal(confData, config.GlobalConfig); ok != nil {
		return errors.New(fmt.Sprintf("%s is wrong format sfType = %s", fileName, sfType))
	}

	// 启动Metrics服务
	metrics.RunMetrics()

	return nil
}
