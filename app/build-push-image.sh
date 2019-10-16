#!/usr/bin/env bash

BASEDIR=$(dirname "$0")
cd "$BASEDIR" || exit

tag=lluuvr/federation-app

docker build -t $tag -f ./Dockerfile .
docker push $tag
