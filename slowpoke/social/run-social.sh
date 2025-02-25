#!/bin/bash

cd $(dirname $0)/..

# python3 test.py -b social -r mix -x composepost --num_exp 5
# python3 test.py -b social -r mix -x usertimeline --num_exp 10
python3 test.py -b social -r mix -x hometimeline --num_exp 10

