package config

type Config struct {
	Port        string `env:"PORT" envDefault:"8080"`
	GitLabURL   string `env:"GITLAB_URL" envDefault:"https://gitlab.com"`
	GitLabToken string `env:"GITLAB_TOKEN" envRequired:"true"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
}
