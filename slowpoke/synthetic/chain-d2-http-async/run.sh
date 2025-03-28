#!/bin/bash
#!/bin/bash

cd $(dirname $0)/../..

exp=chain-d2-http-async

mkdir -p synthetic/$exp/one-service-per-node-results

# Make it reproducible
target_service_random_pairs="0:13976 1:9944 2:1762"
# target_service_random_pairs="0:4446"

for pair in $target_service_random_pairs
do 
    target_service=$(echo $pair | cut -d':' -f1)
    random_seed=$(echo $pair | cut -d':' -f2)

    output_file=synthetic/$exp/one-service-per-node-results/$exp-service$target_service-8-512.log
    
    if [[ -e $output_file ]]; then
        echo "File $output_file already exists. Skipping..."
        continue
    fi

    touch $output_file
    
    python3 test.py -b synthetic \
        -r $exp \
        -x service$target_service \
        --num_exp 10 \
        -c 512 \
        -t 8 \
        --num_req 20000 \
        --clien_cpu_quota 2 \
        --random_seed $random_seed \
        --repetition 3 \
        >$output_file
done
