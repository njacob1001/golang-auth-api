include .env.development

export

build:
	mkdir -p bin && cd cmd/api && go build -o ../../bin

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

docker:
	echo "Ensure run this command with sudo"
	redis-server & docker run -d -p 8001:8001 rumm-api-auth:alplha
