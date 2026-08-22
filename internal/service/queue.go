package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type Message struct {
	ID             string            `json:"id"`
	Type           string            `json:"type"`
	Payload        map[string]any    `json:"payload"`
	Timestamp      time.Time         `json:"timestamp"`
	ReceiptHandler string            `json:"-"`
	Attributes     map[string]string `json:"-"`
}

type QueueService struct {
	SqsClient *sqs.Client
	QueueURL  string
}

func NewQueueService(client *sqs.Client, queueURL string) *QueueService {
	return &QueueService{client, queueURL}
}

func (s *QueueService) SendMessage(
	ctx context.Context,
	msg *Message,
	delaySeconds int32,
) (string, error) {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}

	msg.Timestamp = time.Now().UTC()

	body, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("marshaling message: %w", err)
	}

	input := &sqs.SendMessageInput{
		QueueUrl:     aws.String(s.QueueURL),
		MessageBody:  aws.String(string(body)),
		DelaySeconds: delaySeconds,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"MessageType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(msg.Type),
			},
			"CorrelationId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(msg.ID),
			},
		},
	}

	result, err := s.SqsClient.SendMessage(ctx, input)
	if err != nil {
		return "", fmt.Errorf("sending message: %w", err)
	}

	return *result.MessageId, nil
}
