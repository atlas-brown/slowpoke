#!/bin/bash

cd $(dirname $0)/../..

EXP=dag-cross-http-sync
DIR=synthetic/$EXP/bus-theory-fixed-json-time-based-sleep
mkdir -p $DIR

# config
THREAD=8
CONN=512
NUM_REQ=20000
# POKER_BATCH=20000000
POKER_BATCH_REQ=100
NUM_EXP=10
REPETITION=1

# Make it reproducible
# target_service_random_pairs="2:12895 3:19434 4:1348"
target_service_random_pairs="4:1348"
# 4:1348 has giant processing time

for pair in $target_service_random_pairs
do 
    target_service=$(echo $pair | cut -d':' -f1)
    random_seed=$(echo $pair | cut -d':' -f2)

    output_file=$DIR/$EXP-service$target_service-t$THREAD-c$CONN-req$NUM_REQ-poker_batch_req$POKER_BATCH_REQ-n$NUM_EXP-rep$REPETITION-relative100.log
    
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
        --poker_batch_req $POKER_BATCH_REQ \
        >$output_file
done

