package jobs

import (
	"context"
	"time"
	"weakRaider/internal/domain/manager"
)

type Talents struct {
	sync *manager.TalentSync
}

func NewTalents(sync *manager.TalentSync) *Talents {
	return &Talents{sync: sync}
}

func (j *Talents) Name() string {
	return "dict_talents"
}

func (j *Talents) Depends() []string {
	return nil
}

func (j *Talents) Run(ctx context.Context, _ time.Duration) error {
	return j.sync.Sync(ctx)
}
