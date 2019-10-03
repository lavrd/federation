#!/usr/bin/env bash

BASEDIR=$(dirname "$0")
cd "$BASEDIR" || exit

arangodb_container_name=federation-arangodb
nats_container_name=federation-nats

docker rm -fv $arangodb_container_name $nats_container_name

docker run -d --name $arangodb_container_name -p 8529:8529 -e ARANGO_ROOT_PASSWORD=arbuz arangodb
docker run -d --name $nats_container_name -p 4222:4222 -p 8222:8222 -p 6222:6222 nats --user root --pass arbuz
