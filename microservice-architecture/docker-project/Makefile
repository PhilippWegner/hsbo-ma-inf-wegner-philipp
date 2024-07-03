GRAPHQL_MACHINE_SERVICE_BINARY=graphql-machine-service-app
GRPC_LOGGER_SERVICE_BINARY=graphql-machine-service-app
RESTFUL_BROKER_SERVICE_BINARY=restful-broker-service-app
STATEMACHINE_SERVICE_BINARY=statemachine-service-app

up:
	docker compose up -d

down:
	docker compose down

up_build: down
	docker compose up --build -d

build_local: build_graphql_machine_service build_grpc_logger_service build_restful_broker_service build_statemaschine

build_graphql_machine_service:
	cd ../graphql-machine-service && env GOOS=linux CGO_ENABLED=0 go build -o ${GRAPHQL_MACHINE_SERVICE_BINARY} ./cmd/api

build_grpc_logger_service:
	cd ../grpc-logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${GRPC_LOGGER_SERVICE_BINARY} ./cmd/api

build_restful_broker_service:
	cd ../restful-broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${RESTFUL_BROKER_SERVICE_BINARY} ./cmd/api

build_statemaschine:
	cd ../statemachine-service && env GOOS=linux CGO_ENABLED=0 go build -o ${STATEMACHINE_SERVICE_BINARY} ./cmd/service
