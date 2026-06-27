package itadclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type ITADClient struct {
	apiKey     string
	httpClient *http.Client
	logger     *slog.Logger
	filterTags []string
}

func NewITADClient(apiKey string, logger *slog.Logger, filterTags []string) *ITADClient {
	return &ITADClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     logger,
		filterTags: filterTags,
	}
}

func (c *ITADClient) GetDeals(ctx context.Context) (*Deals, error) {
	base, err := url.Parse("https://api.isthereanydeal.com/deals/v2")
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("key", c.apiKey)
	base.RawQuery = params.Encode()

	requestbody := DealsRequest{
		Limit:   40,
		Shops:   []int{SteamShopID},
		Country: JPCountryCode,
		Filter: DealsFilter{
			Platforms:  PlatformWindows,
			Types:      []GameType{GameTypeIntGame},
			SteamPerc:  FilterRange{Min: MinReviewRating, Max: 100},
			SteamCount: FilterRange{Min: MinReviewCount, Max: 1000000},
			TagsUnion:  c.filterTags,
		},
	}

	requestJson, err := json.Marshal(requestbody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, base.String(), bytes.NewBuffer(requestJson))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("failed to get deals, status code: %d", resp.StatusCode)
	}

	// レスポンスボディを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Deals
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *ITADClient) GetGameDetail(ctx context.Context, id string) (*GameDetail, error) {
	base, err := url.Parse("https://api.isthereanydeal.com/games/info/v2")
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("id", id)
	base.RawQuery = params.Encode()

	htpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, base.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(htpRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // 必ずCloseする

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("failed to get game detail, status code: %d", resp.StatusCode)
	}

	// レスポンスボディを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result GameDetail
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
