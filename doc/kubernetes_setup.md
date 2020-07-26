# kubernetes集群搭建(centos7)

* 修改每台主机的hostname，如果需要的话
	* `hostname <hostname>`
	* 修改/etc/hostname

* 选择一台机器安装ansible，为了便于从一台机器上操作所有机器
	* 安装zsh & oh-my-zsh，为了更方便的使用命令行（可选）

		```
		yum install -y zsh
		yum install -y git
		sh -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
		```
		
	* `yum install -y ansible`
	* 解决错误 `RequestsDependencyWarning: urllib3 (1.22) or chardet (2.2.1) doesn't match a supported version`
 
		```
		pip uninstall -y urllib3
		pip uninstall -y chardet
		pip install requests
		```

	* 禁用command_warnings，在/etc/ansible/ansible.cfg里将`command_warnings = False`前面的#去掉
	* 将所有机器的内网ip按照分组增加到/etc/ansible/hosts，如下:

		```
		[master]
		172.20.102.[208:210]

		[node]
		172.20.102.[211:212]
		```
		
	* 用root账号通过ssh-keygen生成内网无需密码root登录其它服务器，使用默认选项
	* 用ssh-copy-id将生成的id_rsa.pub传送到所有主机的authorized_hosts里，包括本机，如：

		`ssh-copy-id root@172.20.102.208`
	* 验证ansible是否可以登录所有服务器，如下:

		```
		[root@172 ~]# ansible all -m ping -u root
		172.20.102.208 | SUCCESS => {
		    "changed": false, 
		    "ping": "pong"
		}
		...
		```
		
* 更新所有服务器

	`ansible all -u root -m shell -a "yum update -y"`
	
* 所有服务器上安装docker

	```
	ansible all -u root -m shell -a "yum remove docker docker-client docker-client-latest docker-core docker-latest docker-latest-logrotate docker-logrotate docker-selinux docker-engine-selinux docker-engine"
	ansible all -u root -m shell -a "yum install -y yum-utils"
	ansible all -u root -m shell -a "yum install -y device-mapper-persistent-data"
	ansible all -u root -m shell -a "yum install -y lvm2"
	ansible all -u root -m shell -a "yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo"
	ansible all -u root -m shell -a "yum install -y docker-ce"
	ansible all -u root -m shell -a "systemctl enable docker"
	ansible all -u root -m shell -a "systemctl start docker"
	```
	
* 每台机器上添加阿里云的kubernetes repo

	```
	# cat k8srepo.yaml
	cat <<EOF > /etc/yum.repos.d/kubernetes.repo
	[kubernetes]
	name=Kubernetes
	baseurl=http://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
	enabled=1
	gpgcheck=0
	repo_gpgcheck=0
	gpgkey=http://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg
   	       http://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
	EOF

	# ansible-playbook k8srepo.yaml
	```
	
* 安装kubelet, kubeadm, kubectl, ipvsadm

	`ansible all -u root -m shell -a "yum install -y kubelet kubeadm kubectl ipvsadm"`
	
* 禁用所有服务器上的swap

	`ansible all -u root -m shell -a "swapoff -a"`
	
* 允许所有服务器进行转发，因为k8s的NodePort需要在所有服务器之间进行转发

	`ansible all -u root -m shell -a "iptables -P FORWARD ACCEPT"`

* 由于k8s.gcr.io不能访问，需要从本机科学上网docker pull如下几个image

	```
	k8s.gcr.io/kube-proxy-amd64:v1.11.1
	k8s.gcr.io/kube-controller-manager-amd64:v1.11.1
	k8s.gcr.io/kube-scheduler-amd64:v1.11.1
	k8s.gcr.io/kube-apiserver-amd64:v1.11.1
	k8s.gcr.io/coredns:1.1.3
	k8s.gcr.io/etcd-amd64:3.2.18
	k8s.gcr.io/pause:3.1
	```
	
	通过命令一次完成拉取
	
	`while IFS= read -r line; do docker pull $line; done`
	
	然后上传到一台服务器

	`while IFS= read -r line; do docker save <image> | pv | ssh <user>@<server ip> "docker load"; done`
	
	同步到所有k8s服务器，其中$i是为了匹配所有内网ip
	
	`while IFS= read -r line; do for ((i=208;i<213;i++)); do docker save $line | ssh root@172.20.102.$i "docker load"; done; done`

* 在一台master服务器上初始化集群

	`kubeadm init --api-advertise-addresses <本机内网ip> --kubernetes-version=v1.11.1`
	
	注意最后的`kubeadm join`一行，用来在其它服务器加入集群（稍后用）
	
	初始化配置，master上执行如下命令
	
	```
	mkdir -p $HOME/.kube
	sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
	sudo chown $(id -u):$(id -g) $HOME/.kube/config
	```
	
	添加calico网络
	
	`kubectl apply -f https://docs.projectcalico.org/v3.1/getting-started/kubernetes/installation/hosted/kubeadm/1.7/calico.yaml`
	
* 从所有其它服务器执行master上获得的kubeadm join那行命令，里面包含了加入的token

* 执行`kubectl get nodes`验证集群是否成功
