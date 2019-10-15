#!/usr/bin/env bash

KUBE_DIR=$HOME/.kube
VM_IP=$1

if [ ! -d $KUBE_DIR ]; then
  mkdir $KUBE_DIR
fi

ssh -o 'StrictHostKeyChecking no' root@$VM_IP cat /etc/ssh/ssh_host_dsa_key.pub >>/Users/$USER/.ssh/known_hosts

scp -r root@$VM_IP:/root/.kube/config $KUBE_DIR/config
scp -r root@$VM_IP:/root/.minikube/client.key $KUBE_DIR/client.key
scp -r root@$VM_IP:/root/.minikube/client.crt $KUBE_DIR/client.crt
scp -r root@$VM_IP:/root/.minikube/ca.crt $KUBE_DIR/ca.crt

sed -i -- 's/\/root\/.minikube/\/Users\/'$USER'\/.kube/g' $KUBE_DIR/config
