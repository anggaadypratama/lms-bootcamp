package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	cfg "lms-bootcamp/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)


type S3Config struct {
    cfg    *cfg.Config
	s3c *s3.Client
	bucket string
}

func NewS3Config(cfg *cfg.Config) *S3Config {
	awsEndpoint := cfg.AWS_ENDPOINT
	awsRegion := cfg.AWS_REGION
	awsBucket := cfg.AWS_BUCKET

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		log.Fatalf("Cannot load the AWS configs: %s", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(awsEndpoint)
	})

	return &S3Config{
		cfg:   cfg,
		s3c:  client,
		bucket: awsBucket,
	}
}


func (c *S3Config) EnsureBucket(ctx context.Context) error {
	_, err := c.s3c.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(c.bucket)})
	if err == nil {
		return nil 
	}
	_, err = c.s3c.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(c.bucket),
		ACL:    types.BucketCannedACLPrivate,
	})
	if err != nil {
		return fmt.Errorf("create bucket: %w", err)
	}

	waiter := s3.NewBucketExistsWaiter(c.s3c)
	if werr := waiter.Wait(ctx, &s3.HeadBucketInput{Bucket: aws.String(c.bucket)}, 30*time.Second); werr != nil {
		return fmt.Errorf("wait bucket: %w", werr)
	}
	return nil
}

func (c *S3Config) PutObject(ctx context.Context, key, content string, contentType ...string) error {
	_, err := c.s3c.PutObject(ctx, &s3.PutObjectInput{
		Bucket:     aws.String(c.bucket),
		Key:        aws.String(key),
		Body:      strings.NewReader(content),
		// ContentType: aws.String(contentType[0]),
	})
	return err
}

func (c *S3Config) GetObject(ctx context.Context, key string) (string, error) {
	out, err := c.s3c.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}
	defer out.Body.Close()

	b, err := io.ReadAll(out.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
