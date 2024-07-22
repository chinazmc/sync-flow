package id

import (
	"github.com/google/uuid"
	"strings"
	"sync-flow/common"
)

// SfID 获取一个中随机实例ID
// 格式为  "prefix1-[prefix2-][prefix3-]ID"
// 如：flow-1234567890
// 如：func-1234567890
// 如: conn-1234567890
// 如: func-1-1234567890
func SfID(prefix ...string) (sfId string) {

	idStr := strings.Replace(uuid.New().String(), "-", "", -1)
	sfId = formatSfID(idStr, prefix...)

	return
}

func formatSfID(idStr string, prefix ...string) string {
	var sfId string

	for _, fix := range prefix {
		sfId += fix
		sfId += common.SfIdJoinChar
	}

	sfId += idStr

	return sfId
}
