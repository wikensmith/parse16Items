package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/forgoer/openssl"
	"github.com/wikensmith/parse16Items/const_"
	"github.com/wikensmith/parse16Items/structs"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

 

var key = []byte(const_.Param.PrivateKey)

// GetRefRules16 获取16项规则
func GetRefRules16(map1 map[string]interface{}, OfficeNo string) (string, error) {
	
	data := map[string]interface{}{
		"FareBasis":      map1["FareBasis"],
		"Total":          map1["Total"],
		"IssueDate":      map1["IssueDate"],
		"BillingAirline": map1["BillingAirline"],
		"Segments":       map1["Segments"],
		//"TicketNo": ticketNo,
		"User": map[string]interface{}{
			"AppName":     "refund_ticket", // 大师兄端接口的用户名及密码
			"AppPwd":      "refund_ticket",
			"AppCaptcha":  "",     // 传空
			"UserName":    "退票中心", // 用户名， 用于计算流量
			"Department":  "技术中心", // 使用部门， 用于计算流量
			"ConfigGroup": 0,      //  默认0：使用的配置组 0:常用配置  7316358544891:订单系统配置, 当传入ETermConfig时， 此字段会被短路
		},
		"OtherOfficeNo": OfficeNo, // 其他office号， 传入此选项时， 在黑屏中会use 该office号
		//"ETermConfig": map[string]interface{}{
		//  "ETermCode":      p.AsmsAcount, // 黑屏帐号
		//  "ETermPwd":       p.AsmsPwd,    // 黑屏密码
		//  "ETermIP":        "",           // ip
		//  "ETermPort":      "",           // 端口
		//  "IsBigConfig":    "",
		//  "OfficeNo":       config.GetConfig().BackGround.OfficeNo, // CKG177
		//  "Authentication": "",                                     //认证模式  0:地址认证 7316358544891: 密码认证
		//  "KeepAlive":      "",                                     // 保持连接  true:是  false:否
		//},
	}
	dataB, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error in GetRefRules16.Marshal, error: [%s]", err.Error())
	}
	iv := make([]byte, 16)
	dst, _ := openssl.AesCBCEncrypt(dataB, key, iv, openssl.PKCS7_PADDING)
	dataStr := base64.StdEncoding.EncodeToString(dst)
	URL16 := const_.Param.BaseURL + "/Handler/GetTicketRule.ashx"
	response, err := http.Post(URL16, "application/json", strings.NewReader(dataStr))
	if err != nil {
		return "", fmt.Errorf("error in GetRefRules16, error: [%s]", err.Error())
	}
	defer func() { _ = response.Body.Close() }()
	responseB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error in GetRefRules16.ReadAll, error: [%s]", err.Error())
	}
	dataB, err = base64.StdEncoding.DecodeString(string(responseB))
	if err != nil {
		return "", fmt.Errorf("error in GetRefRules16.DecodeString, error: [%s]", err.Error())
	}
	dst, err = openssl.AesCBCDecrypt(dataB, key, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", fmt.Errorf("error in GetRefRules16.AesCBCDecrypt, error: [%s]", err.Error())
	}
	return string(dst), nil
}

func CalcRefundAmount(ticketNo string) (string, error) {
	URL := const_.Param.BaseURL + "/Handler/CalcRefundAmount.ashx"
	data := map[string]interface{}{
		"TicketNo": ticketNo,
		"User": map[string]interface{}{
			"AppName":     "refund_ticket", // 大师兄端接口的用户名及密码
			"AppPwd":      "refund_ticket",
			"AppCaptcha":  "",     // 传空
			"UserName":    "退票中心", // 用户名， 用于计算流量
			"Department":  "技术中心", // 使用部门， 用于计算流量
			"ConfigGroup": 0,      //  默认0：使用的配置组 0:常用配置  7316358544891:订单系统配置, 当传入ETermConfig时， 此字段会被短路
		},
		"OtherOfficeNo": "ckg177", // 其他office号， 传入此选项时， 在黑屏中会use 该office号
		//"ETermConfig": map[string]interface{}{
		//  "ETermCode":      p.AsmsAcount, // 黑屏帐号
		//  "ETermPwd":       p.AsmsPwd,    // 黑屏密码
		//  "ETermIP":        "",           // ip
		//  "ETermPort":      "",           // 端口
		//  "IsBigConfig":    "",
		//  "OfficeNo":       config.GetConfig().BackGround.OfficeNo, // CKG177
		//  "Authentication": "",                                     //认证模式  0:地址认证 7316358544891: 密码认证
		//  "KeepAlive":      "",                                     // 保持连接  true:是  false:否
		//},
	}
	dataB, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error in CalcRefundAmount.Marshal, error: [%s]", err.Error())
	}
	iv := make([]byte, 16)
	dst, _ := openssl.AesCBCEncrypt(dataB, key, iv, openssl.PKCS7_PADDING)
	dataStr := base64.StdEncoding.EncodeToString(dst)
	response, err := http.Post(URL, "application/json", strings.NewReader(dataStr))
	if err != nil {
		return "", fmt.Errorf("error in CalcRefundAmount, error: [%s]", err.Error())
	}
	defer func() { _ = response.Body.Close() }()
	responseB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error in CalcRefundAmount.ReadAll, error: [%s]", err.Error())
	}
	dataB, err = base64.StdEncoding.DecodeString(string(responseB))
	if err != nil {
		return "", fmt.Errorf("error in CalcRefundAmount.DecodeString, error: [%s]", err.Error())
	}
	dst, err = openssl.AesCBCDecrypt(dataB, key, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", fmt.Errorf("error in CalcRefundAmount.AesCBCDecrypt, error: [%s]", err.Error())
	}
	return string(dst), nil
}

// http请求汇率
func GetExchangeRate2(Amount float64, FromCurrency string, OfficeNo string) (s string, err error) {

	URL := const_.Param.BaseURL + "/Handler/ExchangeRate.ashx"
	data := map[string]interface{}{
		"Amount":       Amount,
		"FromCurrency": FromCurrency,
		"ToCurrency":   "CNY",
		"User": map[string]interface{}{
			"AppName":     "refund_ticket", // 大师兄端接口的用户名及密码
			"AppPwd":      "refund_ticket",
			"AppCaptcha":  "",     // 传空
			"UserName":    "退票中心", // 用户名， 用于计算流量
			"Department":  "技术中心", // 使用部门， 用于计算流量
			"ConfigGroup": 0,      //  默认0：使用的配置组 0:常用配置  7316358544891:订单系统配置, 当传入ETermConfig时， 此字段会被短路
		},
		"OtherOfficeNo": OfficeNo, // 其他office号， 传入此选项时， 在黑屏中会use 该office号
		//"ETermConfig": map[string]interface{}{
		//  "ETermCode":      p.AsmsAcount, // 黑屏帐号
		//  "ETermPwd":       p.AsmsPwd,    // 黑屏密码
		//  "ETermIP":        "",           // ip
		//  "ETermPort":      "",           // 端口
		//  "IsBigConfig":    "",
		//  "OfficeNo":       config.GetConfig().BackGround.OfficeNo, // CKG177
		//  "Authentication": "",                                     //认证模式  0:地址认证 7316358544891: 密码认证
		//  "KeepAlive":      "",                                     // 保持连接  true:是  false:否
		//},
	}
	dataB, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error in GetExchangeRate2.Marshal, error: [%s]", err.Error())
	}

	iv := make([]byte, 16)
	dst, _ := openssl.AesCBCEncrypt(dataB, key, iv, openssl.PKCS7_PADDING)
	dataStr := base64.StdEncoding.EncodeToString(dst)
	response, err := http.Post(URL, "application/json", strings.NewReader(dataStr))
	if err != nil {
		return "", fmt.Errorf("error in GetExchangeRate2, error: [%s]", err.Error())
	}
	defer func() { _ = response.Body.Close() }()
	responseB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error in GetExchangeRate2.ReadAll, error: [%s]", err.Error())
	}
	dataB, err = base64.StdEncoding.DecodeString(string(responseB))
	if err != nil {
		return "", fmt.Errorf("error in GetExchangeRate2.DecodeString, error: [%s]", err.Error())
	}
	dst, err = openssl.AesCBCDecrypt(dataB, key, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", fmt.Errorf("error in GetExchangeRate2.AesCBCDecrypt, error: [%s]", err.Error())
	}
	return string(dst), nil
}

// ParseRate 解析汇率值
func ParseRate(s string) (float64, error) {
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return 0, fmt.Errorf("error in ParseRate, error: [%s]", err.Error())
	}
	if d, ok := m["CalcResults"]; ok {
		p := d.(map[string]interface{})["FARES"]
		//fares := strconv.FormatFloat(p.(float64), 'f', 0, 64)
		return p.(float64), nil
	} else {
		return 0, fmt.Errorf("error in ParseRate, error: [%s", "解析汇率值失败")
	}
}

// GetExchangeRateAndParseRate 获取汇率并解析
func GetExchangeRateAndParseRate(Amount float64, FromCurrency string, OfficeNo string) (float64, error) {
	fmt.Println("")
	resString, err := GetExchangeRate2(Amount, FromCurrency, OfficeNo)
	if err != nil {
		return 0, fmt.Errorf("error in GetExchangeRateAndParseRate.GetExchangeRate2, error: [%s]", err.Error())
	}
	fares, err := ParseRate(resString)
	if err != nil {
		return 0, fmt.Errorf("error in GetExchangeRateAndParseRate.ParseRate, error: [%s]", err.Error())
	} else {
		fmt.Println("税率转换:", Amount, FromCurrency, fares, "CNY")
		return fares, nil
	}
}

//  GetTicketNoHistory 查询历史记录是否 误机  查询操作记录中 open  状态时间
func GetTicketNoHistory(ticketNo string, OfficeNo string) (string, error) {
	URL := const_.Param.BaseURL + "/Handler/GetTicketNoHistory.ashx"
	data := map[string]interface{}{
		"TicketNo": ticketNo,
		"User": map[string]interface{}{
			"AppName":     "refund_ticket", // 大师兄端接口的用户名及密码
			"AppPwd":      "refund_ticket",
			"AppCaptcha":  "",     // 传空
			"UserName":    "退票中心", // 用户名， 用于计算流量
			"Department":  "技术中心", // 使用部门， 用于计算流量
			"ConfigGroup": 0,      //  默认0：使用的配置组 0:常用配置  7316358544891:订单系统配置, 当传入ETermConfig时， 此字段会被短路
		},
		"OtherOfficeNo": OfficeNo, // 其他office号， 传入此选项时， 在黑屏中会use 该office号
		//"ETermConfig": map[string]interface{}{
		//  "ETermCode":      p.AsmsAcount, // 黑屏帐号
		//  "ETermPwd":       p.AsmsPwd,    // 黑屏密码
		//  "ETermIP":        "",           // ip
		//  "ETermPort":      "",           // 端口
		//  "IsBigConfig":    "",
		//  "OfficeNo":       config.GetConfig().BackGround.OfficeNo, // CKG177
		//  "Authentication": "",                                     //认证模式  0:地址认证 7316358544891: 密码认证
		//  "KeepAlive":      "",                                     // 保持连接  true:是  false:否
		//},
	}
	dataB, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error in GetRefRules16.Marshal, error: [%s]", err.Error())
	}
	iv := make([]byte, 16)
	dst, _ := openssl.AesCBCEncrypt(dataB, key, iv, openssl.PKCS7_PADDING)
	dataStr := base64.StdEncoding.EncodeToString(dst)
	response, err := http.Post(URL, "application/json", strings.NewReader(dataStr))
	if err != nil {
		return "", fmt.Errorf("error in GetTicketNoHistory, error: [%s]", err.Error())
	}
	defer func() { _ = response.Body.Close() }()
	responseB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error in GetTicketNoHistory.ReadAll, error: [%s]", err.Error())
	}
	dataB, err = base64.StdEncoding.DecodeString(string(responseB))
	if err != nil {
		return "", fmt.Errorf("error in GetTicketNoHistory.DecodeString, error: [%s]", err.Error())
	}
	dst, err = openssl.AesCBCDecrypt(dataB, key, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", fmt.Errorf("error in GetTicketNoHistory.AesCBCDecrypt, error: [%s]", err.Error())
	}
	return string(dst), nil
}

// DoParamsCalcRefundAmount 解析DETR信息
func DoParamsCalcRefundAmount(AirTwoCode string, params *structs.DETRStruct) (map[string]interface{}, map[string]interface{}, error) {
	// 判断是否是同一个运价级别 还是 多个
	requiredClasses := mapset.NewSet()
	for _, v := range params.Data.TripInfos {
		requiredClasses.Add(v.FareBasis)
	}
	var FareBasis string
	// 如果是同一个运价级别 传 该运价级别
	s := requiredClasses.ToSlice()
	if len(s) == 1 {
		FareBasis = s[0].(string)
		FareBasis = strings.Replace(FareBasis, "/CH", "", -1)
	} else {
		FareBasis = ""
	}

	Segments_ := make([]interface{}, 0)
	for _, v := range params.CostInfo.TripList {
		segMap := make(map[string]interface{})
		segMap["Airline"] = v.Airline
		segMap["FlightNum"] = v.FlightNo
		segMap["Cabin"] = v.Cabin
		segMap["FromAirport"] = v.FromAirport
		segMap["ToAirport"] = v.ToAirport
		segMap["DepTime"] = v.FlyDate + ":00"
		depTimeFormat, _ := time.Parse("2006-01-02 15:04:05", v.FlyDate+":00")
		arrTime := depTimeFormat.Add(3 * time.Hour)
		segMap["ArrTime"] = arrTime.Format("2006-01-02 15:04:05")
		Segments_ = append(Segments_, segMap)
	}

	Total_ := params.CostInfo.Price + params.CostInfo.Tax
	// 出票日期
	IssueDate_ := params.CostInfo.IssueDate
	IssueDate_ = strings.Replace(IssueDate_, "T", " ", -1)

	DETRInfo := map[string]interface{}{
		"FareBasis":      FareBasis,
		"Total":          Total_,
		"IssueDate":      IssueDate_,
		"BillingAirline": AirTwoCode,
		"Segments":       Segments_,
	}
	priceAndFee := map[string]interface{}{
		"Price":            params.CostInfo.Price,
		"UsedHistoryPrice": params.UsedHistoryPrice,
		"Taxs":             params.Taxs,
		"UsedFare":         params.UsedFare,
	}
	return DETRInfo, priceAndFee, nil
}


// DoParamsGetRefRules16 接收返回黑屏 指令和结果的列表
func DoParamsGetRefRules16(s string) ([]string, error) {
	params2 := new(structs.ResGetRefRules16)
	if err := json.Unmarshal([]byte(s), params2); err != nil {
		return nil, fmt.Errorf("error in DoParamsGetRefRules16.Unmarshal, error: [%s]", err.Error())
	}
	if params2.Error != 0 {
		return nil, fmt.Errorf("error in DoParamsGetRefRules16, error: [%s]", params2.Message)
	}
	if len(params2.EtermStr) == 0 {
		return nil, fmt.Errorf("error in DoParamsGetRefRules16, error: [%s]", params2.Message)
	}
	return params2.EtermStr, nil
}

var BASEURL = "http://118.178.135.114:9901/api/basedata/getairline/"
var AirPortURL = "http://118.178.135.114:9901/api/basedata/getairport/"

// 查询机场信息
func QueryAirport(data string) (map[string]string, error) {
	URL := AirPortURL + "by_code/" + data
	client := http.Client{}
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("调用数据中心接口查询航司名称创建newrequests失败, 异常信息: %s", err.Error()))
	}
	request.Header.Set("UserName", "refundcenter")
	now := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	request.Header.Set("TimeStamp", now)
	v := MD5("refundcenter:" + now + "@ghdfg6ryehs3gf6hfd")
	request.Header.Set("Sign", v)
	res, err := client.Do(request)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("调用数据中心接口查询航司名称失败, 异常信息: %s", err.Error()))
	}
	defer func() {
		_ = res.Body.Close()
	}()
	resB, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("调用数据中心接口查询航司名称ReadAll异常, 异常信息: %s", err.Error()))
	}
	//resultLst := []map[string]string{}
	resultLst := make([]map[string]string, 0)
	if err := json.Unmarshal(resB, &resultLst); err != nil {
		return nil, errors.New(fmt.Sprintf("调用数据中心接口查询航司名称unmarshal异常, 异常信息: %s, 返回数据: %s", err.Error(), string(resB)))
	}
	if len(resultLst) == 0 {
		return nil, fmt.Errorf("QueryAirport 查询结果为空, 请核实！")
	}
	return resultLst[0], nil
}

// QueryAirLine 通过票号前三位查询航司名称
func QueryAirLine(data string) (map[string]string, error) {
	URL := BASEURL + "by_code3/" + data
	client := http.Client{}
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("调用数据中心接口查询航司名称创建newrequests失败, 异常信息: %s", err.Error()))
	}
	request.Header.Set("UserName", "refundcenter")
	now := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	request.Header.Set("TimeStamp", now)
	v := MD5("refundcenter:" + now + "@ghdfg6ryehs3gf6hfd")
	request.Header.Set("Sign", v)
	res, err := client.Do(request)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("调用数据中心接口查询航司名称失败, 异常信息: %s", err.Error()))
	}
	defer func() {
		_ = res.Body.Close()
	}()
	resB, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("调用数据中心接口查询航司名称ReadAll异常, 异常信息: %s", err.Error()))
	}
	//resultLst := []map[string]string{}
	resultLst := make([]map[string]string, 0)
	if err := json.Unmarshal(resB, &resultLst); err != nil {
		return nil, errors.New(fmt.Sprintf("调用数据中心接口查询航司名称unmarshal异常, 异常信息: %s, 返回数据: %s", err.Error(), string(resB)))
	}
	return resultLst[0], nil
}

func MD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	result := h.Sum(nil)
	return hex.EncodeToString(result)
}

func DoTicketNoSuspended(params *structs.DETRStruct, OfficeNo string) error {
	if params.Data == nil {
		return fmt.Errorf("error in DoTicketNoSuspended, params.Data 为nil")
	}
	if params.Data.TripInfos == nil {
		return fmt.Errorf("error in DoTicketNoSuspended, params.Data.TripInfos 为nil")
	}
	for i := range params.Data.TripInfos {
		if params.Data.TripInfos[i].TicketNoStatus == "SUSPENDED" {
			res, err := TicketNoSuspended(params.Data.TripInfos[i].TicketNo, OfficeNo)
			fmt.Println("解除SUSPENDED状态接受数据:", res, err)
			if err != nil {
				return fmt.Errorf("error in TicketNoSuspended, error: [%s]", err.Error())
			}
			m := make(map[string]interface{})
			err = json.Unmarshal([]byte(res), &m)
			if err != nil {
				return fmt.Errorf("error in DoTicketNoSuspended.Unmarshal, error: [%s]", err.Error())
			}
			if m["Error"].(float64) != 0 {
				return fmt.Errorf("error in DoTicketNoSuspended, error: [%s]", m["Message"])
			} else {
				params.Data.TripInfos[i].TicketNoStatus = "OPEN FOR USE"
			}

			// 遍历所有票号， 如果此票号与解锁票号相等， 则赋值OPEN FOR USE
			for j, v := range params.Data.TripInfos {
				if params.Data.TripInfos[i].TicketNo == v.TicketNo {
					params.Data.TripInfos[j].TicketNoStatus = "OPEN FOR USE"
				}
			}
		}
	}
	return nil
}

// 解锁
func TicketNoSuspended(ticketNo string, OfficeNo string) (string, error) {
	data := map[string]interface{}{
		"TicketNo":    ticketNo,
		"IsSuspended": false,
		"User": map[string]interface{}{
			"AppName":     "refund_ticket", // 大师兄端接口的用户名及密码
			"AppPwd":      "refund_ticket",
			"AppCaptcha":  "",     // 传空
			"UserName":    "退票中心", // 用户名， 用于计算流量
			"Department":  "技术中心", // 使用部门， 用于计算流量
			"ConfigGroup": 0,      //  默认0：使用的配置组 0:常用配置  7316358544891:订单系统配置, 当传入ETermConfig时， 此字段会被短路
		},
		"OtherOfficeNo": OfficeNo, // 其他office号， 传入此选项时， 在黑屏中会use 该office号
		//"ETermConfig": map[string]interface{}{
		//  "ETermCode":      p.AsmsAcount, // 黑屏帐号
		//  "ETermPwd":       p.AsmsPwd,    // 黑屏密码
		//  "ETermIP":        "",           // ip
		//  "ETermPort":      "",           // 端口
		//  "IsBigConfig":    "",
		//  "OfficeNo":       config.GetConfig().BackGround.OfficeNo, // CKG177
		//  "Authentication": "",                                     //认证模式  0:地址认证 7316358544891: 密码认证
		//  "KeepAlive":      "",                                     // 保持连接  true:是  false:否
		//},
	}
	dataB, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error in GetRefRules16.Marshal, error: [%s]", err.Error())
	}
	iv := make([]byte, 16)
	dst, _ := openssl.AesCBCEncrypt(dataB, key, iv, openssl.PKCS7_PADDING)
	dataStr := base64.StdEncoding.EncodeToString(dst)
	URL := const_.Param.BaseURL + "/Handler/TicketNoSuspended.ashx"
	response, err := http.Post(URL, "application/json", strings.NewReader(dataStr))
	if err != nil {
		return "", fmt.Errorf("error in TicketNoSuspended, error: [%s]", err.Error())
	}
	defer func() { _ = response.Body.Close() }()
	responseB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error in TicketNoSuspended.ReadAll, error: [%s]", err.Error())
	}
	dataB, err = base64.StdEncoding.DecodeString(string(responseB))
	if err != nil {
		return "", fmt.Errorf("error in TicketNoSuspended.DecodeString, error: [%s]", err.Error())
	}
	dst, err = openssl.AesCBCDecrypt(dataB, key, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", fmt.Errorf("error in TicketNoSuspended.AesCBCDecrypt, error: [%s]", err.Error())
	}
	return string(dst), nil
}

// 外航 查询历史记录是否 误机  查询操作记录中 open  状态时间
func GetPnrHistory(Pnr string, FlightNo string, FromDate string, OfficeNo string) (string, error) {
	//URL := config.GetConfig().BaseURL + "/Handler/GetPnrHistory.ashx"
	URL := const_.Param.BaseURL + "/Handler/GetPnrHistory.ashx"
	data := map[string]interface{}{
		"Pnr": Pnr,
		"FromDate": FromDate,
		"FlightNo": FlightNo,
		"User": map[string]interface{}{
			"AppName":     "refund_ticket", // 大师兄端接口的用户名及密码
			"AppPwd":      "refund_ticket",
			"AppCaptcha":  "",     // 传空
			"UserName":    "退票中心", // 用户名， 用于计算流量
			"Department":  "技术中心", // 使用部门， 用于计算流量
			"ConfigGroup": 0,      //  默认0：使用的配置组 0:常用配置  7316358544891:订单系统配置, 当传入ETermConfig时， 此字段会被短路
		},
		"OtherOfficeNo": OfficeNo, // 其他office号， 传入此选项时， 在黑屏中会use 该office号
		//"ETermConfig": map[string]interface{}{
		//  "ETermCode":      p.AsmsAcount, // 黑屏帐号
		//  "ETermPwd":       p.AsmsPwd,    // 黑屏密码
		//  "ETermIP":        "",           // ip
		//  "ETermPort":      "",           // 端口
		//  "IsBigConfig":    "",
		//  "OfficeNo":       config.GetConfig().BackGround.OfficeNo, // CKG177
		//  "Authentication": "",                                     //认证模式  0:地址认证 7316358544891: 密码认证
		//  "KeepAlive":      "",                                     // 保持连接  true:是  false:否
		//},
	}
	dataB, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error in GetRefRules16.Marshal, error: [%s]", err.Error())
	}
	iv := make([]byte, 16)
	dst, _ := openssl.AesCBCEncrypt(dataB, key, iv, openssl.PKCS7_PADDING)
	dataStr := base64.StdEncoding.EncodeToString(dst)
	response, err := http.Post(URL, "application/json", strings.NewReader(dataStr))
	if err != nil {
		return "", fmt.Errorf("error in GetTicketNoHistory, error: [%s]", err.Error())
	}
	defer func() { _ = response.Body.Close() }()
	responseB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error in GetTicketNoHistory.ReadAll, error: [%s]", err.Error())
	}
	dataB, err = base64.StdEncoding.DecodeString(string(responseB))
	if err != nil {
		return "", fmt.Errorf("error in GetTicketNoHistory.DecodeString, error: [%s]", err.Error())
	}
	dst, err = openssl.AesCBCDecrypt(dataB, key, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", fmt.Errorf("error in GetTicketNoHistory.AesCBCDecrypt, error: [%s]", err.Error())
	}
	return string(dst), nil
}

func UWingQuery(data map[string]string)  (string, error){
	//URL := "http://192.168.0.79:8058/admin_api/uwing"
	URL := const_.Param.UWURL
	dataB, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error in UWingQuery.Marshal, error: [%s]", err.Error())
	}
	response, err := http.Post(URL, "application/json", bytes.NewReader(dataB))
	if err != nil {
		return "", fmt.Errorf("error in UWingQuery, error: [%s]", err.Error())
	}
	defer func() { _ = response.Body.Close() }()
	responseB, err := ioutil.ReadAll(response.Body)
	if err != nil{
		return "", fmt.Errorf("error in UWingQuery.ReadAll, error: [%s]", err.Error())
	}
	return string(responseB), nil




}