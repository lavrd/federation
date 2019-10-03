#!/usr/bin/env bash

BASEDIR=$(dirname "$0")
cd "$BASEDIR" || exit

NATS_URL=http://127.0.0.1:4222 NATS_USER=root NATS_PASS=arbuz \
  ARANGODB_URL=http://127.0.0.1:8529 ARANGODB_USER=root ARANGODB_PASS=arbuz \
  go run ./app.go
