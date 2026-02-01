package gojinn

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (r *Gojinn) getS3Client(ctx context.Context) (*s3.Client, error) {
	if r.S3Bucket == "" {
		return nil, fmt.Errorf("s3_bucket not configured")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(r.S3Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			r.S3AccessKey,
			r.S3SecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	if r.S3Endpoint != "" {
		cfg.BaseEndpoint = aws.String(r.S3Endpoint)
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	}), nil
}

func (r *Gojinn) s3Put(ctx context.Context, key string, data []byte) error {
	client, err := r.getS3Client(ctx)
	if err != nil {
		return err
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.S3Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	return err
}

func (r *Gojinn) s3Get(ctx context.Context, key string) ([]byte, error) {
	client, err := r.getS3Client(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
