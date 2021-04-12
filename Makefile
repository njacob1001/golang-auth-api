set_develop_env:
	export RUMM_HOST=localhost
	export RUMM_PORT=8080
	export RUMM_SHUTDOWNTIMEOUT=10s
	export RUMM_DBUSER=admin
	export RUMM_DBPASS=admin
	export RUMM_DBHOST=localhost
	export RUMM_DBPORT=5432
	export RUMM_DBNAME=rummdb
	export RUMM_DBSHUTDOWNTIMEOUT=5s
run: set_develop_env
	cd cmd/api && go run main.go

dockerbuild:
	echo "Ensure run this command with sudo"
	docker-compose build

dockerup:
	echo "Ensure run this command with sudo"
	docker-compose up -d

