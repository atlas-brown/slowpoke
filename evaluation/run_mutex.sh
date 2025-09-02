#!/bin/bash

export SLOWPOKE_TOP=${SLOWPOKE_TOP:-$(cd "${BASH_SOURCE%/*}/.." && pwd -P)}

kubectl delete deployments --all
kubectl delete services --all

cd $(dirname $0)
mkdir -p results
time bash mutex/run-no-lock.sh > results/mutex-microbenchmark-no-lock.log
time bash mutex/run-lock.sh > results/mutex-microbenchmark-lock.log

outdir=$(realpath ./results)
draw_script=$(realpath ./draw-mutex.py)

echo "The results are stored in ${outdir}"
echo "To visualize the results, run: "
echo ""
echo "python3 ${draw_script} ${outdir}"
