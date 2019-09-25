#!/usr/bin/env bash

# donwload minikube binary file
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 &&
  chmod +x minikube

# move it to /usr/local/bin and delete
install minikube /usr/local/bin
rm ./minikube

# clear previous minikube state
minikube delete

# start minikube with none driver (minikube runs by docker)
minikube start --vm-driver=none
