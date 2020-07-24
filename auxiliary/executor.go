package auxiliary

import (
	"encoding/json"
	"fmt"
	"github.com/wikensmith/toLogCenter"
	"github.com/wikensmith/parse16Items/structs"
	"github.com/wikensmith/parse16Items/utils"
	"strings"
)

// 把flightNo从constInfo中更新至Data里面
func UpdateInfoForDETR(tempStruct *structs.DETRStruct)  {
	for i, m := range tempStruct.CostInfo.TripList{
		for j, n := range tempStruct.Data.TripInfos{
			if m.FromAirport == n.FromAirport && m.ToAirport == n.ToAirport {
				tempStruct.Data.TripInfos[j].FlightNo = tempStruct.CostInfo.TripList[i].FlightNo
				if n.DepartureTime == "" {
					lst := strings.Split(tempStruct.CostInfo.TripList[i].FlyDate, " ")
					if len(lst) >=  2 {
						tempStruct.Data.TripInfos[j].DepartureTime = strings.Trim(lst[1], " ")
						tempStruct.Data.TripInfos[j].FlightDate = strings.Trim(lst[0], " ")
					}
				}
			}
		}
	}
}


type Executor struct {
	ComingData *structs.ComingData
	Log *toLogCenter.Logger
}

func (e *Executor) New(inputData string, logger toLogCenter.Logger) (*Executor, error){
	e.Log = logger.New()
	e.ComingData = new(structs.ComingData)
	err := json.Unmarshal([]byte(inputData), e.ComingData)
	if err != nil{
		return nil, fmt.Errorf("error in ParseCalc.Unmarshal, error: [%s]", err.Error())
	}
	e.Log.AddField(0, e.ComingData.BuyOrders[0].Passengers[0].PassengerVoyages[0].TicketNo)
	e.Log.AddField(1, e.ComingData.BuyOrders[0].Pnr["PnrCode"])
	e.Log.AddField(2, e.ComingData.Voyages[0].AirLine)
	e.Log.AddField(3, e.ComingData.BuyOrders[0].OfficeNo)
	e.Log.AddField(4, "wiken")
	e.Log.PrintInput(e.ComingData)
	return e, nil
}

func (e *Executor)assignDETR()  {
	for k, b := range e.ComingData.BuyOrders{

		for pk, p := range b.Passengers {
			if data, err := json.Marshal(p.DETR);err != nil {
				fmt.Println(err)
			}else {
				e.ComingData.BuyOrders[k].Passengers[pk].RefundCenterDETR = string(data)
			}
		}
	}
}
// 删除行程中的缺口程, (退票的时候需要重新提取相关信息)
func (e *Executor)deleteEmptyTrips(bk, pk int , tempDETR *structs.DETRStruct)  {
	for k, v :=  range tempDETR.Data.TripInfos {
		if v.FlightNo == ""|| v.FlightNo == "unkown" {
			tempDETR.Data.TripInfos = append(tempDETR.Data.TripInfos[:k], tempDETR.Data.TripInfos[k+1:]...)
		}
	}
	for k, v := range tempDETR.CostInfo.TripList{
		if v.FlightNo  == "" || v.FlightNo == "unkown"{
			tempDETR.CostInfo.TripList = append(tempDETR.CostInfo.TripList[:k], tempDETR.CostInfo.TripList[k+1:]...)
		}
	}
	e.ComingData.BuyOrders[bk].Passengers[pk].DETRNotAll = tempDETR
}


func (e *Executor) Do() (interface{}, error) {
	for bk, b := range e.ComingData.BuyOrders {
		for pk, p := range b.Passengers {
			// 解析detr成结构体
			tempStrct := new(structs.DETRStruct)
			if err := json.Unmarshal([]byte(p.RefundCenterDETR), tempStrct); err != nil {
				fmt.Println("error for parse detr:",err)
			}else {
				UpdateInfoForDETR(tempStrct)
				e.ComingData.BuyOrders[bk].Passengers[pk].DETR = tempStrct
				e.deleteEmptyTrips(bk, pk, tempStrct)
			}
			// 解锁 SUSPENDED
			err := utils.DoTicketNoSuspended(tempStrct, GetOfficeNo(b))
			if err != nil {
				return "", fmt.Errorf("error in DoRefCalc.DoTicketNoSuspended, error: [%s]", err.Error())
			}

			// 获取16项原始数据
			 q := Query16Item{
				Passenger: p,
				detrInfo:  nil,
				OfficeNo:  GetOfficeNo(b),
				Log: e.Log}
			origin16Data, err  := q.Do()
			e.Log.Print(origin16Data)
			if err != nil {
				return nil, fmt.Errorf("error in Executor.Do.Query16Items, 查询黑屏16项数据失败, error:[%s]", err.Error())
			}

			// 匹配模板
			m :=  MatchModel{
				Passengers:        p,
				Origin16ItemsData: origin16Data,
				OfficeNo:GetOfficeNo(b),
				Log: e.Log}
			f, err := m.Do()
			if err != nil {
				return nil, fmt.Errorf("error in Executor.Do.Query16Items,匹配模板失败, error:[%s]", err.Error())
			}

			// 执行计算费用
			if err := f.DoCalc(); err != nil {
				return nil, fmt.Errorf("error in Executor.Do.DoCalc, 计算退票费失败,error: [%s]", err)
			}

		}
	}
	e.assignDETR()
	return e.ComingData, nil
}
