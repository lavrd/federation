#!/usr/bin/env bash

BASEDIR=$(dirname "$0")
cd "$BASEDIR" || exit

federation_network_name=federation
arangodb_container_name=federation-arangodb
redis_container_name=federation-redis

docker rm -fv $arangodb_container_name $redis_container_name
docker network rm $federation_network_name

docker network create $federation_network_name
docker run -d --name $arangodb_container_name --net $federation_network_name -p 8529:8529 -e ARANGO_ROOT_PASSWORD=arbuz arangodb:3.5.1
docker run -d --name $redis_container_name --net $federation_network_name -p 6379:6379 redis:5.0.6-alpine
