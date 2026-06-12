package manager

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"weakRaider/internal/clients/blizzard"
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
	logger      *zerolog.Logger
	talentRepo  *repository.TalentsRepository
	authManager *AuthManager
	blizzard    *blizzard.ClientWithResponses
}

func NewTalentSync(
	logger *zerolog.Logger,
	talentRepo *repository.TalentsRepository,
	authManager *AuthManager,
	blizzard *blizzard.ClientWithResponses,
) *TalentSync {
	return &TalentSync{
		talentRepo:  talentRepo,
		authManager: authManager,
		logger:      logger,
		blizzard:    blizzard,
	}
}

func (t *TalentSync) Sync() error {
	key, err := t.authManager.BlizzardKey()
	if err != nil {
		return fmt.Errorf("Talent sync, BlizzardKey: %w", err)
	}

	re, err := regexp.Compile(`talent-tree\/(\d+)\/playable-specialization\/(\d+)`)
	if err != nil {
		fmt.Errorf("Talent sync, regexp compile: %w", err)
	}

	talTreesResp, err := t.blizzard.GetTalentTreesWithResponse(
		context.Background(),
		&blizzard.GetTalentTreesParams{
			BattlenetNamespace: blizzard.GetTalentTreesParamsBattlenetNamespaceStaticEu,
			Authorization:      key,
		},
	)

	if err != nil {
		return fmt.Errorf("talent sync, get talent trees: %w", err)
	}

	for _, spec := range talTreesResp.JSON200.SpecTalentTrees {
		r := re.FindStringSubmatch(spec.Key.Href)
		if len(r) != 3 {
			log.Err(errors.New("bad regexp")).Msg("sync talents, find substring match")
			continue
		}

		resp, err := t.blizzard.GetTalentsWithResponse(
			context.Background(),
			blizzard.TalentTreeId(r[1]),
			blizzard.SpecId(r[2]),
			&blizzard.GetTalentsParams{
				BattlenetNamespace: blizzard.GetTalentsParamsBattlenetNamespaceStaticEu,
				Authorization:      key,
			},
		)

		if err != nil {
			log.Err(err).Msg("sync talents, get talents")
			continue
		}

		if resp.JSON200 == nil {
			log.Err(ErrInvalidApiResponse).Msg("sync talents, nil json200")
			continue
		}

		res := resp.JSON200

		templateTalents := entity.Talents{
			ClassName: types.NewTranslation(res.PlayableClass.Name),
			ClassID:   res.PlayableClass.Id,
			SpecName:  types.NewTranslation(res.PlayableSpecialization.Name),
			SpecID:    res.PlayableSpecialization.Id,
			Talents:   make(map[int]entity.TalentNode),
		}

		parseTalents(&templateTalents, resp.JSON200.ClassTalentNodes)
		parseTalents(&templateTalents, resp.JSON200.SpecTalentNodes)

		for _, heroTals := range resp.JSON200.HeroTalentTrees {
			classTalents := templateTalents
			if !isPlayable(classTalents.SpecID, heroTals.PlayableSpecializations) {
				continue
			}

			classTalents.HeroSpecID = heroTals.Id
			classTalents.HeroSpecName = types.NewTranslation(heroTals.Name)
			parseTalents(&classTalents, heroTals.HeroTalentNodes)
			fmt.Println(classTalents)
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
				return fmt.Errorf("Talent sync, as choice: %w", err)
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
			return fmt.Errorf("Talent sync, as talent: %w", err)
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
