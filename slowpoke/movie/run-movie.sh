#!/bin/bash

cd $(dirname $0)/..

# python3 test.py -b movie -r mix -x castinfo --num_exp 5 -t 8 -c 512

python3 test.py -b movie -r mix -x moviereviews --num_exp 10 -t 8 -c 1024 >movie/new-results/mix-moviereviews-8-1024-50000-10-3.log

