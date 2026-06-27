package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/azure/azure-functions-golang-worker/sdk"
	"github.com/azure/azure-functions-golang-worker/sdk/bindings"
	"github.com/azure/azure-functions-golang-worker/worker"

	"github.com/ryun-stack/steam-deal-discord-notifier/discordclient"
	"github.com/ryun-stack/steam-deal-discord-notifier/igdbclient"
	"github.com/ryun-stack/steam-deal-discord-notifier/itadclient"
	"github.com/ryun-stack/steam-deal-discord-notifier/openaiclient"
)

func main() {

	app := sdk.FunctionApp()
	app.Timer("steamDealNotifier", SteamDealNotifier, sdk.WithSchedule("0 0 9 * * 1"))
	worker.Start(app)
}

func SteamDealNotifier(ctx context.Context, timerInfo bindings.TimerInfo) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	cfg := LoadConfig()
	if err := cfg.Validate(); err != nil {
		logger.Error("Config validation error", "error", err)
		os.Exit(1)
	}

	apiKey := cfg.ITADAPIKey
	discordWebhookURL := cfg.DiscordWebhookURL

	discordClient, err := discordclient.NewDiscordClient(discordWebhookURL, logger)
	if err != nil {
		logger.Error("Failed to create Discord client", "error", err)
		os.Exit(1)
	}

	client := itadclient.NewITADClient(apiKey, logger, cfg.DealsFilterTags)
	result, err := client.GetDeals(ctx)
	if err != nil || len(result.Results) == 0 {
		discordClient.SendErrorNotificationIfNeeded(err, "ITAD API - GetDeals")
		return err
	}
	deal := result.GetRandomDeal()

	game, err := client.GetGameDetail(ctx, deal.ID)
	if err != nil {
		discordClient.SendErrorNotificationIfNeeded(err, "ITAD API - GetGameDetail")
		return err
	}

	igdbCli := igdbclient.NewIGDBClient(cfg.TwitchClientID, cfg.TwitchClientSecret, cfg.TwitchOAuthEndpoint, logger)
	err = igdbCli.GetToken(ctx)
	if err != nil {
		discordClient.SendErrorNotificationIfNeeded(err, "IGDB API - GetToken")
		return err
	}

	igdb, err := igdbCli.GetGames(ctx, game.Title)
	if err != nil || len(igdb) == 0 {
		discordClient.SendErrorNotificationIfNeeded(err, "IGDB API - GetGames")
		return err
	}

	openaiCli := openaiclient.NewClient(cfg.AzureOpenAIEndpoint, cfg.AzureOpenAIAPIKey, cfg.AOAIChatCompletionModel, logger)

	gameDescription, err := openaiCli.GetGameDescription(ctx, igdb[0].Summary)
	if err != nil {
		discordClient.SendErrorNotificationIfNeeded(err, "OpenAI API - GetGameDescription")
		return err
	}

	obj := discordclient.Notification{
		Title:       game.Title,
		Description: gameDescription,
		Price:       deal.Deal.Price.String(),
		URL:         deal.Deal.URL,
		Thumbnail:   game.Assets.GetBestBanner(),
	}

	err = discordClient.SendDiscordNotification(ctx, obj)
	if err != nil {
		discordClient.SendErrorNotificationIfNeeded(err, "Discord API - SendDiscordNotification")
		return err
	}

	return nil
}
