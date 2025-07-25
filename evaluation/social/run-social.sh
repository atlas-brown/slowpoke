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
# python3 $SLOWPOKE/src/main.py -b social -r mix -x hometimeline -t 4 -c 128 --num_exp 10 --poker_batch 60000000 --num_req 50000 >$FILE

#!/bin/bash

cd $(dirname $0)/..

target=hometimeline
thread=8
conn=1024
repetitions=5
num_req=20000
poker_batch_req=100
num_exp=10
DIR=social/04-09-pokerpp-rm-deadlock
FILE=mix-$target-t$thread-c$conn-r$repetitions-req$num_req-n$num_exp-poker_batch_req$poker_batch_req.log
mkdir -p $DIR
if [ -f $DIR/$FILE ]; then
    echo "File $DIR/$FILE already exists. Skipping test."
    exit 0
fi

python3 $SLOWPOKE_TOP/src/main.py -b social -r mix -x hometimeline --num_exp $num_exp -t $thread -c $conn --poker_batch_req $poker_batch_req --repetition $repetitions --num_req $num_req >$DIR/$FILE
