#!/bin/bash

outfile=$1
outdir=$(realpath $(dirname $outfile))
name=$(basename $outfile)
outfile=$outdir/$name
cd $(dirname "$0")/..

# python3 $SLOWPOKE_TOP/src/main.py -b boutique -x cart -r mix -d 100

target=cart
thread=8
conn=1024
repetitions=1
num_req=80000
poker_batch_req=100
num_exp=10
DIR=boutique/04-15-pokerpp-fixed-deadlock
FILE=mix-$target-t$thread-c$conn-r$repetitions-req$num_req-n$num_exp-poker_batch_req$poker_batch_req.log

kubectl delete deployments --all
kubectl delete services --all
python3 $SLOWPOKE_TOP/src/main.py -b boutique -x $target -r mix -t $thread -c $conn --num_exp $num_exp --repetitions $repetitions --num_req $num_req --poker_batch_req $poker_batch_req >$outfile
