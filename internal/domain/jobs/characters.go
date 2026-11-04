package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"
	"weakRaider/internal/domain/manager"
	"weakRaider/internal/domain/repository"
)

type Characters struct {
	sync   *manager.CharacterSync
	guilds *repository.GuildRepository
}

func NewCharacters(sync *manager.CharacterSync, guilds *repository.GuildRepository) *Characters {
	return &Characters{sync: sync, guilds: guilds}
}

func (j *Characters) Name() string {
	return "characters"
}

func (j *Characters) Depends() []string {
	return nil
}

func (j *Characters) Run(ctx context.Context, _ time.Duration) error {
	all, err := j.guilds.FindAll()
	if err != nil {
		return err
	}
	var failures []error
	for _, guild := range all {
		if err := j.sync.Sync(ctx, &guild); err != nil {
			failures = append(failures, fmt.Errorf("guild %s: %w", guild.Name, err))
		}
	}
	return errors.Join(failures...)
}
