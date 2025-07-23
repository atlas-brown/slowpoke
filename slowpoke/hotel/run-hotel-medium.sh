#!/bin/bash

outfile=$1
outdir=$(realpath $(dirname $outfile))
name=$(basename $outfile)
outfile=$outdir/$name
cd $(dirname $0)/..

target=profile
thread=8
conn=1024
repetitions=1
num_req=10000
poker_batch_req=100
num_exp=10
DIR=hotel/04-08-pokerpp
FILE=mix-$target-t$thread-c$conn-r$repetitions-req$num_req-n$num_exp-poker_batch_req$poker_batch_req-2.log
mkdir -p $DIR

kubectl delete deployments --all
kubectl delete services --all
python3 test.py -b hotel -r mix -x $target --num_exp $num_exp -t $thread -c $conn --poker_batch_req $poker_batch_req --repetition $repetitions --num_req $num_req >$outfile
