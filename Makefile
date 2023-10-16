up:
	docker compose -f docker-compose-dev.yml up -d

down:
	docker compose -f docker-compose-dev.yml down

logs:
	docker compose -f docker-compose-dev.yml logs


EXECUTABLE_NAME=weathersvc

build:
	go build -o ${EXECUTABLE_NAME} cmd/main.go

run: build
	./${EXECUTABLE_NAME}