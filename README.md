# Monitorink

Request a URL processes the results by a pipeline:

## The Producer Pipeline

1. Produce a scheduled http request to a configured URL
2. Process the response by calculating the response time
3. Push the processed data to a provided consumer ( in memeory queue , Apache Kafka etc.)

## The Consumer Pipeline

1. Consumes pushed messages from the Producer pipeline
2. Persists them into any configured store ( in memory store, Prostgres, etc.)

## Packages & Directories

- data: describes the data model that is being produced and consumed
- flyway: coniguration and db migration scripts. flyway is used to do the migrations
- monitor: the producer(heartbeat) & the consumer(collector) pipelines are implemented together with their
  recpective consumer and producer implemenetations:
- pipeline: abstractions for creating pipelines, stages, producers and consumers (based on : https://github.com/PacktPublishing/Hands-On-Software-Engineering-with-Golang/tree/master/Chapter07/pipeline)
- queue: a queue interface and the Kafka and In Memory Queue implementations. Kafka Produce and Consumer are mainly based on the examples from (https://pkg.go.dev/github.com/Shopify/sarama?tab=doc)
- store: a persistent store interface and the Postgres and In Memory Store implementations

## Files

- docker-compose.yml: Runs zoekeeper, kafka and postgres locally.
- env.sh: environment variables needed to run the system.
- Makefile: some targets to setup and run the application

## Testing

`$ make test`

## Running

0. _Requirements:_ docker, docker-compose, >=go1.14
1. Update env.sh contents with the appropriate settings
2. `$ make run-infra`
3. `$ make migrate`
4. `$ make run-producer`
5. `$ make run-consumer`

## Todos

- [ ] Extend tests
- [ ] Add regex matcher stage to pipeline
- [ ] Improve error handling
- [ ] Containerize produce and consumer
- [ ] Add cloud provider support (Aiven, Confluent etc.) for queue and store
- [ ] Support multiple url registrations
- [ ] Create kubernetes deployment manifests
