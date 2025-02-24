#!/bin/bash

cd $(dirname $0)

export SLOWPOKE_DELAY_MICROS_CART=0
export SLOWPOKE_DELAY_MICROS_CHECKOUT=0
export SLOWPOKE_DELAY_MICROS_CURRENCY=0
export SLOWPOKE_DELAY_MICROS_EMAIL=0
export SLOWPOKE_DELAY_MICROS_FRONTEND=0
export SLOWPOKE_DELAY_MICROS_PAYMENT=0
export SLOWPOKE_DELAY_MICROS_PRODUCT_CATALOG=0
export SLOWPOKE_DELAY_MICROS_RECOMMENDATIONS=0
export SLOWPOKE_DELAY_MICROS_SHIPPING=0
export SLOWPOKE_PRERUN=false

request=home
thread=${1:-16}
conn=${2:-512}
rate=${3:-0}
duration=60

bash run.sh $request $thread $conn $duration
# bash run-wrk2.sh $request $thread $conn $duration $rate