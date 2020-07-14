package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/deliergky/monitorink/data"
)

// KafkaConfig is passed to a KafkaQueue to initialize the producer and the consumer
type KafkaConfig struct {
	Brokers string
	Topic   string
}

// KafkaConsumer implements the Queue interface with Kafka specific methods
type KafkaConsumer struct {
	config         KafkaConfig
	consumer       sarama.PartitionConsumer
	producer       sarama.AsyncProducer
	current        data.ResponseData
	producerConfig *sarama.Config
	consumerConfig *sarama.Config
}

// NewKafkaConsumer creates a new kafka queue with the provided parameters
func NewKafkaConsumer(config KafkaConfig, producerConfig *sarama.Config, consumerConfig *sarama.Config) *KafkaConsumer {
	return &KafkaConsumer{
		config:         config,
		producerConfig: producerConfig,
		consumerConfig: consumerConfig,
	}
}

// CreateConsumer initializes a new Kafka PartitionConsumer
func (q *KafkaConsumer) CreateConsumer(ctx context.Context) error {
	master, err := sarama.NewConsumer(strings.Split(q.config.Brokers, ","), q.consumerConfig)
	if err != nil {
		return err
	}
	consumer, err := master.ConsumePartition(q.config.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	q.consumer = consumer

	return nil
}

// CreateProducder initializes a new Kafka AsyncProducer
func (q *KafkaConsumer) CreateProducer(ctx context.Context) error {

	producer, err :=
		sarama.NewAsyncProducer(strings.Split(q.config.Brokers, ","), q.producerConfig)
	if err != nil {
		return err
	}
	q.producer = producer
	return nil
}

// CloseConsumer closes the consumer
func (q *KafkaConsumer) CloseConsumer(ctx context.Context) error {
	return q.consumer.Close()
}

// CloseProducer closes the producer
func (q *KafkaConsumer) CloseProducer(ctx context.Context) {
	q.producer.AsyncClose()
}

func (q *KafkaConsumer) Push(d data.ResponseData) error {

	var err error
	var wg sync.WaitGroup

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	wg.Add(2)
	go func() {
		q.producer.Input() <- &sarama.ProducerMessage{Topic: q.config.Topic, Key: nil, Value: sarama.ByteEncoder(b)}
		wg.Done()
	}()

	go func() {
		select {
		case result := <-q.producer.Successes():
			log.Printf("> message: \"%s\" sent to partition  %d at offset %d\n", result.Value, result.Partition, result.Offset)
		case err = <-q.producer.Errors():
			log.Println("Failed to produce message", err)
		}
		wg.Done()
	}()
	wg.Wait()
	return err
}

func (q *KafkaConsumer) Pop() (data.ResponseData, error) {
	return q.Current(), nil
}

func (q *KafkaConsumer) Current() data.ResponseData {
	return q.current
}

func (q *KafkaConsumer) Next(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	case message := <-q.consumer.Messages():
		d := data.ResponseData{}
		err := json.Unmarshal(message.Value, &d)
		if err != nil {
			log.Printf("Error unmarshalling data %v\n", err)
		}
		q.current = d
		return true
	case err := <-q.consumer.Errors():
		log.Printf("Consumed error %v", err)
		return false
	case <-time.After(60 * time.Second):
		log.Printf("Did not receive any messages\n")
		return false
	}
}
