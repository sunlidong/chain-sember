#!/usr/bin/env bash

### 清空证书文件以及 配置文件
#rm -rf crypto-config/
#rm -rf channel-artifacts/*
SLEEP_SECOND=10
### 生成证书文件
cryptogen generate --config=./crypto-config.yaml
sleep $SLEEP_SECOND
### 生成创世块文件
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block

### 生成通道配置文件
configtxgen -profile TwoOrgsChannelByskt -outputCreateChannelTx ./channel-artifacts/mychannel.tx -channelID mychannelbyskt


### 生成锚节点文件
configtxgen -profile TwoOrgsChannelByskt -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannelbyskt -asOrg Org1MSP
configtxgen -profile TwoOrgsChannelByskt -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID mychannelbyskt -asOrg Org2MSP


### 启动网络
docker-compose -f docker-compose-cli.yaml -p demo up -d


echo "create cert is successful"

