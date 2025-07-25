#!/bin/bash

cd $(dirname $0)

all_prefix="fanout-w7 fanout-w10-http-sync fanout-w10-http-async"

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