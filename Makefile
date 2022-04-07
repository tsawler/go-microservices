FRONT_END_BINARY=frontApp
LOGGER_BINARY=logServiceApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LISTENER_BINARY=listener
MAIL_BINARY=mailerServiceApp
AUTH_VERSION=1.0.0
BROKER_VERSION=1.0.0
LISTENER_VERSION=1.0.2
MAIL_VERSION=1.0.0
LOGGER_VERSION=1.0.0

## up: starts all containers in the background without forcing build
up:
	@echo "Starting docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## down: stop docker compose
down:
	@echo "Stopping docker images..."
	docker-compose down
	@echo "Docker stopped!"

## build_dockerfiles: builds all dockerfile images
build_dockerfiles: build_auth build_broker build_listener build_logger build_mail front_end_linux
	@echo "Building dockerfiles..."
	docker build -f front-end.dockerfile -t tsawler/front-end .
	docker build -f authentication-service.dockerfile -t tsawler/authentication:${AUTH_VERSION} .
	docker build -f broker-service.dockerfile -t tsawler/broker:1.0.0 .
	docker build -f listener-service.dockerfile -t tsawler/listener:1.0.2 .
	docker build -f mail-service.dockerfile -t tsawler/mail:1.0.0 .
	docker build -f logger-service.dockerfile -t tsawler/logger:1.0.0 .

## push_dockerfiles: pushes tagged versions to docker hub
push_dockerfiles: build_dockerfiles
	docker push tsawler/authentication:${AUTH_VERSION}
	docker push tsawler/broker:${BROKER_VERSION}
	docker push tsawler/listener:${LISTENER_VERSION}
	docker push tsawler/mail:${MAIL_VERSION}
	docker push tsawler/logger:${LOGGER_VERSION}
	@echo "Done!"

## front_end_linux: builds linux executable for front end
front_end_linux:
	@echo "Building linux version of front end..."
	cd front-end && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o frontEndLinux ./cmd/web
	@echo "Done!"

## swarm_up: starts the swarm
swarm_up:
	@echo "Starting swarm..."
	docker stack deploy -c swarm.yml myapp

## swarm_down: stops the swarm
swarm_down:
	@echo "Stopping swarm..."
	docker stack rm myapp

## build_auth: builds the authentication binary as a linux executable
build_auth:
	@echo "Building authentication binary.."
	cd authentication-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Authentication binary built!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd logger-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/web
	@echo "Logger binary built!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd broker-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Broker binary built!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd listener-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Listener binary built!"

## build_mail: builds the mail binary as a linux executable
build_mail:
	@echo "Building mailer binary..."
	cd mail-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Mailer binary built!"


## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_auth build_broker build_listener build_logger build_mail
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## auth: stops authentication-service, removes docker image, builds service, and starts it
auth: build_auth
	@echo "Building authentication-service docker image..."
	- docker-compose stop authentication-service
	- docker-compose rm -f authentication-service
	docker-compose up --build -d authentication-service
	docker-compose start authentication-service
	@echo "authentication-service built and started!"

## broker: stops broker-service, removes docker image, builds service, and starts it
broker: build_broker
	@echo "Building broker-service docker image..."
	- docker-compose stop broker-service
	- docker-compose rm -f broker-service
	docker-compose up --build -d broker-service
	docker-compose start broker-service
	@echo "broker-service rebuilt and started!"

## logger: stops logger-service, removes docker image, builds service, and starts it
logger: build_logger
	@echo "Building logger-service docker image..."
	- docker-compose stop logger-service
	- docker-compose rm -f logger-service
	docker-compose up --build -d logger-service
	docker-compose start logger-service
	@echo "broker-service rebuilt and started!"

## mail: stops mail-service, removes docker image, builds service, and starts it
mail: build_mail
	@echo "Building mail-service docker image..."
	- docker-compose stop mail-service
	- docker-compose rm -f mail-service
	docker-compose up --build -d mail-service
	docker-compose start mail-service
	@echo "mail-service rebuilt and started!"

## listener: stops listener-service, removes docker image, builds service, and starts it
listener: build_listener
	@echo "Building listener-service docker image..."
	- docker-compose stop listener-service
	- docker-compose rm -f listener-service
	docker-compose up --build -d listener-service
	docker-compose start listener-service
	@echo "listener-service rebuilt and started!"

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