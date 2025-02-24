#!python3

import os
import subprocess
import copy

TARGET_SERVICE = "cart"
TARGET_PROCESSING_TIME_RANGE = [0, 1000]
TARGET_NUM_EXP = 10
# REQUEST_RATIO = {
#     "cart": 1, 
#     "checkout": 1, 
#     "currency": 1, 
#     "email": 1, 
#     "frontend": 1, 
#     "payment": 1, 
#     "product_catalog": 1, 
#     "recommendations": 1, 
#     "shipping": 1
# }
REQUEST_RATIO = {
    'cart': 0.45175405908969923, 
    'checkout': 0.05106201756720788, 
    'currency': 0.2597231833910035, 
    'email': 0.0, 
    'frontend': 1.0, 
    'payment': 0.05106201756720788, 
    'product_catalog': 0.6091136545115784, 
    'recommendations': 0.0, 
    'shipping': 0.10212403513441576
}
BASELINE_SERVICE_PROCESSING_TIME = {
        "cart":0,
        "checkout":0,
        "currency":0,
        "email":0,
        "frontend":0,
        "payment":0,
        "product_catalog":0,
        "recommendations":0,
        "shipping":0
    }
REQUEST="home"

def exp(service_delay):
    env = os.environ.copy()
    for service, delay in service_delay.items():
        env[f"SLOWPOKE_DELAY_MICROS_{service.upper()}"] = str(delay)
    env["SLOWPOKE_PRERUN"] = "false" # Disable request counting during normal execution
    cmd = f"bash run.sh {REQUEST}"
    print(f"[test.py] Running {REQUEST} request with the following configuration: {service_delay}", flush=True)
    process = subprocess.Popen(cmd, shell=True, env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    print(f"[test.py] Executing {cmd}:")
    throughput = -1
    for line in process.stdout:
        line_output = line.decode().strip()
        if "Requests/sec" in line_output:
            throughput = float(line_output.split()[1])
        print(f"    {line.decode().strip()}", flush=True)
    if process.wait() != 0:
        print(f"Error running {cmd}")
        for line in process.stderr:
            print(f"    {line.decode().strip()}", flush=True)
        raise Exception(f"Error running {cmd}")
    return throughput

def run():
    service_delay = copy.deepcopy(BASELINE_SERVICE_PROCESSING_TIME)
    processing_time_diff = TARGET_PROCESSING_TIME_RANGE[1]-TARGET_PROCESSING_TIME_RANGE[0]
    processing_time_range = range(TARGET_PROCESSING_TIME_RANGE[0], TARGET_PROCESSING_TIME_RANGE[1], processing_time_diff//TARGET_NUM_EXP)

    # groundtruth
    groundtruth = []
    for p_t in processing_time_range:
        service_delay[TARGET_SERVICE] = p_t
        res = exp(service_delay)
        while int(res) == 0:
            print("[test.py] Found 0 throughtput, rerun experiment")
            res = exp(service_delay)
        groundtruth.append(res)
    
    # slowdown
    slowdown = []
    predicted = []
    for p_t in processing_time_range:
        for service in service_delay:
            if service == TARGET_SERVICE:
                service_delay[service] = TARGET_PROCESSING_TIME_RANGE[1]
            else:
                if REQUEST_RATIO[service] == 0:
                    delay = 0
                else:
                    delay = int(
                        ((TARGET_PROCESSING_TIME_RANGE[1] - p_t)*REQUEST_RATIO[TARGET_SERVICE]) / REQUEST_RATIO[service]
                    )
                service_delay[service] = BASELINE_SERVICE_PROCESSING_TIME[service] + delay
        res = exp(service_delay)
        while int(res) == 0:
            print("[test.py] Found 0 throughput, rerun experiment")
            res = exp(service_delay)
        slowdown.append(res)

        try:
            delay = (TARGET_PROCESSING_TIME_RANGE[1] - p_t)*REQUEST_RATIO[TARGET_SERVICE]
            predicted_throughput = 1000000/(1000000/slowdown[-1]-delay)
        except:
            print("[test.py] Error: Division by zero")
            predicted_throughput = -1
        predicted.append(predicted_throughput)
    
    err = [(predicted[i]-groundtruth[i])*100/groundtruth[i] for i in range(len(predicted))]
    print("[test.py] Groundtruth: ", groundtruth, flush=True)
    print("[test.py] Slowdown: ", slowdown, flush=True)
    print("[test.py] Predicted: ", predicted, flush=True)
    print("[test.py] Error percentage: ", err, flush=True)

    return groundtruth, slowdown, predicted, err

if __name__ == "__main__":
    try:
        TARGET_SERVICE = os.sys.argv[1]
        REQUEST = os.sys.argv[2]
    except: 
        print("[test.py] No target service and request specified, using default values")
    print(f"[test.py] Running experiment for {TARGET_SERVICE}...")
    print(f"[test.py] REQUEST_RATIO: {REQUEST_RATIO}")
    print(f"[test.py] BASELINE: {BASELINE_SERVICE_PROCESSING_TIME}")
    print(f"[test.py] TARGET_PROCESSING_TIME_RANGE: {TARGET_PROCESSING_TIME_RANGE}")
    print(f"[test.py] TARGET_NUM_EXP: {TARGET_NUM_EXP}")
    print(f"[test.py] REQUEST: {REQUEST}")
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    groundtruths, slowdowns, predicteds, errs = [], [], [], []
    for i in range(3):
        print(f"[test.py] Running experiment {i}...")
        groundtruth, slowdown, predicted, err = run()
        groundtruths.append(groundtruth)
        slowdowns.append(slowdown)
        predicteds.append(predicted)
        errs.append(err)
    print("[test.py] Summary: ")
    for i in range(len(groundtruths)):
        print(f"[test.py] Result for the experiment {i}: ")
        print(f"    Groundtruth: {groundtruths[i]}")
        print(f"    Slowdown:    {slowdowns[i]}")
        print(f"    Predicted:   {predicteds[i]}")
        print(f"    Error Perc:  {errs[i]}")
