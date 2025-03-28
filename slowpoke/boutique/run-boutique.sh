#!/bin/bash

cd $(dirname "$0")/..

# python3 test.py -b boutique -x cart -r mix -d 100

python3 test.py -b boutique -x cart -r mix -t 8 -c 512 --clien_cpu_quota 8 --num_exp 10 --repetitions 3 --num_req 200000 >boutique/results-one-service-per-node/mix-cart.log