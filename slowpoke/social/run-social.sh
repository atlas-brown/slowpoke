#!/bin/bash

cd $(dirname $0)/..

DIR=social/one-service-per-node-poker-time-spining
mkdir -p $DIR
FILE=$DIR/hometimeline-4-128-50000-poker60000000.log
# Check if exists to avoid overwritting
if [ -e "$FILE" ]; then
    echo "File $FILE already exists. Exiting..."
    exit 1  # Exit with status code 1 to indicate an error
fi
python3 test.py -b social -r mix -x hometimeline -t 4 -c 128 --num_exp 2 --poker_batch 60000000 --num_req 50000 >$FILE
