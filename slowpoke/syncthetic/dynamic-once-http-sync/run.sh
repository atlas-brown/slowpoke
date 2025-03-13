#!/bin/bash

cd $(dirname $0)/../..

exp=dynamic-once-http-sync

mkdir -p syncthetic/$exp/results

target_services="1 2"
for target_service in $target_services
do 
    if [[ -e syncthetic/$exp/results/$exp-service$target_service.log ]]; then
        echo "File syncthetic/$exp/results/$exp-service$target_service.log already exists. Skipping..."
        continue
    fi
    touch syncthetic/$exp/results/$exp-service$target_service.log
    python3 test.py -b syncthetic \
        -r $exp \
        -x service$target_service \
        --num_exp 10 \
        -c 128 \
        -t 2 \
        --num_req 18000 \
        --clien_cpu_quota 2 \
        --random_seed $RANDOM \
        >syncthetic/$exp/results/$exp-service$target_service.log
done
