#!/usr/bin/env bash

BASEDIR=$(dirname "$0")
cd "$BASEDIR" || exit

case "$1" in

up)
  echo "create VMs"
  terraform apply -var-file=./.env.tfvars -auto-approve
  ;;

down)
  echo "remove VMs"
  terraform destroy -var-file=./.env.tfvars -auto-approve
  ;;

refresh)
  echo "refresh VM states"
  terraform refresh -var-file=./.env.tfvars
  ;;

*)
  echo "unknown command"
  ;;

esac
