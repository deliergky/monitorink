package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/deliergky/monitorink/data"
)

// KafkaConsumerGroup implements the Queue interface with Kafka specific methods
type KafkaConsumerGroup struct {
	config         KafkaConfig
	group          sarama.ConsumerGroup
	producer       sarama.AsyncProducer
	handler        msgConsumerGroupHandler
	producerConfig *sarama.Config
	consumerConfig *sarama.Config
}

// NewKafkaConsumerGroup creates a new kafka queue with the provided parameters
func NewKafkaConsumerGroup(config KafkaConfig, producerConfig *sarama.Config, consumerConfig *sarama.Config) *KafkaConsumerGroup {
	return &KafkaConsumerGroup{
		config:         config,
		producerConfig: producerConfig,
		consumerConfig: consumerConfig,
	}
}

// CreateConsumer initializes a new Kafka PartitionConsumer
func (q *KafkaConsumerGroup) CreateConsumer(ctx context.Context) error {
	cg, err := sarama.NewConsumerGroup(strings.Split(q.config.Brokers, ","), q.config.Topic, q.consumerConfig)
	if err != nil {
		return err
	}
	q.group = cg
	q.handler = msgConsumerGroupHandler{
		current: make(chan *sarama.ConsumerMessage),
	}
	return nil
}

// CreateProducder initializes a new Kafka AsyncProducer
func (q *KafkaConsumerGroup) CreateProducer(ctx context.Context) error {

	producer, err :=
		sarama.NewAsyncProducer(strings.Split(q.config.Brokers, ","), q.producerConfig)
	if err != nil {
		return err
	}
	q.producer = producer
	return nil
}

// CloseConsumer closes the consumer
func (q *KafkaConsumerGroup) CloseConsumer(ctx context.Context) error {
	return q.group.Close()
}

// CloseProducer closes the producer
func (q *KafkaConsumerGroup) CloseProducer(ctx context.Context) {
	q.producer.AsyncClose()
}

func (q *KafkaConsumerGroup) Push(d data.ResponseData) error {

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

func (q *KafkaConsumerGroup) Pop() (data.ResponseData, error) {
	return q.Current(), nil
}

func (q *KafkaConsumerGroup) Current() data.ResponseData {
	message := <-q.handler.current

	d := data.ResponseData{}
	err := json.Unmarshal(message.Value, &d)
	if err != nil {
		log.Printf("Error unmarshalling data %v\n", err)
	}
	log.Printf("consumes msg%v\n", d)
	return d
}

func (q *KafkaConsumerGroup) Next(ctx context.Context) bool {
	topics := []string{q.config.Topic}
	// `Consume` should be called inside an infinite loop, when a
	// server-side rebalance happens, the consumer session will need to be
	// recreated to get the new claims
	select {
	case <-ctx.Done():
		return false
	default:
	}

	go func() {
		err := q.group.Consume(ctx, topics, q.handler)
		if err != nil {
			log.Printf("Error consuming from group %v\n", err)
		}
	}()
	return true
}
