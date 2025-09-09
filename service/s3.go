package service

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Service handles S3 operations
type S3Service struct {
	client     *s3.Client
	bucketName string
}

type S3Config struct {
	BucketName  string
	Region      string
	AccessKeyID string
	SecretKey   string
	EndpointURL string
}

func NewS3Service(cfg S3Config) (*S3Service, error) {
	var awsCfg aws.Config
	var err error

	configOptions := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
	}

	if cfg.AccessKeyID != "" && cfg.SecretKey != "" {
		configOptions = append(configOptions, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretKey, ""),
		))
	}

	awsCfg, err = config.LoadDefaultConfig(context.TODO(), configOptions...)
	if err != nil {
		return nil, err
	}

	s3Options := []func(*s3.Options){}

	if cfg.EndpointURL != "" {
		s3Options = append(s3Options, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.EndpointURL)
			o.UsePathStyle = true // Required for R2
		})
	}

	client := s3.NewFromConfig(awsCfg, s3Options...)

	return &S3Service{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}

func (s *S3Service) PresignGetObject(ctx context.Context, key string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func (s *S3Service) PresignMultipleObjects(ctx context.Context, keys []string) (map[string]string, error) {
	results := make(map[string]string)

	for _, key := range keys {
		url, err := s.PresignGetObject(ctx, key)
		if err != nil {
			return nil, err
		}
		results[key] = url
	}

	return results, nil
}
