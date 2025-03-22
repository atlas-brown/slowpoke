#!/bin/bash

cd $(dirname $0)/..

python3 test.py -b hotel -r mix -x profile --num_exp 10 -t 8 -c 512 --clien_cpu 8 >hotel/new-results/mix-profile.log