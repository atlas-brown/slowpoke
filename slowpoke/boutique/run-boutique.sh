#!/bin/bash

cd $(dirname "$0")/..

# python3 test.py -b boutique -x cart -r mix -d 100

target=cart
thread=8
conn=512
repetitions=3
num_req=200000
poker_batch=30000000
num_exp=2
DIR=boutique/results-all-flush-time-based-sleep
FILE=mix-$target-t$thread-c$conn-r$repetitions-req$num_req-n$num_exp-batch$poker_batch-no-flush-fixed-batch-delay-before-processing.log
mkdir -p $DIR
if [ -f $DIR/$FILE ]; then
    echo "File $DIR/$FILE already exists. Skipping test."
    exit 0
fi
python3 test.py -b boutique -x cart -r mix -t $thread -c $conn --num_exp $num_exp --repetitions $repetitions --num_req $num_req --poker_batch $poker_batch >$DIR/$FILE