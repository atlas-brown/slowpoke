#!/bin/bash

cd $(dirname $0)

export PROCESSING_TIME_SERVICE0=0.00097282
export PROCESSING_TIME_SERVICE1=0.00082341
export PROCESSING_TIME_SERVICE2=0.00069813
export PROCESSING_TIME_SERVICE3=0.0006208
export PROCESSING_TIME_SERVICE4=0.00053849
export PROCESSING_TIME_SERVICE5=0.00039561
export SLOWPOKE_DELAY_MICROS_SERVICE0=0.0
export SLOWPOKE_POKER_BATCH_THRESHOLD_SERVICE0=100
export SLOWPOKE_IS_TARGET_SERVICE_SERVICE0=true
export SLOWPOKE_DELAY_MICROS_SERVICE1=0.0
export SLOWPOKE_POKER_BATCH_THRESHOLD_SERVICE1=100
export SLOWPOKE_IS_TARGET_SERVICE_SERVICE1=false
export SLOWPOKE_DELAY_MICROS_SERVICE2=0.0
export SLOWPOKE_POKER_BATCH_THRESHOLD_SERVICE2=100
export SLOWPOKE_IS_TARGET_SERVICE_SERVICE2=false
export SLOWPOKE_DELAY_MICROS_SERVICE3=0.0
export SLOWPOKE_POKER_BATCH_THRESHOLD_SERVICE3=100
export SLOWPOKE_IS_TARGET_SERVICE_SERVICE3=false
export SLOWPOKE_DELAY_MICROS_SERVICE4=0.0
export SLOWPOKE_POKER_BATCH_THRESHOLD_SERVICE4=100
export SLOWPOKE_IS_TARGET_SERVICE_SERVICE4=false
export SLOWPOKE_DELAY_MICROS_SERVICE5=0.0
export SLOWPOKE_POKER_BATCH_THRESHOLD_SERVICE5=100
export SLOWPOKE_IS_TARGET_SERVICE_SERVICE5=false

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
# FILE=100reqbatching-all-delay-order-large.log
# if [[ -e $DIR/$FILE ]]; then
#     echo "File $DIR/$FILE already exists. Skipping..."
#     exit 0
# fi
# env | grep "SLOWPOKE" >$DIR/$FILE
# env | grep "PROCESSING" >>$DIR/$FILE
# echo "bash run.sh $benchmark $request $thread $conn $num_req" >>$DIR/$FILE
# bash run.sh $benchmark $request $thread $conn $num_req >>$DIR/$FILE

YAMLS_DIR=/home/ubuntu/mucache/slowpoke/synthetic/chain-d4-http-sync/yamls
for file in $YAMLS_DIR/*.yaml
do 
    envsub <$file | kubectl delete -f -
    envsub <$file | kubectl apply -f -
done