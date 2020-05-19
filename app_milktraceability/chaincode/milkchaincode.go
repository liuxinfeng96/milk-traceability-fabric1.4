package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

//MilkChaincode 自定义链码
type MilkChaincode struct {
}

//SourceInfo 原料信息
type SourceInfo struct {
	GrassState string `json:"grassState"` //牧草指标
	CowState   string `json:"cowState"`   //奶牛状态指标
	MilkState  string `json:"milkState"`  //生牛乳指标
}

//ProcessInfo 加工信息
type ProcessInfo struct {
	ProteinContent string `json:"proteinContent"` //蛋白质含量
	SterilizeTime  string `json:"sterilizeTime"`  //杀菌时间
	StorageTime    string `json:"storageTime"`    //存储时间
}

//LogInfo 配送信息
type LogInfo struct {
	LogCopName     string `json:"logCopName"`     //物流公司名称
	LogDepartureTm string `json:"logDepartureTm"` //出发时间
	LogArrivalTm   string `json:"logArrivalTm"`   //到达时间
	LogDeparturePl string `json:"logDeparturePl"` //出发地
	LogDest        string `json:"logDest"`        //目的地
	LogMOT         string `json:"logMOT"`         //运送方式
	TempAvg        string `json:"tempAvg"`        //平均温度
}

//MilkInfo 牛奶产品
type MilkInfo struct {
	MilkID          string      `json:"milkID"`          //牛奶ID
	MilkSourceInfo  SourceInfo  `json:"milkSourceInfo"`  //原料信息
	MilkProcessInfo ProcessInfo `json:"milkProcessInfo"` //加工信息
	MilkLogInfo     LogInfo     `json:"milkLogInfo"`     //配送信息
}

//Init 函数
func (m *MilkChaincode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

//Invoke 函数
func (m *MilkChaincode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
	if function == "addSourceInfo" {
		return m.addSourceInfo(APIstub, args)
	} else if function == "addProcessInfo" {
		return m.addProcessInfo(APIstub, args)
	} else if function == "addLogInfo" {
		return m.addLogInfo(APIstub, args)
	} else if function == "queryMilk" {
		return m.queryMilk(APIstub, args)
	} else if function == "queryAllMilks" {
		return m.queryAllMilks(APIstub)
	} else if function == "getHistoryInfo" {
		return m.getHistoryInfo(APIstub, args)
	} else if function == "initLedger" {
		return m.initLedger(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}
func (m *MilkChaincode) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	milks := []MilkInfo{
		MilkInfo{
			MilkID:          "milk1",
			MilkSourceInfo:  SourceInfo{GrassState: "8.8", CowState: "6.7", MilkState: "9.4"},
			MilkProcessInfo: ProcessInfo{ProteinContent: "3.5%", SterilizeTime: "2小时", StorageTime: "5小时"},
			MilkLogInfo:     LogInfo{LogCopName: "顺丰", LogDepartureTm: "2020-4-11", LogArrivalTm: "2020-4-13", LogDeparturePl: "大连", LogDest: "北京", LogMOT: "低温", TempAvg: "5摄氏度"},
		},
		MilkInfo{
			MilkID:          "milk2",
			MilkSourceInfo:  SourceInfo{GrassState: "7.8", CowState: "8.7", MilkState: "9.4"},
			MilkProcessInfo: ProcessInfo{ProteinContent: "3.0%", SterilizeTime: "2小时", StorageTime: "5小时"},
			MilkLogInfo:     LogInfo{LogCopName: "韵达", LogDepartureTm: "2020-4-11", LogArrivalTm: "2020-4-13", LogDeparturePl: "大连", LogDest: "北京", LogMOT: "低温", TempAvg: "5摄氏度"},
		},
		MilkInfo{
			MilkID:          "milk3",
			MilkSourceInfo:  SourceInfo{GrassState: "8.5", CowState: "6.7", MilkState: "9.0"},
			MilkProcessInfo: ProcessInfo{ProteinContent: "3.2%", SterilizeTime: "2小时", StorageTime: "5小时"},
			MilkLogInfo:     LogInfo{LogCopName: "德邦", LogDepartureTm: "2020-4-11", LogArrivalTm: "2020-4-13", LogDeparturePl: "大连", LogDest: "北京", LogMOT: "低温", TempAvg: "5摄氏度"},
		},
	}

	i := 0
	for i < len(milks) {
		fmt.Println("i is ", i)
		milkAsBytes, _ := json.Marshal(milks[i])
		APIstub.PutState(milks[i].MilkID, milkAsBytes)
		fmt.Println("Added", milks[i])
		i = i + 1
	}

	return shim.Success(nil)
}

//添加原料信息
func (m *MilkChaincode) addSourceInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var err error
	var MilkInfos MilkInfo

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments.")
	}
	MilkInfos.MilkID = args[0]
	if MilkInfos.MilkID == "" {
		return shim.Error("MilkID can not be empty.")
	}
	//比较关键：增加之前先要看看账本里是否存在其他信息
	milkAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("查看账本出现错误")
	}
	if milkAsBytes != nil {
		errs := json.Unmarshal(milkAsBytes, &MilkInfos)
		if errs != nil {
			fmt.Println(err)
			return shim.Error("反序列化失败")
		}
	}

	MilkInfos.MilkSourceInfo.GrassState = args[1]
	MilkInfos.MilkSourceInfo.CowState = args[2]
	MilkInfos.MilkSourceInfo.MilkState = args[3]
	SourceInfoJSONasBytes, err := json.Marshal(MilkInfos)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = APIstub.PutState(MilkInfos.MilkID, SourceInfoJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//添加加工信息
func (m *MilkChaincode) addProcessInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var err error
	var MilkInfos MilkInfo
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments.")
	}
	MilkInfos.MilkID = args[0]
	if MilkInfos.MilkID == "" {
		return shim.Error("MilkID can not be empty.")
	}
	//比较关键：增加之前先要看看账本里是否存在其他信息
	milkAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("查看账本出现错误")
	}
	if milkAsBytes == nil {
		return shim.Error("对不起，没有此产品")
	}
	errs := json.Unmarshal(milkAsBytes, &MilkInfos)
	if errs != nil {
		fmt.Println(err)
		return shim.Error("反序列化失败")
	}
	MilkInfos.MilkProcessInfo.ProteinContent = args[1]
	MilkInfos.MilkProcessInfo.SterilizeTime = args[2]
	MilkInfos.MilkProcessInfo.StorageTime = args[3]
	ProcessInfoJSONasBytes, err := json.Marshal(MilkInfos)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(MilkInfos.MilkID, ProcessInfoJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//添加配送信息
func (m *MilkChaincode) addLogInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var err error
	var MilkInfos MilkInfo
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments.")
	}
	MilkInfos.MilkID = args[0]
	if MilkInfos.MilkID == "" {
		return shim.Error("MilkID can not be empty.")
	}
	//比较关键：增加之前先要看看账本里是否存在其他信息
	milkAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("查看账本出现错误")
	}
	if milkAsBytes == nil {
		return shim.Error("对不起，没有此产品")
	}
	errs := json.Unmarshal(milkAsBytes, &MilkInfos)
	if errs != nil {
		fmt.Println(err)
		return shim.Error("反序列化失败")
	}
	MilkInfos.MilkLogInfo.LogCopName = args[1]
	MilkInfos.MilkLogInfo.LogDepartureTm = args[2]
	MilkInfos.MilkLogInfo.LogArrivalTm = args[3]
	MilkInfos.MilkLogInfo.LogDeparturePl = args[4]
	MilkInfos.MilkLogInfo.LogDest = args[5]
	MilkInfos.MilkLogInfo.LogMOT = args[6]
	MilkInfos.MilkLogInfo.TempAvg = args[7]
	LogInfoJSONasBytes, err := json.Marshal(MilkInfos)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(MilkInfos.MilkID, LogInfoJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//查询产品
func (m *MilkChaincode) queryMilk(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	milkAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(milkAsBytes)
}

//查询所有产品
func (m *MilkChaincode) queryAllMilks(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "milk1"
	endKey := "milk99"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllMilks:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//查询历史信息
func (m *MilkChaincode) getHistoryInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var MilkInfos MilkInfo
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	MilkInfos.MilkID = args[0]
	fmt.Printf("- start getHistory: %s\n", MilkInfos.MilkID)
	resultsIterator, err := APIstub.GetHistoryForKey(MilkInfos.MilkID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryInfo returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(MilkChaincode))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
