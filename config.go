package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ITADAPIKey              string
	DiscordWebhookURL       string
	AOAIChatCompletionModel string
	AzureOpenAIEndpoint     string
	AzureOpenAIAPIKey       string
	TwitchClientID          string
	TwitchClientSecret      string
	TwitchOAuthEndpoint     string
	DealsFilterTags         []string
}

func LoadConfig() *Config {
	return &Config{
		ITADAPIKey:              os.Getenv("ITAD_API_KEY"),
		DiscordWebhookURL:       os.Getenv("DISCORD_WEBHOOK_URL"),
		AOAIChatCompletionModel: os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL"),
		AzureOpenAIEndpoint:     os.Getenv("AZURE_OPENAI_ENDPOINT"),
		AzureOpenAIAPIKey:       os.Getenv("AZURE_OPENAI_API_KEY"),
		TwitchClientID:          os.Getenv("TWITCH_CLIENT_ID"),
		TwitchClientSecret:      os.Getenv("TWITCH_CLIENT_SECRET"),
		TwitchOAuthEndpoint:     os.Getenv("TWITCH_OAUTH_ENDPOINT"),
		DealsFilterTags:         parseFilterTags(os.Getenv("ITAD_FILTER_TAGS")),
	}
}

// parseFilterTags はカンマ区切りの文字列をタグスライスに変換する。
// 未設定の場合はデフォルト値を返す。
func parseFilterTags(raw string) []string {
	if raw == "" {
		return []string{"Online Co-Op", "Automation"}
	}
	tags := strings.Split(raw, ",")
	for i, t := range tags {
		tags[i] = strings.TrimSpace(t)
	}
	return tags
}

func (c *Config) Validate() error {
	if c.ITADAPIKey == "" {
		return fmt.Errorf("ITAD_API_KEY is not set")
	}
	if c.DiscordWebhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_URL is not set")
	}
	if c.AOAIChatCompletionModel == "" {
		return fmt.Errorf("AOAI_CHAT_COMPLETIONS_MODEL is not set")
	}
	if c.AzureOpenAIEndpoint == "" {
		return fmt.Errorf("AZURE_OPENAI_ENDPOINT is not set")
	}
	if c.AzureOpenAIAPIKey == "" {
		return fmt.Errorf("AZURE_OPENAI_API_KEY is not set")
	}
	if c.TwitchClientID == "" {
		return fmt.Errorf("TWITCH_CLIENT_ID is not set")
	}
	if c.TwitchClientSecret == "" {
		return fmt.Errorf("TWITCH_CLIENT_SECRET is not set")
	}
	if c.TwitchOAuthEndpoint == "" {
		return fmt.Errorf("TWITCH_OUTH_ENDPOINT is not set")
	}
	return nil
}
