include .env.development

export

build:
	cd cmd/api && go build -o ../../bin

run:
	redis-server & cd cmd/api && go run main.go

runproduction:
	cd bin && ./api

test:
	go test ./...

dockerbuild:
	echo "Ensure run this command with sudo"
	docker-compose build

dockerup:
	echo "Ensure run this command with sudo"
	docker-compose up -d

