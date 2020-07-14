package auxiliary

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/wikensmith/toLogCenter"
	"github.com/wikensmith/parse16Items/structs"
	"github.com/wikensmith/parse16Items/utils"
	"strings"
	"time"
)

type Query16Item struct{
	Passenger *structs.Passengers_
	detrInfo map[string]interface{}
	OfficeNo string
	Log *toLogCenter.Logger
}

// 解析detr数据， 为获取16项组织数据
func (q *Query16Item)parseDETR() error {
	params := q.Passenger.DETR
	AirTwoCode, err := utils.QueryAirLine(q.Passenger.PassengerVoyages[0].TicketNo[:3])
	if err != nil {
	}
	// 判断是否是同一个运价级别 还是 多个
	requiredClasses := mapset.NewSet()
	for _, v := range params.Data.TripInfos{
		requiredClasses.Add(v.FareBasis )
	}
	var FareBasis string
	// 如果是同一个运价级别 传 该运价级别
	s := requiredClasses.ToSlice()
	if len(s) == 1{
		FareBasis = s[0].(string)
		FareBasis = strings.Replace(FareBasis, "/CH", "", -1)
	}else {
		FareBasis = ""
	}
	if params.Error != 0 {
		return fmt.Errorf("error in parseDETR, detr查询状态异常， error msg: [%s], error code: [%d]", params.Message, params.Error)
	}

	Segments_ := make([]interface{}, 0)
	for _, v := range params.CostInfo.TripList{
		segMap := make(map[string]interface{})
		segMap["Airline"] =  v.Airline
		segMap["FlightNum"] = v.FlightNo
		segMap["Cabin"] = v.Cabin
		segMap["FromAirport"] = v.FromAirport
		segMap["ToAirport"] = v.ToAirport
		segMap["DepTime"] = v.FlyDate  + ":00"
		depTimeFormat, _:= time.Parse("2006-01-02 15:04:05", v.FlyDate  + ":00")
		arrTime := depTimeFormat.Add(3*time.Hour)
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
		"BillingAirline": AirTwoCode["iata_code"],
		"Segments":       Segments_,
	}
	q.detrInfo = DETRInfo
	return  nil
}
func (q *Query16Item)Do() (string, error){
	if err := q.parseDETR(); err != nil {
		return "", fmt.Errorf("error in MatchModel.Do.parseDETR, error: [%s]", err.Error())
	}
	return utils.GetRefRules16(q.detrInfo, q.OfficeNo)
}