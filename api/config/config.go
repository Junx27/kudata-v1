package config

type Config struct {
	BaseURLSurvey    string `env:"BASE_URL_SURVEY"`
	BaseURLUser      string `env:"BASE_URL_USER"`
	BaseURLPayment   string `env:"BASE_URL_PAYMENT"`
	BaseURLResponden string `env:"BASE_URL_RESPONDEN"`
	BaseURLImage     string `env:"BASE_URL_MINIO"`
}
