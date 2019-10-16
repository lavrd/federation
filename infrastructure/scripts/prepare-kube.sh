#!/usr/bin/env bash

BASEDIR=$(dirname "$0")
cd "$BASEDIR" || exit

bash ./prepare-remote-kube-access.sh "$@"
bash ./prepare-kube-dashboard.sh
