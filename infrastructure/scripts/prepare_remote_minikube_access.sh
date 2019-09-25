#!/usr/bin/env bash

KUBE_DIR=$HOME/.kube

if [ ! -d $KUBE_DIR ]; then
  mkdir $KUBE_DIR
fi

# TODO make this ip dynamically change after terraform apply or get from flag (can use terraform output ips command)
scp -r root@206.189.100.31:/root/.kube/config $KUBE_DIR/config
scp -r root@206.189.100.31:/root/.minikube/client.key $KUBE_DIR/client.key
scp -r root@206.189.100.31:/root/.minikube/client.crt $KUBE_DIR/client.crt
scp -r root@206.189.100.31:/root/.minikube/ca.crt $KUBE_DIR/ca.crt

sed -i -- 's/\/root\/.minikube/\/Users\/$USER\/.kube/g' $KUBE_DIR/config
