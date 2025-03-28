#!/bin/bash

cd $(dirname $0)/..

DIR=movie/one-service-per-node-poker-time-spinning
mkdir -p $DIR

# config
THREAD=8
CONN=1024
NUM_REQ=50000
POKER_BATCH=60000000
FILE=$DIR/moviereviews-$THREAD-$CONN-$NUM_REQ-poker$POKER_BATCH-full.log
# Check if exists to avoid overwritting
if [ -e "$FILE" ]; then
    echo "File $FILE already exists. Exiting..."
    exit 1  # Exit with status code 1 to indicate an error
fi

# python3 test.py -b movie -r mix -x castinfo --num_exp 5 -t 8 -c 512

python3 test.py -b movie -r mix -x moviereviews --num_exp 10 -t $THREAD -c $CONN --poker_batch $POKER_BATCH >$FILE

