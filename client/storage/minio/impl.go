package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageClient struct {
	client    *minio.Client
	bucket    string
	publicURL string
}

func NewStorageClient(endpoint, accessKey, secretKey, bucket, publicURL string, useSSL bool) (*StorageClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio new client: %w", err)
	}

	return &StorageClient{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
	}, nil
}

func (c *StorageClient) EnsureBucket(ctx context.Context) error {
	exists, err := c.client.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("bucket exists check: %w", err)
	}

	if exists {
		return nil
	}

	err = c.client.MakeBucket(ctx, c.bucket, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("make bucket: %w", err)
	}

	policy := fmt.Sprintf(
		`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`,
		c.bucket,
	)

	err = c.client.SetBucketPolicy(ctx, c.bucket, policy)
	if err != nil {
		return fmt.Errorf("set bucket policy: %w", err)
	}

	return nil
}

func (c *StorageClient) Upload(ctx context.Context, key string, r io.Reader, size int64, contentType string) (string, error) {
	_, err := c.client.PutObject(ctx, c.bucket, key, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("minio put object: %w", err)
	}

	return fmt.Sprintf("%s/%s/%s", c.publicURL, c.bucket, key), nil
}

func (c *StorageClient) Delete(ctx context.Context, key string) error {
	err := c.client.RemoveObject(ctx, c.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("minio remove object: %w", err)
	}

	return nil
}
