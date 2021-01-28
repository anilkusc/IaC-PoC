#!/bin/bash 
setenforce 0 && sed -i "s/^SELINUX\=enforcing/SELINUX\=disabled/g" /etc/selinux/config && sed -i "s/^SELINUX\=permissive/SELINUX\=disabled/g" /etc/selinux/config && systemctl disable firewalld; systemctl stop firewalld; systemctl mask firewalld
yum install openssl-devel  gcc  libffi-devel  wget  nano  git  sshpass  python3 python3-pip  python3-devel -y
mkdir -p /root/master/kubespray
git clone https://github.com/kubernetes-incubator/kubespray.git /root/master/kubespray
pip3 install -r /root/master/kubespray/requirements.txt
ansible-playbook --connection=local -i ./files/hosts.yaml ./files/keys.yaml
ansible-playbook --private-key=/root/master/id_rsa  -i /root/master/files/hosts.yaml /root/master/files/playbook.yaml
ansible-playbook --private-key=/root/master/id_rsa -i /root/master/kubespray/inventory/sample/inventory.ini  --become --become-user=root /root/master/kubespray/cluster.yml
ansible-playbook --private-key=/root/master/id_rsa -i /root/master/files/hosts.yaml /root/master/files/apply.yaml
ansible-playbook --private-key=/root/master/id_rsa -i /root/master/files/hosts.yaml  /root/master/files/federation-consul.yaml   
ansible-playbook --private-key=/root/master/id_rsa -i /root/master/files/hosts.yaml  /root/master/files/elk.yaml
cd ./controller
bash ./generate.sh --namespace default --service nslimiter --secret nslimiter
mkdir /mutate && cp -fr ../controller/* /mutate/
kubectl apply -f manifest.yaml