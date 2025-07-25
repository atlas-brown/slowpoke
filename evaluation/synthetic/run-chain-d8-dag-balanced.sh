#!/bin/bash

cd $(dirname $0)

all_prefix="chain-d8 dag-balanced"

for prefix in $all_prefix
do
    for folder in ./$prefix*
    do
    kubectl get deploy -o name | xargs kubectl delete
    kubectl get svc -o name | xargs kubectl delete 
    bash $folder/run.sh
    # clean up all deployment and services
    done
done