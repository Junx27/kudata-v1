package config

import "github.com/minio/minio-go/v7"

var (
	AppConfig   Config
	MinioClient *minio.Client
)

type Config struct {
	BaseURLSurvey    string `env:"BASE_URL_SURVEY"`
	BaseURLUser      string `env:"BASE_URL_USER"`
	BaseURLPayment   string `env:"BASE_URL_PAYMENT"`
	BaseURLResponden string `env:"BASE_URL_RESPONDEN"`
	BaseURLImage     string `env:"BASE_URL_MINIO"`

	AMQPHost string `env:"AMQP_HOST"`

	MinioHost      string `env:"MINIO_HOST" envDefault:"minio:9000"`
	MinioPort      string `env:"MINIO_PORT"`
	MinioAccessKey string `env:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `env:"MINIO_SECRET_KEY"`
	MinioBucket    string `env:"MINIO_BUCKET"`
	MinioUseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"false"`
}
