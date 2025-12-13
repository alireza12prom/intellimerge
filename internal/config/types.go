package config

type Config struct {
	// Server
	Port          string `env:"PORT"             envDefault:"8080"`
	WebhookPath   string `env:"WEBHOOK_PATH"     envDefault:"/webhook"`
	WebhookSecret string `env:"WEBHOOK_SECRET"   envRequired:"true"`

	// GitLab
	GitLabURL   string `env:"GITLAB_URL"       envDefault:"https://gitlab.com"`
	GitLabToken string `env:"GITLAB_TOKEN"     envRequired:"true"`

	// Jira
	JiraURL      string `env:"JIRA_URL"        envRequired:"true"`
	JiraEmail    string `env:"JIRA_EMAIL"      envRequired:"true"`
	JiraAPIToken string `env:"JIRA_API_TOKEN" envRequired:"true"`

	// LLM
	LLMProvider string `env:"LLM_PROVIDER"     envDefault:"openai"` // openai, anthropic, etc.
	LLMAPIKey   string `env:"LLM_API_KEY"     envRequired:"true"`
	LLMBaseURL  string `env:"LLM_BASE_URL"    envDefault:""`      // Optional, for custom endpoints
	LLMModel    string `env:"LLM_MODEL"       envDefault:"gpt-4"` // Model name

	// Logging
	LogLevel string `env:"LOG_LEVEL"        envDefault:"info"`
}
