#!/bin/bash

cd $(dirname $0)/../..

python3 test.py -b syncthetic -r chain-d2-grpc-sync -x service0 --num_exp 5 -c 64 -t 2 --num_req 18000