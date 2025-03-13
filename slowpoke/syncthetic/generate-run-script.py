#!/usr/bin/env python3

import os
import sys

bash_script = """#!/bin/bash

cd $(dirname $0)/../..

exp=$FOLDER_NAME

mkdir -p syncthetic/$exp/results

target_services="$TARGET_SERVICE"
for target_service in $target_services
do 
    if [[ -e syncthetic/$exp/results/$exp-service$target_service.log ]]; then
        echo "File syncthetic/$exp/results/$exp-service$target_service.log already exists. Skipping..."
        continue
    fi
    touch syncthetic/$exp/results/$exp-service$target_service.log
    python3 test.py -b syncthetic \\
        -r $exp \\
        -x service$target_service \\
        --num_exp 10 \\
        -c 128 \\
        -t 2 \\
        --num_req 18000 \\
        --clien_cpu_quota 2 \\
        --random_seed $RANDOM \\
        >syncthetic/$exp/results/$exp-service$target_service.log
done
"""

def build_bash_script(folder_name, target_service):
    script = bash_script
    script = script.replace("$FOLDER_NAME", folder_name)
    script = script.replace("$TARGET_SERVICE", target_service)
    return script

special_params = {
    "chain-d2": [0, 1, 2],
    "chain-d4": [0, 2, 5],
    "chain-d8": [0, 5, 9],
    "fanout-w3": [0, 2],
    "fanout-w10": [0, 5],
    "dag-balanced": [0, 2, 12],
    "dag-unbalanced": [1, 4, 5],
    "dag-relay": [1, 4, 5],
    "dag-cross": [3, 2, 4],
    "dynamic-once": [1, 2],
    "dynamic-twice": [0, 3, 5],
    "dynamic-cache": [1, 2],
    "dynamic-cycle": [1, 2]
}

if __name__ == "__main__":

    # change to current directory
    os.chdir(os.path.dirname(os.path.abspath(__file__)))

    for dir_ in os.listdir("."):
        if not os.path.isdir(dir_):
            continue
        topology = "-".join(dir_.split("-")[:2])
        if topology not in special_params:
            continue
        script = build_bash_script(str(dir_), " ".join([str(i) for i in special_params[topology]]))
        if os.path.exists(f"{dir_}/run.sh"):
            print(f"❌ Already exists: syncthetic/{dir_}/run.sh")
            continue
        os.makedirs(f"{dir_}", exist_ok=True)
        with open(f"{dir_}/run.sh", "w") as f:
            f.write(script)
        os.chmod(f"{dir_}/run.sh", 0o755)
        print(f"✅Generated script for {dir_} with target services {special_params[topology]}")