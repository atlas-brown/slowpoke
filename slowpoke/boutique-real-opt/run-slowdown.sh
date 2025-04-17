cd $(dirname $0)

DIR=/home/ubuntu/mucache/slowpoke/boutique-real-opt/results
for i in $(seq 1 10); do
    FILE=slowdown-$i.log
    if [ ! -d "$DIR" ]; then
    mkdir -p $DIR
    fi
    if [ -f "$DIR/$FILE" ]; then
    echo "$DIR/$FILE exists."
    continue
    fi
    source <(python3 env.py baseline-$i.log)
    env | grep SLOWPOKE > $DIR/$FILE
    bash /home/ubuntu/mucache/slowpoke/run.sh boutique mix 8 1024 100000 >> $DIR/$FILE
done