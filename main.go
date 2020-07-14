package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/deliergky/monitorink/monitor"
	"github.com/deliergky/monitorink/queue/kafka"
	"github.com/deliergky/monitorink/store"
	"github.com/kelseyhightower/envconfig"
)

func main() {

	runner := flag.String("run", "producer", "producer|consumer")
	flag.Parse()

	kafkaConfig := kafka.KafkaConfig{}
	err := envconfig.Process("kafka", &kafkaConfig)

	if err != nil {
		log.Panicf("Could not parse kafka config %v", err)
	}

	consumerConfig := sarama.NewConfig()
	consumerConfig.Version = sarama.V2_5_0_0

	producerConfig := sarama.NewConfig()
	producerConfig.Version = sarama.V2_5_0_0
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.Return.Errors = true

	kq := kafka.NewKafkaConsumer(kafkaConfig, producerConfig, consumerConfig)
	kg := kafka.NewKafkaConsumerGroup(kafkaConfig, producerConfig, consumerConfig)

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, os.Interrupt)

	if *runner == "producer" {
		err := kq.CreateProducer(ctx)
		if err != nil {
			log.Panicf("Error creating the kafka producer %v", err)
		}
		defer kq.CloseProducer(ctx)

		hbConfig := monitor.HeartbeatConfig{}
		err = envconfig.Process("heartbeat", &hbConfig)
		if err != nil {
			log.Panicf("Could not heartbeat config %v", err)
		}
		consumer := monitor.NewHeartbeatConsumer(kq)
		producer := monitor.NewHeartbeatProducer(hbConfig.Interval)
		hb := monitor.NewHeartbeat(hbConfig, producer, consumer)
		go func() {
			if err := hb.Beat(ctx); err != nil {
				log.Printf("Error processing heartbeat %v", err)
			}
		}()

	} else if *runner == "consumer" {
		postgresConfig := store.PostgresConfig{}
		err := envconfig.Process("postgres", &postgresConfig)

		if err != nil {
			log.Panicf("Could not parse db config %v", err)
		}
		store, err := store.NewPostgresStore(ctx, postgresConfig)
		if err != nil {
			log.Panicf("Error creating db connection %v", err)
		}
		defer store.Close()
		err = kq.CreateConsumer(ctx)
		if err != nil {
			log.Panicf("Error creating the kafka producer %v", err)
		}
		defer func() {
			if err := kq.CloseConsumer(ctx); err != nil {
				log.Printf("Error closing kafka consumer %v", err)
			}
		}()
		consumer := monitor.NewResponseConsumer(store)
		producer := monitor.NewResponseProducer(kq)
		collector := monitor.NewCollector(producer, consumer)
		go func() {
			if err := collector.Collect(ctx); err != nil {
				log.Printf("Error processing collector %v", err)
			}
		}()
	} else {

		postgresConfig := store.PostgresConfig{}
		err := envconfig.Process("postgres", &postgresConfig)

		if err != nil {
			log.Panicf("Could not parse db config %v", err)
		}
		store, err := store.NewPostgresStore(ctx, postgresConfig)
		if err != nil {
			log.Panicf("Error creating db connection %v", err)
		}
		defer store.Close()
		err = kg.CreateConsumer(ctx)
		if err != nil {
			log.Panicf("Error creating the kafka producer %v", err)
		}
		defer func() {
			if err := kg.CloseConsumer(ctx); err != nil {
				log.Printf("Error closing kafka consumer %v", err)
			}
		}()
		consumer := monitor.NewResponseConsumer(store)
		producer := monitor.NewResponseProducer(kg)
		collector := monitor.NewCollector(producer, consumer)
		go func() {
			if err := collector.Collect(ctx); err != nil {
				log.Printf("Error processing collector %v", err)
			}
		}()
	}

	s := <-signals
	log.Printf("Got interrupt signal, now shutting down application  %v", s)
	cancel()

	// giving enough time to process producers and consumers
	cleanupContext, cleanupFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cleanupFunc()

	<-cleanupContext.Done()
	if !errors.Is(ctx.Err(), context.Canceled) {
		log.Printf("Application could not be shutdown in a clean way. Error was %v\n", ctx.Err())
	}

}
