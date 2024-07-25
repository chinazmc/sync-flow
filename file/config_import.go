package file

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync-flow/common"
	"sync-flow/config"
	"sync-flow/flow"
	"sync-flow/sf"
)

type allConfig struct {
	Flows map[string]*config.SfFlowConfig
	Funcs map[string]*config.SfFuncConfig
	Conns map[string]*config.SfConnConfig
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

	if _, ok := all.Flows[flow.FlowName]; ok {
		return errors.New(fmt.Sprintf("%s set repeat flow_id:%s", fileName, flow.FlowName))
	}

	// 加入配置集合中
	all.Flows[flow.FlowName] = flow

	return nil
}

// sfTypeFuncConfigure 解析Function配置文件，yaml格式
func sfTypeFuncConfigure(all *allConfig, confData []byte, fileName string, sfType interface{}) error {
	function := new(config.SfFuncConfig)
	if ok := yaml.Unmarshal(confData, function); ok != nil {
		return errors.New(fmt.Sprintf("%s has wrong format sfType = %s", fileName, sfType))
	}
	if _, ok := all.Funcs[function.FName]; ok {
		return errors.New(fmt.Sprintf("%s set repeat function_id:%s", fileName, function.FName))
	}

	// 加入配置集合中
	all.Funcs[function.FName] = function

	return nil
}

// sfTypeConnConfigure 解析Connector配置文件，yaml格式
func sfTypeConnConfigure(all *allConfig, confData []byte, fileName string, sfType interface{}) error {
	conn := new(config.SfConnConfig)
	if ok := yaml.Unmarshal(confData, conn); ok != nil {
		return errors.New(fmt.Sprintf("%s is wrong format nsType = %s", fileName, sfType))
	}

	if _, ok := all.Conns[conn.CName]; ok {
		return errors.New(fmt.Sprintf("%s set repeat conn_id:%s", fileName, conn.CName))
	}

	// 加入配置集合中
	all.Conns[conn.CName] = conn

	return nil
}

// parseConfigWalkYaml 全盘解析配置文件，yaml格式, 讲配置信息解析到allConfig中
func parseConfigWalkYaml(loadPath string) (*allConfig, error) {

	all := new(allConfig)

	all.Flows = make(map[string]*config.SfFlowConfig)
	all.Funcs = make(map[string]*config.SfFuncConfig)
	all.Conns = make(map[string]*config.SfConnConfig)
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

	all, err := parseConfigWalkYaml(loadPath)
	if err != nil {
		return err
	}

	for flowName, flowConfig := range all.Flows {

		// 构建一个Flow
		newFlow := flow.NewSfFlow(flowConfig)

		for _, fp := range flowConfig.Flows {
			if err := buildFlow(all, fp, newFlow, flowName); err != nil {
				return err
			}
		}

		//将flow添加到FlowPool中
		sf.Pool().AddFlow(flowName, newFlow)
	}

	return nil
}
func buildFlow(all *allConfig, fp config.SfFlowFunctionParam, newFlow sf.Flow, flowName string) error {
	//加载当前Flow依赖的Function
	if funcConfig, ok := all.Funcs[fp.FuncName]; !ok {
		return errors.New(fmt.Sprintf("FlowName [%s] need FuncName [%s], But has No This FuncName Config", flowName, fp.FuncName))
	} else {
		//flow add connector
		if funcConfig.Option.CName != "" {
			// 加载当前Function依赖的Connector
			if connConf, ok := all.Conns[funcConfig.Option.CName]; !ok {
				return errors.New(fmt.Sprintf("FuncName [%s] need ConnName [%s], But has No This ConnName Config", fp.FuncName, funcConfig.Option.CName))
			} else {
				// Function Config 关联 Connector Config
				_ = funcConfig.AddConnConfig(connConf)
			}
		}

		//flow add function
		if err := newFlow.Link(funcConfig, fp.Params); err != nil {
			return err
		}
	}

	return nil
}
