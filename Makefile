FRONT_END_BINARY=frontApp
LOGGER_BINARY=logServiceApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LISTENER_BINARY=listener
MAIL_BINARY=mailerServiceApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## build_auth: builds the authentication binary as a linux executable
build_auth:
	cd authentication-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${AUTH_BINARY} .
	@echo "authentication-service built!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	cd logger-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/web
	@echo "authentication-service built!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	cd broker-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "authentication-service built!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	cd listener-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "authentication-service built!"

## build_mail: builds the mail binary as a linux executable
build_mail:
	cd mail-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${MAIL_BINARY} .
	@echo "authentication-service built!"


## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_auth build_broker build_listener build_logger build_mail
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## auth: stops authentication-service, removes docker image, builds service, and starts it
auth:
	@echo "Building authentication-service docker image..."
	- docker-compose stop authentication-service
	- docker-compose rm -f authentication-service
	cd authentication-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o authApp .
	docker-compose up --build -d authentication-service
	docker-compose start authentication-service
	@echo "authentication-service built and started!"

## broker: stops broker-service, removes docker image, builds service, and starts it
broker:
	@echo "Building broker-service docker image..."
	- docker-compose stop broker-service
	- docker-compose rm -f broker-service
	cd broker-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o brokerApp ./cmd/api
	docker-compose up --build -d broker-service
	docker-compose start broker-service
	@echo "broker-service rebuilt and started!"

## logger: stops logger-service, removes docker image, builds service, and starts it
logger:
	@echo "Building logger-service docker image..."
	- docker-compose stop logger-service
	- docker-compose rm -f logger-service
	cd logger-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o logServiceApp ./cmd/web
	docker-compose up --build -d logger-service
	docker-compose start logger-service
	@echo "broker-service rebuilt and started!"

## mail: stops mail-service, removes docker image, builds service, and starts it
mail:
	@echo "Building mail-service docker image..."
	- docker-compose stop mail-service
	- docker-compose rm -f mail-service
	cd mail-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o mailerServiceApp .
	docker-compose up --build -d mail-service
	docker-compose start mail-service
	@echo "mail-service rebuilt and started!"

## listener: stops listener-service, removes docker image, builds service, and starts it
listener:
	@echo "Building listener-service docker image..."
	- docker-compose stop listener-service
	- docker-compose rm -f listener-service
	cd listener-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o listener .
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
	@cd broker-service && rm -f ${BROKER_BINARY}
	@cd broker-service && go clean
	@cd listener-service && rm -f ${LISTENER_BINARY}
	@cd listener-service && go clean
	@cd authentication-service && rm -f ${AUTH_BINARY}
	@cd authentication-service && go clean
	@cd mail-service && rm -f ${MAIL_BINARY}
	@cd mail-service && go clean
	@cd logger-service && rm -f ${LOGGER_BINARY}
	@cd logger-service && go clean
	@cd front-end && go clean
	@cd front-end && rm -f ${FRONT_END_BINARY}
	@echo "Cleaned!"

## help: displays help
help: Makefile
	@echo " Choose a command:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'