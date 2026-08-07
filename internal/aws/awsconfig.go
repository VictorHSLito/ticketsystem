package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
}

func LoadFromEnv() *Config {
	return &Config{
		Region:          os.Getenv("AWS_REGION"),
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRECT_ACCESS_KEY"),
		Endpoint:        os.Getenv("AWS_ENDPOINT"),
	}
}

func NewSQSClient(ctx context.Context, cfg *Config) (*sqs.Client, error) {
	var opts []func(*config.LoadOptions) error

	opts = append(opts, config.WithRegion(cfg.Region))

	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			),
		))
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	var sqsOpts []func(*sqs.Options)

	if cfg.Endpoint != "" {
		sqsOpts = append(sqsOpts, func(o *sqs.Options) {
			o.BaseEndpoint = &cfg.Endpoint
		})
	}

	return sqs.NewFromConfig(awsCfg, sqsOpts...), nil
}
