package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	localpkg "github.com/cthulhu-platform/filemanager/internal/pkg"
)

type AWSStorage struct {
	Client       *s3.Client
	PresignClient *s3.PresignClient
	BucketName   string
}

type AWSStorageConfig struct {
	AccessKeyID       string
	SecretAccessKey   string
	Endpoint          string
	PresignedEndpoint string // optional; if set, presigned URLs use this (e.g. localhost for browser); server-side uses Endpoint
	Region            string
	BucketName        string
	ForcePathStyle    bool
}

func NewAWSStorage(ctx context.Context, cfg AWSStorageConfig) (*AWSStorage, error) {
	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	creds := credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")

	makeClient := func(endpoint string) (*s3.Client, error) {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if endpoint != "" {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               endpoint,
					HostnameImmutable: true,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})
		newCfg, err := config.LoadDefaultConfig(
			ctx,
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(creds),
			config.WithEndpointResolverWithOptions(customResolver),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS config: %v", err)
		}
		return s3.NewFromConfig(newCfg, func(o *s3.Options) {
			o.UsePathStyle = cfg.ForcePathStyle
		}), nil
	}

	client, err := makeClient(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	presignBase := client
	if cfg.PresignedEndpoint != "" && cfg.PresignedEndpoint != cfg.Endpoint {
		presignBase, err = makeClient(cfg.PresignedEndpoint)
		if err != nil {
			return nil, err
		}
	}
	presignClient := s3.NewPresignClient(presignBase, func(o *s3.PresignOptions) {
		o.Expires = localpkg.PRESIGNED_URL_EXPIRATION
	})

	log.Println("Successfully connected to AWS S3")

	return &AWSStorage{Client: client, PresignClient: presignClient, BucketName: cfg.BucketName}, nil
}

func (s *AWSStorage) Close() error {
	return nil
}

func (s *AWSStorage) PresignPut(ctx context.Context, key string, contentLength int64, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.BucketName),
		Key:           aws.String(key),
		ContentLength: aws.Int64(contentLength),
		ContentType:   aws.String(contentType),
	}
	req, err := s.PresignClient.PresignPutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("presign put object: %w", err)
	}
	return req.URL, nil
}

func (s *AWSStorage) PresignGet(ctx context.Context, key string) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	}
	req, err := s.PresignClient.PresignGetObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("presign get object: %w", err)
	}
	return req.URL, nil
}

func (s *AWSStorage) DeleteObject(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete object %q: %w", key, err)
	}
	return nil
}
