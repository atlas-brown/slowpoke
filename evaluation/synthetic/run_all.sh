#!/bin/bash

cd $(dirname $0)

for dir in $(ls -d */); do
    cd "$dir"

    if [ -f "run.sh" ]; then
        echo "Running $dir"
        time ./run.sh
    fi

    cd ..
done