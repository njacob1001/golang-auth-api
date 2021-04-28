# Rumm api


![Tests and build](https://github.com/njacob1001/rumm-api-alpha/actions/workflows/master.yml/badge.svg?branch=master)
![Tests and build](https://github.com/njacob1001/rumm-api-alpha/actions/workflows/develop.yml/badge.svg?branch=develop)




API for clients data management

For *local development* create  a `.env.deleopment` file en then run:

```shell
make run
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


