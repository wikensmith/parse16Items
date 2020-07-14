package parseClient

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var Data = `{
    "ApiInfo": null,
    "Attachments": [],
    "AuditorName": "国际机票┝正航组┝综合组┝出票/退改┠陈佳顺(9615)",
    "AuditorTime": "2020-06-27T13:47:05.5855356",
    "BuyAmount": -1728,
    "BuyOrders": [
      {
        "Submitter": "",
        "IsVoluntary": false,
        "RefundBuyOrderNo": "202006271240190093",
        "OuterBuyOderId": "",
        "BuyChannel": {
          "Account": "CKG177",
          "AccountAbbr": "CKG177|CKG177",
          "ChannelInt": "",
          "ChannelType": "BOP",
          "ChannelTypeName": "BOP",
          "SaleBuyCapacity": "Buy",
          "OrderCapacity": "It",
          "Name": "BOP_CKG177",
          "NameAbbr": "BOP_CKG177|BOP_CKG177",
          "Remark": "",
          "Organization": "",
          "OrganizationId": "",
          "PaymentId": 166,
          "PaymentName": "德付通-68650323@sohu.com",
          "Id": 24
        },
        "Pnr": {
          "BigPnrCode": "MGPTDF",
          "GdsName": "TravelSky",
          "Id": "1547728",
          "OfficeNo": "CKG177",
          "PnrCode": "KX1JFY",
          "Source": "BySale",
          "SourceNo": "202006121448500055"
        },
        "GdsName": "TravelSky",
        "OfficeNo": "CKG177",
        "Currency": "CNY",
        "ExchangeRate": 1,
        "BuyAmount": -1728,
        "BuyForeignAmount": -1728,
        "ReceivedAmount": 0,
        "RefundType": "",
        "SubmitStatus": "NotSubmit",
        "RefundStatus": "NonePay",
        "Passengers": [
          {
            "PassengerId": null,
            "ChangePassengerId": "",
            "PassengerVoyages": [
              {
                "TicketNo": "731-3391534281",
                "Id": 1517787,
                "VoyageId": 2420227
              }
            ],
            "RefundCenterDETR": "{\"UsedFare\":280.0,\"UsedHistoryPrice\":350.0,\"Taxs\":{\"MO\":27.00,\"WN\":96.00,\"YR\":265.00},\"Data\":{\"PassengerName\":\"XIE/YUECHAO MR\",\"Endorsement\":\"Q/NON-END\",\"OldTicketNo\":\"\",\"TripInfos\":[{\"TripCode\":\"O\",\"TripNo\":\"1\",\"Airline\":\"MF\",\"FlightNo\":\"841\",\"Cabin\":\"R\",\"CabinMark\":null,\"FromCity\":\"HGH\",\"FromAirport\":\"HGH\",\"ToCity\":\"MFM\",\"ToAirport\":\"MFM\",\"FormTerminal\":\"T2\",\"ToTerminal\":\"--\",\"TicketNoStatus\":\"USED/FLOWN\",\"FlightDate\":\"2019-11-25\",\"FareBasis\":\"R14MOC\",\"DepartureTime\":\"13:50\",\"ArrivalTime\":null,\"Luggage\":\"1PC\",\"TicketNo\":\"7313391534281\",\"Pnr\":\"KE22D0\"},{\"TripCode\":\"O\",\"TripNo\":\"2\",\"Airline\":\"MF\",\"FlightNo\":\"\",\"Cabin\":\"R\",\"CabinMark\":null,\"FromCity\":\"MFM\",\"FromAirport\":\"MFM\",\"ToCity\":\"HGH\",\"ToAirport\":\"HGH\",\"FormTerminal\":\"\",\"ToTerminal\":\"\",\"TicketNoStatus\":\"OPEN FOR USE\",\"FlightDate\":\"\",\"FareBasis\":\"R14MOC\",\"DepartureTime\":\"\",\"ArrivalTime\":\"\",\"Luggage\":\"1PC\",\"TicketNo\":\"7313391534281\",\"Pnr\":null}],\"itinerary\":null},\"CostInfo\":{\"Currency\":\"CNY\",\"ROEValue\":7.1468,\"NUCValue\":76.94,\"CNFee\":90.00,\"YQFee\":0.0,\"YRFee\":530.00,\"Taxs\":{\"CN\":90.00,\"MO\":27.00,\"WN\":96.00,\"YR\":530.00},\"Price\":550.00,\"Tax\":743.00,\"AgencyFee\":0.0,\"EXCH\":null,\"CONJTKT\":\"\",\"IssueDate\":\"2019-11-22T00:00:00\",\"ThreeSideAgreementNo\":null,\"ProductCode\":null,\"Pnr\":\"KE22D0\",\"TripList\":[{\"Airline\":\"MF\",\"FromAirport\":\"HGH\",\"ToAirport\":\"MFM\",\"Share\":false,\"FlightNo\":\"841\",\"Cabin\":\"R\",\"FlyDate\":\"2019-11-25 13:50\"},{\"Airline\":\"MF\",\"FromAirport\":\"MFM\",\"ToAirport\":\"HGH\",\"Share\":false,\"FlightNo\":\"842\",\"Cabin\":\"R\",\"FlyDate\":\"2019-11-27 17:40\"}],\"TripPriceList\":[{\"Airline\":\"MF\",\"FromCity\":\"HGH\",\"ToCity\":\"MFM\",\"Share\":false,\"FlightNo\":\"841\",\"Cabin\":\"R\",\"FlyDate\":\"2019-11-25 13:50\",\"QValue\":0.0,\"SValue\":0.0,\"OtherValue\":null,\"Value\":38.47,\"Mileage\":0},{\"Airline\":\"MF\",\"FromCity\":\"MFM\",\"ToCity\":\"HGH\",\"Share\":false,\"FlightNo\":\"842\",\"Cabin\":\"R\",\"FlyDate\":\"2019-11-27 17:40\",\"QValue\":0.0,\"SValue\":0.0,\"OtherValue\":null,\"Value\":38.47,\"Mileage\":0}]},\"Error\":0,\"Message\":\"获取要退的税详情成功\",\"EtermStr\":[\">DETR TN/7313391534281\",\"DETR TN/7313391534281\\rISSUED BY: XIAMEN AIRLINES           ORG/DST: HGH/HGH                 BSP-I\\rE/R: Q/NON-END Q/FARE RESTRICTIONS APPLY\\rTOUR CODE:                                                   \\rPASSENGER: XIE/YUECHAO MR\\rEXCH:                               CONJ TKT: \\rO FM:1HGH MF     841  R 25NOV 1350 OK R14MOC           /09DEC9 1PC USED/FLOWN\\r     T2-- RL:MLV0W2  /KE22D01E \\rO TO:2MFM MF    OPEN  R OPEN          R14MOC           /09DEC9 1PC OPEN FOR USE\\r          RL:                  \\r  TO: HGH\\rFC: A  25NOV19HGH MF MFM38.47MF HGH38.47NUC76.94END ROE7.146800 XT 96.00WN5\\r30.00YR\\rFARE:           CNY  550.00|FOP:CASH\\rTAX:            CNY 90.00CN|OI: \\rTAX:            CNY 27.00MO|\\rTAX:            CNY 96.00WN|FOR ALL TAXES: DETR:TN/731-3391534281,X\\rTOTAL:          CNY 1293.00|TKTN: 731-3391534281\",\">XS FSI/MF//.22NOV19\\rS MF 841R25NOV HGH1350 1650MFM0S\",\"SI/MF//.22NOV19   \\rS MF   841R25NOV HGH1350 1650MFM0S  \\r01 ROWMOC                705 CNY                    INCL TAX\\r***** 历史运价计算仅供参考 *****\\r*SYSTEM DEFAULT-CHECK EQUIPMENT/OPERATING CARRIER   \\r*ATTN PRICED ON 13JUL20*0923\\r HGH\\r MFM ROWMOC                            NVB        NVA25NOV20 1PC\\rFARE  CNY     350   \\rTAX   CNY      90CN CNY     265YR   \\rTOTAL CNY     705   \\r25NOV19HGH MF MFM48.97NUC48.97END ROE7.146800   \\rENDOS *Q/NON-END\\rENDOS *Q/FARE RESTRICTIONS APPLY\\r*AUTO BAGGAGE INFORMATION AVAILABLE - SEE FSB   \\r*COMMISSION VALIDATED - DATA SOURCE TRAVELSKY   \\rTKT/TL25NOV19*1350  \\rCOMMISSION  0.00 PERCENT OF GROSS   \\rFSKY/1E/R3KUNCSIVLCVU44/FCC=D/\",\"SI/MF//.22NOV19   \\rS MF   841R25NOV HGH1350 1650MFM0S  \\r01 ROWMOC                705 CNY                    INCL TAX\\r***** 历史运价计算仅供参考 *****\\r*SYSTEM DEFAULT-CHECK EQUIPMENT/OPERATING CARRIER   \\r*ATTN PRICED ON 13JUL20*0923\\r HGH\\r MFM ROWMOC                            NVB        NVA25NOV20 1PC\\rFARE  CNY     350   \\rTAX   CNY      90CN CNY     265YR   \\rTOTAL CNY     705   \\r25NOV19HGH MF MFM48.97NUC48.97END ROE7.146800   \\rENDOS *Q/NON-END\\rENDOS *Q/FARE RESTRICTIONS APPLY\\r*AUTO BAGGAGE INFORMATION AVAILABLE - SEE FSB   \\r*COMMISSION VALIDATED - DATA SOURCE TRAVELSKY   \\rTKT/TL25NOV19*1350  \\rCOMMISSION  0.00 PERCENT OF GROSS   \\rFSKY/1E/R3KUNCSIVLCVU44/FCC=D/\"],\"EtermTraffic\":4}",
            "Name": "JIANG/LIBO",
            "PassengerType": "Adult",
            "CertType": "Passport",
            "CertNo": "EF9082940",
            "SaleTicketPrice": -1495,
            "SaleTicketTax": -238,
            "SaleInsurePrice": 0,
            "SaleChangePrice": 0,
            "SaleAgencyFee": 0,
            "SaleRebate": 0,
            "SaleCashBack": 0,
            "SaleAdjustPrice": 0,
            "SaleServiceFee": 0,
            "SaleAdjustServiceFee": 0,
            "SaleDelayFee": 0,
            "SaleOtherFee": 0,
            "SaleSum": -1733,
            "BuyTicketPrice": -1490,
            "BuyTicketTax": -238,
            "BuyAgencyFee": 0,
            "SaleRefundFee": 0,
            "BuyRebate": 0,
            "BuyCashBack": 0,
            "BuyInsurePrice": 0,
            "BuyChangePrice": 0,
            "BuyAdjustPrice": 0,
            "BuyServiceFee": 0,
            "BuyAdjustServiceFee": 0,
            "BuyDelayFee": 0,
            "BuyRefundFee": 0,
            "BuyOtherFee": 0,
            "BuySum": -1728,
            "BuyTaxDetail": "$id:7,G3:82,HK:110,I5:46",
            "Id": 104243
          }
        ],
        "Paylogs": "",
        "CreationTime": "2020-06-27T12:40:19.5234388",
        "Id": 75432,
        "ExtensionData": "{\"TicketTime\":\"2020-06-12 14:53\",\"BuyRefundReason\":\"因航班取消延误，申请全退\",\"SubmitMsg\":\"申请提交__2020-06-28 09:43:21┠ 信息中心-坤昌┝研发部┝YATP组┠文康(7921)@_@\",\"SubmitUser\":\"信息中心-坤昌┝研发部┝YATP组┠文康(7921)\",\"SubmitTime\":\"2020-06-28 09:43:21\"}"
      }
    ],
    "BuyPnr": "",
    "CancelReason": null,
    "Contact": {
      "Email": "qunaer@qq.com",
      "Name": "去哪儿固定联系人",
      "Phone": "15730076283",
      "Tel": "888888"
    },
    "CreateDepartment": null,
    "CreationTime": "2020-06-27T12:40:19.4765575",
    "CreatorName": "Api_ImportOrder(Api_ImportOrder)",
    "DeliveryInfo": "",
    "ExchangeRate": 1,
    "Id": 75423,
    "IntlFlag": true,
    "IsAbandon": false,
    "IsVoluntary": false,
    "KeepSeat": true,
    "NeedDelivery": false,
    "OrganizationUnitName": "国际机票┝正航组┝去哪儿组",
    "OuterSaleOderId": "ysi200612144738810eb0ab",
    "Passengers": [
      {
        "AirRax": 0,
        "BuyAdjustPrice": 0,
        "BuyAdjustServiceFee": 0,
        "BuyAgencyFee": 0,
        "BuyAirRax": 0,
        "BuyCashBack": 0,
        "BuyChangePrice": 0,
        "BuyDelayFee": 0,
        "BuyInsurePrice": 0,
        "BuyOtherFee": 0,
        "BuyRebate": 0,
        "BuyRefundFee": 0,
        "BuyServiceFee": 0,
        "BuySum": -1728,
        "BuyTaxDetail": "$id:7,G3:82,HK:110,I5:46",
        "BuyTicketPrice": -1490,
        "BuyTicketTax": -238,
        "CertNo": "EF9082940",
        "CertType": "Passport",
        "Id": 104243,
        "Name": "JIANG/LIBO",
        "PassengerType": "Adult",
        "PassengerVoyages": [
          {
            "Id": 1517787,
            "TicketNo": "8513397350173",
            "VoyageId": 2420227
          }
        ],
        "SaleAdjustPrice": 0,
        "SaleAdjustServiceFee": 0,
        "SaleAgencyFee": 0,
        "SaleCashBack": 0,
        "SaleChangePrice": 0,
        "SaleDelayFee": 0,
        "SaleInsurePrice": 0,
        "SaleOtherFee": 0,
        "SaleRebate": 0,
        "SaleRefundFee": 0,
        "SaleServiceFee": 0,
        "SaleSum": -1733,
        "SaleTicketPrice": -1495,
        "SaleTicketTax": -238
      }
    ],
    "PassengersName": "JIANG/LIBO",
    "Paylogs": [],
    "Pnr": {
      "BigPnrCode": "MGPTDF",
      "GdsName": "TravelSky",
      "Id": "1547728",
      "OfficeNo": "CKG177",
      "PnrCode": "KX1JFY",
      "Source": "BySale",
      "SourceNo": "202006121448500055"
    },
    "ProfitCenters": {
      "Code": "00030002",
      "DepartmentCodeList": "|202|120|",
      "DepartmentNameList": "票务国际机票┝正航组┝去哪儿组,畅游国际机票┝正航组┝去哪儿组",
      "Id": "201811081627290003",
      "Name": "国际正航去哪儿组",
      "Params": {
        "AsmsAcount": "9690",
        "AsmsPwd": "For123",
        "Dh": "",
        "Sj": "18983760075",
        "UserId": "9639",
        "Xm": "叶旺",
        "Yx": "",
        "id": null
      },
      "Remark": "去哪国际、去哪特惠、去哪畅游国际",
      "SubsetList": null,
      "id": "201811081627290003"
    },
    "Reason": "因航班取消延误，申请全退",
    "ReceivedAmount": 0,
    "RefundAmount": 0,
    "RefundOrderNo": "202006271240190093",
    "RefundStatus": "NonePay",
    "Remark": "",
    "SaleAmount": -1733,
    "SaleCurrency": "CNY",
    "SaleForeignAmount": -1733,
    "SalePnr": "KX1JFY",
    "Status": "Checked",
    "TicketNo": "8513397350173",
    "Voyages": [
      {
        "Airline": "HX",
        "ArrAirport": "PVG",
        "ArrTerminal": "T1",
        "ArrTime": "2020-07-08T21:55:00",
        "Cabin": "L",
        "DepAirport": "HKG",
        "DepTerminal": "T1",
        "DepTime": "2020-07-08T19:05:00",
        "FlightNo": "HX232",
        "Id": 2420227,
        "SegmentRph": 1
      }
    ]
  }`

var key = []byte("cqyshkgsetermsdk")

//  发送请求查询黑屏并计算税金
func DoQuery(ticketNo string,  officeNo string) (string, error) {
	URL := "http://eterm.70168.com:8081"+ "/Handler/CalcRefundAmount.ashx"
	data := map[string]interface{}{
		"TicketNo": ticketNo,
		"User": map[string]interface{}{
			"AppName":     "refund_ticket", // 大师兄端接口的用户名及密码
			"AppPwd":      "refund_ticket",
			"AppCaptcha":  "",     // 传空
			"UserName":    "退票中心", // 用户名， 用于计算流量
			"Department":  "技术中心", // 使用部门， 用于计算流量
			"ConfigGroup": 0,      //  默认0：使用的配置组 0:常用配置  1:订单系统配置, 当传入ETermConfig时， 此字段会被短路
		},
		"OtherOfficeNo": officeNo, // 其他office号， 传入此选项时， 在黑屏中会use 该office号
		//"ETermConfig": map[string]interface{}{
		//	"ETermCode":      p.AsmsAcount, // 黑屏帐号
		//	"ETermPwd":       p.AsmsPwd,    // 黑屏密码
		//	"ETermIP":        "",           // ip
		//	"ETermPort":      "",           // 端口
		//	"IsBigConfig":    "",
		//	"OfficeNo":       config.GetConfig().BackGround.OfficeNo, // CKG177
		//	"Authentication": "",                                     //认证模式  0:地址认证 1: 密码认证
		//	"KeepAlive":      "",                                     // 保持连接  true:是  false:否
		//},
	}
	dataB, err := json.Marshal(data)
	if err != nil {
		return "", errors.New(fmt.Sprintf("marshal 黑屏请求参数失败, 异常信息: %s", err.Error()))
	}
	cipherData := AESEncrypt(dataB, key)
	//var d []byte
	//base64.StdEncoding.Encode(d, cipherData)
	d := base64.StdEncoding.EncodeToString(cipherData)
	response, err := http.Post(URL, "application/json", strings.NewReader(d))

	if err != nil {
		return "", err
	}
	defer func() { _ = response.Body.Close() }()
	responseB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	//var bb []byte
	n, err := base64.StdEncoding.DecodeString(string(responseB))
	//fmt.Println("here:", n, err)
	de := AESDecrypt(n, key)
	return string(de), nil
}
