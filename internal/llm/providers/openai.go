package providers

import (
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenAIProvider struct {
	client *openai.Client
}

func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &OpenAIProvider{client: &client}
}
