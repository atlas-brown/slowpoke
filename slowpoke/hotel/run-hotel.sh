#!/bin/bash

cd $(dirname $0)/..
DIR=hotel/one-service-per-node-poker-time-spinning
mkdir -p $DIR
FILE=$DIR/profile-8-512-50000-poker60000000-long-timeout-full.log
# Check if exists to avoid overwritting
if [ -e "$FILE" ]; then
    echo "File $FILE already exists. Exiting..."
    exit 1  # Exit with status code 1 to indicate an error
fi

python3 test.py -b hotel -r mix -x profile --num_exp 10 -t 8 -c 512 --clien_cpu 8 --poker_batch 60000000 >$FILE