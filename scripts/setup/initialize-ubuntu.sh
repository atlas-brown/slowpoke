#!/bin/bash

nodes="yizhengx@ms0236.utah.cloudlab.us"
counter=0

for node in $nodes; do
	if [[ $counter -eq 0 ]]; then
		control=$node
		echo "Copying files to node $node"
		scp -o StrictHostKeyChecking=no init_control.sh $node:/users/yizhengx
		
		echo "Initializing node $node"
		ssh -tt -o StrictHostKeyChecking=no $node bash init_control.sh
		# ssh -tt -o StrictHostKeyChecking=no $node 'git clone https://<github_token>@github.com/CASP-Systems-BU/causal-profiling-microbenchmark.git && git checkout <your_branch>'
		echo "done with $counter"
	else
		echo "Copying files to node $node"
		scp -o StrictHostKeyChecking=no init_worker.sh $node:/users/yizhengx
		
		echo "Initializing node $node"
		ssh -tt -o StrictHostKeyChecking=no $node bash init_worker.sh
	fi	
	echo ""
	counter=$((counter+1))
done

TOKEN_PREFIX=$(ssh -tt -o StrictHostKeyChecking=no $control sudo kubeadm token create --print-join-command)
counter=0

for node in $nodes; do
	if [[ $counter -ne 0 ]]; then
		TOKEN_SUFFIX="--node-name worker$counter --cri-socket unix:///var/run/cri-dockerd.sock"
		TOKEN=$(echo "sudo $TOKEN_PREFIX $TOKEN_SUFFIX" | sed 's/\r//')
		
		echo "Node $node is joining the cluster"
		ssh -tt -o StrictHostKeyChecking=no $node eval "$TOKEN"
	fi	
	counter=$((counter+1))
done

echo "Checking nodes in the control plane"
ssh -tt -o StrictHostKeyChecking=no $control kubectl get nodes