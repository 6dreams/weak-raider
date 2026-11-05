package jobs

import (
	domainJobs "weakRaider/internal/domain/jobs"
	"weakRaider/internal/domain/manager"
	"weakRaider/internal/domain/repository"
)

// Registry список возможных задач сервиса, зависимости регистрируются раньше.
func Registry(seasons *manager.SeasonSync, characters *manager.CharacterSync, talents *manager.TalentSync, logs *manager.WarcraftLogs, guilds *repository.GuildRepository) []domainJobs.Job {
	return []domainJobs.Job{
		domainJobs.NewSeasons(seasons),
		domainJobs.NewInstances(seasons),
		domainJobs.NewCharacters(characters, guilds),
		domainJobs.NewTalents(talents),
		domainJobs.NewWarcraftLogs(logs),
	}
}
