package objectstorage

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"contestr/internal/configs"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const pdfContentType = "application/pdf"

type Client struct {
	bucket        string
	publicBaseURL string
	s3            *s3.Client
	presign       *s3.PresignClient
}

func NewClient(cfg *configs.Config) (*Client, error) {
	osc := cfg.ObjectStorage
	if !osc.Enabled() {
		return nil, fmt.Errorf("object storage is not configured")
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(osc.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			osc.AccessKeyID,
			osc.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(osc.Endpoint)
		o.UsePathStyle = true
	})

	return &Client{
		bucket:        osc.Bucket,
		publicBaseURL: strings.TrimRight(osc.PublicBaseURL, "/"),
		s3:            s3Client,
		presign:       s3.NewPresignClient(s3Client),
	}, nil
}

func (c *Client) PutObject(ctx context.Context, objectKey string, body io.Reader, size int64) error {
	_, err := c.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(c.bucket),
		Key:           aws.String(objectKey),
		Body:          body,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(pdfContentType),
	})
	return err
}

func (c *Client) DeleteObject(ctx context.Context, objectKey string) error {
	if strings.TrimSpace(objectKey) == "" {
		return nil
	}
	_, err := c.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(objectKey),
	})
	return err
}

func (c *Client) HeadObject(ctx context.Context, objectKey string) (int64, error) {
	out, err := c.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return 0, err
	}
	if out.ContentLength == nil {
		return 0, nil
	}
	return *out.ContentLength, nil
}

func (c *Client) PublicURL(objectKey string) string {
	return c.publicBaseURL + "/" + strings.TrimPrefix(objectKey, "/")
}

func (c *Client) PresignPut(ctx context.Context, objectKey string, expires time.Duration) (string, error) {
	out, err := c.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(objectKey),
		ContentType: aws.String(pdfContentType),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", err
	}
	return out.URL, nil
}

func ObjectKey(contestID int, fileID string) string {
	return fmt.Sprintf("contests/%d/problems/%s.pdf", contestID, fileID)
}
