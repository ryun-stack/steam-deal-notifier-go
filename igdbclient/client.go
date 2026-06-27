package igdbclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type IGDBClient struct {
	clientId     string
	clientSecret string
	endpoint     string
	token        string
	httpClient   *http.Client
	logger       *slog.Logger
}

func NewIGDBClient(clientId, clientSecret, endpoint string, logger *slog.Logger) *IGDBClient {
	return &IGDBClient{
		clientId:     clientId,
		clientSecret: clientSecret,
		endpoint:     endpoint,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		logger:       logger,
	}
}

func (c *IGDBClient) GetToken(ctx context.Context) error {
	params := url.Values{}
	params.Set("client_id", c.clientId)
	params.Set("client_secret", c.clientSecret)
	params.Set("grant_type", "client_credentials")

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get token, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	c.token = data.AccessToken
	return nil
}

func (c *IGDBClient) GetGames(ctx context.Context, title string) ([]IGDBGameInfo, error) {

	if c.token == "" {
		return nil, fmt.Errorf("token is not set")
	}

	base, err := url.Parse("https://api.igdb.com/v4/games")
	if err != nil {
		return nil, err
	}

	escapedTitle := strings.ReplaceAll(title, `\`, `\\`)
	escapedTitle = strings.ReplaceAll(escapedTitle, `"`, `\"`)
	requestbodyRaw := fmt.Sprintf(`fields name,summary; search "%s"; where version_parent = null;`, escapedTitle)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base.String(), strings.NewReader(requestbodyRaw))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.clientId)
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
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

	var result []IGDBGameInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
