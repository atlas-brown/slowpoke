#!/bin/bash

cd $(dirname $0)

all_prefix="fanout-w10 fanout-w7"

for prefix in $all_prefix
do
    for folder in ./$prefix*
    do
    kubectl get deploy -o name | xargs kubectl delete
    kubectl get svc -o name | grep service/service | xargs kubectl delete 
    bash $folder/run.sh
    # clean up all deployment and services
    done
done