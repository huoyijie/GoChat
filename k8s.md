## 关闭Swap

```bash
sudo swapoff -a
# 取消挂载swap分区，在swap行首添加# 
例如 #/swap.img      none    swap    sw      0       0
sudo vim /etc/fstab
```

## 安装Docker

```bash
# 更新安装源列表
sudo apt-get update
# 安装前提必须文件
sudo apt-get install ca-certificates curl gnupg lsb-release
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
# 添加docker源
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
# 安装dokcer和容器运行时等
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin
# 运行docker，并加入开机启动
sudo systemctl enable --now docker
# 测试安装是否成功，如果安装成功会看到hello world输出
sudo docker run hello-world
```

## 使用阿里镜像源安装Kubeadm

```bash
sudo apt-get install -y apt-transport-https

sudo curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add - 

sudo cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF

sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo systemctl enable --now kubelet
```

## 安装Kubernetes

### 生成并编辑containerd配置文件

```bash
containerd config default | sudo tee /etc/containerd/config.toml
# 使用阿里源替换不可访问的国外源
sudo sed -i 's/registry.k8s.io/registry.aliyuncs.com\/google_containers/g' /etc/containerd/config.toml
```

在config.toml的[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]条目下，将SystemdCgroup = false改为SystemdCgroup = true
注意，config.toml有多处SystemdCgroup配置，不要改错了位置。

```bash
sudo systemctl enable containerd
sudo systemctl restart containerd
```

### 使用Kubeadm部署Kubernetes

```bash
# 根据环境配置你的--pod-network-cidr的值，不能与已有的网络重复
sudo kubeadm init --image-repository registry.aliyuncs.com/google_containers --pod-network-cidr=10.10.0.0/16
```

### 复制配置文件到当前用户目录

```bash
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

### 安装Calico

```bash
kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.24.5/manifests/tigera-operator.yaml
wget https://raw.githubusercontent.com/projectcalico/calico/v3.24.5/manifests/custom-resources.yaml
# 编辑custom-resources.yaml，确保cidr配置的值与之前--pod-network-cidr的值相同
kubectl create -f custom-resources.yaml
```

如果你无法访问以上两个Github文件，可以从这里下载

### 查看结果

以下命令需要显示所有的pod处于Running状态

```bash
watch kubectl get -A pods
```

### 安装Dashboard

如果只有一个节点，那么安装Dashboard前需要运行以下命令

```bash
kubectl taint nodes --all node-role.kubernetes.io/control-plane-
```

### 开始部署Dashboard

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml
```

[参考](https://www.jianshu.com/p/7eedcaa05b02)