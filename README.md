# chain-sember
###### fabric 案例


### 项目介绍

#### 1.  目录结构介绍

├── base ********* docker继承的yaml文件

├── build.sh ********* 初始化网络脚本(创建通道，节点加入通道,安装链码，实例化链码，初始化链码)

├── chaincode ********* 链码(存放链码的位置)

├── channel-artifacts ********* 配置文件存放目录(genesis.block,channel.tx,锚节点文件)

├── clear.sh ********* 清除容器以及挂在卷脚本(清除环境)

├── configtx.yaml ********* configtx配置文件(生成 创世快配置文件，通道配置文件，以及锚节点)

├── crypto-config ********* 证书目录

├── crypto-config.yaml ********* 证书配置文件(生成证书,网络的基石，所有的开始)

├── docker-compose-cli.yaml ********* Docker 容器配置文件( 1 orderer 4 peer)

├── init.sh ********* 生成配置文件脚本(生成:证书文件，创世块文件，通道配置文件,锚节点文件)

└── scripts ********* 脚本目录(目前为空，暂时没用到)


#### 2.  重要参数信息

##### 网络结构： 2个组织 ， 一个组织各有2个节点
##### 组织名称 ： huawei xiaomi
##### 通道名称：mychannelbyskt
##### 链码名称：AssetToChain_realty
##### 链码版本：1.0

#### 3. 容器名称信息

###### cli
###### orderer.cpu.com
###### peer0.org1.huawei.com
###### peer1.org1.huawei.com
###### peer0.org2.xiaomi.com
###### peer1.org2.xiaomi.com
 
 

#### 4. 启动顺序 

##### 启动前提
1. 执行过clear.sh 确保环境干净
2. 删除  channel-artifacts 目录下的配置文件
3. 删除  crypto-config 目录下的证书
4. 所有执行操作在sudo模式下 

##### 启动

1.  执行 init.sh,生成证书文件以及配置文件 ： ./init.sh
2.  执行 build.sh 初始化网络 ：./build.sh
3.  查看网络 ：docker ps -a

##### 清空网络

1. 执行清空网络脚本 ：./clear.sh(以下2步骤可以写到clear.sh脚本中)
2. 删除  channel-artifacts 目录下的配置文件
3. 删除  crypto-config 目录下的证书
