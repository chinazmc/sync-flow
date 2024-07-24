package sf

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync-flow/common"
	"sync-flow/log"
)

var _poolOnce sync.Once

// sfPool 用于管理全部的Function和Flow配置的池子
type sfPool struct {
	fnRouter funcRouter   // 全部的Function管理路由
	fnLock   sync.RWMutex // fnRouter 锁

	flowRouter flowRouter   // 全部的flow对象
	flowLock   sync.RWMutex // flowRouter 锁

	cInitRouter connInitRouter // 全部的Connector初始化路由
	ciLock      sync.RWMutex   // cInitRouter 锁

	cTree      connTree             //全部Connector管理路由
	connectors map[string]Connector // 全部的Connector对象
	cLock      sync.RWMutex         // cTree 锁
}

// 单例
var _pool *sfPool

// Pool 单例构造
func Pool() *sfPool {
	_poolOnce.Do(func() {
		//创建sfPool对象
		_pool = new(sfPool)

		// fnRouter初始化
		_pool.fnRouter = make(funcRouter)

		// flowRouter初始化
		_pool.flowRouter = make(flowRouter)
		// connTree初始化
		_pool.cTree = make(connTree)
		_pool.cInitRouter = make(connInitRouter)
		_pool.connectors = make(map[string]Connector)
	})

	return _pool
}
func (pool *sfPool) AddFlow(name string, flow Flow) {
	pool.flowLock.Lock()
	defer pool.flowLock.Unlock()

	if _, ok := pool.flowRouter[name]; !ok {
		pool.flowRouter[name] = flow
	} else {
		errString := fmt.Sprintf("Pool AddFlow Repeat FlowName=%s\n", name)
		panic(errString)
	}

	log.GetLogger().InfoF("Add FlowRouter FlowName=%s\n", name)
}

func (pool *sfPool) GetFlow(name string) Flow {
	pool.flowLock.RLock()
	defer pool.flowLock.RUnlock()

	if flow, ok := pool.flowRouter[name]; ok {
		return flow
	} else {
		return nil
	}
}

// FaaS 注册 Function 计算业务逻辑, 通过Function Name 索引及注册
func (pool *sfPool) FaaS(fnName string, f FaaS) {
	pool.fnLock.Lock()
	defer pool.fnLock.Unlock()

	if _, ok := pool.fnRouter[fnName]; !ok {
		pool.fnRouter[fnName] = f
	} else {
		errString := fmt.Sprintf("SfPoll FaaS Repeat FuncName=%s", fnName)
		panic(errString)
	}

	log.GetLogger().InfoF("Add SfPool FuncName=%s", fnName)
}

// CallFunction 调度 Function
func (pool *sfPool) CallFunction(ctx context.Context, fnName string, flow Flow) error {

	if f, ok := pool.fnRouter[fnName]; ok {
		return f(ctx, flow)
	}

	log.GetLogger().ErrorFX(ctx, "FuncName: %s Can not find in SfPool, Not Added.\n", fnName)

	return errors.New("FuncName: " + fnName + " Can not find in NsPool, Not Added.")
}

// CaaSInit 注册Connector初始化业务
func (pool *sfPool) CaaSInit(cname string, c ConnInit) {
	pool.ciLock.Lock() // 写锁
	defer pool.ciLock.Unlock()

	if _, ok := pool.cInitRouter[cname]; !ok {
		pool.cInitRouter[cname] = c
	} else {
		errString := fmt.Sprintf("SfPool Reg CaaSInit Repeat CName=%s\n", cname)
		panic(errString)
	}

	log.GetLogger().InfoF("Add SfPool CaaSInit CName=%s", cname)
}

// CallConnInit 调度 ConnInit
func (pool *sfPool) CallConnInit(conn Connector) error {
	pool.ciLock.RLock() // 读锁
	defer pool.ciLock.RUnlock()

	init, ok := pool.cInitRouter[conn.GetName()]

	if !ok {
		panic(errors.New(fmt.Sprintf("init connector cname = %s not reg..", conn.GetName())))
	}

	return init(conn)
}

// CaaS 注册Connector Call业务
func (pool *sfPool) CaaS(cname string, fname string, mode common.SfMode, c CaaS) {
	pool.cLock.Lock() // 写锁
	defer pool.cLock.Unlock()

	if _, ok := pool.cTree[cname]; !ok {
		//cid 首次注册，不存在，创建二级树NsConnSL
		pool.cTree[cname] = make(connSL)

		//初始化各类型FunctionMode
		pool.cTree[cname][common.S] = make(connFuncRouter)
		pool.cTree[cname][common.L] = make(connFuncRouter)
	}

	if _, ok := pool.cTree[cname][mode][fname]; !ok {
		pool.cTree[cname][mode][fname] = c
	} else {
		errString := fmt.Sprintf("CaaS Repeat CName=%s, FName=%s, Mode =%s\n", cname, fname, mode)
		panic(errString)
	}

	log.GetLogger().InfoF("Add SfPool CaaS CName=%s, FName=%s, Mode =%s", cname, fname, mode)
}

// CallConnector 调度 Connector
func (pool *sfPool) CallConnector(ctx context.Context, flow Flow, conn Connector, args interface{}) error {
	fn := flow.GetThisFunction()
	fnConf := fn.GetConfig()
	mode := common.SfMode(fnConf.FMode)

	if callback, ok := pool.cTree[conn.GetName()][mode][fnConf.FName]; ok {
		return callback(ctx, conn, fn, flow, args)
	}

	log.GetLogger().ErrorFX(ctx, "CName:%s FName:%s mode:%s Can not find in SfPool, Not Added.\n", conn.GetName(), fnConf.FName, mode)

	return errors.New(fmt.Sprintf("CName:%s FName:%s mode:%s Can not find in SfPool, Not Added.", conn.GetName(), fnConf.FName, mode))
}
