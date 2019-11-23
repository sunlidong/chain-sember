package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"sort"
	"strings"
)

// 底层资产合约
type UnderlyingChaincode struct {
}

//	logger
var logger = shim.NewLogger("UnderlyingChaincode")

// 资产的类型 || 房地产资产类型
const (
	DATATYPEKYC string = "datakyc" // kyc
	DATATYPEBYC string = "databyc" // byc
)

// 资产操作类型
const (
	assetOperationType_ADD    string = "ADD"
	assetOperationType_DELETE string = "DELETE"
	assetOperationType_UPDATE string = "UPDATE"
)

// 初始化操作
func (c *UnderlyingChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// 定义的各种操作
func (c *UnderlyingChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Debug(function, "--入参：　", args)
	var resp pb.Response
	switch function {
	// -------------------新增数据的时候 -----------------------
	case uploadAsset:
		resp = c.uploadAsset(stub, args)
	case updateAsset:
		resp = c.updateAsset(stub, args)
	case getAssetByID:
		// 根据assetID,updateTime 获取当时的对象
		resp = c.getAssetByID(stub, args)
	case getAssetByUpdateTime:
		// 根据assetID,updateTime 获取当时的对象
		resp, _ = c.getAssetByUpdateTime(stub, args)
	case getAssetList:
		// 根据assetID 获取关联信息
		resp = c.getAssetList(stub, args)
	default:
		resp = shim.Error(fmt.Sprintf("方法未定义:- %s", function))
	}
	logger.Debug(function, "--响应：　", " \n status:", resp.Status, " \n  Message:", resp.Message, " \n  Payload:", string(resp.Payload))
	return resp
}

const keySplitKey string = "+"
const resultSplitKey string = "=="

//	查询||根据Key获取数据
func (c *UnderlyingChaincode) getAssetByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if err := checkArgsForCount(args, 1); err != nil {
		return shim.Error(err.Error())
	}
	currentAssetID := args[0]
	assetAsBytes, err := stub.GetState(currentAssetID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if assetAsBytes == nil {
		return shim.Error("the asset is not existed!")
	}

	return shim.Success(assetAsBytes)
}

//	查询||根据Key获取关联数据
func (c *UnderlyingChaincode) getAssetList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	BycData := getByc()
	//BycDataList := getBycList()
	KycDataList := getKycList()
	//
	if err := checkArgsForCount(args, 1); err != nil {
		return shim.Error(err.Error())
	}

	currentAssetID := args[0]

	// 获取标签表
	assetAsBytes, err := stub.GetState(BYC + currentAssetID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if assetAsBytes == nil {
		return shim.Error("the asset is not existed!")
	}

	// 序列化 到结构体
	json.Unmarshal(assetAsBytes, &BycData)
	//获取 Sun	list

	SunList, stateSun := getAssetSunList(stub, BycData.BycMes)
	if stateSun {
		*KycDataList = append(*KycDataList, *SunList...)
	}

	// 反序列化
	Mardata, err := json.Marshal(*KycDataList)
	if err != nil {
		return shim.Error(err.Error())
	}

	//	返回
	return shim.Success(Mardata)
}

//	查询||根据Key+Time获取数据
func (c *UnderlyingChaincode) getAssetByUpdateTime(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
	if err := checkArgsForCount(args, 2); err != nil {
		return shim.Error(err.Error()), ""
	}
	currentAssetID := args[0]
	updateTime := args[1]

	assetMap := make(map[string][]byte)
	assetTxIDMap := make(map[string]string)

	//1 迭代查找历史数据
	resultsIterator, err := stub.GetHistoryForKey(currentAssetID)
	if err != nil {
		return shim.Error(err.Error()), ""
	}
	defer resultsIterator.Close()
	if !resultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("the assetID : < %s > has not any history info", currentAssetID)), ""
	}
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()

		if err != nil {
			return shim.Error(err.Error()), ""
		}
		//2 查找更新时间　并存入map中
		assetUpdateTime := getUpdateTime(string(response.Value))
		if len(assetUpdateTime) > 0 {
			assetMap[assetUpdateTime] = response.Value
			assetTxIDMap[assetUpdateTime] = response.TxId
		} else {
			return shim.Error("can not find value about field of  updateDate in key:" + currentAssetID), ""
		}
	}
	keys := []string{}
	//3 根据日期返回对象
	for key, value := range assetMap {
		keys = append(keys, key)
		if strings.EqualFold(key, updateTime) {
			return shim.Success(value), assetTxIDMap[key]
		}
	}
	//4 记录首次操作时间
	sort.Strings(keys)
	startTime := keys[0]

	//5 否则　查找最近日期对象
	keys = append(keys, updateTime)
	sort.Strings(keys)
	for index, time := range keys {
		if strings.EqualFold(time, updateTime) {
			//6 判断是否为头  如果是头的话　就有问题
			if index == 0 {
				return shim.Error(fmt.Sprintf("the  assetID :%s  startTime is  %s  ,but receive : %s", currentAssetID, startTime, updateTime)), ""
			}
			return shim.Success(assetMap[keys[index-1]]), assetTxIDMap[keys[index-1]]
		}
	}
	return shim.Error(fmt.Sprintf("failed to find asset history info  by  assetID : %s , updateTime : %s ", currentAssetID, updateTime)), ""
}

//	新增|| 数据上链
func (c *UnderlyingChaincode) uploadAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0　　资产类型
	// 1　　资产结构体json
	// 2　　标签结构体
	// 参数检查
	if err := checkArgs(args); err != nil {
		return shim.Error(err.Error())
	}

	// 资产的上链
	if _, err := uploadAssetInternal(stub, args); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(stub.GetTxID()))
}

//	更新|| 数据上链更新
func (c *UnderlyingChaincode) updateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 参数检查
	if err := checkArgs(args); err != nil {
		return shim.Error(err.Error())
	}
	assetType := args[0]
	// 获取上链结构体对象
	assetStruct, err := getAssetStructByType(assetType)
	if err != nil {
		return shim.Error(err.Error())
	}
	// 资产的更新
	if _, err := updateAssetInternal(stub, args, assetStruct); err != nil {
		return shim.Error(err.Error())
	}

	// 资产关联
	assetStructrel, err := getAssetStructByType(BYC)
	if err != nil {
		return shim.Error(err.Error())
	}
	// 资产关联  权限  数据
	if _, err := updateAssetInternalRel(stub, args, assetStructrel); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(stub.GetTxID()))
}

// 操作资产前的　参数检查
func checkArgs(args []string) error {
	return checkArgsForCount(args, 2)
}

//	校验||参数检查
func checkArgsRec(args []string) error {
	return checkArgsForCount(args, 4)
}

//	校验||参数检查 4
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

//	结构体||根据类型返回对应结构体
func getAssetStructByType(assetType string) (interface{}, error) {
	var assetStruct interface{}
	switch assetType {

	case KYC:
		assetStruct = &Kyc{}
	case BYC:
		assetStruct = &Byc{}
	default:
		return nil, errors.New("The assetType is not supported:" + assetType)
	}

	return assetStruct, nil
}

//	结构体||根据类型返回对应结构体[]
func getAssetStructByTypelist(assetType string) (interface{}, error) {
	var assetStruct interface{}
	switch assetType {
	//	01. 房地产主表
	case DATATYPEKYC:
		assetStruct = &[]Kyc{}
	case DATATYPEBYC:
		assetStruct = &[]Byc{}
	default:
		return nil, errors.New("The assetType is not supported:" + assetType)
	}
	return assetStruct, nil
}

// 获取 Kyc
func getKyc() *Kyc {
	data := Kyc{}
	return &data
}

// 获取 []Kyc
func getKycList() *[]Kyc {
	data := []Kyc{}
	return &data
}

// 获取 Byc
func getByc() *Byc {
	data := Byc{}
	return &data
}

// 获取 []Byc
func getBycList() *[]Byc {
	data := []Byc{}
	return &data
}

// 获取相关数据
func getAssetSunList(stub shim.ChaincodeStubInterface, List []BycMes) (Kyc *[]Kyc, state bool) {
	//
	KycList := getKycList()

	if len(List) > 0 {
		//
		for k, _ := range List {
			if List[k].ParBycID != "" {
				KycData := getKyc()
				//查询关联 元数据   kyc
				assetAsBytes, err := stub.GetState(KYC + List[k].ParBycID)
				if err != nil {
					return nil, false
				}
				if assetAsBytes == nil {
					return nil, false
				}
				// 序列化
				json.Unmarshal(assetAsBytes, &KycData)
				*KycList = append(*KycList, *KycData)
			}
		}
	}
	return KycList, true
}

//	主函数|| main
func main() {
	err := shim.Start(new(UnderlyingChaincode))
	if err != nil {
		fmt.Printf("Error starting ProductChaincode - %s", err)
	}
}
