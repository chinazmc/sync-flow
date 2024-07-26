package flow

import (
	"context"
	"errors"
	"sync"
	"sync-flow/common"
	"sync-flow/config"
	"sync-flow/conn"
	"sync-flow/function"
	"sync-flow/id"
	"sync-flow/log"
	"sync-flow/sf"
)

// SfFlow 用于贯穿整条流式计算的上下文环境
type SfFlow struct {
	// 基础信息
	Id   string               // Flow的分布式实例ID(用于SfFlow内部区分不同实例)
	Name string               // Flow的可读名称
	Conf *config.SfFlowConfig // Flow配置策略

	// Function列表
	Funcs          map[string]sf.Function // 当前flow拥有的全部管理的全部Function对象, key: FunctionName
	FlowHead       sf.Function            // 当前Flow所拥有的Function列表表头
	FlowTail       sf.Function            // 当前Flow所拥有的Function列表表尾
	flock          sync.RWMutex           // 管理链表插入读写的锁
	ThisFunction   sf.Function            // Flow当前正在执行的SfFunction对象
	ThisFunctionId string                 // 当前执行到的Function ID (策略配置ID)
	PrevFunctionId string                 // 当前执行到的Function 上一层FunctionID(策略配置ID)

	// Function列表参数
	funcParams map[string]config.FParam // flow在当前Function的自定义固定配置参数,Key:function的实例NsID, value:FParam
	fplock     sync.RWMutex             // 管理funcParams的读写锁

	buffer common.SfRowArr  // 用来临时存放输入字节数据的内部Buf, 一条数据为interface{}, 多条数据为[]interface{} 也就是SfBatch
	data   common.SfDataMap // 流式计算各个层级的数据源
	inPut  common.SfRowArr  // 当前Function的计算输入数据
	// SfFlow Action
	action sf.Action // 当前Flow所携带的Action动作
	abort  bool      // 是否中断Flow
}

// NewSfFlow 创建一个SfFlow.
func NewSfFlow(conf *config.SfFlowConfig) sf.Flow {
	flow := new(SfFlow)
	// 实例Id
	flow.Id = id.SfID(common.SfIdTypeFlow)

	// 基础信息
	flow.Name = conf.FlowName
	flow.Conf = conf

	// Function列表
	flow.Funcs = make(map[string]sf.Function)
	flow.funcParams = make(map[string]config.FParam)

	flow.data = make(common.SfDataMap)
	return flow
}

// Link 将Function链接到Flow中
// fConf: 当前Function策略
// fParams: 当前Flow携带的Function动态参数
func (flow *SfFlow) Link(fConf *config.SfFuncConfig, fParams config.FParam) error {
	// 创建Function
	f := function.NewSfFunction(flow, fConf)

	if fConf.Option.CName != "" {
		// 当前Function有Connector关联，需要初始化Connector实例

		// 获取Connector配置
		connConfig, err := fConf.GetConnConfig()
		if err != nil {
			panic(err)
		}

		// 创建Connector对象
		connector := conn.NewSfConnector(connConfig)

		// 初始化Connector, 执行Connector Init 方法
		if err = connector.Init(); err != nil {
			panic(err)
		}

		// 关联Function实例和Connector实例关系
		_ = f.AddConnector(connector)
	}

	// Flow 添加 Function
	if err := flow.appendFunc(f, fParams); err != nil {
		return err
	}

	return nil
}

// appendFunc 将Function添加到Flow中, 链表操作
func (flow *SfFlow) appendFunc(function sf.Function, fParam config.FParam) error {

	if function == nil {
		return errors.New("AppendFunc append nil to List")
	}

	flow.flock.Lock()
	defer flow.flock.Unlock()

	if flow.FlowHead == nil {
		// 首次添加节点
		flow.FlowHead = function
		flow.FlowTail = function

		function.SetN(nil)
		function.SetP(nil)

	} else {
		// 将function插入到链表的尾部
		function.SetP(flow.FlowTail)
		function.SetN(nil)

		flow.FlowTail.SetN(function)
		flow.FlowTail = function
	}

	//将Function Name 详细Hash对应关系添加到flow对象中
	flow.Funcs[function.GetConfig().FName] = function

	//先添加function 默认携带的Params参数
	params := make(config.FParam)
	for key, value := range function.GetConfig().Option.Params {
		params[key] = value
	}

	//再添加flow携带的function定义参数(重复即覆盖)
	for key, value := range fParam {
		params[key] = value
	}

	// 将得到的FParams存留在flow结构体中，用来function业务直接通过Hash获取
	// key 为当前Function的SfId，不用Fid的原因是为了防止一个Flow添加两个相同策略Id的Function
	flow.funcParams[function.GetId()] = params

	return nil
}

// Run 启动SfFlow的流式计算, 从起始Function开始执行流
func (flow *SfFlow) Run(ctx context.Context) error {

	var fn sf.Function

	fn = flow.FlowHead
	// 重置 abort
	flow.abort = false //  每次进入调度，要重置abort状态

	if flow.Conf.Status == int(common.FlowDisable) {
		//flow被配置关闭
		return nil
	}
	// 因为此时还没有执行任何Function, 所以PrevFunctionId为FirstVirtual 因为没有上一层Function
	flow.PrevFunctionId = common.FunctionIdFirstVirtual

	// 提交数据流原始数据
	if err := flow.commitSrcData(ctx); err != nil {
		return err
	}
	//流式链式调用
	for fn != nil && flow.abort != true {

		// ========= 数据流 新增 ===========
		// flow记录当前执行到的Function 标记
		fid := fn.GetId()
		flow.ThisFunction = fn
		flow.ThisFunctionId = fid

		// 得到当前Function要处理与的源数据
		if inputData, err := flow.getCurData(); err != nil {
			log.GetLogger().ErrorFX(ctx, "flow.Run(): getCurData err = %s\n", err.Error())
			return err
		} else {
			flow.inPut = inputData
		}
		// ========= 数据流 新增 ===========

		if err := fn.Call(ctx, flow); err != nil {
			//Error
			return err
		} else {
			//Success
			fn, err = flow.dealAction(ctx, fn)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func (flow *SfFlow) GetName() string {
	return flow.Name
}

func (flow *SfFlow) GetThisFunction() sf.Function {
	return flow.ThisFunction
}

func (flow *SfFlow) GetThisFuncConf() *config.SfFuncConfig {
	return flow.ThisFunction.GetConfig()
}

// GetConnector 得到当前正在执行的Function的Connector
func (flow *SfFlow) GetConnector() (sf.Connector, error) {
	if conn := flow.ThisFunction.GetConnector(); conn != nil {
		return conn, nil
	} else {
		return nil, errors.New("GetConnector(): Connector is nil")
	}
}

// GetConnConf 得到当前正在执行的Function的Connector的配置
func (flow *SfFlow) GetConnConf() (*config.SfConnConfig, error) {
	if conn := flow.ThisFunction.GetConnector(); conn != nil {
		return conn.GetConfig(), nil
	} else {
		return nil, errors.New("GetConnConf(): Connector is nil")
	}
}

func (flow *SfFlow) GetConfig() *config.SfFlowConfig {
	return flow.Conf
}

// GetFuncConfigByName 得到当前Flow的配置
func (flow *SfFlow) GetFuncConfigByName(funcName string) *config.SfFuncConfig {
	if f, ok := flow.Funcs[funcName]; ok {
		return f.GetConfig()
	} else {
		log.GetLogger().ErrorF("GetFuncConfigByName(): Function %s not found", funcName)
		return nil
	}
}

// Next 当前Flow执行到的Function进入下一层Function所携带的Action动作
func (flow *SfFlow) Next(acts ...sf.ActionFunc) error {

	// 加载Function FaaS 传递的 Action动作
	flow.action = sf.LoadActions(acts)

	return nil
}
