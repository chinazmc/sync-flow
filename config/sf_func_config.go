package config

import (
	"errors"
	"sync-flow/common"
	"sync-flow/log"
)

// FParam 在当前Flow中Function定制固定配置参数类型
type FParam map[string]string

// SfSource 表示当前Function的业务源
type SfSource struct {
	Name string   `yaml:"name"` //本层Function的数据源描述
	Must []string `yaml:"must"` //source必传字段
}

// SfFuncOption 可选配置
type SfFuncOption struct {
	CName        string `yaml:"cname"`           //连接器Connector名称
	RetryTimes   int    `yaml:"retry_times"`     //选填,Function调度重试(不包括正常调度)最大次数
	RetryDuriton int    `yaml:"return_duration"` //选填,Function调度每次重试最大时间间隔(单位:ms)
	Params       FParam `yaml:"default_params"`  //选填,在当前Flow中Function定制固定配置参数
}

// SfFuncConfig 一个SfFunction策略配置
type SfFuncConfig struct {
	SfType string       `yaml:"sfType"`
	FName  string       `yaml:"fname"`
	FMode  string       `yaml:"fmode"`
	Source SfSource     `yaml:"source"`
	Option SfFuncOption `yaml:"option"`

	connConf *SfConnConfig
}

func (fConf *SfFuncConfig) AddConnConfig(cConf *SfConnConfig) error {
	if cConf == nil {
		return errors.New("SfConnConfig is nil")
	}

	// Function需要和Connector进行关联
	fConf.connConf = cConf

	// Connector需要和Function进行关联
	_ = cConf.WithFunc(fConf)
	// 更新Function配置中的CName
	fConf.Option.CName = cConf.CName
	return nil
}

func (fConf *SfFuncConfig) GetConnConfig() (*SfConnConfig, error) {
	if fConf.connConf == nil {
		return nil, errors.New("SfFuncConfig.connConf not set")
	}

	return fConf.connConf, nil
}

// NewFuncConfig 创建一个Function策略配置对象, 用于描述一个SfFunction信息
func NewFuncConfig(
	funcName string, mode common.SfMode,
	source *SfSource, option *SfFuncOption) *SfFuncConfig {

	config := new(SfFuncConfig)
	config.FName = funcName

	if source == nil {
		defaultSource := SfSource{
			Name: "unNamedSource",
		}
		source = &defaultSource
		log.GetLogger().ErrorF("funcName NewConfig Error, source is nil, funcName = %s\n", funcName)
	}
	config.Source = *source

	config.FMode = string(mode)

	//FunctionS 和 L 需要必传SfConnector参数,原因是S和L需要通过Connector进行建立流式关系
	//if mode == common.Save || mode == common.L {
	//	if option == nil {
	//		log.GetLogger().ErrorF("Funcion S/L need option->Cid\n")
	//		return nil
	//	} else if option.CName == "" {
	//		log.GetLogger().ErrorF("Funcion S/L need option->Cid\n")
	//		return nil
	//	}
	//}

	if option != nil {
		config.Option = *option
	}

	return config
}
