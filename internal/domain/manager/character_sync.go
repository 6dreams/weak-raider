package manager

import (
	"context"
	"weakRaider/internal/domain/clients"
	"weakRaider/internal/domain/entity"
	"weakRaider/internal/domain/repository"

	"github.com/rs/zerolog/log"
)

type CharacterSync struct {
	characterRepo *repository.CharacterRepository
	audit         *clients.Audit
}

func NewCharacterSync(
	characterRepo *repository.CharacterRepository,
	audit *clients.Audit,
) *CharacterSync {
	return &CharacterSync{characterRepo: characterRepo, audit: audit}
}

func (c *CharacterSync) Sync(ctx context.Context, guild *entity.Guild) error {
	characters, err := c.characterRepo.FindAll(guild)
	if err != nil {
		return err
	}

	auditCharacters, err := c.audit.GetCharacters(ctx, guild.WowauditKey)
	if err != nil {
		return err
	}

	for _, v := range auditCharacters {
		_, ok := characters[v.Name+v.Realm]
		character := &entity.Character{
			Name:    v.Name,
			Class:   v.Class,
			Rank:    v.Rank,
			Role:    v.Role,
			Realm:   v.Realm,
			Note:    v.Note,
			GuildId: guild.ID,
			Guild:   guild,
		}

		if !ok {
			if err := c.characterRepo.Create(character); err != nil {
				log.Err(err).Msgf("failed create new character: %v", character.Name+" "+character.Realm)
			}
			continue
		}

		if err := c.characterRepo.Update(character); err != nil {
			log.Err(err).Msgf("failed update character: %v", character.Name+" "+character.Realm)
		}
	}
	return nil
}
