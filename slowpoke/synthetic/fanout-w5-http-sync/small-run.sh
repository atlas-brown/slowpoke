#!/bin/bash

cd /home/ubuntu/mucache/slowpoke


export SLOWPOKE_POKER_BATCH_THRESHOLD=10000000
export SLOWPOKE_PRERUN=false
export PROCESSING_TIME_SERVICE0=0.000709 # base
export PROCESSING_TIME_SERVICE1=0.000010 # base
export PROCESSING_TIME_SERVICE2=0.000032 # base
export PROCESSING_TIME_SERVICE3=0.000012 # base
export PROCESSING_TIME_SERVICE4=0.000050 # base
export PROCESSING_TIME_SERVICE5=0.000020 # base
export PROCESSING_TIME_SERVICE6=0.000070 # base
# export SLOWPOKE_DELAY_MICROS_SERVICE0=0
# export SLOWPOKE_DELAY_MICROS_SERVICE1=1000
# export SLOWPOKE_DELAY_MICROS_SERVICE2=200
# export SLOWPOKE_DELAY_MICROS_SERVICE3=300
# export SLOWPOKE_DELAY_MICROS_SERVICE4=400
# export SLOWPOKE_DELAY_MICROS_SERVICE5=100
# export SLOWPOKE_DELAY_MICROS_SERVICE6=0
export SLOWPOKE_DELAY_MICROS_SERVICE0=0
export SLOWPOKE_DELAY_MICROS_SERVICE1=1500
export SLOWPOKE_DELAY_MICROS_SERVICE2=502
export SLOWPOKE_DELAY_MICROS_SERVICE3=740
export SLOWPOKE_DELAY_MICROS_SERVICE4=560
export SLOWPOKE_DELAY_MICROS_SERVICE5=670
export SLOWPOKE_DELAY_MICROS_SERVICE6=680



cd /home/ubuntu/mucache/slowpoke
DIR=/home/ubuntu/mucache/slowpoke/synthetic/fanout-w5-http-sync/bus-theory-microbenchmark
mkdir -p $DIR

benchmark=synthetic
request=fanout-w5-http-sync
thread=1
conn=1
num_req=2000
# FILE=service1-base.log
# FILE=combined_no_rlock-poker$threadshold.log
# FILE=s1_no_rlock-base.log
# FILE=100reqbatching-all-delay-fixed-context-2proc-fixed-time.log
# if [[ -e $DIR/$FILE ]]; then
#     echo "File $DIR/$FILE already exists. Skipping..."
#     exit 0
# fi
env | grep "SLOWPOKE" 
env | grep "PROCESSING"
echo "bash run.sh $benchmark $request $thread $conn $num_req"
bash run.sh $benchmark $request $thread $conn $num_req
