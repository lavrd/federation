#!/usr/bin/env bash

BASEDIR=$(dirname "$0")
cd "$BASEDIR" || exit

# https://kubernetes.io/docs/tasks/tools/install-minikube/
# https://kubernetes.io/docs/setup/learning-environment/minikube/

curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 &&
  chmod +x minikube

install minikube /usr/local/bin
rm ./minikube

minikube delete
minikube start --vm-driver=none
