#!/usr/bin/env bash

###  需要修改参数
GO_CC_NAME=("AssetToChain_realty" "AssetToChain_increment")
GO_CC_SRC_PATH=("github.com/chaincode/chaincodeRealty/dingchain"  "github.com/chaincode/chaincodeRealty/increment")
CC_VERSION="1.0"

############################## 参数列表 ##############################

### 通道名称
CHANNEL_NAME="mychannelbyskt"

### 通道.tx
CHANNEL_NAMETX="mychannel"
### 生成通道配置文件
#CHANNEL_NAMETX =>configtxgen -profile TwoOrgsChannelByskt -outputCreateChannelTx ./channel-artifacts/mychannel.tx -channelID mychannelbyskt

### 睡眠时间
SLEEP_SECOND=10


DOMAIN_NAME="com"

ORDERER_ADDRESS="orderer.cpu.com:7050"
ORG_NAME=("org1" "org2")
CC_VERSION="1.0"
TLS_PATH="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/"
ORDERER_TLS_PATH="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/"
ORDERER_CAFILE="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cpu.com/orderers/orderer.cpu.com/tls/ca.crt"


### org1
HUAWEI_TLSROOTCERTFILE="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.huawei.com/peers/peer0.org1.huawei.com/tls/ca.crt"

### org2
XIAOMI_TLSROOTCERTFILE="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.xiaomi.com/peers/peer0.org2.xiaomi.com/tls/ca.crt"
### 开始
echo
echo " ____    _____      _      ____    _____ "
echo "/ ___|  |_   _|    / \    |  _ \  |_   _|"
echo "\___ \    | |     / _ \   | |_) |   | |  "
echo " ___) |   | |    / ___ \  |  _ <    | |  "
echo "|____/    |_|   /_/   \_\ |_| \_\   |_|  "
echo
echo "Build your Server......."
echo
###

###
get_mspid() {
    local org=$1
    case "$org" in
        org1)
            echo "Org1MSP"
            ;;
        org2)
            echo "Org2MSP"
            ;;
        *)
            echo "error org name $org"
            exit 1
            ;;
    esac
}

get_msp_config_path() {

    local org=$1
    local peer=$2
    local com=$3

    if [[ "$org" = "Org1" ]] && [[ "$org" = "Org2" ]]; then
        echo "error org name $org"
        exit 1
    fi

    if [[ "$peer" = "peer0" ]] && [[ "$peer" = "peer1" ]]; then
        echo "error peer name $peer"
        exit 1
    fi

    echo "${TLS_PATH}peerOrganizations/$org.$com.com/users/Admin@$org.$com.com/msp"

}

get_peer_address() {
    local org=$1
    local peer=$2
    local port=$3
    local com=$4
    if [[ "$org" != "org1" ]] && [[ "$org" != "org2" ]]; then
        echo "error org name $org"
        exit 1
    fi

    echo "${peer}.${org}.${com}.${DOMAIN_NAME}:$port"
}

get_peer_tls_cert(){
    local org=$1
    local peer=$2
    local type=$3
    local com=$4
    if [[ "$org" != "org1" ]] && [[ "$org" != "org2" ]]; then
        echo "error org name $org"
        exit 1
    fi

    echo "${TLS_PATH}peerOrganizations/${org}.${com}.com/peers/${peer}.${org}.${com}.com/tls/$type"

}

get_orderer_tls_cert(){
    local org=$1
    if [[ "$org" != "orderer" ]] && [[ "$org" != "orderer1" ]]; then
        echo "error org name $org"
        exit 1
    fi

    echo "${ORDERER_TLS_PATH}ordererOrganizations/cpu.com/orderers/orderer.cpu.com/tls/tlsintermediatecerts/tls-localhost-7055.pem"
}

### 第一步：创建通道
channel_create() {
    local channel=$1
    local org="org1"
    local peer="peer0"
    local port="7051"
    local cert="server.crt"
    local key="server.key"
    local rootcert="ca.crt"
    local orderer="orderer"

    docker exec \
        -e "CORE_PEER_LOCALMSPID=Org1MSP" \
        -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.huawei.com/users/Admin@org1.huawei.com/msp" \
        -e "CORE_PEER_ADDRESS=peer0.org1.huawei.com:7051" \
        -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/crypto/peerOrganizations/org1.huawei.com/peers/peer0.org1.huawei.com/tls/server.crt" \
        -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/crypto/peerOrganizations/org1.huawei.com/peers/peer0.org1.huawei.com/tls/server.key" \
        -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/crypto/peerOrganizations/org1.huawei.com/peers/peer0.org1.huawei.com/tls/ca.crt" \
        cli \
        peer channel create -o $ORDERER_ADDRESS -c $channel  -f ./channel-artifacts/$CHANNEL_NAMETX.tx --tls true  --cafile $ORDERER_CAFILE

       echo "第一步：创建通道完成"
}

### 节点加入通道
channel_join() {
    local channel=$1
    local org=$2
    local peer=$3
    local port=$4
    local cert=$5
    local key=$6
    local rootcert=$7
    local com=$8
    ###
    docker exec \
        -e "CORE_PEER_LOCALMSPID=$(get_mspid $org)"\
        -e "CORE_PEER_MSPCONFIGPATH=$(get_msp_config_path $org $peer $com)"\
        -e "CORE_PEER_ADDRESS=$(get_peer_address $org $peer $port $com)"\
        -e "CORE_PEER_TLS_CERT_FILE=$(get_peer_tls_cert $org $peer $cert $com)"\
        -e "CORE_PEER_TLS_KEY_FILE=$(get_peer_tls_cert $org $peer $key $com)"\
        -e "CORE_PEER_TLS_ROOTCERT_FILE=$(get_peer_tls_cert $org $peer $rootcert $com)"\
        cli \
        peer channel join -b $channel.block

     echo "********************$org...$peer join channel $channel successful***************"
}

### 安装流程
install_and_instantiate_one() {
    local lang=$1
    local cc_name=($2)
    local cc_src_path=($3)

    ### install
    chaincode_install $CHANNEL_NAME  "org1" "peer0" "7051" "server.crt" "server.key" "ca.crt" "orderer" ${cc_name[0]} ${cc_src_path[0]} $lang "org1"  "huawei"
    chaincode_install $CHANNEL_NAME  "org1" "peer1" "8051" "server.crt" "server.key" "ca.crt" "orderer" ${cc_name[0]} ${cc_src_path[0]} $lang "org1"   "huawei"
    chaincode_install $CHANNEL_NAME  "org2" "peer0" "9051" "server.crt" "server.key" "ca.crt" "orderer" ${cc_name[0]} ${cc_src_path[0]} $lang "org2"    "xiaomi"
    chaincode_install $CHANNEL_NAME  "org2" "peer1" "10051" "server.crt" "server.key" "ca.crt" "orderer" ${cc_name[0]} ${cc_src_path[0]} $lang "org2"   "xiaomi"

    ### instantiate
    chaincode_instantiate   $CHANNEL_NAME "org1" "peer0" "7051" "server.crt" "server.key" "ca.crt" "orderer" ${cc_name[0]} ${cc_src_path[0]} $lang "org1" "huawei"

    #### invoke
    sleep $SLEEP_SECOND

     chaincode_invoke $CHANNEL_NAME "org1" "org1" "peer0" "7051" "server.crt" "server.key" "ca.crt" "orderer" ${cc_name[0]} "org2" "org2" "peer0" "9051"  '{"function":"","Args":[""]}' "huawei"
}
### 链码安装
chaincode_install() {
    local channel=$1
    local org=$2
    local peer=$3
    local port=$4
    local cert=$5
    local key=$6
    local rootcert=$7
    local orderer=$8
    local cc_name=$9
    local cc_src_path=${10}
    local lang=${11}
    local Org=${12}
    local com=${13}


    docker exec \
        -e "CORE_PEER_LOCALMSPID=$(get_mspid $org)" \
        -e "CORE_PEER_MSPCONFIGPATH=$(get_msp_config_path $org $peer $com)" \
        -e "CORE_PEER_ADDRESS=$(get_peer_address $org $peer $port $com)" \
        -e "CORE_PEER_TLS_CERT_FILE=$(get_peer_tls_cert $org $peer $cert $com)" \
        -e "CORE_PEER_TLS_KEY_FILE=$(get_peer_tls_cert $org $peer $key $com)" \
        -e "CORE_PEER_TLS_ROOTCERT_FILE=$(get_peer_tls_cert $org $peer $rootcert $com)" \
        cli \
        peer chaincode install -n $cc_name  -v $CC_VERSION -l $lang -p $cc_src_path

       echo "********************$org...$peer install chaincode $cc_name successful***************"
}

### 链码实例化
chaincode_instantiate() {
    local channel=$1
    local org=$2
    local peer=$3
    local port=$4
    local cert=$5
    local key=$6
    local rootcert=$7
    local orderer=$8
    local cc_name=$9
    local cc_src_path=${10}
    local lang=${11}
    local Org=${12}
    local com=${13}

    docker exec \
        -e "CORE_PEER_LOCALMSPID=$(get_mspid $org)" \
        -e "CORE_PEER_MSPCONFIGPATH=$(get_msp_config_path $org $peer $com)" \
        -e "CORE_PEER_ADDRESS=$(get_peer_address $org $peer $port $com)" \
        -e "CORE_PEER_TLS_CERT_FILE=$(get_peer_tls_cert $org $peer $cert $com)" \
        -e "CORE_PEER_TLS_KEY_FILE=$(get_peer_tls_cert $org $peer $key $com)" \
        -e "CORE_PEER_TLS_ROOTCERT_FILE=$(get_peer_tls_cert $org $peer $rootcert $com)" \
        cli \
        peer chaincode instantiate -o $ORDERER_ADDRESS --tls true --cafile $ORDERER_CAFILE -C $CHANNEL_NAME -n $cc_name -l golang -v 1.0 -c '{"Args":[""]}' -P 'OR ('\''Org1MSP.member'\'','\''Org2MSP.member'\'')'
      echo "********************$org...$peer instantiate chaincode $cc_name successful***************"
}

### 链码初始化
chaincode_invoke() {
    local channel=$1
    local org1=$2
    local Org1=$3
    local peer=$4
    local port=$5
    local cert=$6
    local key=$7
    local rootcert=$8
    local orderer=$9
    local cc_name=${10}
    local org2=${11}
    local Org2=${12}
    local Org2peer=${13}
    local Org2port=${14}
    local  cmd=${15}
    local  com=${16}

    docker exec \
        -e "CORE_PEER_LOCALMSPID=$(get_mspid $org1)" \
        -e "CORE_PEER_MSPCONFIGPATH=$(get_msp_config_path $Org1 $peer $com)" \
        -e "CORE_PEER_ADDRESS=$(get_peer_address $org1 $peer $port $com)" \
        -e "CORE_PEER_TLS_CERT_FILE=$(get_peer_tls_cert $org1 $peer $cert $com)" \
        -e "CORE_PEER_TLS_KEY_FILE=$(get_peer_tls_cert $org1 $peer $key $com)" \
        -e "CORE_PEER_TLS_ROOTCERT_FILE=$(get_peer_tls_cert $org1 $peer $rootcert $com)" \
        cli \
        peer chaincode invoke -o $ORDERER_ADDRESS --tls true --cafile $ORDERER_CAFILE -C $CHANNEL_NAME -n $cc_name --peerAddresses peer0.org1.huawei.com:7051 --tlsRootCertFiles  $HUAWEI_TLSROOTCERTFILE   --peerAddresses peer0.org2.xiaomi.com:9051 --tlsRootCertFiles $XIAOMI_TLSROOTCERTFILE -c  '{"function":"","Args":[""]}'

    echo "**********************************invoke chaincode*******$cc_name************************************************"
}

###########################################   脚本开始执行 ##########################################


### 创建通道
channel_create $CHANNEL_NAME

### 节点加入网络
channel_join $CHANNEL_NAME "org1" "peer0" "7051" "server.crt" "server.key" "ca.crt" "huawei"
channel_join $CHANNEL_NAME "org1" "peer1" "8051" "server.crt" "server.key" "ca.crt" "huawei"
channel_join $CHANNEL_NAME "org2" "peer0" "9051" "server.crt" "server.key" "ca.crt" "xiaomi"
channel_join $CHANNEL_NAME "org2" "peer1" "10051" "server.crt" "server.key" "ca.crt" "xiaomi"


### 安装链码
install_and_instantiate_one "golang" "${GO_CC_NAME[*]}" "${GO_CC_SRC_PATH[*]}"
echo
echo " _____   _   _   ____   "
echo "| ____| | \ | | |  _ \  "
echo "|  _|   |  \| | | | | | "
echo "| |___  | |\  | | |_| | "
echo "|_____| |_| \_| |____/  "
echo

