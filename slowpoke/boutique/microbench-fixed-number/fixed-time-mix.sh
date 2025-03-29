#!/bin/bash

ubuntu_client=$(kubectl get pod | grep ubuntu-client- | cut -f 1 -d " ") 
times="30 60 90 120 150"
for i in {1..5}
do
    time=$(echo $times | cut -d " " -f $i)
    echo "[run.sh] Running the actual test $i with duration $time"
    kubectl exec $ubuntu_client -- /wrk/wrk -t8 -c512 -d${time}s -L -s /wrk/scripts/online-boutique/mix.lua http://frontend:80
    sleep 30
done