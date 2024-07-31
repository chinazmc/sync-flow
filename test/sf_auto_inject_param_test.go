package test

import (
	"context"
	"sync-flow/file"
	"sync-flow/sf"
	"sync-flow/test/caas"
	"sync-flow/test/faas"
	"sync-flow/test/proto"
	"testing"
)

func TestAutoInjectParamWithConfig(t *testing.T) {
	ctx := context.Background()

	sf.Pool().FaaS("AvgStuScore", faas.AvgStuScore)
	sf.Pool().FaaS("PrintStuAvgScore", faas.PrintStuAvgScore)
	sf.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("./load_conf/"); err != nil {
		panic(err)
	}

	// 2. 获取Flow
	flow1 := sf.Pool().GetFlow("StuAvg")
	if flow1 == nil {
		panic("flow1 is nil")
	}

	// 3. 提交原始数据
	_ = flow1.CommitRow(&faas.AvgStuScoreIn{
		StuScores: proto.StuScores{
			StuId:  100,
			Score1: 1,
			Score2: 2,
			Score3: 3,
		},
	})
	_ = flow1.CommitRow(faas.AvgStuScoreIn{
		StuScores: proto.StuScores{
			StuId:  100,
			Score1: 1,
			Score2: 2,
			Score3: 3,
		},
	})

	// 提交原始数据（json字符串）
	_ = flow1.CommitRow(`{"stu_id":101}`)

	// 4. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}
