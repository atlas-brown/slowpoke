#!/bin/bash

cd $(dirname $0)/..

python3 test.py -b hotel -r mix -x profile --num_exp 10 -t 8 -c 1024 --clien_cpu 8 >hotel/new-results/mix-profile-8-1024.log