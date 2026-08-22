package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Handler func(ctx context.Context, msg *Message) error

type Consumer struct {
	client            *sqs.Client
	queueURL          string
	handler           Handler
	maxMessages       int32
	waitTimeSeconds   int32
	visibilityTimeout int32
	workerCount       int
}

type ConsumerConfig struct {
	QueueURL          string
	MaxMessages       int32
	VisibilityTimeout int32
	WaitTimeSeconds   int32
	WorkerCount       int
}

func NewConsumer(client *sqs.Client, cfg ConsumerConfig, handler Handler) *Consumer {
	if cfg.MaxMessages <= 0 || cfg.MaxMessages > 10 {
		cfg.MaxMessages = 10
	}

	if cfg.VisibilityTimeout <= 0 {
		cfg.VisibilityTimeout = 30
	}

	if cfg.WaitTimeSeconds <= 0 {
		cfg.WaitTimeSeconds = 20
	}

	if cfg.WorkerCount <= 0 {
		cfg.WorkerCount = 5
	}

	return &Consumer{
		client:            client,
		queueURL:          cfg.QueueURL,
		maxMessages:       cfg.MaxMessages,
		waitTimeSeconds:   cfg.WaitTimeSeconds,
		workerCount:       cfg.WorkerCount,
		visibilityTimeout: cfg.VisibilityTimeout,
		handler:           handler,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	msgChan := make(chan *Message, c.workerCount*2)

	var wg sync.WaitGroup

	for i := 0; i < c.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			c.worker(ctx, workerID, msgChan)
		}(i)
	}

	go func() {
		c.poll(ctx, msgChan)
		close(msgChan)
	}()

	wg.Wait()
	return ctx.Err()
}

func (c *Consumer) poll(ctx context.Context, msgChan chan *Message) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			messages, err := c.receiveMessages(ctx)
			if err != nil {
				log.Printf("Error receiving messages: %v", err)
				time.Sleep(3 * time.Second)
				continue
			}

			for _, msg := range messages {
				select {
				case msgChan <- msg:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

func (c *Consumer) worker(ctx context.Context, id int, msgChan <-chan *Message) {
	for msg := range msgChan {
		select {
		case <-ctx.Done():
			return
		default:
			c.processMessage(ctx, msg)
		}
	}
}

func (c *Consumer) receiveMessages(ctx context.Context) ([]*Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(c.queueURL),
		MaxNumberOfMessages:   c.maxMessages,
		VisibilityTimeout:     c.visibilityTimeout,
		WaitTimeSeconds:       c.waitTimeSeconds,
		MessageAttributeNames: []string{"All"},
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameAll,
		},
	}

	result, err := c.client.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("receiving message: %w", err)
	}

	messages := make([]*Message, 0, len(result.Messages))

	for _, sqsMessage := range result.Messages {
		if sqsMessage.Body == nil || sqsMessage.ReceiptHandle == nil {
			continue
		}

		var msg Message
		if err := json.Unmarshal([]byte(*sqsMessage.Body), &msg); err != nil {
			log.Printf("Error unmashilling message %s: %v", *sqsMessage.MessageId, err)
			continue
		}

		msg.ReceiptHandler = *sqsMessage.ReceiptHandle
		msg.Attributes = sqsMessage.Attributes
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (c *Consumer) processMessage(ctx context.Context, msg *Message) {
	timeout := time.Duration(c.visibilityTimeout-5) * time.Second
	if timeout <= 0 {
		timeout = time.Second
	}

	processCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if c.handler == nil {
		log.Printf("No handler configured for message %s", msg.ID)
		return
	}

	if err := c.handler(processCtx, msg); err != nil {
		log.Printf("Error processing message %s: %v", msg.ID, err)
		return
	}

	if err := c.deleteMessage(ctx, msg.ReceiptHandler); err != nil {
		log.Printf("Error deleting message %s: %v", msg.ID, err)
	}
}

func (c *Consumer) deleteMessage(ctx context.Context, receiptHandle string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := c.client.DeleteMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("deleting message: %w", err)
	}

	return nil
}
