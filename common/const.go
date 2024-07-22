package common

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
	// V 为校验特征的SfFunction,
	// 主要进行数据的过滤，验证，字段梳理，幂等等前置数据处理
	V SfMode = "Verify"

	// S 为存储特征的SfFunction,
	// S会通过NsConnector进行将数据进行存储，数据的临时声明周期为Ns Window
	S SfMode = "Save"

	// L 为加载特征的SfFunction，
	// L会通过SfConnector进行数据加载，通过该Function可以从逻辑上与对应的S Function进行并流
	L SfMode = "Load"

	// C 为计算特征的SfFunction,
	// C会通过SfFlow中的数据计算，生成新的字段，将数据流传递给下游S进行存储，或者自己也已直接通过SfConnector进行存储
	C SfMode = "Calculate"

	// E 为扩展特征的SfFunction，
	// 作为流式计算的自定义特征Function，如，Notify 调度器触发任务的消息发送，删除一些数据，重置状态等。
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
