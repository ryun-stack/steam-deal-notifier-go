package discordclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type DiscordClient struct {
	webhookURL url.URL
	httpClient *http.Client
	logger     *slog.Logger
}

func NewDiscordClient(webhookURL string, logger *slog.Logger) (*DiscordClient, error) {
	parsedURL, err := url.Parse(webhookURL)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook URL: %w", err)
	}
	return &DiscordClient{
		webhookURL: *parsedURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     logger,
	}, nil
}

func (c *DiscordClient) SendDiscordNotification(ctx context.Context, notification Notification) error {

	discordMessage := DiscordMessage{
		Embeds: []DiscordEmbed{
			{
				Title:       notification.Title,
				Description: fmt.Sprintf("%s\n価格: %s", notification.Description, notification.Price),
				URL:         notification.URL,
				Thumbnail: EmbedImage{
					URL: notification.Thumbnail,
				},
			},
		},
	}
	payloadJson, err := json.Marshal(discordMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL.String(), bytes.NewBuffer(payloadJson))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send notification, status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *DiscordClient) SendErrorNotificationIfNeeded(receivedError error, where string) error {
	if receivedError == nil {
		return nil
	}

	// 呼び出し元のctxがキャンセル・タイムアウト済みでも通知できるよう独立したコンテキストを使用する
	// タイムアウトはHTTPClientの30秒設定に委ねる
	notifyCtx := context.Background()

	discordMessage := DiscordMessage{
		Embeds: []DiscordEmbed{
			{
				Title:       fmt.Sprintf("Error Occurred in %s", where),
				Description: fmt.Sprintf("An error occurred: %v", receivedError),
				Color:       0xFF0000, // Red color for error
			},
		},
	}
	payloadJson, err := json.Marshal(discordMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	request, err := http.NewRequestWithContext(notifyCtx, http.MethodPost, c.webhookURL.String(), bytes.NewBuffer(payloadJson))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send notification, status code: %d", resp.StatusCode)
	}

	return nil
}
