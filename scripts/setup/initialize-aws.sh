#!/bin/bash
set -euo pipefail

CLUSTER_DIR=$(realpath $1)
DIR="$(dirname "$(realpath "$0")")"
cd $DIR
key=${CLUSTER_DIR}/slowpoke-expr.pem
user=ubuntu
nodes=$(cat ${CLUSTER_DIR}/ec2_ips)
counter=0

for node in $nodes; do
	if [[ $counter -eq 0 ]]; then
		control=$user@$node
		echo "Copying files to node $node"
		scp -o StrictHostKeyChecking=no -i $key init_control.sh $user@$node:/home/$user
		
		echo "Initializing node $node"
		ssh -tt -o StrictHostKeyChecking=no -i $key $user@$node bash init_control.sh
		echo "done with $counter"
	else
		echo "Copying files to node $node"
		scp -o StrictHostKeyChecking=no -i $key init_worker.sh $user@$node:/home/$user
		
		echo "Initializing node $node"
		ssh -tt -o StrictHostKeyChecking=no -i $key $user@$node bash init_worker.sh &
	fi	
	echo ""
	counter=$((counter+1))
done
wait

TOKEN_PREFIX=$(ssh -tt -o StrictHostKeyChecking=no -i $key $control sudo kubeadm token create --print-join-command)
counter=0

for node in $nodes; do
	if [[ $counter -ne 0 ]]; then
		worker_counter=$((counter-1))
		TOKEN_SUFFIX="--node-name worker$worker_counter --cri-socket unix:///var/run/cri-dockerd.sock"
		TOKEN=$(echo "sudo $TOKEN_PREFIX $TOKEN_SUFFIX" | sed 's/\r//')
		
		echo "Node $node is joining the cluster"
		ssh -tt -o StrictHostKeyChecking=no -i $key $user@$node eval "$TOKEN"
	fi	
	counter=$((counter+1))
done

echo "Checking nodes in the control plane"
ssh -tt -o StrictHostKeyChecking=no -i $key $control kubectl get nodes
