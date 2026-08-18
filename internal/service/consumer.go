package service

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Consumer struct {
	client          *sqs.Client
	queueURL        string
	maxMessages     int32
	waitTimeSeconds int32
	workerCount     int
}

type ConsumerConfig struct {
	QueueURL          string
	MaxMessages       int32
	VisibilityTimeout int32
	WaitTimeSeconds   int32
	WorkerCount       int
}

func NewConsumer(client *sqs.Client, cfg ConsumerConfig) *Consumer {
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
		client:          client,
		queueURL:        cfg.QueueURL,
		maxMessages:     cfg.MaxMessages,
		waitTimeSeconds: cfg.WaitTimeSeconds,
		workerCount:     cfg.WorkerCount,
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

func (c *Consumer) poll(ctx context.Context, msgChan <-chan *Message) {
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

// TODO: Implement this function logic
func (c *Consumer) processMessage(ctx context.Context, msg *Message) {
}
