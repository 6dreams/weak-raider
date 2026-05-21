package manager

import (
	"context"
	"encoding/base64"
	"fmt"
	"weakRaider/internal/clients"
	"weakRaider/internal/clients/blizzard"
	"weakRaider/internal/clients/warcraftlogs"
	"weakRaider/internal/config"
)

type AuthManager struct {
	config   *config.Config
	blizzard *blizzard.ClientWithResponses
	logs     *warcraftlogs.ClientWithResponses
	keys     *clients.ApiKeys
}

func NewAuthManager(config *config.Config, blizzard *blizzard.ClientWithResponses, logs *warcraftlogs.ClientWithResponses, keys *clients.ApiKeys) *AuthManager {
	return &AuthManager{
		config:   config,
		blizzard: blizzard,
		logs:     logs,
		keys:     keys,
	}
}

// BlizzardKey Возвращает ключ для API Blizzard, который сразу можно использовать в заголовке авторизации
func (a *AuthManager) BlizzardKey() (string, error) {
	if a.keys.Blizzard == nil {
		resp, err := a.blizzard.AuthorizeWithFormdataBodyWithResponse(
			context.TODO(),
			&blizzard.AuthorizeParams{
				Authorization: fmt.Sprintf(
					"Basic %s",
					base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(
						"%s:%s",
						a.config.Auth.Blizzard.Client,
						a.config.Auth.Blizzard.Secret,
					))),
				),
			},
			blizzard.AuthorizeFormdataRequestBody{
				GrantType: "client_credentials",
			},
		)

		if err != nil || resp.JSON200 == nil {
			return "", err
		}

		key := fmt.Sprintf("Bearer %s", resp.JSON200.AccessToken)
		a.keys.Blizzard = &key
	}

	return *a.keys.Blizzard, nil
}

func (a *AuthManager) LogsKey() (string, error) {
	// todo: impl me

	return "", nil
}
