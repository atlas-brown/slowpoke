#!/bin/bash

cd $(dirname $0)/..

target=moviereviews
thread=8
conn=1024
repetitions=2
num_req=20000
poker_batch_req=100
num_exp=10
DIR=movie/04-08-pokerpp
FILE=mix-$target-t$thread-c$conn-r$repetitions-req$num_req-n$num_exp-poker_batch_req$poker_batch_req.log
mkdir -p $DIR
if [ -f $DIR/$FILE ]; then
    echo "File $DIR/$FILE already exists. Skipping test."
    exit 0
fi

python3 test.py -b movie -r mix -x moviereviews --num_exp $num_exp -t $thread -c $conn --poker_batch_req $poker_batch_req --repetition $repetitions --num_req $num_req >$DIR/$FILE