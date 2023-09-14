package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
	"net/http"
)

type Contract_in_company struct {
	ContractName    string `json:"contractname"`
	ContractContent string `json:"contractcontent"`
	CreaterName     string `json:"creatername"`
	CreaterSign     string `json:"creatersign"`
	CreateTime      string `json:"creatertime"`
}
type Contract_name struct {
	ContractName string `json:"contractname"`
}

var (
	SDK           *fabsdk.FabricSDK
	channelClient *channel.Client
	channelName   = "mychannel"
	chaincodeName = "fabcar"
	orgName       = "Org1"
	orgAdmin      = "Admin"
	org1Peer0     = "peer0.org1.example.com"
	org2Peer0     = "peer0.org2.example.com"
)

func ChannelExecute(funcName string, args [][]byte) (channel.Response, error) {
	var err error
	configPath := "./config.yaml"
	configProvider := config.FromFile(configPath)
	SDK, err = fabsdk.New(configProvider)
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}
	ctx := SDK.ChannelContext(channelName, fabsdk.WithOrg(orgName), fabsdk.WithUser(orgAdmin))
	channelClient, err = channel.New(ctx)
	response, err := channelClient.Execute(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         funcName,
		Args:        args,
	})
	if err != nil {
		return response, err
	}
	SDK.Close()
	return response, nil
}
func main() {
	r := gin.Default()
	r.GET("/queryAllCars", func(c *gin.Context) {
		var result channel.Response
		result, err := ChannelExecute("queryAllCars", [][]byte{})
		fmt.Println(result)
		if err != nil {
			log.Fatalf("Failed to evaluate transaction: %s\n", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Query All Success",
			"result":  string(result.Payload),
		})
	})
	r.POST("/StartContract", func(c *gin.Context) {
		var contract Contract_in_company
		c.BindJSON(&contract)
		var result channel.Response
		result, err := ChannelExecute("StartContract", [][]byte{[]byte(contract.ContractName), []byte(contract.ContractContent), []byte(contract.CreaterName), []byte(contract.CreaterSign)})
		fmt.Println(result)
		if err != nil {
			log.Fatalf("Failed to evaluate transaction: %s\n", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Create Success",
			"result":  string(result.Payload),
		})
	})
	rr.GET("/QueryALLContract_incompany", func(c *gin.Context) {
		var result channel.Response
		result, err := ChannelExecute("QueryALLContract_incompany", [][]byte{})
		fmt.Println(result)
		if err != nil {
			log.Fatalf("Failed to evaluate transaction: %s\n", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "Query All Success",
			"result":  string(result.Payload),
		})
	})

	r.Run(":9099")
}
