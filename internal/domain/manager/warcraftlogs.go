package manager

import (
	"context"
	"fmt"
	"net/http"
	"weakRaider/internal/clients/warcraftlogs"
	"weakRaider/internal/clients/warcraftlogsql"

	"github.com/Khan/genqlient/graphql"
)

type authTransport struct {
	authManager *AuthManager
	wrapped     http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	key, err := t.authManager.LogsKey()

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", key)

	return t.wrapped.RoundTrip(req)
}

type WarcraftLogs struct {
	auth *AuthManager
}

func NewWarcraftLogs(auth *AuthManager) *WarcraftLogs {
	return &WarcraftLogs{
		auth: auth,
	}
}

func (w *WarcraftLogs) Sync() error {
	httpClient := http.Client{
		Transport: &authTransport{
			authManager: w.auth,
			wrapped:     http.DefaultTransport,
		},
	}

	client := graphql.NewClient(warcraftlogs.ServerUrlHttpswwwWarcraftlogsComapiv2client, &httpClient)

	rankings, err := warcraftlogsql.GetEncounterRankings(context.TODO(), client, 3176, "Hunter", "BeastMastery")
	if err != nil {
		return err
	}

	fmt.Println("Encounter Name:", rankings.WorldData.Encounter.Name)

	return nil
}
