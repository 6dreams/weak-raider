package manager

import (
	"context"
	"encoding/base64"
	"fmt"
	"weakRaider/internal/clients"
	"weakRaider/internal/clients/warcraftlogs"
	"weakRaider/internal/config"
)

type AuthManager struct {
	config *config.Config
	wl     *warcraftlogs.ClientWithResponses
	logs   *warcraftlogs.ClientWithResponses
	keys   *clients.ApiKeys
}

func NewAuthManager(
	config *config.Config,
	wl *warcraftlogs.ClientWithResponses,
	keys *clients.ApiKeys,
) *AuthManager {
	return &AuthManager{
		config: config,
		wl:     wl,
		keys:   keys,
	}
}

func (a *AuthManager) LogsKey(ctx context.Context) (string, error) {
	if a.keys.Logs == nil {
		resp, err := a.wl.AuthorizeWithFormdataBodyWithResponse(
			ctx,
			&warcraftlogs.AuthorizeParams{
				Authorization: fmt.Sprintf(
					"Basic %s",
					base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(
						"%s:%s",
						a.config.Auth.Logs.Client,
						a.config.Auth.Logs.Secret,
					))),
				),
			},
			warcraftlogs.AuthorizeFormdataRequestBody{
				GrantType: "client_credentials",
			},
		)

		if err != nil || resp.JSON200 == nil {
			return "", err
		}

		key := fmt.Sprintf("Bearer %s", resp.JSON200.AccessToken)
		a.keys.Logs = &key
	}

	return *a.keys.Logs, nil
}
