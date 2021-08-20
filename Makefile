include .env.development

export

build:
	cd cmd/api && go build main.go

run:
	redis-server & cd cmd/api && go run main.go

runproduction:
	cd cmd/api && ./main

test:
	go test ./...

dockerbuild:
	echo "Ensure run this command with sudo"
	docker-compose build

dockerup:
	echo "Ensure run this command with sudo"
	docker-compose up -d

