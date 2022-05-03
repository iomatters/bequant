# bequant

## Description
- **collector** - daemon running in the background to collect and record prices from the cryptocompare API
- **endpoint** - webSocket API endpoint to interface prices pulled from the cryptocompare API

PostgreSQL is being used as a backend database.

## Prerequisites
- docker
- docker-compose
- make

## Customization
You may want to use your own PostgreSQL setup, if that is the case, modify pg_dump.sql to use alternative database user and password, also modify both collector/config.toml and endpoint/config.toml respectivlly to reflect the changes.

Otherwise, use the defaults.

## Installation
1. Build collector image
```
make -C collector/ docker.build
```
2. Build endpoint image
```
make -C endpoint/ docker.build
```
3. Get containers up & running
```
docker-compose up -d
```
4. Respore PostgreSQL dump
```
cat pg_dump.sql | docker exec -i postgres psql -Upostgres
```
5. Make sure containers are up & running
```
docker ps && docker logs collector && docker logs endpoint
```
## How To
Identify the server running your containers and specify it within wsclient.html. For that, locate the line `var socket = new WebSocket("ws://localhost:8080/price` in wsclient.html and modify it accordingly using your server name, e.g. `var socket = new WebSocket("ws://MYSERVER:8080/price`.

Open wsclient.html in a browser and give it a try. 
