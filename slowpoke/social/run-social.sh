#!/bin/bash

cd $(dirname $0)/..

python3 test.py -b social -r mix -x composepost -d10 --num_exp 2 