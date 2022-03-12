## docker_up_build: Build all projects and start docker compose
docker_up_build:
	@echo "Starting docker images..."
	docker-compose up --build -d
	@echo "Docker started!"

## docker_up: start docker images
docker_up:
	@echo "Starting docker images..."
	docker-compose up -d
	@echo "Docker started!"

## docker_down: Stop docker compose
docker_down:
	@echo "Stopping docker images..."
	docker-compose down
	@echo "Docker stopped!"

test:
	@echo "Testing..."
	go test -v ./...
