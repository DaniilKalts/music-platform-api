package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/DaniilKalts/music-platform-api/internal/config"
)

type Client struct {
	minio     *minio.Client
	bucket    string
	publicURL string
}

func NewClient(ctx context.Context, cfg *config.S3) (*Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("check bucket: %w", err)
	}

	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("make bucket: %w", err)
		}

		policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, cfg.Bucket)
		if err := client.SetBucketPolicy(ctx, cfg.Bucket, policy); err != nil {
			return nil, fmt.Errorf("set bucket policy: %w", err)
		}
	}

	return &Client{
		minio:     client,
		bucket:    cfg.Bucket,
		publicURL: strings.TrimRight(cfg.PublicURL, "/"),
	}, nil
}

func (c *Client) Upload(ctx context.Context, filename string, reader io.Reader, size int64, contentType string) (string, error) {
	_, err := c.minio.PutObject(ctx, c.bucket, filename, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("put object: %w", err)
	}

	return fmt.Sprintf("%s/%s/%s", c.publicURL, c.bucket, filename), nil
}

func (c *Client) Delete(ctx context.Context, fileURL string) error {
	u, err := url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	key := strings.TrimPrefix(u.Path, "/")
	key = strings.TrimPrefix(key, c.bucket+"/")

	if err := c.minio.RemoveObject(ctx, c.bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("remove object: %w", err)
	}

	return nil
}
