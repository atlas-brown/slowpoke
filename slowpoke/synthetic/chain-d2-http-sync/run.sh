#!/bin/bash

cd $(dirname $0)/../..

exp=chain-d2-http-sync
DIR=synthetic/$exp/one-service-per-node-poker
mkdir -p $DIR

# Make it reproducible
target_service_random_pairs="0:4446 1:7748 2:22717"
# target_service_random_pairs="0:4446"

for pair in $target_service_random_pairs
do 
    target_service=$(echo $pair | cut -d':' -f1)
    random_seed=$(echo $pair | cut -d':' -f2)

    output_file=$DIR/$exp-service$target_service.log
    
    if [[ -e $output_file ]]; then
        echo "File $output_file already exists. Skipping..."
        continue
    fi

    touch $output_file
    
    python3 test.py -b synthetic \
        -r $exp \
        -x service$target_service \
        --num_exp 10 \
        -c 256 \
        -t 4 \
        --num_req 20000 \
        --clien_cpu_quota 2 \
        --random_seed $random_seed \
        --repetition 3 \
        --poker_batch 40000000 \
        >$output_file
done
