package v1

import (
	bc "application/blockchain"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SellingRequestBody struct {
	ObjectOfSale string  `json:"objectOfSale"` //销售对象(正在出售的房地产RealEstateID)
	Seller       string  `json:"seller"`       //发起销售人、卖家(卖家AccountId)
	Price        float64 `json:"price"`        //价格
	SalePeriod   int     `json:"salePeriod"`   //智能合约的有效期(单位为天)
}

type SellingByBuyRequestBody struct {
	ObjectOfSale string `json:"objectOfSale"` //销售对象(正在出售的房地产RealEstateID)
	Seller       string `json:"seller"`       //发起销售人、卖家(卖家AccountId)
	Buyer        string `json:"buyer"`        //买家(买家AccountId)
}

type SellingListQueryRequestBody struct {
	Seller string `json:"seller"` //发起销售人、卖家(卖家AccountId)
}

type SellingListQueryByBuyRequestBody struct {
	Buyer string `json:"buyer"` //买家(买家AccountId)
}

type UpdateSellingRequestBody struct {
	ObjectOfSale string `json:"objectOfSale"` //销售对象(正在出售的房地产RealEstateID)
	Seller       string `json:"seller"`       //发起销售人、卖家(卖家AccountId)
	Buyer        string `json:"buyer"`        //买家(买家AccountId)
	Status       string `json:"status"`       //需要更改的状态
}

func CreateSelling(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SellingRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ObjectOfSale == "" || body.Seller == "" {
		appG.Response(http.StatusBadRequest, "失败", "ObjectOfSale销售对象和Seller发起销售人不能为空")
		return
	}
	if body.Price <= 0 || body.SalePeriod <= 0 {
		appG.Response(http.StatusBadRequest, "失败", "Price价格和SalePeriod智能合约的有效期(单位为天)必须大于0")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfSale))
	bodyBytes = append(bodyBytes, []byte(body.Seller))
	bodyBytes = append(bodyBytes, []byte(strconv.FormatFloat(body.Price, 'E', -1, 64)))
	bodyBytes = append(bodyBytes, []byte(strconv.Itoa(body.SalePeriod)))
	//调用智能合约
	resp, err := bc.ChannelExecute("createSelling", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func CreateSellingByBuy(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SellingByBuyRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ObjectOfSale == "" || body.Seller == "" || body.Buyer == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数不能为空")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfSale))
	bodyBytes = append(bodyBytes, []byte(body.Seller))
	bodyBytes = append(bodyBytes, []byte(body.Buyer))
	//调用智能合约
	resp, err := bc.ChannelExecute("createSellingByBuy", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func QuerySellingList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SellingListQueryRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.Seller != "" {
		bodyBytes = append(bodyBytes, []byte(body.Seller))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("querySellingList", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func QuerySellingListByBuyer(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SellingListQueryByBuyRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.Buyer == "" {
		appG.Response(http.StatusBadRequest, "失败", "必须指定买家AccountId查询")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.Buyer))
	//调用智能合约
	resp, err := bc.ChannelQuery("querySellingListByBuyer", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func UpdateSelling(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(UpdateSellingRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ObjectOfSale == "" || body.Seller == "" || body.Status == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数不能为空")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfSale))
	bodyBytes = append(bodyBytes, []byte(body.Seller))
	bodyBytes = append(bodyBytes, []byte(body.Buyer))
	bodyBytes = append(bodyBytes, []byte(body.Status))
	//调用智能合约
	resp, err := bc.ChannelExecute("updateSelling", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

//新建
type Contract_in_company struct {
	ContractName    string `json:"contractname"`
	ContractContent string `json:"contractcontent"`
	CreaterName     string `json:"creatername"`
	CreaterSign     string `json:"creatersign"`
	CreateTime      string `json:"creatertime"`
	CompanyName     string `json:"companyname"`
}
type Contract_among_company struct {
	ContractName       string `json:"contractname"`
	ContractContent    string `json:"contractcontent"`
	CreaterCompanyName string `json:"creatercompanyname"`
	CreaterCompanySign string `json:"creatercompanysign"`
	SignTime           string `json:"signtime"`
}

type Contract_request struct {
	ContractName string `json:"contractname"`
	CompanyName  string `json:"companyname"`
}
type Contract_sanction struct {
	ContractName string `json:"contractname"`
	CompanyName  string `json:"companyname"`
	CompanySign  string `json:"companysign"`
}

func StartContract(c *gin.Context) {
	//companyname string, contractname string, departmentname string, signature string, contractcontent string
	appG := app.Gin{C: c}
	body := new(Contract_in_company)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ContractName == "" || body.ContractContent == "" || body.CompanyName == "" || body.CreaterName == "" || body.CreaterSign == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数不能为空")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ContractName))
	bodyBytes = append(bodyBytes, []byte(body.CreaterName))
	bodyBytes = append(bodyBytes, []byte(body.CreaterSign))
	bodyBytes = append(bodyBytes, []byte(body.ContractContent))
	//调用智能合约
	resp, err := bc.ChannelExecute_new("StartContract", body.CompanyName, bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func QueryContract_incompany(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(Contract_request)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ContractName == "" || body.CompanyName == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数不能为空")
		return
	}
	var bodyBytes [][]byte

	bodyBytes = append(bodyBytes, []byte(body.ContractName))

	//调用智能合约
	resp, err := bc.ChannelQuery_new("QueryContract_incompany", body.CompanyName, bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

//审核合同并上传
func ContractSanction_upload(c *gin.Context) {

	appG := app.Gin{C: c}
	body := new(Contract_sanction)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ContractName == "" || body.CompanyName == "" || body.CompanySign == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数不能为空")
		return
	}
	var bodyBytes_1 [][]byte
	bodyBytes_1 = append(bodyBytes_1, []byte(body.CompanyName))
	bodyBytes_1 = append(bodyBytes_1, []byte(body.CompanySign))
	bodyBytes_1 = append(bodyBytes_1, []byte(body.ContractName))

	//调用智能合约
	resp_1, err := bc.ChannelQuery_new("ContractSanction", body.CompanyName, bodyBytes_1)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	/*var data_1 map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data_1); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}

	contract_among_company_new := Contract_among_company{
		ContractName:       data_1["ContractName"].(string),
		ContractContent:    data_1["ContractContent"].(string),
		CreaterCompanyName: data_1["company_name"].(string),
		CreaterCompanySign: data_1["signature"].(string),
		SignTime:         data_1["signTime"].(string),
	}*/

	var bodyBytes_2 [][]byte
	bodyBytes_2 = append(bodyBytes_2, resp_1.Payload)
	//定义公司间通道为"Among_Companies"
	resp_2, err := bc.ChannelExecute_new("ContractSanction_upload", "Among_Companies", bodyBytes_2)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp_2.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}


func QueryContract_amongcompany_1(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(Contract_request)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ContractName == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数不能为空")
		return
	}
	var bodyBytes [][]byte

	bodyBytes = append(bodyBytes, []byte(body.ContractName))

	//调用智能合约
	resp, err := bc.ChannelQuery_new("QueryContract_amongcompany", "Among_Companies", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func QueryContract_state(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(Contract_request)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ContractName == "" || body.CompanyName == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数不能为空")
		return
	}
	var bodyBytes [][]byte

	bodyBytes = append(bodyBytes, []byte(body.ContractName))

	//调用智能合约
	resp1, err1 := bc.ChannelQuery_new("QueryContract_incompany", body.CompanyName, bodyBytes)
	if err1 != nil {
		appG.Response(http.StatusInternalServerError, "失败", err1.Error())
		return
	}

	//调用智能合约
	resp2, err2 := bc.ChannelQuery_new("QueryContract_amongcompany", "Among_Companies", bodyBytes)
	if err2 != nil {
		appG.Response(http.StatusInternalServerError, "失败", err2.Error())
		return
	}
	// 反序列化json
	var data1 []map[string]interface{}
	if err1 = json.Unmarshal(bytes.NewBuffer(resp1.Payload).Bytes(), &data1); err1 != nil {
		appG.Response(http.StatusInternalServerError, "失败", err1.Error())
		return
	}
	var data2 []map[string]interface{}
	if err2 = json.Unmarshal(bytes.NewBuffer(resp2.Payload).Bytes(), &data2); err2 != nil {
		appG.Response(http.StatusInternalServerError, "失败", err2.Error())
		return
	}
	var result string

	if len(data1) == 0 && len(data2) == 0 {
		result = "no such contract"
	} else if len(data1) == 0 {
		result = "other company's contract"
	} else if len(data2) == 0 {
		result = "Unsigned contract"
	}else{
		result = "signed contract"
	}

	appG.Response(http.StatusOK, "成功", result)
}
