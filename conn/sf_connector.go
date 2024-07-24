package conn

import (
	"context"
	"sync"
	"sync-flow/common"
	"sync-flow/config"
	"sync-flow/id"
	"sync-flow/sf"
)

type SfConnector struct {
	// Connector ID
	CId string
	// Connector Name
	CName string
	// Connector Config
	Conf *config.SfConnConfig

	// Connector Init
	onceInit sync.Once
}

// NewSfConnector 根据配置策略创建一个SfConnector
func NewSfConnector(config *config.SfConnConfig) *SfConnector {
	conn := new(SfConnector)
	conn.CId = id.SfID(common.SfIdTypeConnnector)
	conn.CName = config.CName
	conn.Conf = config

	return conn
}

// Init 初始化Connector所关联的存储引擎链接等
func (conn *SfConnector) Init() error {
	var err error

	//一个Connector只能执行初始化业务一次
	conn.onceInit.Do(func() {
		err = sf.Pool().CallConnInit(conn)
	})

	return err
}

// Call 调用Connector 外挂存储逻辑的读写操作
func (conn *SfConnector) Call(ctx context.Context, flow sf.Flow, args interface{}) error {
	if err := sf.Pool().CallConnector(ctx, flow, conn, args); err != nil {
		return err
	}

	return nil
}

func (conn *SfConnector) GetName() string {
	return conn.CName
}

func (conn *SfConnector) GetConfig() *config.SfConnConfig {
	return conn.Conf
}

func (conn *SfConnector) GetId() string {
	return conn.CId
}
