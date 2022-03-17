FRONT_END_BINARY=frontApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## auth: stops authentication-service, removes docker image, builds service, and starts it
auth:
	@echo "Building authentication-service docker image..."
	docker-compose stop authentication-service
	docker-compose rm -f authentication-service
	docker-compose up --build -d authentication-service
	docker-compose start authentication-service
	@echo "authentication-service built and started!"

## broker: stops broker-service, removes docker image, builds service, and starts it
broker:
	@echo "Building broker-service docker image..."
	docker-compose stop broker-service
	docker-compose rm -f broker-service
	docker-compose up --build -d broker-service
	docker-compose start broker-service
	@echo "broker-service rebuilt and started!"

## logger: stops logger-service, removes docker image, builds service, and starts it
logger:
	@echo "Building logger-service docker image..."
	docker-compose stop logger-service
	docker-compose rm -f logger-service
	docker-compose up --build -d logger-service
	docker-compose start logger-service
	@echo "broker-service rebuilt and started!"

## mail: stops mail-service, removes docker image, builds service, and starts it
mail:
	@echo "Building mail-service docker image..."
	docker-compose stop mail-service
	docker-compose rm -f mail-service
	docker-compose up --build -d mail-service
	docker-compose start mail-service
	@echo "mail-service rebuilt and started!"

## listener: stops listener-service, removes docker image, builds service, and starts it
listener:
	@echo "Building listener-service docker image..."
	docker-compose stop listener-service
	docker-compose rm -f listener-service
	docker-compose up --build -d listener-service
	docker-compose start listener-service
	@echo "listener-service rebuilt and started!"

## down: stop docker compose
down:
	@echo "Stopping docker images..."
	docker-compose down
	@echo "Docker stopped!"

## start: starts the front end
start:
	@echo "Starting front end"
	cd front-end && go build -o ${FRONT_END_BINARY} ./cmd/web
	cd front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

## test: runs all tests
test:
	@echo "Testing..."
	go test -v ./...

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@cd authentication-service && go clean
	@cd front-end && go clean
	@cd front-end && rm -f ${FRONT_END_BINARY}
	@cd logger-service && go clean
	@cd listener-service && go clean
	@echo "Cleaned!"

## help: displays help
help: Makefile
	@echo " Choose a command:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'