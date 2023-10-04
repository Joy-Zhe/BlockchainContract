package routers

import (
	v1 "application/api/v1"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由信息
func InitRouter() *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/hello", v1.Hello)
		apiV1.POST("/queryAccountList", v1.QueryAccountList)
		apiV1.POST("/createRealEstate", v1.CreateRealEstate)
		apiV1.POST("/queryRealEstateList", v1.QueryRealEstateList)
		apiV1.POST("/createSelling", v1.CreateSelling)
		apiV1.POST("/createSellingByBuy", v1.CreateSellingByBuy)
		apiV1.POST("/querySellingList", v1.QuerySellingList)
		apiV1.POST("/querySellingListByBuyer", v1.QuerySellingListByBuyer)
		apiV1.POST("/updateSelling", v1.UpdateSelling)
		apiV1.POST("/createDonating", v1.CreateDonating)
		apiV1.POST("/queryDonatingList", v1.QueryDonatingList)
		apiV1.POST("/queryDonatingListByGrantee", v1.QueryDonatingListByGrantee)
		apiV1.POST("/updateDonating", v1.UpdateDonating)
		apiV1.POST("/StartContract", v1.StartContract)//部门节点调用，创建合同
		apiV1.POST("/QueryContract_incompany", v1.QueryContract_incompany)//部门节点和老板节点调用，查询公司内合同
		apiV1.POST("/ContractSanction_upload", v1.ContractSanction_upload)//老板节点调用，签署合同并上传到公司间channel
		apiV1.POST("/QueryContract_amongcompany", v1.QueryContract_amongcompany_1)//老板节点调用，查询公司间合同
		apiV1.POST("/QueryContract_state", v1.QueryContract_state)//老板节点调用，查询某合同签署情况（在发起未签署、已经签署、为其他公司的合同、无此合同
	}
	return r
}
