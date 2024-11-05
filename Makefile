include .env

build:
	go build ./cmd/.

test:
	go test -count=1 -cover ./...

docker-build:
	docker build .

run:
	docker-compose up -d

lint:
	golangci-lint run -E gofmt --skip-dirs=./vendor --deadline=10m
