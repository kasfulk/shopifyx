package functions

import (
	"context"
	"io"
	"shopifyx/configs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var credentialProvider = func(cfg configs.Config) aws.CredentialsProviderFunc {
	return func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     cfg.S3ID,
			SecretAccessKey: cfg.S3SecretKey,
		}, nil
	}
}

func newS3Uploader(cfg configs.Config) *manager.Uploader {
	client := s3.New(s3.Options{
		Region:      "ap-southeast-1",
		Credentials: credentialProvider(cfg),
	})

	return manager.NewUploader(client)
}

type ImageUploader struct {
	uploader *manager.Uploader
}

func NewImageUploader(cfg configs.Config) *ImageUploader {
	return &ImageUploader{
		uploader: newS3Uploader(cfg),
	}
}

func (i *ImageUploader) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	result, err := i.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String("sprint-bucket-public-read"),
		Key:    aws.String(filename),
		Body:   file,
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}
