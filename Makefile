lint:
	golangci-lint --deadline=5m run --enable=goimports --enable=golint --color=always --out-format=tab  ./...
migrate:
	docker run --rm -v $(shell pwd)/flyway/sql:/flyway/sql flyway/flyway:6.5.0-alpine migrate -url=$(FLYWAY_URL) -user=$(FLYWAY_USER) -password=$(FLYWAY_PASSWORD)
test:
	go test ./...
run-consumer:
	source env.sh && go run main.go -run=consumer
run-producer:
	source env.sh && go run main.go
run-infra:
	source env.sh && docker-compose up
stop-infra:
	docker-compose rm -vfs
build-producer:
	docker build -f Dockerfile.producer -t monitorink-producer:0.1.0 .
build-consumer:
	docker build -f Dockerfile.consumer -t monitorink-consumer:0.1.0 .
