package config

type Config struct {
	BaseURLSurvey string `env:"BASE_URL_SURVEY"`
	BaseURLUser   string `env:"BASE_URL_USER"`
}
