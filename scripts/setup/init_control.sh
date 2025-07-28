#!/bin/bash
CRI_DOCKER_VER=0.3.1

docker_install () {
    # Add Docker's official GPG key:
    sudo apt-get update
    sudo apt-get -y install ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    # Add the repository to Apt sources:
    echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
    $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
    sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # Not needed for Ubuntu
# sudo systemctl start docker
}

cri_dockerd_install () {
    wget https://github.com/Mirantis/cri-dockerd/releases/download/v${CRI_DOCKER_VER}/cri-dockerd-${CRI_DOCKER_VER}.amd64.tgz
    tar xvf cri-dockerd-${CRI_DOCKER_VER}.amd64.tgz
    sudo mv cri-dockerd/cri-dockerd /usr/local/bin/

    wget https://raw.githubusercontent.com/Mirantis/cri-dockerd/master/packaging/systemd/cri-docker.service
    wget https://raw.githubusercontent.com/Mirantis/cri-dockerd/master/packaging/systemd/cri-docker.socket
    sudo mv cri-docker.socket cri-docker.service /etc/systemd/system/
    sudo sed -i -e 's,/usr/bin/cri-dockerd,/usr/local/bin/cri-dockerd,' /etc/systemd/system/cri-docker.service

    sudo systemctl daemon-reload
    sudo systemctl enable cri-docker.service
    sudo systemctl enable --now cri-docker.socket
}

kube_install () {
    sudo apt-get update
    # apt-transport-https may be a dummy package; if so, you can skip that package
    sudo apt-get install -y apt-transport-https ca-certificates curl
    # If the folder `/etc/apt/keyrings` does not exist, it should be created before the curl command, read the note below.
    # sudo mkdir -p -m 755 /etc/apt/keyrings
    curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.29/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
    # This overwrites any existing configuration in /etc/apt/sources.list.d/kubernetes.list
    echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.29/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
    sudo apt-get update
    sudo apt-get install -y kubectl kubelet kubeadm
    sudo systemctl enable --now kubelet
}

istio_install () {
    curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.18.0  sh -
    cd istio* && cd bin
    sudo mv istioctl /usr/local/bin
    cd ../../
    rm -rf istio*
}

init_cluster () {
    sudo swapoff -a
    sudo kubeadm init --cri-socket unix:///var/run/cri-dockerd.sock --pod-network-cidr=192.168.0.0/16
    mkdir -p $HOME/.kube
    sudo cp /etc/kubernetes/admin.conf $HOME/.kube/config
    sudo chown "$(id -u):$(id -g)" $HOME/.kube/config


    # # Setup networking with Calico
    # kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.25.0/manifests/tigera-operator.yaml
    # kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.25.0/manifests/custom-resources.yaml
    kubectl apply -f https://github.com/weaveworks/weave/releases/download/v2.8.1/weave-daemonset-k8s.yaml
}

init_dependencies () {
    sudo apt-get -y install python3
    sudo apt-get -y install python3-pip python3-matplotlib python3-progress
    sudo apt-get -y install screen
    sudo apt-get -y install nginx
    # enable anyone to drop things in shared directory
    sudo chmod 777 /var/www/html
    # sudo pip3 install --upgrade pip
    # sudo pip3 install pandas
    # sudo python3 -m pip install --upgrade Pillow
    # sudo pip3 install matplotlib
    # sudo pip3 install progress
    kubectl apply -f https://raw.githubusercontent.com/pythianarora/total-practice/master/sample-kubernetes-code/metrics-server.yaml
}

echo '* libraries/restart-without-asking boolean true' | sudo debconf-set-selections
docker_install
cri_dockerd_install
kube_install
sudo modprobe br_netfilter
sudo systemctl stop firewalld 
init_cluster
istio_install
init_dependencies
