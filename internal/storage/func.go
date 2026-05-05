package storage

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Store struct {
	Client *s3.Client
}

func (r *R2Store) UploadFile(ctx context.Context, bucketName string, key string,contentType string, file io.Reader) error {
	_, err := r.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
		ContentType: aws.String(contentType),

	})
	return err
}

func DeleteFile(ctx context.Context, client *s3.Client, bucketName string, key string) error {
	_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	return err
}

func (r *R2Store) GetDownloadURL(ctx context.Context, bucketName, fileName string) (string, error) {
	presignClient := s3.NewPresignClient(r.Client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	}, s3.WithPresignExpires(time.Hour*1)) // URL valid for 1 hour

	if err != nil {
		return "", err
	}

	return request.URL, nil
}
