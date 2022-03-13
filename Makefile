## docker_up_build: Build all projects and start docker compose
up:
	@echo "Starting docker images..."
	docker-compose up --build -d
	@echo "Docker started!"

## docker_down: Stop docker compose
down:
	@echo "Stopping docker images..."
	docker-compose down
	@echo "Docker stopped!"

## start: starts the front end
start:
	@echo "Starting front end"
	cd front-end && go build -o frontApp ./cmd/web
	cd front-end && ./frontApp &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./frontApp"
	@echo "Stopped front end!"

## restart_broker: rebuilds and restarts broker-service
restart_broker:
	@echo "Stopping broker service"
	docker-compose build broker-service && docker-compose up -d
	@echo "Restarted!"

## restart_auth: rebuilds and restarts authentication-service
restart_auth:
	@echo "Stopping authentication service"
	docker-compose build authentication-service && docker-compose up -d
	@echo "Restarted!"

## restart_listener: rebuilds and restarts queue-listener-service
restart_listener:
	@echo "Stopping queue listener service"
	docker-compose build queue-listener-service && docker-compose up -d
	@echo "Restarted!"

## restart_logger: rebuilds and restarts logger-service
restart_logger:
	@echo "Stopping logger service"
	docker-compose build logger-service && docker-compose up -d
	@echo "Restarted!"

## test: runs all tests
test:
	@echo "Testing..."
	go test -v ./...

## help: displays help
help: Makefile
	@echo " Choose a command:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'