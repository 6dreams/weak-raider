package manager

import (
	"context"
	"weakRaider/internal/clients/wowaudit"
	"weakRaider/internal/domain/entity"
	"weakRaider/internal/domain/repository"
)

type CharacterSync struct {
	characterRepo *repository.CharacterRepository
	wowAudit      *wowaudit.ClientWithResponses
}

func NewCharacterSync(
	characterRepo *repository.CharacterRepository,
	wowAudit *wowaudit.ClientWithResponses,
) *CharacterSync {
	return &CharacterSync{characterRepo: characterRepo, wowAudit: wowAudit}
}

func (c *CharacterSync) Sync(guild *entity.Guild) error {
	characters, err := c.characterRepo.FindAll(guild) //может кэшировать персов? каждый раз выгружать эту гору персонажей из бд не перебор?
	if err != nil {
		return err
	}

	expResp, err := c.wowAudit.GetCharactersWithResponse(
		context.Background(),
		&wowaudit.GetCharactersParams{Authorization: guild.WowauditKey},
	)
	if err != nil {
		return err
	}
	if expResp.JSON200 == nil {
		return ErrInvalidApiResponse
	}

	for _, v := range *expResp.JSON200 {
		_, ok := characters[v.Name+v.Realm]
		character := &entity.Character{
			Name:    v.Name,
			Class:   v.Class,
			Rank:    v.Rank,
			Role:    v.Role,
			Realm:   v.Realm,
			GuildId: guild.ID,
			Guild:   guild,
		}

		if v.Note != nil {
			character.Note = *v.Note
		}

		if !ok {
			if err := c.characterRepo.Create(character); err != nil {
				return err
			}
			continue
		}

		if err := c.characterRepo.Update(character); err != nil {
			return err
		}
	}
	return nil
}
