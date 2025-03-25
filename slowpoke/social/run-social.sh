#!/bin/bash

cd $(dirname $0)/..

FILE=social/one-service-per-node-poker/hometimeline-4-128-50000-poker50000000-full.log
# Check if exists to avoid overwritting
if [ -e "$FILE" ]; then
    echo "File $FILE already exists. Exiting..."
    exit 1  # Exit with status code 1 to indicate an error
fi
python3 test.py -b social -r mix -x hometimeline -t 4 -c 128 --num_exp 10 --poker_batch 50000000 --num_req 50000 >$FILE
