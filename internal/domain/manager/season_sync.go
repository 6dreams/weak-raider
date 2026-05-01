package manager

import (
	"weakRaider/internal/clients/blizzard"
	"weakRaider/internal/domain/repository"
)

type SeasonSync struct {
	seasonRepo  *repository.SeasonRepository
	blizzard    *blizzard.ClientWithResponses
	authManager *AuthManager
}

func NewSeasonSync(seasonRepo *repository.SeasonRepository, blizzard *blizzard.ClientWithResponses, auth *AuthManager) *SeasonSync {
	return &SeasonSync{
		seasonRepo:  seasonRepo,
		blizzard:    blizzard,
		authManager: auth,
	}
}

func (s *SeasonSync) Sync() error {
	_, err := s.seasonRepo.FindAll()
	if err != nil {
		return err
	}

	//s.blizzard.

	return nil
}
