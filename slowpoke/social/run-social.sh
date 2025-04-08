#!/bin/bash

# cd $(dirname $0)/..

# DIR=social/one-service-per-node-poker-time-spining
# mkdir -p $DIR
# FILE=$DIR/hometimeline-4-128-50000-poker60000000-full.log
# # Check if exists to avoid overwritting
# if [ -e "$FILE" ]; then
#     echo "File $FILE already exists. Exiting..."
#     exit 1  # Exit with status code 1 to indicate an error
# fi
# python3 test.py -b social -r mix -x hometimeline -t 4 -c 128 --num_exp 10 --poker_batch 60000000 --num_req 50000 >$FILE

#!/bin/bash

cd $(dirname $0)/..

target=profile
thread=8
conn=256
repetitions=3
num_req=50000
poker_batch=40000000
poker_batch_req=50
num_exp=10
DIR=social/04-07-all-flush-time-based-sleep
FILE=mix-$target-t$thread-c$conn-r$repetitions-req$num_req-n$num_exp-poker_batch_req$poker_batch_req.log
mkdir -p $DIR
if [ -f $DIR/$FILE ]; then
    echo "File $DIR/$FILE already exists. Skipping test."
    exit 0
fi

python3 test.py -b social -r mix -x hometimeline --num_exp $num_exp -t $thread -c $conn --poker_batch_req $poker_batch_req --repetition $repetitions --num_req $num_req >$DIR/$FILE
