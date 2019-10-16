#!/usr/bin/env bash

BASEDIR=$(dirname "$0")
cd "$BASEDIR" || exit

# https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/

kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta4/aio/deploy/recommended.yaml
cd ../kube/ || exit
kubectl apply -f ./dashboard-admin-user.yml
dashboard_secret=$(kubectl -n kubernetes-dashboard describe secret "$(kubectl -n kubernetes-dashboard get secret | grep admin-user | awk '{print $1}')")
rm ./dashboard-secret || true && echo "$dashboard_secret" >./dashboard-secret
echo "access dashboard link: http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/"
