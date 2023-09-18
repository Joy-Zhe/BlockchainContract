package main

import (
	"bytes"
	"chaincode/model"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func initTest(t *testing.T) *shim.MockStub {
	scc := new(BlockChainRealEstate)
	stub := shim.NewMockStub("ex01", scc)
	checkInit(t, stub, [][]byte{[]byte("init")})
	return stub
}

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) pb.Response {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
	return res
}

// 测试链码初始化
func TestBlockChainRealEstate_Init(t *testing.T) {
	stub := initTest(t)
	fmt.Println(fmt.Sprintf("1、初始化\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("InitLedger"),
		}).Payload)))

	fmt.Println(fmt.Sprintf("2、上传合同\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("ContractSanction_upload"),
			[]byte("{\"contractname\":\"NEW1\"," +
				"\"contractcontent\":\"6b34\",\"creatercompanyname\":" +
				"\"taobao\",\"creatercompanysign\":\"mayun\",\"creatertime\":" +
				"\"2023-09-15 20:51:15\"}"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("3、查询合同\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("QueryContract_amongcompany"),
		}).Payload)))
}
