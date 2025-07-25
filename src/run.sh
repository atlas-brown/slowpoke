#!/bin/bash

export SLOWPOKE_TOP=${SLOWPOKE_TOP:-$(cd "${BASH_SOURCE%/*}/.." && pwd -P)}
echo "SLOWPOKE_TOP is set to $SLOWPOKE_TOP"

cd $(dirname $0)

benchmark=${1:-boutique}
request=${2:-home}
thread=${3:-16}
conn=${4:-512}
TOTAL_REQ=${5:-50000}

duration=60

YAML_PATH=$SLOWPOKE_TOP/evaluation/$benchmark/yamls
if [[ $benchmark == "synthetic" ]]; then
    YAML_PATH=$SLOWPOKE_TOP/$benchmark/$request/yamls
fi

supported_benchmarks=("boutique" "social" "movie" "hotel" "synthetic")

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

check_connectivity_all(){
    echo "[run.sh] Checking heartbeat for all services"
    # if "grpc" in request, ignore
    if [[ $request == *grpc* ]]; then
        sleep 5
        return 0
    fi
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
}

fix_req_num() {
    local benchmark=$1
    local client=$2
    counter=$((TOTAL_REQ / thread))
    PER_THREAD_COUNTER=$counter envsubst < $SLOWPOKE_TOP/client/fix_req_n.lua > /tmp/temp_fix_req_n.lua
    kubectl cp /tmp/temp_fix_req_n.lua ${client}:/wrk/fix_req_n.lua
    rm /tmp/temp_fix_req_n.lua  # clean up
    if [[ $benchmark == *"boutique"* ]]; then
        kubectl exec ${client} -- /bin/sh -c "cat /wrk/fix_req_n.lua >> /wrk/scripts/online-boutique/${request}.lua"
        return
    fi
}

warmup_and_speed() {
    local client=$1
    local thread=$2
    local conn=$3
    local host=$4
    local script="-s $5"
    if [[ -z $5 ]]; then
        script=""
    fi
    output=$(kubectl exec $client -- /wrk/wrk -t${thread} -c${conn} --timeout 3s -d3s -L $script $host)
    echo $output
}

run_test() {
    local benchmark=$1
    local ubuntu_client=$(kubectl get pod | grep ubuntu-client- | cut -f 1 -d " ") 

    if [[ $benchmark != "boutique" && $benchmark != "synthetic" ]]; then
        echo "[run.sh] Starting the rust proxy first for $benchmark"
        kubectl exec $ubuntu_client -- bash -c "/mucache/proxy/target/release/proxy ${benchmark} &"
        sleep 3
    fi

    echo "[run.sh] Running warmup test" 
    if [[ $benchmark == "boutique" ]]; then
       # run the load generator
        echo "[run.sh] /wrk/wrk -t${thread} -c${conn} --timeout 3s -d3s -L -s /wrk/scripts/online-boutique/${request}.lua http://frontend:80"
        output=$(kubectl exec $ubuntu_client -- /wrk/wrk -t${thread} -c${conn} --timeout 3s -d3s -L -s /wrk/scripts/online-boutique/${request}.lua http://frontend:80)
    elif [[ $benchmark == "synthetic" ]]; then
        echo "[run.sh] /wrk/wrk -t${thread} -c${conn} --timeout 3s -d3s -L http://service0:80/endpoint1"
        output=$(kubectl exec $ubuntu_client -- /wrk/wrk -t${thread} -c${conn} --timeout 3s -d3s -L http://service0:80/endpoint1)
    else 
        echo "[run.sh] /wrk/wrk -t${thread} -c${conn} --timeout 3s -d3s -L http://localhost:3000"
        output=$(kubectl exec $ubuntu_client -- /wrk/wrk -t${thread} -c${conn} --timeout 3s -d3s -L http://localhost:3000)
    fi
    echo "$output"

    # get the speed of the warmup test and estimate the duration
    speed=$(echo "$output" | grep "Requests/sec:" | awk '{print $2}')
    duration=$(awk -v r="$TOTAL_REQ" -v s="$speed" 'BEGIN{print (1.5*r/s < 200) ? int(1.5*r/s) : 200}')
    echo "[run.sh] Speed is $speed, duration is $duration"

    echo "[run.sh] Fix the request number."
    fix_req_num $benchmark $ubuntu_client

    echo "[run.sh] Running the actual test"
    sleep 10
    if [[ $benchmark == "boutique" ]]; then
        echo "[run.sh] /wrk/wrk --timeout 20s -t${thread} -c${conn} -d${duration}s -L -s /wrk/scripts/online-boutique/${request}.lua http://frontend:80"
        kubectl exec $ubuntu_client -- /wrk/wrk --timeout 20s -t${thread} -c${conn} -d${duration}s -L -s /wrk/scripts/online-boutique/${request}.lua http://frontend:80
    elif [[ $benchmark == "synthetic" ]]; then
        echo "[run.sh] /wrk/wrk --timeout 20s -t${thread} -c${conn} -d${duration}s -L -s /wrk/fix_req_n.lua http://service0:80/endpoint1"
        # kubectl exec $ubuntu_client -- /wrk/wrk --timeout 20s -t${thread} -c${conn} -d20s -L http://service0:80/endpoint1
        kubectl exec $ubuntu_client -- /wrk/wrk --timeout 20s -t${thread} -c${conn} -d${duration}s -L -s /wrk/fix_req_n.lua http://service0:80/endpoint1
    else
        echo "[run.sh] /wrk/wrk --timeout 20s -t${thread} -c${conn} -d${duration}s -s /wrk/fix_req_n.lua -L http://localhost:3000"
        kubectl exec $ubuntu_client -- /wrk/wrk --timeout 20s -t${thread} -c${conn} -d${duration}s -s /wrk/fix_req_n.lua -L http://localhost:3000
    fi
}

populate() {
    local ubuntu_client=$(kubectl get pod | grep ubuntu-client- | cut -f 1 -d " ") 
    local benchmark=$1
    if [[ $benchmark == "boutique" ]]; then
        echo "[run.sh] No population needed for $benchmark"
        return
    fi
    if [[ $benchmark == "hotel" || $benchmark == "movie" ]]; then
        echo "[run.sh] Copying $SLOWPOKE_TOP/evaluation/$benchmark/analysis.txt to $ubuntu_client:/analysis.txt"
        kubectl cp $SLOWPOKE_TOP/evaluation/$benchmark/data/analysis.txt $ubuntu_client:/analysis.txt
        echo "[run.sh] Finished populating $benchmark"
        return
    fi
    echo "[run.sh] Populating social benchmark"
    bash $SLOWPOKE_TOP/evaluation/$benchmark/populate.sh 
    echo "[run.sh] Copying $SLOWPOKE_TOP/evaluation/$benchmark/analysis.txt to $ubuntu_client:/analysis.txt"
    kubectl cp $SLOWPOKE_TOP/evaluation/$benchmark/analysis.txt $ubuntu_client:/analysis.txt
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
kubectl delete -f $YAML_PATH --ignore-not-found=true
kubectl delete -f client.yaml --ignore-not-found=true
# wait for all pods to be deleted
echo "[run.sh] Waiting for all pods to be deleted"
while [[ $(kubectl get pods | grep -v -E 'STATUS' | wc -l) -ne 0 ]]; do
    sleep 1
done

# deploy all services
echo "[run.sh] Deploying all services"
for file in $(ls -d $YAML_PATH/*.yaml)
do
    envsubst < $file | kubectl apply -f - 
done

kubectl get pod | grep ubuntu-client- 
if [ $? -ne 0 ]
then
    echo "[run.sh] Client pod not found, deploying client"
    envsubst < $SLOWPOKE_TOP/client/client.yaml | kubectl apply -f -
fi

# wait until all pods are ready by checking the log to see if the "server started" message is printed
echo "[run.sh] Waiting for all pods to be running"
while [[ $(kubectl get pods | grep -v -E 'Running|Completed|STATUS' | wc -l) -ne 0 ]]; do
  sleep 1
done
while [[ $(kubectl get pods | grep -v -E '1/1|STATUS' | wc -l) -ne 0 ]]; do
  sleep 1
done
echo "[run.sh] All pods are running"


check_connectivity_all


if [[ $benchmark != "synthetic" ]]; then
    populate $benchmark
fi

sleep 5

run_test $benchmark &
pid=$!

# sleep 0.9*duration
sleep $(echo "$duration*0.8" | bc -l)
echo "[run.sh] Checking the resource usage"
kubectl top pods

wait $pid
status=$?
echo "[run.sh] Test finished with status $status"
exit $status
