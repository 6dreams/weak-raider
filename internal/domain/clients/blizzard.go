package clients

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	apiClients "weakRaider/internal/clients"
	"weakRaider/internal/clients/blizzard"
	"weakRaider/internal/config"
)

type Expansion struct {
	Id   int
	Name string
}

type TalentTree struct {
	TreeId string
	SpecId string
	Name   map[string]string
}

var talentTreePattern = regexp.MustCompile(`talent-tree/(\d+)/playable-specialization/(\d+)(?:[/?#]|$)`)

type ExpansionInstance struct {
	Id   int
	Name map[string]string
}

type ExpansionDetails struct {
	Id       int
	Name     map[string]string
	Dungeons []ExpansionInstance
	Raids    []ExpansionInstance
}

type InstanceEncounter struct {
	Id   int
	Name map[string]string
}

type BlizzardClient struct {
	client    *blizzard.ClientWithResponses
	battleNet *blizzard.ClientWithResponses
	config    *config.Config
	keys      *apiClients.ApiKeys
}

func NewBlizzardClient(client, battleNet *blizzard.ClientWithResponses, config *config.Config, keys *apiClients.ApiKeys) *BlizzardClient {
	return &BlizzardClient{client: client, battleNet: battleNet, config: config, keys: keys}
}

// getAuthKey возвращает токен Blizzard с префиксом Bearer для заголовка авторизации.
func (c *BlizzardClient) getAuthKey(ctx context.Context) (string, error) {
	if c.keys.Blizzard == nil {
		credentials := c.config.Auth.Blizzard.Client + ":" + c.config.Auth.Blizzard.Secret
		response, err := c.battleNet.AuthorizeWithFormdataBodyWithResponse(
			ctx,
			&blizzard.AuthorizeParams{
				Authorization: "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials)),
			},
			blizzard.AuthorizeFormdataRequestBody{GrantType: "client_credentials"},
		)
		if err != nil {
			return "", fmt.Errorf("authorize Blizzard: %w", err)
		}
		if response.JSON200 == nil {
			return "", fmt.Errorf("invalid Blizzard authorization response (HTTP %d)", response.StatusCode())
		}

		key := "Bearer " + response.JSON200.AccessToken

		c.keys.Blizzard = &key
	}
	return *c.keys.Blizzard, nil
}

// GetTalentTrees возвращает идентификаторы деревьев талантов, специализаций и их названия.
func (c *BlizzardClient) GetTalentTrees(ctx context.Context) ([]TalentTree, error) {
	key, err := c.getAuthKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("get Blizzard authorization: %w", err)
	}

	response, err := c.client.GetTalentTreesWithResponse(ctx, &blizzard.GetTalentTreesParams{
		Authorization:      key,
		BattlenetNamespace: blizzard.GetTalentTreesParamsBattlenetNamespaceStaticEu,
	})
	if err != nil {
		return nil, fmt.Errorf("get talent trees: %w", err)
	}
	if response.JSON200 == nil {
		return nil, fmt.Errorf("invalid talent trees response (HTTP %d)", response.StatusCode())
	}

	trees := make([]TalentTree, 0, len(response.JSON200.SpecTalentTrees))
	for _, tree := range response.JSON200.SpecTalentTrees {
		match := talentTreePattern.FindStringSubmatch(tree.Key.Href)
		if len(match) != 3 {
			continue
		}
		trees = append(trees, TalentTree{TreeId: match[1], SpecId: match[2], Name: tree.Name})
	}
	return trees, nil
}

// GetTalents возвращает данные талантов из ответа Blizzard API.
func (c *BlizzardClient) GetTalents(ctx context.Context, treeId, specId string) (*blizzard.TalentsResponse, error) {
	key, err := c.getAuthKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("get Blizzard authorization: %w", err)
	}

	response, err := c.client.GetTalentsWithResponse(
		ctx, treeId, specId,
		&blizzard.GetTalentsParams{
			Authorization:      key,
			BattlenetNamespace: blizzard.GetTalentsParamsBattlenetNamespaceStaticEu,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get talents: %w", err)
	}
	if response.JSON200 == nil {
		return nil, fmt.Errorf("invalid talents response (HTTP %d)", response.StatusCode())
	}
	return response.JSON200, nil
}

// GetExpansions возвращает разделы журнала приключений.
func (c *BlizzardClient) GetExpansions(ctx context.Context) (map[string]Expansion, error) {
	key, err := c.getAuthKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("get Blizzard authorization: %w", err)
	}

	response, err := c.client.JournalExpansionIndexWithResponse(
		ctx,
		&blizzard.JournalExpansionIndexParams{
			Authorization:      key,
			BattlenetNamespace: blizzard.JournalExpansionIndexParamsBattlenetNamespaceStaticEu,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get journal expansions: %w", err)
	}
	if response.JSON200 == nil {
		return nil, fmt.Errorf("invalid journal expansions response (HTTP %d)", response.StatusCode())
	}

	expansions := make(map[string]Expansion, len(response.JSON200.Tiers))
	for _, tier := range response.JSON200.Tiers {
		expansions[tier.Name] = Expansion{Id: tier.Id, Name: tier.Name}
	}
	return expansions, nil
}

// GetExpansionInstances возвращает сведения о разделе журнала приключений и его инстансах.
func (c *BlizzardClient) GetExpansionInstances(ctx context.Context, id int) (*ExpansionDetails, error) {
	key, err := c.getAuthKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("get Blizzard authorization: %w", err)
	}

	response, err := c.client.JournalExpansionByIdWithResponse(
		ctx,
		strconv.Itoa(id),
		&blizzard.JournalExpansionByIdParams{
			Authorization:      key,
			BattlenetNamespace: blizzard.JournalExpansionByIdParamsBattlenetNamespaceStaticEu,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get journal expansion %d: %w", id, err)
	}
	if response.JSON200 == nil {
		return nil, fmt.Errorf("invalid journal expansion %d response (HTTP %d)", id, response.StatusCode())
	}

	data := response.JSON200
	return &ExpansionDetails{
		Id:       data.Id,
		Name:     data.Name,
		Dungeons: expansionInstances(data.Dungeons),
		Raids:    expansionInstances(data.Raids),
	}, nil
}

func expansionInstances(instances *blizzard.EJSeasonInstance) []ExpansionInstance {
	if instances == nil {
		return nil
	}
	result := make([]ExpansionInstance, 0, len(*instances))
	for _, instance := range *instances {
		result = append(result, ExpansionInstance{Id: instance.Id, Name: instance.Name})
	}
	return result
}

// GetInstanceEncounters возвращает боссов инстанса из журнала приключений.
func (c *BlizzardClient) GetInstanceEncounters(ctx context.Context, id int) ([]InstanceEncounter, error) {
	key, err := c.getAuthKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("get Blizzard authorization: %w", err)
	}

	response, err := c.client.JournalInstanceWithResponse(
		ctx,
		strconv.Itoa(id),
		&blizzard.JournalInstanceParams{
			Authorization:      key,
			BattlenetNamespace: blizzard.JournalInstanceParamsBattlenetNamespaceStaticEu,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get instance %d encounters: %w", id, err)
	}
	if response.JSON200 == nil {
		return nil, fmt.Errorf("invalid instance %d response (HTTP %d)", id, response.StatusCode())
	}

	encounters := make([]InstanceEncounter, 0, len(response.JSON200.Encounters))
	for _, encounter := range response.JSON200.Encounters {
		encounters = append(encounters, InstanceEncounter{Id: encounter.Id, Name: encounter.Name})
	}
	return encounters, nil
}

// GetCurrentKeystoneSeason возвращает идентификатор текущего сезона Mythic+.
func (c *BlizzardClient) GetCurrentKeystoneSeason(ctx context.Context) (int, error) {
	key, err := c.getAuthKey(ctx)
	if err != nil {
		return 0, fmt.Errorf("get Blizzard authorization: %w", err)
	}

	response, err := c.client.MythicKeystoneSeasonsIndexWithResponse(
		ctx,
		&blizzard.MythicKeystoneSeasonsIndexParams{
			Authorization:      key,
			BattlenetNamespace: blizzard.MythicKeystoneSeasonsIndexParamsBattlenetNamespaceDynamicEu,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("get current Mythic+ season: %w", err)
	}
	if response.JSON200 == nil {
		return 0, fmt.Errorf("invalid Mythic+ seasons response (HTTP %d)", response.StatusCode())
	}
	season := response.JSON200.CurrentSeason
	if season == nil {
		return 0, errors.New("no current Mythic+ season found")
	}

	return season.Id, nil
}
