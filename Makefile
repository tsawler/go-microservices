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

start:
	@echo "Starting front end"
	cd front-end && go build -o frontApp ./cmd/web
	cd front-end && ./frontApp &

stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./frontApp"
	@echo "Stopped front end!"

restart_broker:
	@echo "Stopping broker service"
	docker-compose build broker-service && docker-compose up -d
	@echo "Restarted!"

restart_auth:
	@echo "Stopping authentication service"
	docker-compose build authentication-service && docker-compose up -d
	@echo "Restarted!"

test:
	@echo "Testing..."
	go test -v ./...
