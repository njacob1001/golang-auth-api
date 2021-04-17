include .env.development

export

run:
	redis-server & cd cmd/api && go run main.go

test:
	go test ./...

dockerbuild:
	echo "Ensure run this command with sudo"
	docker-compose build

dockerup:
	echo "Ensure run this command with sudo"
	docker-compose up -d

