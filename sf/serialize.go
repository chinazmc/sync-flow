package sf

import (
	"reflect"
	"sync-flow/common"
	"sync-flow/serialize"
)

// Serialize 数据序列化接口
type Serialize interface {
	// UnMarshal 用于将 SfRowArr 反序列化为指定类型的值。
	UnMarshal(common.SfRowArr, reflect.Type) (reflect.Value, error)
	// Marshal 用于将指定类型的值序列化为 SfRowArr。
	Marshal(interface{}) (common.SfRowArr, error)
}

// defaultSerialize SfFlow提供的默认序列化实现(开发者可以自定义)
var defaultSerialize = &serialize.DefaultSerialize{}

// isSerialize 判断传递进来的 paramType 是否实现了 Serialize 接口
func isSerialize(paramType reflect.Type) bool {
	return paramType.Implements(reflect.TypeOf((*Serialize)(nil)).Elem())
}
