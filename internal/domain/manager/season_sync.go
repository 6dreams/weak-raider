package manager

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"weakRaider/internal/clients/raidbots"
	"weakRaider/internal/domain/clients"
	"weakRaider/internal/domain/entity"
	"weakRaider/internal/domain/entity/types"
	"weakRaider/internal/domain/repository"

	"github.com/rs/zerolog"
)

var ErrInvalidApiResponse = errors.New("invalid api response, empty or non-parsable body")
var ErrNothingNoUpdate = errors.New("nothing to update")

type SeasonSync struct {
	Logger         *zerolog.Logger
	seasonRepo     *repository.SeasonRepository
	instanceRepo   *repository.InstanceRepository
	encounterRepo  *repository.EncounterRepository
	blizzardClient *clients.BlizzardClient
	audit          *clients.Audit
	raidBots       *raidbots.ClientWithResponses
}

func NewSeasonSync(
	Logger *zerolog.Logger,
	seasonRepo *repository.SeasonRepository,
	instanceRepo *repository.InstanceRepository,
	encounterRepo *repository.EncounterRepository,
	blizzardClient *clients.BlizzardClient,
	audit *clients.Audit,
	raidBots *raidbots.ClientWithResponses,
) *SeasonSync {
	return &SeasonSync{
		Logger:         Logger,
		seasonRepo:     seasonRepo,
		instanceRepo:   instanceRepo,
		encounterRepo:  encounterRepo,
		blizzardClient: blizzardClient,
		audit:          audit,
		raidBots:       raidBots,
	}
}

func (s *SeasonSync) SyncInstances(ctx context.Context, period time.Duration) error {
	instances, err := s.instanceRepo.FindWithInstances()
	if err != nil {
		return err
	}

	if len(instances) == 0 {
		return ErrNothingNoUpdate
	}

	for _, instance := range instances {
		if !s.isInstanceRequireUpdate(&instance, period) {
			continue
		}
		encounters, err := s.blizzardClient.GetInstanceEncounters(ctx, instance.ID)

		if err != nil {
			s.Logger.Error().Msgf("[InstanceSync] failed update instance `%v`: %v", instance.ID, err)
			continue
		}

		for _, encData := range encounters {
			encounter := instance.GetEncounter(encData.Id)
			if encounter == nil {
				encounter = &entity.Encounter{
					Name:       types.NewTranslation(encData.Name),
					Instance:   &instance,
					BlizzardId: encData.Id,
				}

				err = s.encounterRepo.Upsert(encounter)
				if err != nil {
					s.Logger.Error().Msgf("[InstanceSync] failed upsert encounter `%v`: %v", instance.ID, err)
				}
			} else {
				encounter.Name = types.NewTranslation(encData.Name)
				err = s.encounterRepo.Update(encounter)
				if err != nil {
					s.Logger.Error().Msgf("[InstanceSync] failed update encounter `%v`: %v", instance.ID, err)
				}
			}
		}
		if err := s.instanceRepo.Update(&instance); err != nil {
			s.Logger.Error().Msgf("[InstanceSync] failed update instance: `%v`: %v", instance.ID, err)
		}
	}

	return nil
}

func (s *SeasonSync) Sync(ctx context.Context, period time.Duration) error {
	seasons, err := s.seasonRepo.FindAll()
	if err != nil {
		return err
	}

	instances, err := s.instanceRepo.FindAll()
	if err != nil {
		return err
	}

	currentSeason := getCurrentSeason(&seasons)
	if currentSeason != nil && !s.isSeasonUpdateRequired(currentSeason, period) && !isWowAuditRequired(currentSeason) && len(instances) > 0 {
		return nil
	}

	expansions, err := s.blizzardClient.GetExpansions(ctx)
	if err != nil {
		return err
	}

	seasonId, err := s.blizzardClient.GetCurrentKeystoneSeason(ctx)
	if err != nil {
		return err
	}

	expansion, ok := expansions["Current Season"]
	if !ok {
		return errors.New("no current expansionId found")
	}

	season, ok := seasons[seasonId]
	updateSeason := false

	if !ok {
		season = entity.Season{
			BlizzardId:  seasonId,
			ExpansionId: expansion.Id,
			IsCurrent:   true,
		}

		updateSeason = true
	}

	currentWasChanged := false
	if !season.IsCurrent {
		season.IsCurrent = true
		updateSeason = true
		currentWasChanged = true
	}

	if isWowAuditRequired(&season) {
		auditSeason, err := s.audit.GetCurrentSeason(ctx)
		if err != nil {
			return err
		}
		season.Name = auditSeason.Name
		season.WowAuditId = sql.NullInt64{
			Int64: int64(auditSeason.Id),
			Valid: true,
		}
		updateSeason = true
	}

	if updateSeason {
		if err := s.seasonRepo.Update(&season); err != nil {
			return err
		}
	}

	for _, seasonTemp := range seasons {
		if seasonTemp.IsCurrent && season.ID != seasonTemp.ID {
			seasonTemp.IsCurrent = false
			if err := s.seasonRepo.Update(&seasonTemp); err != nil {
				return err
			}
		}
	}

	if !s.isSeasonUpdateRequired(&season, period) && !currentWasChanged && len(instances) > 0 {
		return nil
	}

	validInstanceIds, err := s.getSeasonInstances(ctx)
	if err != nil {
		return err
	}

	expansionInstances, err := s.blizzardClient.GetExpansionInstances(ctx, season.ExpansionId)
	if err != nil {
		return err
	}

	s.storeInstances(validInstanceIds, &instances, &season, expansionInstances.Dungeons, false)
	s.storeInstances(validInstanceIds, &instances, &season, expansionInstances.Raids, true)

	return nil
}

func (s *SeasonSync) storeInstances(validInstances map[int]bool, instances *entity.InstanceMap, season *entity.Season, api []clients.ExpansionInstance, isRaid bool) {
	for _, apiInstance := range api {
		instance, exists := (*instances)[apiInstance.Id]
		_, valid := validInstances[apiInstance.Id]

		if !exists && valid {
			instance = entity.Instance{
				ID:       apiInstance.Id,
				Name:     types.NewTranslation(apiInstance.Name),
				SeasonId: season.ID,
				IsRaid:   isRaid,
			}

			err := s.instanceRepo.Upsert(&instance)
			if err == nil {
				(*instances)[instance.ID] = instance
			}
		}
	}
}

func getCurrentSeason(seasons *entity.SeasonMap) *entity.Season {
	for _, season := range *seasons {
		if season.IsCurrent {
			return &season
		}
	}

	return nil
}

func isWowAuditRequired(season *entity.Season) bool {
	return !season.WowAuditId.Valid || season.Name == ""
}

func (s *SeasonSync) isInstanceRequireUpdate(instance *entity.Instance, period time.Duration) bool {
	return len(instance.Encounters) == 0 || time.Since(instance.UpdatedAt) >= period
}

func (s *SeasonSync) isSeasonUpdateRequired(season *entity.Season, period time.Duration) bool {
	return time.Since(season.UpdatedAt) >= period
}

func (s *SeasonSync) getSeasonInstances(ctx context.Context) (map[int]bool, error) {
	instances := make(map[int]bool, 0)

	data, err := s.raidBots.GetInstancesWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	if data.JSON200 == nil {
		return nil, ErrInvalidApiResponse
	}

	for _, instance := range *data.JSON200 {
		// Mythic Plus Dungeons.
		// У них есть вменяемый тип.
		if instance.Type == "mplus-chest" {
			for _, encounter := range instance.Encounters {
				instances[encounter.Id] = true
			}
		}

		// Рейды, у них есть тип, но с этим же типом есть сомнительные группы.
		if instance.Id > 0 && instance.Type == "raid" {
			instances[instance.Id] = true
		}
	}

	return instances, nil
}
