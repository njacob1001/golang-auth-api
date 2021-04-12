include .env

export

run:
	cd cmd/api && go run main.go

dockerbuild:
	echo "Ensure run this command with sudo"
	docker-compose build

dockerup:
	echo "Ensure run this command with sudo"
	docker-compose up -d

