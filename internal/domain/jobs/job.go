package jobs

import (
	"context"
	"time"
)

type Job interface {
	Name() string
	Depends() []string
	Run(context.Context, time.Duration) error
}
