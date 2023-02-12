# GoChat

## v0.3 Release

* ui
* ...

### Run

* server

![gochat-server](docs/images/gochat-server.gif)

* client (huoyijie)

![gochat-huoyijie](docs/images/gochat-huoyijie.gif)

* client (jack)

![gochat-jack](docs/images/gochat-jack.gif)

### Features

![gochat-features-uml](docs/images/gochat-features-uml.svg)

### Diagrams

* lib

![gochat-lib-uml](docs/images/gochat-lib-uml.svg)

* server

![gochat-server-uml](docs/images/gochat-server-uml.svg)

* client

![gochat-client-uml](docs/images/gochat-client-uml.svg)

* sequence

![gochat-sequence-uml](docs/images/gochat-sequence-uml.svg)

## Docker

```bash
# work dir
cd server

# build executable
go build -o target/gochat-server

# build docker image
docker build -t huoyijie/gochat-server:latest .

# run docker c
docker run -it -v "$(pwd)"/target:/root/.gochat huoyijie/gochat-server:latest

# open container's shell
docker exec -it af2e58909af8 /bin/bash
```

# Kubeadm

* 安装 containerd/runc/cni-plugins

```bash
# The containerd.io package contains runc too, but does not contain CNI plugins.
sudo apt install -y containerd.io

# 安装 CNI plugins，从 https://github.com/containernetworking/plugins/releases 下载
sudo mkdir -p /opt/cni/bin
sudo tar Cxzvf /opt/cni/bin  cni-plugins-linux-amd64-v1.2.0.tgz

# 编辑 /etc/containerd/config.toml，取消禁止 cri 集成插件，并配置 systemd cgroup 驱动

#version = 2
#[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
#  SystemdCgroup = true
sudo systemctl restart containerd
```

* 安装 kubeadm/kubelet/kubectl

[参考这里](https://kubernetes.io/zh-cn/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)

切换软件源到国内镜像

* 初始化集群

```bash
sudo crictl --runtime-endpoint unix:///var/run/containerd/containerd.sock pull registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.6

sudo ctr -n k8s.io image tag registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.6 registry.k8s.io/pause:3.6

sudo kubeadm init  --v=6  --image-repository='registry.cn-hangzhou.aliyuncs.com/google_containers' --ignore-preflight-errors=all --pod-network-cidr=10.244.0.0/16
```

Your Kubernetes control-plane has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

Alternatively, if you are the root user, you can run:

  export KUBECONFIG=/etc/kubernetes/admin.conf

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

Then you can join any number of worker nodes by running the following on each as root:

kubeadm join 172.21.0.16:6443 --token fnegf9.vg14engc3y0llf6g \
	--discovery-token-ca-cert-hash sha256:30bcb402141023db100ce6a07eb59df6563ffc342568169cfb6c46511cac11cc

* 安装 Pod 网络插件

```bash
kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml
sudo systemctl restart kubelet containerd
```

## v0.4 todo

* tls
* emoji
* send file
* group chat