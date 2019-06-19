[TOC]

# ZDHsys

Platform for automatic installation of operating system based on golang and saltstack

The main implementation of binding MAC address, independent installation of the operating system, change the host name, IP address

The project is still being updated...

## Installation

### Installing golang(linux)

- Download address：[https://golang.google.cn/dl/](https://golang.google.cn/dl/)

- package ：go1.12.6.linux-amd64.tar.gz

- uncompress

  ```shell
  # tar -C /usr/local/ -zxf go1.11.2.linux-amd64.tar.gz
  ```

- Setting environment variables

  ```shell
  # echo export PATH=$PATH:/usr/local/go/bin/ >> /etc/profile
  # source /etc/profile
  # go env\
  ```

### Installing beego

```shell
# yum install -y git
# go get github.com/astaxie/beego
# go get github.com/beego/bee
# go get github.com/astaxie/beego/orm
# echo export PATH=$PATH:/root/go/bin >> /etc/profile
# source /etc/profile
# bee
```

### Installing other packages

- Json parsing library
- mysql driver
- call python for golang

```shell
# go get github.com/bitly/go-simplejson
# go get github.com/go-sql-driver/mysql
# yum install -y python-devel
# go get github.com/sbinet/go-python
# mkdir -p /root/go/src/golang.org/x
# cd $GOPATH/src/golang.org/x
# git clone https://github.com/golang/crypto.git
```

### Installing mysql

- download：[mysql](https://dev.mysql.com/downloads/mysql/) 

- create a remote login user

  ```mysql
  mysql> grant all on *.* to "root"@"192.168.100.111" identified by "WWW.1.com";
  mysql> flush privileges;
  ```

- create database

  ```mysql
  mysql> create database test charset=utf8;
  ```

### Installing cobbler 

- slightly

### Installing ZDHsys

- attention：The path of ZDHsys is  `/root/go/src/ZDHsys`

- run the script

  ```shell
  bash /root/go/src/ZDHsys/setup.sh
  ```

- register and login

  ![注册](http://chuantu.xyz/t6/702/1560917837x2918527082.png)

  ![登陆](http://chuantu.xyz/t6/702/1560917934x2362407012.png)

- create  server（Bind system to MAC address）

  ![录入主机](http://chuantu.xyz/t6/702/1560925278x3752237043.png)

  ![mac](http://chuantu.xyz/t6/702/1560925396x3752237043.png)

  ![sss](http://chuantu.xyz/t6/702/1560926226x1709417317.png)

  ![](http://chuantu.xyz/t6/702/1560926259x1033347913.png)

  ![](http://chuantu.xyz/t6/702/1560926288x1709417317.png)

- then  a system will be create in the cobbler
