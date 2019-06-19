#!/bin/bash
#

systemctl stop firewalld
sed -i '/^SELINUX=/c \SELINUX=disabled' /etc/selinux/config

echo "正在安装软件请稍后！！！"

yum install -y epel-release &>/dev/null
yum install -y net-tools python python26 python-jinja2 salt-master salt-api salt-minion pyOpenSSL &>/dev/null

clear

if [ ! -f /etc/salt/master ];then
clear
read -t 5 -p "salt-master 没有安装成功，请检查安装源（需要连接网络）！"
return
fi

rpm -ql salt-master > /dev/null
if [ $? = 0 ];then
read -t 5 -p "salt-master 已经安装成功！正在安装salt-minion"
clear
fi


if [ ! -f /etc/salt/minion ];then
clear
read -t 5 -p "salt-minion 没有安装成功，请检查安装源（需要连接网络）！"
return
fi

rpm -ql salt-minion >/dev/null
if [ $? = 0 ];then
read -t 5 -p "salt-minion 已经安装成功！"
fi


sed -i '/#interface:/s/^#//' /etc/salt/master
sed -i '670s/^#//' /etc/salt/master
sed -i '671s/^#//' /etc/salt/master
sed -i '406s/^#//' /etc/salt/master
sed -i '407s/^#//' /etc/salt/master
sed -i '408s/^#//' /etc/salt/master
sed -i '409s/^#//' /etc/salt/master
sed -i '/#     - \/srv\/salt\/dev\/services/s/#     - \/srv\/salt\/dev\/services/     - \/srv\/salt\/dev\//' /etc/salt/master
ipnum=`ifconfig | grep "inet"|sed -n 1p |awk '{print $2}'`
sed -i "/^interface/s/0.0.0.0/$ipnum/" /etc/salt/master
sed -i '/#auto_accept/s/^#//' /etc/salt/master
sed -i '/auto_accept/s/False/True/' /etc/salt/master
sed -i '/#default_include/s/#default/default/g' /etc/salt/master

[ ! -d /etc/salt/master.d ] && mkdir /etc/salt/master.d

cat /etc/passwd|grep salter &>/dev/null
if [ $? != 0 ];then
useradd -M -s /sbin/nologin salter
echo "salter" | passwd salter --stdin
fi
echo "salter用户创建完毕！"

cat << EOF > /etc/salt/master.d/api.conf
rest_cherrypy:
  port: 8000
  ssl_crt: /etc/pki/tls/certs/localhost.crt
  ssl_key: /etc/pki/tls/certs/localhost.key
EOF
cat << EOF > /etc/salt/master.d/eauth.conf
external_auth:
  pam:
    saltss:
      - .*
      - '@wheel'
      - '@runner'
EOF

systemctl start salt-master
systemctl enable salt-master

echo "master已经启动！"

read -p "请设置master IP：" serverip
sed -i '/#master:/s/^#//' /etc/salt/minion
sed -i "/^master:/s/salt/$serverip/" /etc/salt/minion
read -p "设定客户端名称编号：" clientip
sed -i '/#id/s/^#//' /etc/salt/minion
sed -i "/^id/s/id:/id: $clientip/" /etc/salt/minion
sed -i '/#log_file/s/^#//' /etc/salt/minion
sed -i '/#key_logfile/s/^#//' /etc/salt/minion

systemctl start salt-minion
systemctl enable salt-minion

echo "minion已经启动！"


CERT(){
echo "开始生成自签名证书，请稍等......"
echo "master: $ipadd" >> /etc/salt/minion
echo "id: $hostname" >> /etc/salt/minion
sed -i '/#log_file/s/^#//' /etc/salt/minion
sed -i '/#key_logfile/s/^#//' /etc/salt/minion
systemctl restart salt-master &>/dev/null
systemctl restart salt-minion &>/dev/null
salt-call tls.create_self_signed_cert
if [ $? != 0 ];then
salt-call tls.create_self_signed_cert
fi
read -t 4 -p "证书生成完毕。"
clear
}


clear
ipnum=`ip ad|grep global|awk -F"[ /]+" '{print $3}'|wc -l`
if [ $ipnum = 1 ];then
ipadd=`ip ad|grep global|awk -F"[ /]+" '{print $3}'`
hostname=`hostname`
CERT

sed -i "s/192.168.100.111/$ipadd/g" /root/go/src/ZDHsys/conf/app.conf
sed -i "/salt_username/d" /root/go/src/ZDHsys/conf/app.conf
sed -i "/salt_password/d" /root/go/src/ZDHsys/conf/app.conf
echo "salt_username = salter" >> /root/go/src/ZDHsys/conf/app.conf
echo "salt_password = salter" >> /root/go/src/ZDHsys/conf/app.conf
sed -i "s/192.168.100.111/$ipadd/" /root/go/src/ZDHsys/models/model.go
cd /root/go/src/ZDHsys
bee run



else
ip ad|grep global|awk -F"[ /]+" '{print $3}'
read -p "系统有多个IP地址，请选择： " ipadd
CERT

sed -i "s/192.168.100.111/$ipadd/g" /root/go/src/ZDHsys/conf/app.conf
sed -i "/salt_username/d" /root/go/src/ZDHsys/conf/app.conf
sed -i "/salt_password/d" /root/go/src/ZDHsys/conf/app.conf
echo "salt_username = salter" >> /root/go/src/ZDHsys/conf/app.conf
echo "salt_password = salter" >> /root/go/src/ZDHsys/conf/app.conf
sed -i "s/192.168.100.111/$ipadd/" /root/go/src/ZDHsys/models/model.go

cd /root/go/src/ZDHsys
bee run

fi



systemctl restart salt-minion &>/dev/null
systemctl restart salt-master &>/dev/null
systemctl start salt-api &>/dev/null
if [ $? != 0 ];then
systemctl restart salt-api
fi
systemctl enable salt-api &>/dev/null
ps aux |grep salt-master
ps aux |grep salt-api
ps aux |grep SST
echo "ZDHsys以安装成功，请在ZDHsys工作目录运行bee run,启动项目，打开http://$ipadd:8080 访问页面！"
clear


