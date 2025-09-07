package s3

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const (
	// S3BucketName is the name of the s3 bucket
	MediaBucketName = "remont_ai_media_storage"
	projectID       = "remontai-app"
)

// S3Repository is a repository for s3
type S3Repository interface {
	UploadFile(ctx context.Context, file []byte, filename string) error
	DeleteFile(ctx context.Context, filename string) error
	IsFileExist(ctx context.Context, filename string) (bool, error)
}

// NewRepository creates a new s3 repository
func NewRepository(ctx context.Context, s3Credentials string) (S3Repository, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(s3Credentials)))
	if err != nil {
		return nil, err
	}

	return &repository{
		client: client,
	}, nil
}

type repository struct {
	client *storage.Client
}

// UploadFile uploads a file to s3
func (r *repository) UploadFile(ctx context.Context, file []byte, filename string) error {
	bucket := r.client.Bucket(MediaBucketName)
	obj := bucket.Object(filename)
	w := obj.NewWriter(ctx)
	w.Write(file)
	w.Close()

	return nil
}

func (r *repository) DeleteFile(ctx context.Context, filename string) error {
	bucket := r.client.Bucket(MediaBucketName)
	obj := bucket.Object(filename)
	return obj.Delete(ctx)
}

func (r *repository) IsFileExist(ctx context.Context, filename string) (bool, error) {
	bucket := r.client.Bucket(MediaBucketName)
	obj := bucket.Object(filename)
	_, err := obj.Attrs(ctx)
	if err != nil {
		return false, nil
	}
	return true, nil
}
