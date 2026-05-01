package manager

import (
	"context"
	"errors"
	"strconv"
	"weakRaider/internal/clients/blizzard"
	"weakRaider/internal/domain/entity"
	"weakRaider/internal/domain/repository"
)

var ErrInvalidApiResponse = errors.New("invalid api response, empty or non-parsable body")

type SeasonSync struct {
	seasonRepo   *repository.SeasonRepository
	instanceRepo *repository.InstanceRepository
	blizzard     *blizzard.ClientWithResponses
	authManager  *AuthManager
}

func NewSeasonSync(
	seasonRepo *repository.SeasonRepository,
	instanceRepo *repository.InstanceRepository,
	blizzard *blizzard.ClientWithResponses,
	auth *AuthManager,
) *SeasonSync {
	return &SeasonSync{
		seasonRepo:   seasonRepo,
		instanceRepo: instanceRepo,
		blizzard:     blizzard,
		authManager:  auth,
	}
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
	if !ok {
		season := &entity.Season{
			BlizzardId: seasonId,
			Name:       season.Name,
		}

		if err := s.seasonRepo.Upsert(season); err != nil {
			return err
		}
	}

	apiInstances, err := s.blizzard.JournalExpansionByIdWithResponse(
		context.TODO(),
		strconv.Itoa(seasonId),
		&blizzard.JournalExpansionByIdParams{
			Authorization:      key,
			BattlenetNamespace: blizzard.ServerUrlHttpseuApiBlizzardCom,
		},
	)

	if err != nil {
		return err
	}
	if apiInstances.JSON200 == nil {
		return ErrInvalidApiResponse
	}

	s.storeInstances(&instances, &season, apiInstances.JSON200.Dungeons, false)
	s.storeInstances(&instances, &season, apiInstances.JSON200.Raids, true)

	return nil
}

func (s *SeasonSync) storeInstances(instances *entity.InstanceMap, season *entity.Season, api *blizzard.EJSeasonInstance, isRaid bool) {
	if api == nil {
		return
	}

	for _, apiInstance := range *api {
		instance, ok := (*instances)[apiInstance.Id]

		if !ok {
			instance = entity.Instance{
				ID:       apiInstance.Id,
				Name:     apiInstance.Name,
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
