package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type store interface {
	UploadFile(ctx context.Context, client *s3.Client, bucketName string, key string, file io.Reader) error
	DeleteFile(ctx context.Context, client *s3.Client, bucketName string, key string, file io.Reader) error
	GetDownloadURL(ctx context.Context, client *s3.Client, bucketName string, fileName string) (string, error)
}
