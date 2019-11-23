# chain-sember
###### fabric 案例


### 项目介绍

#### 1.  结构介绍

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

