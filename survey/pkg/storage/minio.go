package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"survey/internal/config"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func UploadImageToMinio(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	filename := uuid.New().String() + ext
	objectName := "surveys/" + filename

	_, err := config.MinioClient.PutObject(context.Background(),
		config.AppConfig.MinioBucket,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", err
	}

	// Buat URL untuk akses file
	url := fmt.Sprintf("http://%s:%s/%s/%s",
		config.AppConfig.MinioHost,
		config.AppConfig.MinioPort,
		config.AppConfig.MinioBucket,
		objectName,
	)

	return url, nil
}
