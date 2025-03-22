#!/bin/bash

cd $(dirname $0)/..

# python3 test.py -b social -r mix -x composepost --num_exp 5
python3 test.py -b social -r mix -x hometimeline -t 8 -c 512 --num_exp 10 >social/new-results/hometimeline-8-512-100000.log
# python3 test.py -b social -r mix -x hometimeline --num_exp 10

