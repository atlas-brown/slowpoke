#!/bin/bash

cd $(dirname $0)

request=${1:-home}
thread=${2:-16}
conn=${3:-512}
duration=${4:-60}
rate=${5:-1000}
echo "[run.sh] Running $request with $thread threads, $conn connections, and $duration seconds"

# delete all services
echo "[run.sh] Deleting all services"
kubectl delete -f yamls/ --ignore-not-found=true
# wait for all pods to be deleted
echo "[run.sh] Waiting for all pods to be deleted"
while [[ $(kubectl get pods | grep -v -E 'STATUS|ubuntu' | wc -l) -ne 0 ]]; do
    sleep 1
done

kubectl get pod | grep ubuntu-client- 
if [ $? -ne 0 ]
then
    echo "[run.sh] Client pod not found, deploying client"
    kubectl apply -f client.yaml  
fi
# # deploy client
# echo "[run.sh] Deploying client"
# kubectl apply -f client.yaml

echo "[run.sh] Waiting for client to be ready"
ubuntu_client=$(kubectl get pod | grep ubuntu-client- | cut -f 1 -d " ") 
while true
do 
    kubectl logs $ubuntu_client | grep "Init" > /dev/null
    if [ $? -eq 0 ]
    then
        break
    fi
    sleep 1
done
echo "[run.sh] Client is ready"

# deploy all services
echo "[run.sh] Deploying all services"
for file in $(ls -d yamls/*.yaml)
do
    envsubst < $file | kubectl apply -f - 
done

# wait until all pods are ready by checking the log to see if the "server started" message is printed
echo "[run.sh] Waiting for all pods to be running"
while [[ $(kubectl get pods | grep -v -E 'Running|Completed|STATUS' | wc -l) -ne 0 ]]; do
  sleep 1
done

echo "[run.sh] Waiting for all pods to be ready"
while true
do
    res=$(kubectl get pods | cut -f 1 -d " " | grep -vE "ubuntu|NAME|redis")
    check=0
    echo $res
    IFS=$'\n' read -rd '' -a array <<< "$res"
    for value in "${array[@]:1:${#array[@]}-1}"
    do
        kubectl logs $value | grep "Server started" > /dev/null
        if [ $? -ne 0 ]
        then
            check=1
            # echo "waiting for $value"
            sleep 1
            break
        fi
    done
    if [ $check -eq 0 ]
    then
        echo "[run.sh] All pods are ready"
        break
    fi
done

echo "[run.sh] Checking heartbeat for all services"
while true
do
    res=$(kubectl get svc | cut -f 1 -d " " | grep -vE "kube")
    check=0
    IFS=$'\n' read -rd '' -a array <<< "$res"
    for value in "${array[@]:1:${#array[@]}-1}"
    do
        kubectl exec $ubuntu_client -- curl $value:80/heartbeat --max-time 1 | grep Heartbeat > /dev/null
        if [ $? -ne 0 ]
        then
            echo "[run.sh] $value is not reachable"
            check=1
            sleep 1
            break
        fi
    done
    if [ $check -eq 0 ]
    then
        echo "[run.sh] All services are ready"
        break
    fi
done
sleep 10

# run the load generator
echo "[run.sh] Running warmup test"
kubectl exec $ubuntu_client -- /wrk2/wrk -t${thread} -c${conn} -R${rate} -d3s -L -s /wrk/scripts/online-boutique/home.lua http://frontend:80 &
sleep 10
echo "[run.sh] /wrk2/wrk -t${thread} -c${conn} -d${duration}s -L -s /wrk/scripts/online-boutique/home.lua http://frontend:80"
kubectl exec $ubuntu_client -- /wrk2/wrk -t${thread} -c${conn} -R${rate} -d${duration}s -L -s /wrk/scripts/online-boutique/home.lua http://frontend:80 &
pid=$!

sleep 50
echo "[run.sh] Checking the resource usage"
kubectl top pods

wait $pid