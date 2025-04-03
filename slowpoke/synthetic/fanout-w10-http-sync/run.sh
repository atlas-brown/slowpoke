#!/bin/bash

cd $(dirname $0)/../..

EXP=fanout-w10-http-sync
DIR=synthetic/$EXP/locker-correction-norlock
mkdir -p $DIR

# config
THREAD=8
CONN=512
NUM_REQ=40000
POKER_BATCH=0
NUM_EXP=2
REPETITION=1

# Make it reproducible
target_service_random_pairs="0:29513 5:30007"

for pair in $target_service_random_pairs
do 
    target_service=$(echo $pair | cut -d':' -f1)
    random_seed=$(echo $pair | cut -d':' -f2)

    output_file=$DIR/$EXP-service$target_service-t$THREAD-c$CONN-req$NUM_REQ-poker$POKER_BATCH-n$NUM_EXP-rep$REPETITION.log
    
    if [[ -e $output_file ]]; then
        echo "File $output_file already exists. Skipping..."
        continue
    fi

    touch $output_file
    
    python3 test.py -b synthetic \
        -r $EXP \
        -x service$target_service \
        --num_exp $NUM_EXP \
        -c $CONN \
        -t $THREAD \
        --num_req $NUM_REQ \
        --random_seed $random_seed \
        --repetition $REPETITION \
        --poker_batch $POKER_BATCH \
        >$output_file
done
