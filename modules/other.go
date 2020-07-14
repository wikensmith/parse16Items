package modules

import (
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/wikensmith/toLogCenter"
	"github.com/wikensmith/parse16Items/structs"
	"github.com/wikensmith/parse16Items/utils"
	"regexp"
	"strings"
	"time"
	"strconv"
)

type Calc interface {
	DoCalc() (error)  // 执行模板计算
	Name()string  // 返回模板名称
	//m.Passengers, m.etermMap, m.OfficeNo, m.Log
	New(
		Passengers *structs.Passengers_,
		EtermMap map[string]string,
		OfficeNo string,
		Log *toLogCenter.Logger) Calc
}

// getHighestFee 如果有多个退票费， 转换汇率成CNY，返回最高价
func GetHighestFee(refundFeeLst [][]string) (maxPrice float64, currency string, err error) {
	// [["100", "CNY"], ["CNY", "200"]]
	if len(refundFeeLst) == 1 {
		if len(refundFeeLst[0]) == 2 {
			l := refundFeeLst[0]
			if l[0][1] <= 90 && l[0][1] >= 65 {
				maxPrice, _ = strconv.ParseFloat(l[1], 10)
				currency = l[0]
			} else {
				maxPrice, _ = strconv.ParseFloat(l[0], 10)
				currency = l[1]
			}
			if currency != "CNY" {
				tempPrice, err1 := utils.GetExchangeRateAndParseRate(maxPrice, currency, "CKG177")
				if err1 != nil {
					err =  fmt.Errorf("error in getHighestFee.GetExchangeRateAndParseRate error:[%s]","转换汇率失败")
					return
				}
				maxPrice = tempPrice
				currency = "CNY"
				return
			}

		}
	}
	for _, v := range refundFeeLst {
		if len(v) == 2 {
			// 是大写字母
			if v[0][1] <= 90 && v[0][1] >= 65 {
				if v[0] != "CNY" {
					price, _ := strconv.ParseFloat(v[1], 10)
					priceCNY, err := utils.GetExchangeRateAndParseRate(price, v[0], "CKG177")
					if err != nil {
						return 0, "", fmt.Errorf("error in getHighestFee.GetExchangeRateAndParseRate error:[%s]",
							"转换汇率失败")
					}
					if price > maxPrice {
						maxPrice = priceCNY
						currency = "CNY"
					}
				} else {
					price, _ := strconv.ParseFloat(v[1], 10)
					if price > maxPrice {
						maxPrice = price
						currency = v[0]
					}
				}
			}
		}
	}
	if maxPrice == 0 {
		err = fmt.Errorf("error in getHighestFee, error: [%s: %#v]", "获取最大退票费失败, ", refundFeeLst)
		return
	}
	return
}

// ISMissedFlight 是否误机
func ISMissedFlight(
	ticketNo string,
	OfficeNo string,
	params *structs.DETRStruct,
	NoShowRuleTime int) (bool, error) {

	// 国内航司调用此接口， 国际航司调用其他接口
	TickHistoryStr, err := utils.GetTicketNoHistory(ticketNo, OfficeNo)
	//fmt.Println("接收历史操作记录数据...", TickHistoryStr, err)
	// {"Error":0,"Message":"succuss","EtermStr":[">DETR TN/826-5582086910,H","DETR TN/826-5582086910,H\rNAME: JIANG/JUNTIAN MR TKTN:8265582086910\rIATA OFFC: 08316954 ISSUED: 20JUL19 RVAL: 00\r  7 2  08AUG/0912/9940     EOTU CHG FLT FROM GS7955/08AUG19/U/TSNSVO TO\r                                GSOPEN/OPEN/U/TSNSVO\r  6 2  08AUG/0912/9940     EOTU MKG RL JTGPCR1E CLEARED\r  5 2  08AUG/0912/9940     EOTU RES RL MLSRNW   CLEARED\r  4 1  07AUG/2156/9940     EOTU CHG FLT FROM GS7854/07AUG19/Y/HGHTSN TO\r                                GSOPEN/OPEN/Y/HGHTSN\r  3 1  07AUG/2156/9940     EOTU MKG RL JTGPCR1E CLEARED\r  2 1  07AUG/2156/9940     EOTU RES RL MLSRNW   CLEARED\r  1    20JUL/1121/52666    TRMK YAKE+CKG177+DEV-16"],"EtermTraffic":2}U
	if err != nil {
		return false, fmt.Errorf("error in ISMissedFlight.GetTicketNoHistory, error: [%s]", err.Error())
	}

	TickHistoryMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(TickHistoryStr), &TickHistoryMap)
	if err != nil {
		return false, fmt.Errorf("error in ISMissedFlight.Unmarshal, error: [%s]", err.Error())
	}

	if TickHistoryMap["Error"].(float64) != 0 {
		return false, fmt.Errorf("error in ISMissedFlight, error: [%s]", TickHistoryMap["Message"])
	}

	TickHistoryStr = strings.Replace(TickHistoryStr, "\\r", "", -1)
	TickHistoryStr = strings.Replace(TickHistoryStr, " ", "", -1)
	//reg := regexp.MustCompile(`([\d]{2})([A-Z]{3})/([\d]{4})/\d+EOTUCHGFLTFROM[A-Z]{2}[\d]*/[\d]{2}[A-Z]{3}([\d]{2})/[A-Z]+/[A-Z]*/OPEN/[A-Z]+/([A-Z]{6})`)

	// 取消的占座时间
	reg := regexp.MustCompile(`([\d]{2})([A-Z]{3})/([\d]{4})/\d+[A-Z]{4}CHGFLTFROM[A-Z0-9]{2}[\d]*/[\d]{2}[A-Z]{3}([\d]{2})/[A-Z]+/[A-Z]{6}TO[A-Z0-9]{2}OPEN/[A-Z0-9]{4,7}/[A-Z]+/([A-Z]{6})`)
	subLst := reg.FindAllStringSubmatch(TickHistoryStr, -1)
	var MonthMap = map[string]string{
		"JAN": "01",
		"FEB": "02",
		"MAR": "03",
		"APR": "04",
		"MAY": "05",
		"JUN": "06",
		"JUL": "07",
		"AUG": "08",
		"SEP": "09",
		"OCT": "10",
		"NOV": "11",
		"DEC": "12",
	}
	// 取消机位时间, 每一个航段都要找
	OPENTimeMap := make(map[string]string)
	for _, k := range subLst {
		OPENTimeMap[k[5]] = "20" + k[4] + "-" + MonthMap[k[2]] + "-" + k[1] + " " + k[3][:2] + ":" + k[3][2:] + ":00"
	}
	fmt.Println("")
	fmt.Println("取消机位OPEN时间解析...", OPENTimeMap)
	// DETR 中遍历出起飞时间， {"PEKCAD":"2020-07-07 12:00:00"}
	DETRTimeMap := make(map[string]string)
	for i := range params.Data.TripInfos {
		if params.Data.TripInfos[i].TicketNoStatus == "OPEN FOR USE" {
			TripInfosFromAirport := params.Data.TripInfos[i].FromAirport
			TripInfosToAirport := params.Data.TripInfos[i].ToAirport
			for p := range params.CostInfo.TripList {
				TripListFromAirport := params.CostInfo.TripList[p].FromAirport
				TripListToAirport := params.CostInfo.TripList[p].ToAirport
				if TripInfosFromAirport == TripListFromAirport && TripInfosToAirport == TripListToAirport {
					t, err := ConvertTimeToBeijing(TripInfosFromAirport, params.CostInfo.TripList[p].FlyDate + ":00")
					if err != nil {
						return false, fmt.Errorf("error in ISMissedFlight.ConvertTimeToBeijing, 时区转换失败， 请重试！, error: [%s]", err.Error())
					}
					DETRTimeMap[TripInfosFromAirport+TripInfosToAirport] = t
				}
			}
		}
	}
	fmt.Println("航班起飞DETR时间解析...", DETRTimeMap)
	FromToAirport := make([]string, 0)
	for k := range DETRTimeMap {
		FromToAirport = append(FromToAirport, k)
	}

	// 误机时间规定
	for _, portNode := range FromToAirport {
		if _, ok := OPENTimeMap[portNode]; !ok {
			return false, fmt.Errorf("error in ISMissedFlight, error: [%s]", "判断是否误机失败, 请确认该航段 "+OPENTimeMap[portNode]+"是否取消机位置于OPEN")
		}
		OPENTime, err := time.Parse("2006-01-02 15:04:05", OPENTimeMap[portNode])
		if err != nil {
			return false, fmt.Errorf("error in ISMissedFlight.Parse, error: [%s]", err.Error())
		}
		DETRTime, err := time.Parse("2006-01-02 15:04:05", DETRTimeMap[portNode])
		if err != nil {
			return false, fmt.Errorf("error in ISMissedFlight.Parse, error: [%s]", err.Error())
		}
		dur, _ := time.ParseDuration(fmt.Sprintf("-%dh", NoShowRuleTime))
		DETRTime = DETRTime.Add(dur)
		if OPENTime.After(DETRTime) {
			return true, nil // 误机
		}
	}
	return false, nil // 没有误机
}

// 获取不重复的票号列表
func GetTicketNoLst(passengers_ *structs.Passengers_) []string {
	ticketNoLst := make([]string, 0)
	s := mapset.NewSet()
	for _, pv := range passengers_.PassengerVoyages {
		s.Add(pv.TicketNo)
	}
	for _, v := range s.ToSlice() {
		ticketNoLst = append(ticketNoLst, v.(string))
	}
	return ticketNoLst
}

//  是否部份使用, 部份使用返回true
func IsPartlyUsed(detrStruct *structs.DETRStruct) bool {
	for _, v := range detrStruct.Data.TripInfos {
		if v.TicketNoStatus == "USED FLOWN" || // 客票已使用
			v.TicketNoStatus == "USED/FLOWN" ||
			v.TicketNoStatus == "LIFT/BOARDED" { //  已经换取登机牌
			return true
		}
	}
	return false

}
// ConvertTimeToBeijing 转换时区至北京时间, 传入机场三字码, 2020-07-13 12:10:00
func ConvertTimeToBeijing(AirportCode, timeStr string) (t string, err error) {
	tempTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return "", fmt.Errorf("error in ConvertTimeToBeijing,时间转换失败, 请确定订单详情是否正确, error: [%s]", err.Error())
	}
	airPortMap, err := utils.QueryAirport(AirportCode)
	if err != nil {
		return "", err
	}
	if airPortMap == nil {
		return "", fmt.Errorf("error in ConvertTimeToBeijing 查询机场信息失败, 查询结果为空, 请重试！")
	}
	if timeZoneStr, ok := airPortMap["time_zone"]; !ok {
		return "", fmt.Errorf("error in ConvertTimeToBeijing.airPortMap['time_zone'] 获取时区失败, 请重试！")
	}else {
		fmt.Println(AirportCode, "的时区为: ", timeZoneStr)
		zone, err := strconv.Atoi(timeZoneStr)
		if err != nil {
			return "",fmt.Errorf("error in ConvertTimeToBeijing.Atoi 转换时区字符串为int失败, 请重试！ error: [%s]", timeZoneStr)
		}
		d, err := time.ParseDuration(strconv.Itoa(zone/60 - 8)+"h")
		if err != nil {
			return "",fmt.Errorf("error in ConvertTimeToBeijing.ParseDuration 转换时区字符串为int失败, 请重试！ error: [%s]", timeZoneStr)
		}
		tempTime = tempTime.Add(-d)
		t = tempTime.Format("2006-01-02 15:04:05")
		return t, nil
	}
}