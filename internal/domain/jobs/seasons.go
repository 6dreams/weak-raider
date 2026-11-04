package jobs

import (
	"context"
	"time"
	"weakRaider/internal/domain/manager"
)

type Seasons struct {
	sync *manager.SeasonSync
}

func NewSeasons(sync *manager.SeasonSync) *Seasons {
	return &Seasons{sync: sync}
}

func (j *Seasons) Name() string {
	return "dict_seasons"
}

func (j *Seasons) Depends() []string {
	return nil
}

func (j *Seasons) Run(ctx context.Context, period time.Duration) error {
	return j.sync.Sync(ctx, period)
}
