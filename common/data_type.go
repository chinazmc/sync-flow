package common

// SfRow 一行数据
type SfRow interface{}

// SfRowArr 一次业务的批量数据
type SfRowArr []SfRow

/*
		SfDataMap 当前Flow承载的全部数据，
	   	key	:  数据所在的Function ID
	    value: 对应的SfRow
*/
type SfDataMap map[string]SfRowArr
