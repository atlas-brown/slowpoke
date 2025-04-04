#!/bin/bash

cd $(dirname $0)


export SLOWPOKE_POKER_BATCH_THRESHOLD=10000000
export PROCESSING_TIME_SERVICE0=0.0 # base
export PROCESSING_TIME_SERVICE1=0.0 # base
export PROCESSING_TIME_SERVICE2=0.0 # base
export PROCESSING_TIME_SERVICE3=0.0 # base
export PROCESSING_TIME_SERVICE4=0.0 # base
export PROCESSING_TIME_SERVICE5=0.0 # base
export SLOWPOKE_DELAY_MICROS_SERVICE0=0
export SLOWPOKE_DELAY_MICROS_SERVICE1=1000
export SLOWPOKE_DELAY_MICROS_SERVICE2=200
export SLOWPOKE_DELAY_MICROS_SERVICE3=300
export SLOWPOKE_DELAY_MICROS_SERVICE4=400
export SLOWPOKE_DELAY_MICROS_SERVICE5=100
# export SLOWPOKE_DELAY_MICROS_SERVICE0=0
# export SLOWPOKE_DELAY_MICROS_SERVICE1=0
# export SLOWPOKE_DELAY_MICROS_SERVICE2=0
# export SLOWPOKE_DELAY_MICROS_SERVICE3=0
# export SLOWPOKE_DELAY_MICROS_SERVICE4=0
# export SLOWPOKE_DELAY_MICROS_SERVICE5=0

cd /home/ubuntu/mucache/slowpoke
DIR=/home/ubuntu/mucache/slowpoke/synthetic/chain-d4-http-sync/bus-theory-microbenchmark
mkdir -p $DIR

benchmark=synthetic
request=chain-d4-http-sync
thread=1
conn=1
num_req=2000
# FILE=service1-base.log
# FILE=combined_no_rlock-poker$threadshold.log
# FILE=s1_no_rlock-base.log
FILE=100reqbatching-all-delay-order-large.log
if [[ -e $DIR/$FILE ]]; then
    echo "File $DIR/$FILE already exists. Skipping..."
    exit 0
fi
env | grep "SLOWPOKE" >$DIR/$FILE
env | grep "PROCESSING" >>$DIR/$FILE
echo "bash run.sh $benchmark $request $thread $conn $num_req" >>$DIR/$FILE
bash run.sh $benchmark $request $thread $conn $num_req >>$DIR/$FILE
