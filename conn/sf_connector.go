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
	// SfConnector的自定义临时数据
	metaData map[string]interface{}
	// 管理metaData的读写锁
	mLock sync.RWMutex
}

// NewSfConnector 根据配置策略创建一个SfConnector
func NewSfConnector(config *config.SfConnConfig) *SfConnector {
	conn := new(SfConnector)
	conn.CId = id.SfID(common.SfIdTypeConnnector)
	conn.CName = config.CName
	conn.Conf = config
	conn.metaData = make(map[string]interface{})
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
func (conn *SfConnector) Call(ctx context.Context, flow sf.Flow, args interface{}) (interface{}, error) {
	var result interface{}
	var err error

	result, err = sf.Pool().CallConnector(ctx, flow, conn, args)
	if err != nil {
		return nil, err
	}

	return result, nil
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

// GetMetaData 得到当前Connector的临时数据
func (conn *SfConnector) GetMetaData(key string) interface{} {
	conn.mLock.RLock()
	defer conn.mLock.RUnlock()

	data, ok := conn.metaData[key]
	if !ok {
		return nil
	}

	return data
}

// SetMetaData 设置当前Connector的临时数据
func (conn *SfConnector) SetMetaData(key string, value interface{}) {
	conn.mLock.Lock()
	defer conn.mLock.Unlock()

	conn.metaData[key] = value
}
