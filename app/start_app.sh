#!/usr/bin/env bash

BASEDIR=$(dirname "$0")
cd "$BASEDIR" || exit

case "$1" in
consumer)
  NATS_URL=http://127.0.0.1:4222 NATS_USER=root NATS_PASS=arbuz \
    ARANGODB_URL=http://127.0.0.1:8529 ARANGODB_USER=root ARANGODB_PASS=arbuz \
    go run ./app.go -http :7777 -neighbor http://127.0.0.1:8888/ -consumer
  ;;
producer)
  NATS_URL=http://127.0.0.1:4222 NATS_USER=root NATS_PASS=arbuz \
    ARANGODB_URL=http://127.0.0.1:8529 ARANGODB_USER=root ARANGODB_PASS=arbuz \
    go run ./app.go -http :7777 -neighbor http://127.0.0.1:8888/ -producer
  ;;
*)
  echo "starting app error: incorrect node type"
  ;;
esac
