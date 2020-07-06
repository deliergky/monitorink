migrate:
	docker run --rm -v $(shell pwd)/flyway/sql:/flyway/sql -v $(shell pwd)/flyway/conf:/flyway/conf flyway/flyway:6.5.0-alpine migrate
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
