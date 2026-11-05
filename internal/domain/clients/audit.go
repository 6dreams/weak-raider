package clients

import (
	"context"
	"fmt"
	"weakRaider/internal/clients/wowaudit"
)

// Audit предоставляет доступ к данным wow_audit.
type Audit struct {
	client        *wowaudit.ClientWithResponses
	authorization string
}

// AuditSeason содержит сведения о сезоне из wow_audit.
type AuditSeason struct {
	Id   int
	Name string
}

// AuditCharacter содержит сведения о персонаже из wow_audit.
type AuditCharacter struct {
	BlizzardId int
	Name       string
	Class      string
	Rank       string
	Role       string
	Realm      string
	Note       string
}

func NewAudit(client *wowaudit.ClientWithResponses, authorization string) *Audit {
	return &Audit{client: client, authorization: authorization}
}

// GetCurrentSeason возвращает текущий сезон из wow_audit.
func (a *Audit) GetCurrentSeason(ctx context.Context) (*AuditSeason, error) {
	response, err := a.client.GetPeriodWithResponse(ctx, &wowaudit.GetPeriodParams{
		Authorization: a.authorization,
	})
	if err != nil {
		return nil, fmt.Errorf("get wow_audit current season: %w", err)
	}
	if response.JSON200 == nil {
		return nil, fmt.Errorf("invalid wow_audit period response (HTTP %d)", response.StatusCode())
	}
	season := response.JSON200.CurrentSeason
	return &AuditSeason{Id: season.Id, Name: season.Name}, nil
}

// GetCharacters возвращает персонажей wow_audit по ключу гильдии.
func (a *Audit) GetCharacters(ctx context.Context, guildKey string) ([]AuditCharacter, error) {
	response, err := a.client.GetCharactersWithResponse(ctx, &wowaudit.GetCharactersParams{
		Authorization: guildKey,
	})
	if err != nil {
		return nil, fmt.Errorf("get wow_audit characters: %w", err)
	}
	if response.JSON200 == nil {
		return nil, fmt.Errorf("invalid wow_audit characters response (HTTP %d)", response.StatusCode())
	}

	characters := make([]AuditCharacter, 0, len(*response.JSON200))
	for _, data := range *response.JSON200 {
		character := AuditCharacter{
			BlizzardId: data.BlizzardId,
			Name:       data.Name,
			Class:      data.Class,
			Rank:       data.Rank,
			Role:       data.Role,
			Realm:      data.Realm,
		}
		if data.Note != nil {
			character.Note = *data.Note
		}
		characters = append(characters, character)
	}
	return characters, nil
}
