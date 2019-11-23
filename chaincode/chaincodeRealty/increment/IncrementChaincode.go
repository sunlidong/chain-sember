/*
 * @Author  : yangqingwei
 * @Date    : 2018-08-13 09:33:22
 * @Describe: 增量合约、记录每次更新都涉及了哪些key
 *            基础资产ID-更新时间-资产类型-资产key
 *
 */
package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func main() {
	err := shim.Start(new(IncrementChaincode))
	if err != nil {
		fmt.Printf("Error starting IncrementChaincode - %s", err)
	}
}

const keySplitKey string = "+"
const resultSplitKey string = "=="

// 保理资产合约
type IncrementChaincode struct {
}

var logger = shim.NewLogger("IncrementChaincode")

// 初始化操作
func (c *IncrementChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// 定义的各种操作
func (c *IncrementChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Debug(function, "--入参：　", args)
	var resp pb.Response
	switch function {
	case "addIncrementInfo":
		// 添加增量信息
		resp = c.addIncrementInfo(stub, args)
	case "getIncrementInfo":
		// 添加增量信息
		resp = c.getIncrementInfo(stub, args)
	default:
		resp = shim.Error(fmt.Sprintf("方法未定义:- %s", function))
	}
	logger.Debug(function, "--响应：　", " \n status:", resp.Status, " \n  Message:", resp.Message, " \n  Payload:", string(resp.Payload))
	return resp
}

// 获取增量信息
// 0 assetID 　基础资产ID
// 1 updateTime 更新时间
func (c *IncrementChaincode) getIncrementInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if err := checkArgsForCount(args, 2); err != nil {
		return shim.Error(err.Error())
	}
	assetID := args[0]
	updateTime := args[1]

	//　根据范围模糊查找数据
	startKey := fmt.Sprintf("%s+%s+", assetID, updateTime)
	endKey := fmt.Sprintf("%s+%s0", assetID, updateTime)
	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	// 判空
	defer resultsIterator.Close()
	if !resultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("there are nothing about assetID: %s , updateTime:  %s", assetID, updateTime))
	}
	// 组装数据　给调用者
	var buffer bytes.Buffer
	isFirst := true
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if !isFirst {
			buffer.WriteString(resultSplitKey)
		}
		buffer.WriteString(queryResponse.Key)
		isFirst = false
	}
	return shim.Success(buffer.Bytes())
}

// 添加增量信息　　每次都是不同的　key--value
// 0 assetID 　基础资产ID
// 1 updateTime
// 2 assetType　(1发票、2合同、3其他附件信息、4底层资产、5保理资产、6基础资产)（ 10代表资产无变化、同时assetKey为assetUUID）
// 3 assetKey
// todo 　需要增加资产无变化的情况
func (c *IncrementChaincode) addIncrementInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if err := checkArgsForCount(args, 4); err != nil {
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}
	// 特殊判断底层资产的子资产的情况
	// 类型如果是　1 2 3 　assetKey =  底层资产ID+该类型ID
	specialAssetType := "1+2+3"
	isContains := strings.Contains(specialAssetType, args[2])
	if isContains {
		types := strings.Split(args[3], keySplitKey)
		if len(types) != 2 {
			return shim.Error("assetKey should be \"parentID-currentAssetKey\" ")
		}
		for _, assetType := range types {
			if len(assetType) == 0 {
				return shim.Error("assetKey should  not be empty!")
			}
		}
	}
	// 组建　key value
	key := fmt.Sprintf("%s+%s+%s+%s", args[0], args[1], args[2], args[3])
	if err := stub.PutState(key, []byte{0x00}); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(stub.GetTxID()))
}

// 检查参数　个数以及是否为空
func checkArgsForCount(args []string, count int) error {
	if len(args) != count {
		return fmt.Errorf("Incorrect number of arguments. Expecting :  %v", count)
	}
	// 验空
	for index := 0; index < count; index++ {
		if len(args[index]) <= 0 {
			return fmt.Errorf("index :%v  argument must be a non-empty string", index)
		}
	}
	return nil
}
