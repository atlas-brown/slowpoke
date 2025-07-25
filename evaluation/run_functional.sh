#!/bin/bash

export SLOWPOKE_TOP=${SLOWPOKE_TOP:-$(dirname "$(realpath "$0")")/..}

cd $(dirname $0)
mkdir -p results
bash boutique/run-boutique-tiny.sh results/boutique_tiny.log