#!/bin/bash

# Change this!!
target_service_random_pairs="0:24737 3:1492 5:19356" # Make it reproducible

EXP="$(basename "$(dirname "$(realpath "$0")")")"

cd $(dirname $0)/../../..
DIR=./evaluation/results
mkdir -p $DIR
# config
THREAD=8
CONN=512
NUM_REQ=20000
POKER_BATCH_REQ=100
NUM_EXP=5
REPETITION=1

for pair in $target_service_random_pairs
do 
    kubectl delete deployments --all
    kubectl delete services --all
    target_service=$(echo $pair | cut -d':' -f1)
    random_seed=$(echo $pair | cut -d':' -f2)

    output_file=$DIR/$EXP-service$target_service-t$THREAD-c$CONN-req$NUM_REQ-poker_batch_req$POKER_BATCH_REQ-n$NUM_EXP-rep$REPETITION.log
    
    if [[ -e $output_file ]]; then
        echo "File $output_file already exists. Skipping..."
        continue
    fi

    touch $output_file
    
    python3 src/main.py -b synthetic \
        -r $EXP \
        -x service$target_service \
        --num_exp $NUM_EXP \
        -c $CONN \
        -t $THREAD \
        --num_req $NUM_REQ \
        --random_seed $random_seed \
        --repetition $REPETITION \
        --poker_batch_req $POKER_BATCH_REQ \
        >$output_file
    kubectl delete deployments --all
    kubectl delete services --all
done
