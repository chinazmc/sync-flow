package common

import "time"

// SfIdType 用户生成SfId的字符串前缀
const (
	SfIdTypeFlow       = "flow"
	SfIdTypeConnnector = "conn"
	SfIdTypeFunction   = "func"
	SfIdTypeGlobal     = "global"
	SfIdJoinChar       = "-"
)
const (
	// FunctionIdFirstVirtual 为首结点Function上一层虚拟的Function ID
	FunctionIdFirstVirtual = "FunctionIdFirstVirtual"
	// FunctionIdLastVirtual 为尾结点Function下一层虚拟的Function ID
	FunctionIdLastVirtual = "FunctionIdLastVirtual"
)

type SfMode string

const (
	// V 为校验特征的SfFunction, 主要进行数据的过滤，验证，字段梳理，幂等等前置数据处理
	V SfMode = "Verify"

	// S 为存储特征的SfFunction, S会通过SfConnector进行将数据进行存储. S Function 会通过SfConnector进行数据存储,具备相同Connector的Function在逻辑上可以进行并流
	S SfMode = "Save"

	// L 为加载特征的SfFunction，L会通过SfConnector进行数据加载，L Function 会通过SfConnector进行数据读取，具备相同Connector的Function可以从逻辑上与对应的S Function进行并流
	L SfMode = "Load"

	// C 为计算特征的SfFunction, 可以生成新的字段，计算新的值，进行数据的聚合，分析等
	C SfMode = "Calculate"

	// E 为扩展特征的SfFunction，作为流式计算的自定义特征Function，也同时是SfFlow当前流中的最后一个Function，概念类似Sink。
	E SfMode = "Expand"
)

// SfConnType represents the type of SfConnector
type SfConnType string

const (
	// REDIS is the type of Redis
	REDIS SfConnType = "redis"
	// MYSQL is the type of MySQL
	MYSQL SfConnType = "mysql"
	// KAFKA is the type of Kafka
	KAFKA SfConnType = "kafka"
	// TIDB is the type of TiDB
	TIDB SfConnType = "tidb"
	// ES is the type of Elasticsearch
	ES SfConnType = "es"
)

// SfOnOff  Whether to enable the Flow
type SfOnOff int

const (
	// FlowEnable Enabled
	FlowEnable SfOnOff = 1
	// FlowDisable Disabled
	FlowDisable SfOnOff = 0
)

// cache
const (
	// DeFaultFlowCacheCleanUp SfFlow中Flow对象Cache缓存默认的清理内存时间
	DeFaultFlowCacheCleanUp = 5 //单位 min
	// DefaultExpiration 默认GoCahce时间 ，永久保存
	DefaultExpiration time.Duration = 0
)

// metrics
const (
	METRICS_ROUTE string = "/metrics"

	LABEL_FLOW_NAME     string = "flow_name"
	LABEL_FLOW_ID       string = "flow_id"
	LABEL_FUNCTION_NAME string = "func_name"
	LABEL_FUNCTION_MODE string = "func_mode"

	COUNTER_SFFLOW_DATA_TOTAL_NAME string = "sfflow_data_total"
	COUNTER_SFFLOW_DATA_TOTAL_HELP string = "SfFlow全部Flow的数据总量"

	GANGE_FLOW_DATA_TOTAL_NAME string = "flow_data_total"
	GANGE_FLOW_DATA_TOTAL_HELP string = "SfFlow各个FlowID数据流的数据数量总量"

	GANGE_FLOW_SCHE_CNTS_NAME string = "flow_schedule_cnts"
	GANGE_FLOW_SCHE_CNTS_HELP string = "SfFlow各个FlowID被调度的次数"

	GANGE_FUNC_SCHE_CNTS_NAME        string = "func_schedule_cnts"
	GANGE_FUNC_SCHE_CNTS_HELP        string = "SfFlow各个Function被调度的次数"
	HISTOGRAM_FUNCTION_DURATION_NAME string = "func_run_duration"
	HISTOGRAM_FUNCTION_DURATION_HELP string = "Function执行耗时"

	HISTOGRAM_FLOW_DURATION_NAME string = "flow_run_duration"
	HISTOGRAM_FLOW_DURATION_HELP string = "Flow执行耗时"
)
