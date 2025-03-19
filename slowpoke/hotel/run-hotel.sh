#!/bin/bash

cd $(dirname $0)/..

python3 test.py -b hotel -r mix -x profile --num_exp 10 -t 4 -c 256

