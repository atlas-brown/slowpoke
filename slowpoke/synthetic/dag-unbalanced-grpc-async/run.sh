#!/bin/bash

# Change this!!
target_service_random_pairs="1:2560 4:8264 5:6318" # Make it reproducible

cd $(dirname $0)/../..

EXP=$(dirname $0 | xargs basename)
echo $EXP
DIR=synthetic/$EXP/04-09-pokerpp-rm-deadlock-0maxconn
mkdir -p $DIR

# config
THREAD=8
CONN=512
NUM_REQ=20000
POKER_BATCH_REQ=100
NUM_EXP=10
REPETITION=5

for pair in $target_service_random_pairs
do 
    target_service=$(echo $pair | cut -d':' -f1)
    random_seed=$(echo $pair | cut -d':' -f2)

    output_file=$DIR/$EXP-service$target_service-t$THREAD-c$CONN-req$NUM_REQ-poker_batch_req$POKER_BATCH_REQ-n$NUM_EXP-rep$REPETITION.log
    
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
