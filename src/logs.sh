cd $(dirname $0)
rm -rf logs
mkdir -p logs
log() {
    res=$(kubectl get pods | cut -f 1 -d " ")
    IFS=$'\n' read -rd '' -a array <<< "$res"
    for value in "${array[@]:1:${#array[@]}-1}"
    do
        kubectl logs $value > logs/$value.log
    done
}
log