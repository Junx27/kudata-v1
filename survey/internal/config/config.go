package config

import (
	"github.com/minio/minio-go/v7"
)

var (
	AppConfig   Config
	MinioClient *minio.Client
)

type Config struct {
	Port string `env:"PORT" envDefault:"8000"`
	Env  string `env:"ENV" envDefault:"dev"`

	DatabaseUsername string `env:"DB_USERNAME"`
	DatabasePassword string `env:"DB_PASSWORD"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabasePort     string `env:"DB_PORT"`
	DatabaseName     string `env:"DB_NAME"`

	MigrationPath string `env:"MIGRATION_PATH" envDefault:"/app/migrations"`

	AMQPHost string `env:"AMQP_HOST"`

	MinioHost      string `env:"MINIO_HOST" envDefault:"minio:9000"`
	MinioPort      string `env:"MINIO_PORT"`
	MinioAccessKey string `env:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `env:"MINIO_SECRET_KEY"`
	MinioBucket    string `env:"MINIO_BUCKET"`
	MinioUseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"false"`
}
