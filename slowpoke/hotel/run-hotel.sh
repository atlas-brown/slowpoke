#!/bin/bash

cd $(dirname $0)/..

target=profile
thread=8
conn=512
repetitions=2
num_req=40000
poker_batch=48000000
num_exp=10
DIR=hotel/one-service-per-node-locker-correction-norlock
FILE=mix-$target-t$thread-c$conn-r$repetitions-req$num_req-n$num_exp-batch$poker_batch-again.log
mkdir -p $DIR
if [ -f $DIR/$FILE ]; then
    echo "File $DIR/$FILE already exists. Skipping test."
    exit 0
fi

python3 test.py -b hotel -r mix -x profile --num_exp $num_exp -t $thread -c $conn --poker_batch $poker_batch --repetition $repetitions --num_req $num_req >$DIR/$FILE