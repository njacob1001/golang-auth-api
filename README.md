# Rumm api authentication microservice

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/njacob1001/go-hexagonal-baseline)](https://github.com/njacob1001/rumm-api-alpha)
![Tests and build master](https://github.com/njacob1001/rumm-api-auth/actions/workflows/master.yml/badge.svg)
![Tests and build develop](https://github.com/njacob1001/rumm-api-auth/actions/workflows/develop.yml/badge.svg)




API for clients data management

For *local development* create  a `.env.delopment` file en then run:

```shell
make run
```

### Test project

```shell
make test
```

for deploy in docker container create a `.env` file en then run:

```shell
sudo docker-compose build
sudo docker-compose up -d
```

shutdown docker servers run:

```shell
sudo docker-compose down
```

### Monitoring containers

This project uses portainer for container monitoring, you can go to localhost:8000


