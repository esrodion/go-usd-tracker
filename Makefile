include .env

build:
	go build ./cmd/.

test:
	go test -count=1 -cover ./...

docker-build:
	docker build .

run:
	docker-compose up -d
	docker-compose run --rm --service-ports app ./app -default

lint:
	golangci-lint run --exclude-dirs=./vendor
