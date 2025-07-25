#!/bin/bash

cd $(dirname $0)
mkdir -p results
time bash ./slowpoke/boutique/run-boutique-medium.sh results/boutique_medium.log
time bash ./slowpoke/hotel/run-hotel-medium.sh results/hotel_medium.log
time bash ./slowpoke/social/run-social-medium.sh results/social_medium.log
time bash ./slowpoke/movie/run-movie-medium.sh results/movie_medium.log

outdir=$(realpath ./results)
draw_script=$(realpath ./draw.py)

echo "The results are stored in ${outdir}"
echo "To visualize the results, run: "
echo ""
echo "python3 ${draw_script} ${outdir}"
