#!/usr/bin/env bash

# https://kubernetes.io/docs/tasks/tools/install-minikube/
# https://kubernetes.io/docs/setup/learning-environment/minikube/
# https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/

curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 &&
  chmod +x minikube

install minikube /usr/local/bin
rm ./minikube

minikube delete
minikube start --vm-driver=none

kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta4/aio/deploy/recommended.yaml
cd ../../kube/ && kubectl apply -f ./dashboard_admin_user.yml
dashboard_secret=$(kubectl -n kubernetes-dashboard describe secret "$(kubectl -n kubernetes-dashboard get secret | grep admin-user | awk '{print $1}')")
echo "$dashboard_secret" >./dashboard-secret
