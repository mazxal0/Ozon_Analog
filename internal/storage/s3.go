package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	Client   *minio.Client
	Bucket   string
	Endpoint string
}

func NewMinioStorage() (*MinioStorage, error) {
	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "minio123", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	storage := &MinioStorage{
		Client:   client,
		Bucket:   "products",
		Endpoint: "http://localhost:9000",
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, storage.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		client.MakeBucket(ctx, storage.Bucket, minio.MakeBucketOptions{})
	}

	return storage, nil
}

func (s *MinioStorage) Upload(ctx context.Context, objectName, filePath string) (string, error) {
	_, err := s.Client.FPutObject(ctx, s.Bucket, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	// Публичный URL
	return s.Endpoint + "/" + s.Bucket + "/" + objectName, nil
}
