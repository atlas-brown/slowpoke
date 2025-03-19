#!/bin/bash

cd $(dirname $0)/../..

exp=chain-d2-grpc-async

mkdir -p synthetic/$exp/results

target_services="0 1 2"
for target_service in $target_services
do 
    if [[ -e synthetic/$exp/results/$exp-service$target_service.log ]]; then
        echo "File synthetic/$exp/results/$exp-service$target_service.log already exists. Skipping..."
        continue
    fi
    touch synthetic/$exp/results/$exp-service$target_service.log
    python3 test.py -b synthetic \
        -r $exp \
        -x service$target_service \
        --num_exp 10 \
        -c 128 \
        -t 2 \
        --num_req 18000 \
        --clien_cpu_quota 2 \
        --random_seed $RANDOM \
        >synthetic/$exp/results/$exp-service$target_service.log
done
