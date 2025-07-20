#!/bin/bash

cd $(dirname $0)
mkdir -p results
bash ./slowpoke/boutique/run-boutique-tiny.sh results/boutique_tiny.log
