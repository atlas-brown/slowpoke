#!/bin/bash

export SLOWPOKE_TOP=${SLOWPOKE_TOP:-$(cd "${BASH_SOURCE%/*}/.." && pwd -P)}

cd $(dirname $0)
mkdir -p results
kubectl delete deployments --all
kubectl delete services --all
bash boutique/run-boutique-tiny.sh results/boutique_tiny.log