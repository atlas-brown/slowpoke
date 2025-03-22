#!/bin/bash

cd $(dirname $0)/..

# python3 test.py -b movie -r mix -x castinfo --num_exp 5 -t 8 -c 512

python3 test.py -b movie -r mix -x moviereviews --num_exp 10 -t 4 -c 256 >movie/new-results/mix-moviereviews-4-256.log

