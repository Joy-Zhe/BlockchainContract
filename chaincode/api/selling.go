package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// CreateSelling 发起销售
func CreateSelling(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 4 {
		return shim.Error("参数个数不满足")
	}
	objectOfSale := args[0]
	seller := args[1]
	price := args[2]
	salePeriod := args[3]
	if objectOfSale == "" || seller == "" || price == "" || salePeriod == "" {
		return shim.Error("参数存在空值")
	}
	// 参数数据格式转换
	var formattedPrice float64
	if val, err := strconv.ParseFloat(price, 64); err != nil {
		return shim.Error(fmt.Sprintf("price参数格式转换出错: %s", err))
	} else {
		formattedPrice = val
	}
	var formattedSalePeriod int
	if val, err := strconv.Atoi(salePeriod); err != nil {
		return shim.Error(fmt.Sprintf("salePeriod参数格式转换出错: %s", err))
	} else {
		formattedSalePeriod = val
	}
	//判断objectOfSale是否属于seller
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, model.RealEstateKey, []string{seller, objectOfSale})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("验证%s属于%s失败: %s", objectOfSale, seller, err))
	}
	var realEstate model.RealEstate
	if err = json.Unmarshal(resultsRealEstate[0], &realEstate); err != nil {
		return shim.Error(fmt.Sprintf("CreateSelling-反序列化出错: %s", err))
	}
	//判断记录是否已存在，不能重复发起销售
	//若Encumbrance为true即说明此房产已经正在担保状态
	if realEstate.Encumbrance {
		return shim.Error("此房地产已经作为担保状态，不能重复发起销售")
	}
	createTime, _ := stub.GetTxTimestamp()
	selling := &model.Selling{
		ObjectOfSale:  objectOfSale,
		Seller:        seller,
		Buyer:         "",
		Price:         formattedPrice,
		CreateTime:    time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
		SalePeriod:    formattedSalePeriod,
		SellingStatus: model.SellingStatusConstant()["saleStart"],
	}
	// 写入账本
	if err := utils.WriteLedger(selling, stub, model.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将房子状态设置为正在担保状态
	realEstate.Encumbrance = true
	if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	sellingByte, err := json.Marshal(selling)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(sellingByte)
}

// CreateSellingByBuy 参与销售(买家购买)
func CreateSellingByBuy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 3 {
		return shim.Error("参数个数不满足")
	}
	objectOfSale := args[0]
	seller := args[1]
	buyer := args[2]
	if objectOfSale == "" || seller == "" || buyer == "" {
		return shim.Error("参数存在空值")
	}
	if seller == buyer {
		return shim.Error("买家和卖家不能同一人")
	}
	//根据objectOfSale和seller获取想要购买的房产信息，确认存在该房产
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, model.RealEstateKey, []string{seller, objectOfSale})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s获取想要购买的房产信息失败: %s", objectOfSale, seller, err))
	}
	//根据objectOfSale和seller获取销售信息
	resultsSelling, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellingKey, []string{seller, objectOfSale})
	if err != nil || len(resultsSelling) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s获取销售信息失败: %s", objectOfSale, seller, err))
	}
	var selling model.Selling
	if err = json.Unmarshal(resultsSelling[0], &selling); err != nil {
		return shim.Error(fmt.Sprintf("CreateSellingBuy-反序列化出错: %s", err))
	}
	//判断selling的状态是否为销售中
	if selling.SellingStatus != model.SellingStatusConstant()["saleStart"] {
		return shim.Error("此交易不属于销售中状态，已经无法购买")
	}
	//根据buyer获取买家信息
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{buyer})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("buyer买家信息验证失败%s", err))
	}
	var buyerAccount model.Account
	if err = json.Unmarshal(resultsAccount[0], &buyerAccount); err != nil {
		return shim.Error(fmt.Sprintf("查询buyer买家信息-反序列化出错: %s", err))
	}
	if buyerAccount.UserName == "管理员" {
		return shim.Error(fmt.Sprintf("管理员不能购买%s", err))
	}
	//判断余额是否充足
	if buyerAccount.Balance < selling.Price {
		return shim.Error(fmt.Sprintf("房产售价为%f,您的当前余额为%f,购买失败", selling.Price, buyerAccount.Balance))
	}
	//将buyer写入交易selling,修改交易状态
	selling.Buyer = buyer
	selling.SellingStatus = model.SellingStatusConstant()["delivery"]
	if err := utils.WriteLedger(selling, stub, model.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
		return shim.Error(fmt.Sprintf("将buyer写入交易selling,修改交易状态 失败%s", err))
	}
	createTime, _ := stub.GetTxTimestamp()
	//将本次购买交易写入账本,可供买家查询
	sellingBuy := &model.SellingBuy{
		Buyer:      buyer,
		CreateTime: time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
		Selling:    selling,
	}
	if err := utils.WriteLedger(sellingBuy, stub, model.SellingBuyKey, []string{sellingBuy.Buyer, sellingBuy.CreateTime}); err != nil {
		return shim.Error(fmt.Sprintf("将本次购买交易写入账本失败%s", err))
	}
	sellingBuyByte, err := json.Marshal(sellingBuy)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	//购买成功，扣取余额，更新账本余额，注意，此时需要卖家确认收款，款项才会转入卖家账户，此处先扣除买家的余额
	buyerAccount.Balance -= selling.Price
	if err := utils.WriteLedger(buyerAccount, stub, model.AccountKey, []string{buyerAccount.AccountId}); err != nil {
		return shim.Error(fmt.Sprintf("扣取买家余额失败%s", err))
	}
	// 成功返回
	return shim.Success(sellingBuyByte)
}

// QuerySellingList 查询销售(可查询所有，也可根据发起销售人查询)(发起的)(供卖家查询)
func QuerySellingList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var sellingList []model.Selling
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellingKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var selling model.Selling
			err := json.Unmarshal(v, &selling)
			if err != nil {
				return shim.Error(fmt.Sprintf("QuerySellingList-反序列化出错: %s", err))
			}
			sellingList = append(sellingList, selling)
		}
	}
	sellingListByte, err := json.Marshal(sellingList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QuerySellingList-序列化出错: %s", err))
	}
	return shim.Success(sellingListByte)
}

// QuerySellingListByBuyer 根据参与销售人、买家(买家AccountId)查询销售(参与的)(供买家查询)
func QuerySellingListByBuyer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("必须指定买家AccountId查询"))
	}
	var sellingBuyList []model.SellingBuy
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellingBuyKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var sellingBuy model.SellingBuy
			err := json.Unmarshal(v, &sellingBuy)
			if err != nil {
				return shim.Error(fmt.Sprintf("QuerySellingListByBuyer-反序列化出错: %s", err))
			}
			sellingBuyList = append(sellingBuyList, sellingBuy)
		}
	}
	sellingBuyListByte, err := json.Marshal(sellingBuyList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QuerySellingListByBuyer-序列化出错: %s", err))
	}
	return shim.Success(sellingBuyListByte)
}

// UpdateSelling 更新销售状态（买家确认、买卖家取消）
func UpdateSelling(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 4 {
		return shim.Error("参数个数不满足")
	}
	objectOfSale := args[0]
	seller := args[1]
	buyer := args[2]
	status := args[3]
	if objectOfSale == "" || seller == "" || status == "" {
		return shim.Error("参数存在空值")
	}
	if buyer == seller {
		return shim.Error("买家和卖家不能同一人")
	}
	//根据objectOfSale和seller获取想要购买的房产信息，确认存在该房产
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, model.RealEstateKey, []string{seller, objectOfSale})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s获取想要购买的房产信息失败: %s", objectOfSale, seller, err))
	}
	var realEstate model.RealEstate
	if err = json.Unmarshal(resultsRealEstate[0], &realEstate); err != nil {
		return shim.Error(fmt.Sprintf("UpdateSellingBySeller-反序列化出错: %s", err))
	}
	//根据objectOfSale和seller获取销售信息
	resultsSelling, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellingKey, []string{seller, objectOfSale})
	if err != nil || len(resultsSelling) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s获取销售信息失败: %s", objectOfSale, seller, err))
	}
	var selling model.Selling
	if err = json.Unmarshal(resultsSelling[0], &selling); err != nil {
		return shim.Error(fmt.Sprintf("UpdateSellingBySeller-反序列化出错: %s", err))
	}
	//根据buyer获取买家购买信息sellingBuy
	var sellingBuy model.SellingBuy
	//如果当前状态是saleStart销售中，是不存在买家的
	if selling.SellingStatus != model.SellingStatusConstant()["saleStart"] {
		resultsSellingByBuyer, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellingBuyKey, []string{buyer})
		if err != nil || len(resultsSellingByBuyer) == 0 {
			return shim.Error(fmt.Sprintf("根据%s获取买家购买信息失败: %s", buyer, err))
		}
		for _, v := range resultsSellingByBuyer {
			if v != nil {
				var s model.SellingBuy
				err := json.Unmarshal(v, &s)
				if err != nil {
					return shim.Error(fmt.Sprintf("UpdateSellingBySeller-反序列化出错: %s", err))
				}
				if s.Selling.ObjectOfSale == objectOfSale && s.Selling.Seller == seller && s.Buyer == buyer {
					//还必须判断状态必须为交付中,防止房子已经交易过，只是被取消了
					if s.Selling.SellingStatus == model.SellingStatusConstant()["delivery"] {
						sellingBuy = s
						break
					}
				}
			}
		}
	}
	var data []byte
	//判断销售状态
	switch status {
	case "done":
		//如果是买家确认收款操作,必须确保销售处于交付状态
		if selling.SellingStatus != model.SellingStatusConstant()["delivery"] {
			return shim.Error("此交易并不处于交付中，确认收款失败")
		}
		//根据seller获取卖家信息
		resultsSellerAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{seller})
		if err != nil || len(resultsSellerAccount) != 1 {
			return shim.Error(fmt.Sprintf("seller卖家信息验证失败%s", err))
		}
		var accountSeller model.Account
		if err = json.Unmarshal(resultsSellerAccount[0], &accountSeller); err != nil {
			return shim.Error(fmt.Sprintf("查询seller卖家信息-反序列化出错: %s", err))
		}
		//确认收款,将款项加入到卖家账户
		accountSeller.Balance += selling.Price
		if err := utils.WriteLedger(accountSeller, stub, model.AccountKey, []string{accountSeller.AccountId}); err != nil {
			return shim.Error(fmt.Sprintf("卖家确认接收资金失败%s", err))
		}
		//将房产信息转入买家，并重置担保状态
		realEstate.Proprietor = buyer
		realEstate.Encumbrance = false
		//realEstate.RealEstateID = stub.GetTxID() //重新更新房产ID
		if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		//清除原来的房产信息
		if err := utils.DelLedger(stub, model.RealEstateKey, []string{seller, objectOfSale}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		//订单状态设置为完成，写入账本
		selling.SellingStatus = model.SellingStatusConstant()["done"]
		selling.ObjectOfSale = realEstate.RealEstateID //重新更新房产ID
		if err := utils.WriteLedger(selling, stub, model.SellingKey, []string{selling.Seller, objectOfSale}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		sellingBuy.Selling = selling
		if err := utils.WriteLedger(sellingBuy, stub, model.SellingBuyKey, []string{sellingBuy.Buyer, sellingBuy.CreateTime}); err != nil {
			return shim.Error(fmt.Sprintf("将本次购买交易写入账本失败%s", err))
		}
		data, err = json.Marshal(sellingBuy)
		if err != nil {
			return shim.Error(fmt.Sprintf("序列化购买交易的信息出错: %s", err))
		}
		break
	case "cancelled":
		data, err = closeSelling("cancelled", selling, realEstate, sellingBuy, buyer, stub)
		if err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		break
	case "expired":
		data, err = closeSelling("expired", selling, realEstate, sellingBuy, buyer, stub)
		if err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		break
	default:
		return shim.Error(fmt.Sprintf("%s状态不支持", status))
	}
	return shim.Success(data)
}

// closeSelling 不管是取消还是过期，都分两种情况
// 1、当前处于saleStart销售状态
// 2、当前处于delivery交付中状态
func closeSelling(closeStart string, selling model.Selling, realEstate model.RealEstate, sellingBuy model.SellingBuy, buyer string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	switch selling.SellingStatus {
	case model.SellingStatusConstant()["saleStart"]:
		selling.SellingStatus = model.SellingStatusConstant()[closeStart]
		//重置房产信息担保状态
		realEstate.Encumbrance = false
		if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
			return nil, err
		}
		if err := utils.WriteLedger(selling, stub, model.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
			return nil, err
		}
		data, err := json.Marshal(selling)
		if err != nil {
			return nil, err
		}
		return data, nil
	case model.SellingStatusConstant()["delivery"]:
		//根据buyer获取卖家信息
		resultsBuyerAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{buyer})
		if err != nil || len(resultsBuyerAccount) != 1 {
			return nil, err
		}
		var accountBuyer model.Account
		if err = json.Unmarshal(resultsBuyerAccount[0], &accountBuyer); err != nil {
			return nil, err
		}
		//此时取消操作，需要将资金退还给买家
		accountBuyer.Balance += selling.Price
		if err := utils.WriteLedger(accountBuyer, stub, model.AccountKey, []string{accountBuyer.AccountId}); err != nil {
			return nil, err
		}
		//重置房产信息担保状态
		realEstate.Encumbrance = false
		if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
			return nil, err
		}
		//更新销售状态
		selling.SellingStatus = model.SellingStatusConstant()[closeStart]
		if err := utils.WriteLedger(selling, stub, model.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
			return nil, err
		}
		sellingBuy.Selling = selling
		if err := utils.WriteLedger(sellingBuy, stub, model.SellingBuyKey, []string{sellingBuy.Buyer, sellingBuy.CreateTime}); err != nil {
			return nil, err
		}
		data, err := json.Marshal(sellingBuy)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, nil
	}
}

//新增
//初始化
func InitLedger(stub shim.ChaincodeStubInterface) pb.Response {
	contract_in_company1 := &model.Contract_in_company{
		ContractName: "Toyota", ContractContent: "Prius", CreaterName: "blue", CreaterSign: "Tomoko", CreateTime: "2006-01-02 15:04:05", CompanyName: "A",
	}
	contract_in_company2 := &model.Contract_in_company{
		ContractName: "Toyo", ContractContent: "Pri", CreaterName: "bl", CreaterSign: "Tomo", CreateTime: "2006-01-02 15:04:05", CompanyName: "A",
	}
	if err := utils.WriteLedger(contract_in_company1, stub, "ContractName", []string{contract_in_company1.ContractName}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	if err := utils.WriteLedger(contract_in_company2, stub, "ContractName", []string{contract_in_company2.ContractName}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}

	contractByte2, err := json.Marshal(contract_in_company2)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(contractByte2)
}

//创建合同
func StartContract(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error(fmt.Sprintf("参数长度不满足"))
	}
	contractname := args[0]
	departmentname := args[1]
	signature := args[2]
	contractcontent := args[3]
	companyname := args[4]
	if departmentname == "" || signature == "" || contractcontent == "" || contractname == "" || companyname == "" {
		return shim.Error(fmt.Sprintf("输入存在空值"))
	}

	//判断私钥签名是否属于此部门
	/*
		未完成
	*/

	//时间戳
	timeUnix := time.Now().Unix() //时间戳
	createTime := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	contract_in_company := model.Contract_in_company{
		ContractName:    contractname,
		ContractContent: contractcontent,
		CreaterName:     departmentname,
		CreaterSign:     signature,
		CreateTime:      createTime,
		CompanyName:     companyname,
	}
	if err := utils.WriteLedger(contract_in_company, stub, "ContractName", []string{contract_in_company.ContractName}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}

	contractByte, err := json.Marshal(contract_in_company)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	return shim.Success(contractByte)
}

//按照合同名查询&查询全部
func QueryContract_incompany(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	results, err := utils.GetStateByPartialCompositeKeys2(stub, "ContractName", args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	var contractList []model.Contract_in_company
	for _, v := range results {
		if v != nil {
			var contract model.Contract_in_company
			err := json.Unmarshal(v, &contract)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryContract_incompany-反序列化出错: %s", err))
			}
			contractList = append(contractList, contract)
		}
	}

	contractListByte, err := json.Marshal(contractList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryContract_incompany-序列化出错: %s", err))
	}
	return shim.Success(contractListByte)

}
//创建合同
func ContractSanction_upload(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	contract_among_company_Byte := args[0]
	fmt.Printf("arg:%v\n", args[0])
	/*companyname := args[1]
	companysignature := args[2]*/

	var contract_among_company model.Contract_among_company
	err := json.Unmarshal([]byte(contract_among_company_Byte), &contract_among_company)

	/*contractname := contract_among_company.ContractName
	contractcontent := contract_among_company.ContractContent*/
	//判断私钥签名是否属于此部门
	/*
		未完成
	*/

	//时间戳
	/*timeUnix := time.Now().Unix() //时间戳
	signtime := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	contract_among_company := model.Contract_among_company{
		ContractName:       contractname,
		ContractContent:    contractcontent,
		CreaterCompanyName: companyname,
		CreaterCompanySign: companysignature,
		CreateTime:         createTime,
	}
	//更新时间
	contract_among_company.SignTime = signtime*/
	if err := utils.WriteLedger(contract_among_company, stub, "ContractName", []string{contract_among_company.ContractName}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}

	contractByte, err := json.Marshal(contract_among_company)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}

	return shim.Success(contractByte)
}

//按照合同名查询&查询全部
func QueryContract_amongcompany(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	results, err := utils.GetStateByPartialCompositeKeys2(stub, "ContractName", args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	var contractList []model.Contract_among_company
	for _, v := range results {
		if v != nil {
			var contract model.Contract_among_company
			err := json.Unmarshal(v, &contract)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryContract_amongcompany-反序列化出错: %s", err))
			}
			contractList = append(contractList, contract)
		}
	}

	contractListByte, err := json.Marshal(contractList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryContract_amongcompany-序列化出错: %s", err))
	}
	return shim.Success(contractListByte)

}

// 审核，生成放入跨公司区块的文本【需要调用公司间区块的链码ContractSanction_upload上传】
func ContractSanction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error(fmt.Sprintf("参数长度不满足"))
	}
	company_name := args[0]
	signature := args[1]
	contractname := args[2]
	var query_args []string
	query_args = append(query_args, contractname)
	results, err := utils.GetStateByPartialCompositeKeys2(stub, "ContractName", query_args)
	if len(results) != 1 {
		return shim.Error(fmt.Sprintf("出现重复的合同名"))
	}
	v := results[0]
	var contract model.Contract_in_company
	err1 := json.Unmarshal(v, &contract)
	if err1 != nil {
		return shim.Error(fmt.Sprintf("反序列化出错: %s", err1))
	}

	timeUnix := time.Now().Unix() //时间戳
	createtime := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")

	contract_among_company := model.Contract_among_company{
		ContractName:       contract.ContractName,
		ContractContent:    contract.ContractContent,
		CreaterCompanyName: company_name,
		CreaterCompanySign: signature,
		SignTime:           createtime,
	}
	contractByte, err := json.Marshal(contract_among_company)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化出错: %s", err))
	}
	fmt.Printf("arg:%v\n", contractByte)
	return shim.Success(contractByte)
}

