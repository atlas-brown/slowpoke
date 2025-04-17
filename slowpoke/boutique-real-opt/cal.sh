cd $(dirname $0)

DIR=/home/ubuntu/mucache/slowpoke/social-real-opt/results
for i in $(seq 1 10); do
    echo "===="
    python3 env.py baseline-$i.log x $i
done


