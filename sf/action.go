package sf

// Action SfFlow执行流程Actions
type Action struct {
	// DataReuseEnable 是否复用上层Function数据
	DataReuseEnable bool
	// EntryNextForceEnable 为忽略上述默认规则，没有数据强制进入下一层Function
	EntryNextForceEnable bool
	// FlowAbortEnable 终止Flow的执行
	FlowAbortEnable bool
	// FuncJumpEnable 跳转到指定Function继续执行
	FuncJumpEnable string
}

// ActionFunc SfFlow Functional Option 类型
type ActionFunc func(ops *Action)

// LoadActions 加载Actions，依次执行ActionFunc操作函数
func LoadActions(acts []ActionFunc) Action {
	action := Action{}

	if acts == nil {
		return action
	}

	for _, act := range acts {
		act(&action)
	}

	return action
}

// ActionAbort 终止Flow的执行
func ActionAbort(action *Action) {
	action.FlowAbortEnable = true
}

// ActionDataReuse Next复用上层Function数据Option
func ActionDataReuse(act *Action) {
	act.DataReuseEnable = true
}

// ActionForceEntryNext 强制进入下一层
func ActionForceEntryNext(act *Action) {
	act.FlowAbortEnable = true
}

// ActionJumpFunc 会返回一个ActionFunc函数，并且会将funcName赋值给Action.JumpFunc
// (注意：容易出现Flow循环调用，导致死循环)
func ActionJumpFunc(funcName string) ActionFunc {
	return func(act *Action) {
		act.FuncJumpEnable = funcName
	}
}
