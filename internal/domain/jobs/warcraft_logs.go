package jobs

import (
	"context"
	"time"
	"weakRaider/internal/domain/manager"
)

type WarcraftLogs struct {
	sync *manager.WarcraftLogs
}

func NewWarcraftLogs(sync *manager.WarcraftLogs) *WarcraftLogs {
	return &WarcraftLogs{sync: sync}
}

func (j *WarcraftLogs) Name() string {
	return "warcraft_logs"
}

func (j *WarcraftLogs) Depends() []string {
	return nil
}

func (j *WarcraftLogs) Run(ctx context.Context, _ time.Duration) error {
	return j.sync.Sync(ctx)
}
