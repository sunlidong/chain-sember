#!/usr/bin/env bash

### 生成配置文件
$ cryptogen generate --config=./crypto-config.yaml

### 生成创世块文件
$ configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
### 生成通道配置文件
$ configtxgen -profile TwoOrgsChannelByskt -outputCreateChannelTx ./channel-artifacts/mychannel.tx -channelID mychannelbyskt

### 生成锚节点文件
$ configtxgen -profile TwoOrgsChannelByskt -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannelbyskt -asOrg Org1MSP
$ configtxgen -profile TwoOrgsChannelByskt -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID mychannelbyskt -asOrg Org2MSP

### 启动网络
$ docker-compose -f docker-compose-cli.yaml -p demo up -d

### 进入 容器
$ docker exec -it cli bash

### 查询类

### 以下子命令所操作的peer可以通过环境变量查看(操作在容器中操作)
$ echo $CORE_PEER_ADDRESS
$ echo $CORE_PEER_MSPCONFIGPATH
$ echo $ CORE_PEER_LOCALMSPID

### 查看peer所加入的channel
$ peer channel list

### 查看channel信息
$ peer channel getinfo -c mychannel # 这里的mychannel为channel的名称

### 查看peer已安装的chaincode
$ peer chaincode list --installed

### 查看某channel上已实例化的chaincode
$ peer chaincode list --instantiated -C mychannel # 这里的mychannel为channel的名称

### 查看peer的运行状态
$ peer node status

### 查看容器日志
$ docker logs  容器ID

