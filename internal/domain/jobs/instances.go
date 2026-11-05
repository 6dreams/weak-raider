package jobs

import (
	"context"
	"time"
	"weakRaider/internal/domain/manager"
)

type Instances struct {
	sync *manager.SeasonSync
}

func NewInstances(sync *manager.SeasonSync) *Instances {
	return &Instances{sync: sync}
}

func (j *Instances) Name() string {
	return "dict_instances"
}

func (j *Instances) Depends() []string {
	return []string{"dict_seasons"}
}

func (j *Instances) Run(ctx context.Context, period time.Duration) error {
	return j.sync.SyncInstances(ctx, period)
}
