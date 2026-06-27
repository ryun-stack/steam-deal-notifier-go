package openaiclient

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

type Client struct {
	openAIClient        openai.Client
	DeploymentName      string
	azureOpenAIEndpoint string
	azureOpenAIAPIKey   string
	logger              *slog.Logger
}

func NewClient(endpoint, apiKey, deploymentName string, logger *slog.Logger) *Client {
	return &Client{
		openAIClient: openai.NewClient(
			option.WithBaseURL(endpoint),
			option.WithAPIKey(apiKey),
		),
		DeploymentName:      deploymentName,
		azureOpenAIEndpoint: endpoint,
		azureOpenAIAPIKey:   apiKey,
		logger:              logger,
	}
}

func (c *Client) GetGameDescription(ctx context.Context, gameSummary string) (string, error) {

	resp, err := c.openAIClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(c.DeploymentName),
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfSystem: &openai.ChatCompletionSystemMessageParam{
					Content: openai.ChatCompletionSystemMessageParamContentUnion{
						OfString: openai.String("あなたはゲームの専門家で、ゲームを紹介して多くの人にプレイしてもらうことを生業としています"),
					},
				},
			},
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(fmt.Sprintf("以降のゲーム紹介を100文字で要約して。%s", gameSummary)),
					},
				},
			},
		},
	})

	if err != nil {
		c.logger.Error("failed to get game description", "error", err)
		return "", err
	}

	for _, choice := range resp.Choices {
		if choice.Message.Content != "" {
			return choice.Message.Content, nil
		}
	}
	return "", nil
}
