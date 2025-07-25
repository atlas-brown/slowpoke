#!/bin/bash

export SLOWPOKE_TOP=${SLOWPOKE_TOP:-$(dirname "$(realpath "$0")")/..}

cd $(dirname $0)
mkdir -p results
time bash boutique/run-boutique-medium.sh results/boutique_medium.log
time bash hotel/run-hotel-medium.sh results/hotel_medium.log
time bash social/run-social-medium.sh results/social_medium.log
time bash movie/run-movie-medium.sh results/movie_medium.log

outdir=$(realpath ./results)
draw_script=$(realpath ./draw.py)

echo "The results are stored in ${outdir}"
echo "To visualize the results, run: "
echo ""
echo "python3 ${draw_script} ${outdir}"
