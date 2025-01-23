package util

import (
	"context"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetS3Client(ctx context.Context, config Config) (*s3.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(config.AwsRegion),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.AwsAccessKey,
				config.AwsSecretKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("AWS設定のロードに失敗しました: %w", err)
	}

	return s3.NewFromConfig(cfg), nil
}
