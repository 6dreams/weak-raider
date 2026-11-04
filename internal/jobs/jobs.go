package jobs

import (
	"context"
	"fmt"
	"time"
	"weakRaider/internal/config"
	domainJobs "weakRaider/internal/domain/jobs"

	"github.com/rs/zerolog"
)

type PlannedJob struct {
	Job     domainJobs.Job
	Period  time.Duration
	Depends []*PlannedJob
	next    time.Time
	success bool
}

// Build применяет интервалы к зарегистрированным работам и проверяет их зависимости.
func Build(cfg config.Jobs, registry []domainJobs.Job) ([]*PlannedJob, error) {
	defaultPeriod, err := getPeriod(cfg.Default)
	if err != nil {
		return nil, fmt.Errorf("jobs.default: %w", err)
	}
	result := make([]*PlannedJob, 0, len(registry))
	byName := make(map[string]*PlannedJob, len(registry))
	for _, job := range registry {
		name := job.Name()
		if _, exists := byName[name]; exists {
			return nil, fmt.Errorf("duplicate job %q", name)
		}
		planned := &PlannedJob{Job: job}
		for _, dependency := range job.Depends() {
			dep, exists := byName[dependency]
			if !exists {
				return nil, fmt.Errorf("job %s: dependency %q must be registered earlier", name, dependency)
			}
			planned.Depends = append(planned.Depends, dep)
		}
		period := defaultPeriod
		if configured := cfg.Entries[name].Period; configured != "" {
			period, err = getPeriod(configured)
			if err != nil {
				return nil, fmt.Errorf("job %s: %w", name, err)
			}
		}
		planned.Period = period
		byName[name] = planned
		result = append(result, planned)
	}
	for name := range cfg.Entries {
		if _, exists := byName[name]; !exists {
			return nil, fmt.Errorf("unknown job %q", name)
		}
	}
	return result, nil
}

func getPeriod(value string) (time.Duration, error) {
	period, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}
	if period <= 0 {
		return 0, fmt.Errorf("period must be positive")
	}
	return period, nil
}

// Run проверяет список раз в секунду. Первый запуск — сразу, следующие — через
// Period после завершения попытки. Ошибка зависимости блокирует работу до её успешного повтора.
func Run(ctx context.Context, list []*PlannedJob, logger *zerolog.Logger) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		for _, planned := range list {
			if ctx.Err() != nil {
				return
			}
			planned.run(ctx, logger)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (planned *PlannedJob) run(ctx context.Context, logger *zerolog.Logger) {
	for _, dependency := range planned.Depends {
		if !dependency.success {
			planned.success = false
			return
		}
	}
	if time.Now().Before(planned.next) {
		return
	}
	logger.Info().Str("job", planned.Job.Name()).Msg("job started")
	err := planned.Job.Run(ctx, planned.Period)
	planned.next = time.Now().Add(planned.Period)
	planned.success = err == nil
	if err != nil {
		logger.Error().Err(err).Str("job", planned.Job.Name()).Msg("job failed")
	} else {
		logger.Info().Str("job", planned.Job.Name()).Msg("job completed")
	}
}
