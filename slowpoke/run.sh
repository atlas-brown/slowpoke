#!/bin/bash

cd $(dirname $0)

benchmark=${1:-boutique}
request=${2:-home}
thread=${3:-16}
conn=${4:-512}
duration=${5:-60}

supported_benchmarks=("boutique" "social" "movie")

check_benchmark_supported() {
    local benchmark=$1
    for b in "${supported_benchmarks[@]}"; do
        if [[ $b == $benchmark ]]; then
            return 0
        fi
    done
    return 1
}

check_connectivity() {
    local pod_name=$1
    local service_name=$2
    # if service name is the same as the pod name (prefix), skip the check
    if [[ $1 == $2* ]]; then
        return 0
    fi
    # if the pod is ubuntu client
    if [[ $pod_name == *"ubuntu-client"* ]]; then
        kubectl exec $pod_name -- curl $service_name:80/heartbeat --max-time 1 | grep Heartbeat > /dev/null
        return $?
    fi
    kubectl exec $pod_name -- sh -c "(echo -e \"GET /heartbeat HTTP/1.1\r\nHost: $service_name\r\nConnection: close\r\n\r\n\") \
        | nc -w 1 $service_name 80" | grep Heartbeat > /dev/null
    return $?
}

fix_req_num() {
    local benchmark=$1
    local client=$2
    kubectl cp ./fix_req_n.lua ${client}:/wrk
    kubectl exec ${client} -- /bin/sh -c "cat /wrk/fix_req_n.lua >> ${benchmark}"
}

run_test() {
    local benchmark=$1
    local ubuntu_client=$(kubectl get pod | grep ubuntu-client- | cut -f 1 -d " ") 
    if [[ $benchmark == "boutique" ]]; then
        # run the load generator
        echo "[run.sh] Running warmup test" 
        echo "[run.sh] /wrk/wrk -t${thread} -c${conn} -d3s -L -s /wrk/scripts/online-boutique/${request}.lua http://frontend:80"
        kubectl exec $ubuntu_client -- /wrk/wrk -t${thread} -c${conn} -d3s -L -s /wrk/scripts/online-boutique/${request}.lua http://frontend:80
	fix_req_num "/wrk/scripts/online-boutique/${request}.lua" $ubuntu_client
        sleep 10
        echo "[run.sh] /wrk/wrk -t${thread} -c${conn} -d${duration}s -L -s /wrk/scripts/online-boutique/${request}.lua http://frontend:80"
        kubectl exec $ubuntu_client -- /wrk/wrk -t${thread} -c${conn} -d${duration}s -L -s /wrk/scripts/online-boutique/${request}.lua http://frontend:80
    else
        echo "[run.sh] Starting the rust proxy first for $benchmark"
        kubectl exec $ubuntu_client -- bash -c "/mucache/proxy/target/release/proxy ${benchmark} &"
        sleep 3
        echo "[run.sh] Running warmup test"
        echo "[run.sh] /wrk/wrk -t${thread} -c${conn} -d3s -L http://localhost:3000"
        kubectl exec $ubuntu_client -- /wrk/wrk -t${thread} -c${conn} -d3s http://localhost:3000
        sleep 10
        echo "[run.sh] /wrk/wrk -t${thread} -c${conn} -d${duration}s -L http://localhost:3000"
        kubectl exec $ubuntu_client -- /wrk/wrk -t${thread} -c${conn} -d${duration}s -L http://localhost:3000
    fi
}

populate() {
    local benchmark=$1
    if [[ $benchmark == "boutique" ]]; then
        echo "[run.sh] No population needed for $benchmark"
        return
    fi
    local ubuntu_client=$(kubectl get pod | grep ubuntu-client- | cut -f 1 -d " ") 
    echo "[run.sh] Populating social benchmark"
    bash $benchmark/populate.sh 
    echo "[run.sh] Copying $benchmark/analysis.txt to $ubuntu_client:/analysis.txt"
    kubectl cp $benchmark/data/analysis.txt $ubuntu_client:/analysis.txt
    echo "[run.sh] Finished populating $benchmark"
}

check_benchmark_supported $benchmark
if [ $? -ne 0 ]; then
    echo "[run.sh] Benchmark $benchmark is not supported"
    exit 1
fi

echo "[run.sh] Running benchmark $benchmark with request $request, thread $thread, conn $conn, duration $duration"

# delete all services
echo "[run.sh] Deleting all services"
kubectl delete -f $benchmark/yamls/ --ignore-not-found=true
kubectl delete -f client.yaml --ignore-not-found=true
# wait for all pods to be deleted
echo "[run.sh] Waiting for all pods to be deleted"
while [[ $(kubectl get pods | grep -v -E 'STATUS' | wc -l) -ne 0 ]]; do
    sleep 1
done

# deploy all services
echo "[run.sh] Deploying all services"
for file in $(ls -d $benchmark/yamls/*.yaml)
do
    envsubst < $file | kubectl apply -f - 
done

kubectl get pod | grep ubuntu-client- 
if [ $? -ne 0 ]
then
    echo "[run.sh] Client pod not found, deploying client"
    kubectl apply -f client.yaml  
fi

# wait until all pods are ready by checking the log to see if the "server started" message is printed
echo "[run.sh] Waiting for all pods to be running"
while [[ $(kubectl get pods | grep -v -E 'Running|Completed|STATUS' | wc -l) -ne 0 ]]; do
  sleep 1
done

# echo "[run.sh] Waiting for all pods to be ready"
# while true
# do
#     res=$(kubectl get pods | cut -f 1 -d " " | grep -vE "|NAME|redis")
#     check=0
#     IFS=$'\n' read -rd '' -a array <<< "$res"
#     for value in "${array[@]:1:${#array[@]}-1}"
#     do
#         kubectl logs $value | grep "Server started" > /dev/null
#         if [ $? -ne 0 ]
#         then
#             check=1
#             # echo "waiting for $value"
#             sleep 1
#             break
#         fi
#     done
#     if [ $check -eq 0 ]
#     then
#         echo "[run.sh] All pods are ready"
#         break
#     fi
# done

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

populate $benchmark
sleep 5

run_test $benchmark &
pid=$!

# sleep 0.9*duration
sleep $(echo "$duration*0.9" | bc -l)
echo "[run.sh] Checking the resource usage"
kubectl top pods

wait $pid
status=$?
echo "[run.sh] Test finished with status $status"
exit $status
