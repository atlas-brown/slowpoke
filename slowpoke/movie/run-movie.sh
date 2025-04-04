#!/bin/bash

cd $(dirname $0)/..

target=moviereviews
thread=8
conn=800
repetitions=5
num_req=50000
poker_batch=60000000
num_exp=10
DIR=movie/locker-correction-norlock-num-based-batching
FILE=mix-$target-t$thread-c$conn-r$repetitions-req$num_req-n$num_exp-batch$poker_batch-0reqbatch.log
mkdir -p $DIR
if [ -f $DIR/$FILE ]; then
    echo "File $DIR/$FILE already exists. Skipping test."
    exit 0
fi

python3 test.py -b movie -r mix -x moviereviews --num_exp $num_exp -t $thread -c $conn --poker_batch $poker_batch --repetition $repetitions --num_req $num_req >$DIR/$FILE