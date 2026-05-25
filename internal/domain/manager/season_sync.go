package manager

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"
	"weakRaider/internal/clients/blizzard"
	"weakRaider/internal/clients/raidbots"
	"weakRaider/internal/clients/wowaudit"
	"weakRaider/internal/config"
	"weakRaider/internal/domain/entity"
	"weakRaider/internal/domain/entity/types"
	"weakRaider/internal/domain/repository"

	"github.com/rs/zerolog"
)

var ErrInvalidApiResponse = errors.New("invalid api response, empty or non-parsable body")
var ErrNothingNoUpdate = errors.New("nothing to update")

type SeasonSync struct {
	config        *config.Config
	Logger        *zerolog.Logger
	seasonRepo    *repository.SeasonRepository
	instanceRepo  *repository.InstanceRepository
	encounterRepo *repository.EncounterRepository
	blizzard      *blizzard.ClientWithResponses
	wowaudit      *wowaudit.ClientWithResponses
	raidBots      *raidbots.ClientWithResponses
	authManager   *AuthManager
}

func NewSeasonSync(
	config *config.Config,
	Logger *zerolog.Logger,
	seasonRepo *repository.SeasonRepository,
	instanceRepo *repository.InstanceRepository,
	encounterRepo *repository.EncounterRepository,
	blizzard *blizzard.ClientWithResponses,
	wowaudit *wowaudit.ClientWithResponses,
	raidBots *raidbots.ClientWithResponses,
	auth *AuthManager,
) *SeasonSync {
	return &SeasonSync{
		config:        config,
		Logger:        Logger,
		seasonRepo:    seasonRepo,
		instanceRepo:  instanceRepo,
		encounterRepo: encounterRepo,
		blizzard:      blizzard,
		wowaudit:      wowaudit,
		raidBots:      raidBots,
		authManager:   auth,
	}
}

func (s *SeasonSync) SyncInstances() error {
	instances, err := s.instanceRepo.FindWithInstances()
	if err != nil {
		return err
	}

	if len(instances) == 0 {
		return ErrNothingNoUpdate
	}

	key, err := s.authManager.BlizzardKey()
	if err != nil {
		return err
	}

	for _, instance := range instances {
		if s.isInstanceRequireUpdate(&instance) {
			data, err := s.blizzard.JournalInstanceWithResponse(context.TODO(), strconv.Itoa(instance.ID), &blizzard.JournalInstanceParams{
				Authorization:      key,
				BattlenetNamespace: blizzard.JournalInstanceParamsBattlenetNamespaceStaticEu,
			})

			if err != nil || data.JSON200 == nil {
				s.Logger.Error().Msgf("[InstanceSync] failed update instance `%v`: %v", instance.ID, err)
				continue
			}

			for _, encData := range data.JSON200.Encounters {
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
	}
	key, err := s.authManager.BlizzardKey()
	if err != nil {
		return err
	}

	for _, instance := range instances {
		if s.isInstanceRequireUpdate(&instance) {
			data, err := s.blizzard.JournalInstanceWithResponse(context.TODO(), strconv.Itoa(instance.ID), &blizzard.JournalInstanceParams{
				Authorization:      key,
				BattlenetNamespace: blizzard.JournalInstanceParamsBattlenetNamespaceStaticEu,
			})

			if err != nil || data.JSON200 == nil {
				s.Logger.Error().Msgf("[InstanceSync] failed update instance `%v`: %v", instance.ID, err)
				continue
			}

			for _, encData := range data.JSON200.Encounters {
				encounter := instance.GetEncounter(encData.Id)
				if encounter == nil {
					translation := types.NewTranslation(encData.Name)
					encounter = &entity.Encounter{
						Name:       types.NewTranslation(encData.Name),
						Instance:   &instance,
						BlizzardId: encData.Id,
					}

					err = s.encounterRepo.Upsert(encounter)
					if err != nil {
						s.Logger.Error().Msgf("[InstanceSync] failed update instance `%v`: %v", instance.ID, err)
					}
				}
			}
		}
	}

	return nil
}

func (s *SeasonSync) Sync() error {
	seasons, err := s.seasonRepo.FindAll()
	if err != nil {
		return err
	}

	instances, err := s.instanceRepo.FindAll()
	if err != nil {
		return err
	}

	currentSeason := getCurrentSeason(&seasons)
	if currentSeason != nil && s.isSeasonUpdateRequired(currentSeason) && !isWowAuditRequired(currentSeason) && len(instances) > 0 {
		return nil
	}

	key, err := s.authManager.BlizzardKey()
	if err != nil {
		return err
	}

	expResp, err := s.blizzard.JournalExpansionIndexWithResponse(
		context.TODO(),
		&blizzard.JournalExpansionIndexParams{
			Authorization:      key,
			BattlenetNamespace: blizzard.JournalExpansionIndexParamsBattlenetNamespaceStaticEu,
		},
	)

	if err != nil {
		return err
	}
	if expResp.JSON200 == nil {
		return ErrInvalidApiResponse
	}

	seasonId := 0
	for _, tier := range expResp.JSON200.Tiers {
		if tier.Name == "Current Season" {
			seasonId = tier.Id
		}
	}

	if seasonId == 0 {
		return errors.New("no current seasonId found")
	}

	season, ok := seasons[seasonId]
	updateSeason := false

	if !ok {
		season = entity.Season{
			BlizzardId: seasonId,
			Name:       season.Name,
			IsCurrent:  true,
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
		audit, err := s.wowaudit.GetPeriodWithResponse(context.TODO(), &wowaudit.GetPeriodParams{
			Authorization: s.config.Auth.WowAudit,
		})

		if err == nil && audit.JSON200 != nil {
			season.Name = audit.JSON200.CurrentSeason.Name
			season.WowAuditId = sql.NullInt64{
				Int64: int64(audit.JSON200.CurrentSeason.Id),
				Valid: true,
			}
			updateSeason = true
		}
	}

	if updateSeason {
		if err := s.seasonRepo.Update(&season); err != nil {
			return err
		}
	}

	for _, seasonTemp := range seasons {
		if seasonTemp.IsCurrent && season.ID != seasonTemp.ID {
			seasonTemp.IsCurrent = false
			// ignore errors
			_ = s.seasonRepo.Update(&season)
		}
	}

	if !s.isSeasonUpdateRequired(&season) && !currentWasChanged && len(instances) > 0 {
		return nil
	}

	validInstanceIds, err := s.getSeasonInstances()
	if err != nil {
		return err
	}

	apiInstances, err := s.blizzard.JournalExpansionByIdWithResponse(
		context.TODO(),
		strconv.Itoa(seasonId),
		&blizzard.JournalExpansionByIdParams{
			Authorization:      key,
			BattlenetNamespace: blizzard.JournalExpansionByIdParamsBattlenetNamespaceStaticEu,
		},
	)

	if err != nil {
		return err
	}
	if apiInstances.JSON200 == nil {
		return ErrInvalidApiResponse
	}

	s.storeInstances(validInstanceIds, &instances, &season, apiInstances.JSON200.Dungeons, false)
	s.storeInstances(validInstanceIds, &instances, &season, apiInstances.JSON200.Raids, true)

	return nil
}

func (s *SeasonSync) storeInstances(validInstances map[int]bool, instances *entity.InstanceMap, season *entity.Season, api *blizzard.EJSeasonInstance, isRaid bool) {
	if api == nil {
		return
	}

	for _, apiInstance := range *api {
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
	return !season.WowAuditId.Valid || "" == season.Name
}

func (s *SeasonSync) isInstanceRequireUpdate(instance *entity.Instance) bool {
	return len(instance.Encounters) == 0 || !instance.UpdatedAt.Add(s.config.Tickers.SyncSeasons).After(time.Now())
}

func (s *SeasonSync) isSeasonUpdateRequired(season *entity.Season) bool {
	return !season.UpdatedAt.Add(s.config.Tickers.SyncSeasons).After(time.Now())
}

func (s *SeasonSync) getSeasonInstances() (map[int]bool, error) {
	instances := make(map[int]bool, 0)

	data, err := s.raidBots.GetInstancesWithResponse(context.TODO())
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
