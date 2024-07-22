package test

import (
	"context"
	"sync-flow/log"
	"testing"
)

func TestSfLogger(t *testing.T) {
	ctx := context.Background()

	log.GetLogger().InfoFX(ctx, "GetLogger InfoFX")
	log.GetLogger().ErrorFX(ctx, "GetLogger ErrorFX")
	log.GetLogger().DebugFX(ctx, "GetLogger DebugFX")

	log.GetLogger().InfoF("GetLogger InfoF")
	log.GetLogger().ErrorF("GetLogger ErrorF")
	log.GetLogger().DebugF("GetLogger DebugF")
}
