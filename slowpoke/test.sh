# env

env | grep "PROCESSING"

benchmark=${1:-boutique}
request=${2:-home}
thread=${3:-16}
conn=${4:-512}
# duration=${5:-60}
duration=0

YAML_PATH=$benchmark/yamls
if [[ $benchmark == "syncthetic" ]]; then
    YAML_PATH=$benchmark/$request/yamls
fi


# envsubst < $YAML_PATH/service0.yaml 
# envsubst < $YAML_PATH/service0.yaml | kubectl apply -f -
for file in $YAML_PATH/*;
do
    envsubst < $file | kubectl apply -f -
done