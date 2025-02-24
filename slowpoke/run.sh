#!/bin/bash

cd $(dirname $0)

request=${1:-home}
thread=${2:-16}
conn=${3:-512}
duration=${4:-60}

check_connectivity() {
    local pod_name=$1
    local service_name=$2
    # if the pod is ubuntu client
    if [[ $pod_name == *"ubuntu-client"* ]]; then
        kubectl exec $pod_name -- curl $service_name:80/heartbeat --max-time 1 | grep Heartbeat > /dev/null
        return $?
    fi
    kubectl exec $pod_name -- sh -c "(echo -e \"GET /heartbeat HTTP/1.1\r\nHost: $service_name\r\nConnection: close\r\n\r\n\") \
        | nc -w 1 $service_name 80" | grep Heartbeat > /dev/null
    return $?
}

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
    all_connected=1
    for pod in $(kubectl get pods | grep -v -E 'NAME' | cut -f 1 -d " ")
    do
        for service in $(kubectl get svc | grep -v -E 'NAME|kube' | cut -f 1 -d " ")
        do
            check_connectivity $pod $service
            if [ $? -ne 0 ]
            then
                echo "[run.sh] $pod cannot connect to $service"
                all_connected=0
                break
            fi
        done
        if [ $all_connected -eq 0 ]
        then
            break
        fi
    done
    if [ $all_connected -eq 1 ]
    then
        echo "[run.sh] All pods can connect to all services"
        break
    fi
done

# run the load generator
echo "[run.sh] Running warmup test"
kubectl exec $ubuntu_client -- /wrk/wrk -t${thread} -c${conn} -d3s -L -s /wrk/scripts/online-boutique/${request}.lua http://frontend:80
sleep 10
echo "[run.sh] /wrk/wrk -t${thread} -c${conn} -d${duration}s -L -s /wrk/scripts/online-boutique/${request}.lua http://frontend:80"
kubectl exec $ubuntu_client -- /wrk/wrk -t${thread} -c${conn} -d${duration}s -L -s /wrk/scripts/online-boutique/${request}.lua http://frontend:80 &
pid=$!

# sleep 0.9*duration
sleep $(echo "$duration*0.9" | bc -l)
echo "[run.sh] Checking the resource usage"
kubectl top pods

wait $pid