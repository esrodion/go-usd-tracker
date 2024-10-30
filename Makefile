include .env

PROJECT_PATH=$(CURDIR)

build:
	go build ./cmd/.

test:
	go test -count=1 -cover ./...

docker-build:
	docker build .

run:
	docker-compose up -d

lint:
	docker run --rm --volume="${PROJECT_PATH}:/goserver" -w /goserver golangci/golangci-lint:v1.59-alpine golangci-lint run -E gofmt --skip-dirs=./vendor --deadline=10m