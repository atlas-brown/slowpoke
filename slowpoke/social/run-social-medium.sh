#!/bin/bash

outfile=$1
outdir=$(realpath $(dirname $outfile))
name=$(basename outfile)
outfile=$outdir/$name
cd $(dirname $0)/..

target=hometimeline
thread=8
conn=1024
repetitions=1
num_req=30000
poker_batch_req=100
num_exp=5
DIR=social/04-09-pokerpp-rm-deadlock
FILE=mix-$target-t$thread-c$conn-r$repetitions-req$num_req-n$num_exp-poker_batch_req$poker_batch_req.log
mkdir -p $DIR
if [ -f $DIR/$FILE ]; then
    echo "File $DIR/$FILE already exists. Skipping test."
    exit 0
fi

kubectl delete deployments --all
kubectl delete services --all
python3 test.py -b social -r mix -x hometimeline --num_exp $num_exp -t $thread -c $conn --poker_batch_req $poker_batch_req --repetition $repetitions --num_req $num_req >$outfile
