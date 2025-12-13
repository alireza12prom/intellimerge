package llm

import "github.com/alireza12prom/intellimerge/internal/llm/providers"

func NewClient(provider string, apiKey string) Provider {
	var p Provider

	switch provider {
	case "openai":
		p = providers.NewOpenAIProvider(apiKey, "gpt-4o")
	default:
		p = providers.NewOpenAIProvider(apiKey, "gpt-4o")
	}

	return p
}
