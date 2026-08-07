package aws

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type QueueManager struct {
	client *sqs.Client
}

func NewQueueManager(client *sqs.Client) *QueueManager {
	return &QueueManager{client: client}
}

func (m *QueueManager) CreateStandardQueue(ctx context.Context, name string, visibilityTimeOut, messageRetention int) (string, error) {
	input := &sqs.CreateQueueInput{
		QueueName: aws.String(name),
		Attributes: map[string]string{
			// Time in seconds that a message is hidden after being received
			"VisibilityTimeout": strconv.Itoa(visibilityTimeOut),
			// Time in seconds to retain messages
			"MessageRetentionPeriod":        strconv.Itoa(messageRetention),
			"ReceiveMessageWaitTimeSeconds": "20",
		},
	}

	result, err := m.client.CreateQueue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("creating standard queue: %v", err)
	}

	return *result.QueueUrl, nil
}

func (m *QueueManager) GetQueueUrl(ctx context.Context, name string) (string, error) {
	input := &sqs.GetQueueUrlInput{
		QueueName: aws.String(name),
	}

	result, err := m.client.GetQueueUrl(ctx, input)
	if err != nil {
		return "", fmt.Errorf("getting queue URL: %w", err)
	}

	return *result.QueueUrl, nil
}
