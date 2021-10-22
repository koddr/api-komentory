.PHONY: clean lint security critic test build run

APP_NAME = api
BUILD_DIR = $(PWD)/build

clean:
	rm -rf ./build

lint:
	golangci-lint run ./... --timeout 2m

security:
	gosec -exclude=G307 ./...

critic:
	gocritic check -enableAll ./...

test: lint security critic
	go test -cover ./...

build: clean test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: build
	$(BUILD_DIR)/$(APP_NAME)

docker.run: docker.network docker.postgres docker.redis

docker.network:
	docker network inspect dev-network >/dev/null 2>&1 || \
	docker network create -d bridge dev-network

docker.postgres:
	docker run --rm -d \
		--name dev-postgres \
		--network dev-network \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=postgres \
		-v ${HOME}/dev-postgres/data/:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres:13

docker.redis:
	docker run --rm -d \
		--name dev-redis \
		--network dev-network \
		-p 6379:6379 \
		redis

docker.stop: docker.stop.postgres docker.stop.redis

docker.stop.postgres:
	docker stop dev-postgres

docker.stop.redis:
	docker stop dev-redis
