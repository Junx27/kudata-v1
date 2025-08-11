package config

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioHost *minio.Client

func InitMinio() {
	endpoint := AppConfig.MinioHost
	accessKeyID := AppConfig.MinioAccessKey
	secretAccessKey := AppConfig.MinioSecretKey
	useSSL := AppConfig.MinioUseSSL

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	MinioClient = client
	if err := ensureBucket(AppConfig.MinioBucket); err != nil {
		log.Fatalf("Failed to ensure bucket: %v", err)
	}
	if err := setPublicReadPolicy(AppConfig.MinioBucket); err != nil {
		log.Fatalf("Failed to set public policy: %v", err)
	}

	log.Printf("✅ MinIO client initialized to %s", endpoint)
}

func setPublicReadPolicy(bucket string) error {
	policy := fmt.Sprintf(`
{
  "Version": "2012-10-17",
  "Statement": [
	{
	  "Effect": "Allow",
	  "Principal": "*",
	  "Action": ["s3:GetObject"],
	  "Resource": ["arn:aws:s3:::%s/*"]
	}
  ]
}`, bucket)

	return MinioClient.SetBucketPolicy(context.Background(), bucket, policy)
}

func ensureBucket(bucket string) error {
	exists, err := MinioClient.BucketExists(context.Background(), bucket)
	if err != nil {
		return err
	}

	if !exists {
		err := MinioClient.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		log.Printf("Created new bucket: %s", bucket)
	} else {
		log.Printf("Bucket already exists: %s", bucket)
	}

	return nil
}
