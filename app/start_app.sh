#!/usr/bin/env bash

BASEDIR=$(dirname "$0")
cd "$BASEDIR" || exit

case "$1" in
consumer)
  REDIS_URL=redis://127.0.0.1:6379 \
    ARANGODB_URL=http://127.0.0.1:8529 ARANGODB_USER=root ARANGODB_PASS=arbuz \
    go run ./app.go -consumer
  ;;
producer)
  REDIS_URL=redis://127.0.0.1:6379 \
    ARANGODB_URL=http://127.0.0.1:8529 ARANGODB_USER=root ARANGODB_PASS=arbuz \
    go run ./app.go -producer
  ;;
*)
  echo "starting app error: incorrect node type"
  ;;
esac
