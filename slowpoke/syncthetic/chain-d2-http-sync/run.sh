#!/bin/bash

cd $(dirname $0)/../..

# touch syncthetic/chain-d2-http-sync/results/chain-d2-http-sync-service0.txt

# python3 test.py -b syncthetic -r chain-d2-http-sync -x service0 --num_exp 10 -c 128 -t 2 --num_req 18000 \
#     >syncthetic/chain-d2-http-sync/results/chain-d2-http-sync-service0.txt

touch syncthetic/chain-d2-http-sync/results/chain-d2-http-sync-service1.txt

python3 test.py -b syncthetic -r chain-d2-http-sync -x service1 --num_exp 10 -c 128 -t 2 --num_req 18000 \
    >syncthetic/chain-d2-http-sync/results/chain-d2-http-sync-service1.txt
