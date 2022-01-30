# operation-fire-quasar


# TEST & DEBUG

go test ./...  -coverprofile=c.out


example via curl: 
$ curl -X POST -H "Content-Type: application/json" -d @topSecret_test1_request.json http://localhost:8080/topsecret/


# docker commands

## Builds docker image with a tag 
docker build . --tag samplemoduleuser:1.0.0


## Run app w/ redis via docker-compose


## Runs app with docker (standalone)
docker run -it --rm -p 8081:3001 --name samplemoduleuser-running samplemoduleuser:1.0.0

## starts a redis-server (standalone)
docker run --name some-redis -d redis:6.0-alpine redis-server --save 60 1 --loglevel warning

## Connect to local redis (standalone), to use redis-cli from console
docker network create redis-ntk
docker network connect redis-ntk some-redis
docker run -it --network redis-ntwk --rm redis:6.0-alpine redis-cli -h some-redis

## clean docker volumes
docker volume ls
docker volume rm VOLUME-NAME

## Starts app and redis with docker-compose
docker-compose up
docker-compose down
