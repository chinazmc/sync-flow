package config

import (
	"errors"
	"fmt"
	"sync-flow/common"
)

// SfConnConfig SfConnector 策略配置
type SfConnConfig struct {
	//配置类型
	SfType string `yaml:"sfType"`
	//唯一描述标识
	CName string `yaml:"cname"`
	//基础存储媒介地址
	AddrString string `yaml:"addrs"`
	//存储媒介引擎类型"Mysql" "Redis" "Kafka"等
	Type common.SfConnType `yaml:"type"`
	//一次存储的标识：如Redis为Key名称、Mysql为Table名称,Kafka为Topic名称等
	Key string `yaml:"key"`
	//配置信息中的自定义参数
	Params map[string]string `yaml:"params"`
	//存储读取所绑定的NsFuncionID
	Load []string `yaml:"load"`
	Save []string `yaml:"save"`
}

// NewConnConfig 创建一个SfConnector策略配置对象, 用于描述一个SfConnector信息
func NewConnConfig(cName string, addr string, t common.SfConnType, key string, param FParam) *SfConnConfig {
	strategy := new(SfConnConfig)
	strategy.CName = cName
	strategy.AddrString = addr

	strategy.Type = t
	strategy.Key = key
	strategy.Params = param

	return strategy
}

// WithFunc Connector与Function进行关系绑定
func (cConfig *SfConnConfig) WithFunc(fConfig *SfFuncConfig) error {

	switch common.SfMode(fConfig.FMode) {
	case common.S:
		cConfig.Save = append(cConfig.Save, fConfig.FName)
	case common.L:
		cConfig.Load = append(cConfig.Load, fConfig.FName)
	default:
		return errors.New(fmt.Sprintf("Wrong SfMode %s", fConfig.FMode))
	}

	return nil
}
