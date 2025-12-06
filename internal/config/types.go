package config

type Config struct {
	// Server
	Port string `env:"PORT"             envDefault:"8080"`

	// GitLab
	GitLabURL   string `env:"GITLAB_URL"       envDefault:"https://gitlab.com"`
	GitLabToken string `env:"GITLAB_TOKEN"     envRequired:"true"`

	// Logging
	LogLevel string `env:"LOG_LEVEL"        envDefault:"info"`
}
