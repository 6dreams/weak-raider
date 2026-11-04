package manager

import (
	"context"
	"fmt"
	"maps"
	"weakRaider/internal/clients/blizzard"
	"weakRaider/internal/domain/clients"
	"weakRaider/internal/domain/entity"
	"weakRaider/internal/domain/entity/types"
	"weakRaider/internal/domain/repository"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	choose int = 2
)

type TalentSync struct {
	logger         *zerolog.Logger
	talentRepo     *repository.TalentsRepository
	blizzardClient *clients.BlizzardClient
}

func NewTalentSync(
	logger *zerolog.Logger,
	talentRepo *repository.TalentsRepository,
	blizzardClient *clients.BlizzardClient,
) *TalentSync {
	return &TalentSync{
		talentRepo:     talentRepo,
		logger:         logger,
		blizzardClient: blizzardClient,
	}
}

func (t *TalentSync) Sync(ctx context.Context) error {
	talentTrees, err := t.blizzardClient.GetTalentTrees(ctx)

	if err != nil {
		return fmt.Errorf("talent sync, get talent trees: %w", err)
	}

	for _, spec := range talentTrees {
		res, err := t.blizzardClient.GetTalents(ctx, spec.TreeId, spec.SpecId)

		if err != nil {
			log.Err(err).Msg("sync talents, get talents")
			continue
		}

		templateTalents := entity.Talents{
			ClassName: types.NewTranslation(res.PlayableClass.Name),
			ClassID:   res.PlayableClass.Id,
			SpecName:  types.NewTranslation(res.PlayableSpecialization.Name),
			SpecID:    res.PlayableSpecialization.Id,
			Talents:   make(map[int]entity.TalentNode),
		}

		parseTalents(&templateTalents, res.ClassTalentNodes)
		parseTalents(&templateTalents, res.SpecTalentNodes)

		for _, heroTals := range res.HeroTalentTrees {
			classTalents := templateTalents
			classTalents.Talents = maps.Clone(templateTalents.Talents)
			if !isPlayable(classTalents.SpecID, heroTals.PlayableSpecializations) {
				continue
			}

			classTalents.HeroSpecID = heroTals.Id
			classTalents.HeroSpecName = types.NewTranslation(heroTals.Name)
			parseTalents(&classTalents, heroTals.HeroTalentNodes)
			err := t.talentRepo.Upsert(&classTalents)
			if err != nil {
				log.Err(err).Msg("sync talents, upsert classTalents")
				continue
			}
		}
	}

	return nil
}

func isPlayable(spec int, playable []blizzard.PlayableSpecialization) bool {
	for _, v := range playable {
		if v.Id == spec {
			return true
		}
	}
	return false
}

func parseTalents(talents *entity.Talents, nodes blizzard.TalentNodesSet) error {
	for _, node := range nodes {
		if node.NodeType.Id == choose {
			tals, err := node.Ranks.AsChoiceOfTalents()
			if err != nil {
				return fmt.Errorf("talent sync, as choice: %w", err)
			}
			for _, tal := range tals {
				for _, v := range tal.ChoiceOfTooltips {
					talents.Talents[v.SpellTooltip.Spell.Id] = entity.TalentNode{
						Ranks: len(tal.ChoiceOfTooltips),
						Name:  types.NewTranslation(v.SpellTooltip.Spell.Name),
					}
				}
			}
			continue
		}

		tals, err := node.Ranks.AsTalent()
		if err != nil {
			return fmt.Errorf("talent sync, as talent: %w", err)
		}
		for _, tal := range tals {
			talents.Talents[tal.Tooltip.SpellTooltip.Spell.Id] = entity.TalentNode{
				Ranks: len(tals),
				Name:  types.NewTranslation(tal.Tooltip.SpellTooltip.Spell.Name),
			}
		}
	}
	return nil
}
