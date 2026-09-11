package storage

import (
	"context"
	"io"
)

type Store interface {
	UploadFile(ctx context.Context, bucketName string, key string, contentType string, file io.Reader) error
	DeleteFile(ctx context.Context, bucketName string, key string) error
	GetDownloadURL(ctx context.Context, bucketName string, fileName string) (string, error)
}
