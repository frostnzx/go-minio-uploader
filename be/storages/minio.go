package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

func ConnectMinio() (*minio.Client, error) {
	ctx := context.Background()
	host := os.Getenv("MINIO_HOST")
	port := os.Getenv("MINIO_PORT")
	accessKeyID := os.Getenv("MINIO_ACCESSKEY")
	secretAccessKey := os.Getenv("MINIO_SECRETKEY")
	useSSL := false

	endpoint := fmt.Sprintf("%s:%s", host, port)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// create image-collection bucket
	bucketName := os.Getenv("MINIO_BUCKET")
	location := "us-east-1"

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Successfully created bucket: %s\n", bucketName)
	}
	// create csv bucket
	bucketCsvName := os.Getenv("MINIO_BUCKET_CSV")
	location = "us-east-1"

	exists, err = minioClient.BucketExists(ctx, bucketCsvName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketCsvName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Successfully created bucket: %s\n", bucketCsvName)
	}
	// No need to log anything if the bucket already exists
	return minioClient, nil
}
