package modules

import (
	"fmt"
	"github.com/wikensmith/toLogCenter"
	"github.com/wikensmith/parse16Items/structs"
	"regexp"
)

//ModelNo1 << -------的模板
type ModelNo1 struct {
	Passengers      *structs.Passengers_
	EtermMap        map[string]string
	ticketNoLst     []string
	OfficeNo        string
	ProcessInfo     *structs.ProcessInfo
	Log             *toLogCenter.Logger
}

func (m *ModelNo1)Name() string {
	return "ModelNo1"
}

func (n *ModelNo1) init() {
	n.ticketNoLst = make([]string, 0)
	n.ProcessInfo = new(structs.ProcessInfo)
}

func (n *ModelNo1) getRefundFee(tempMap map[string]string) ([][]string, error) {
	//UsedFlag := IsPartlyUsed(n.Passengers.DETR)
	resultLst := make([][]string, 0)
	for _, v := range tempMap {
		reg := regexp.MustCompile(`CANCELLATION\s<<.*?CHARGE ([A-Z]{3})(\d+)(\.*\d*).*?FOR REFUND/CANCEL`)
		subLst01 := reg.FindAllStringSubmatch(v, -1)
		if len(subLst01) != 0 {
			resultLst = append(resultLst, []string{subLst01[0][1], subLst01[0][2] + subLst01[0][3]})
		}
	}
	// "TICKET IS NON-REFUNDABLE IN CASE OF REFUND"
	return resultLst, nil
}

// 在16项中匹配需要的信息map
func (n *ModelNo1) get16ItemInfo(tempMap map[string]string) (err error) {

	// 1.匹配Noshow fee
	//UsedFlag := IsPartlyUsed(n.Passengers.DETR)
	resultLst := make([][]string, 0)
	for _, v := range tempMap {
		reg := regexp.MustCompile(`CANCELLATION\s<<.*?CHARGE.*?\s<<.*?CHARGE ([A-Z]{3})(\d+)(\.*\d*) FOR NO-SHOW`)
		subLst := reg.FindAllStringSubmatch(v, -1)
		if len(subLst) != 0 {
			resultLst = append(resultLst, []string{subLst[0][1], subLst[0][2] + subLst[0][3]})
		}
	}
	if len(resultLst) != 0 {
		n.ProcessInfo.NoShowFee, n.ProcessInfo.Currency, err = getHighestFee(resultLst)
		// currency都是CNY
	}
	// 2.匹配是否允许退票
	n.ProcessInfo.IsNotPermitted = false

	// 3.确定是否退燃油税规则
	// 没有说退燃油锐的问题， 全退

	// 4. 匹配NoShowRuleTime
	// 无
	return nil
}

func (n *ModelNo1) calc() error {

	detrStruct := n.Passengers.DETR
	// 如果 UsedFare 存在  UsedHistoryPrice 不存在  UsedHistoryPrice 取 UsedFare的值
	if detrStruct.UsedFare != 0 && detrStruct.UsedHistoryPrice == 0 {
		n.ProcessInfo.UsedHistoryPrice = detrStruct.UsedFare
	}else {
		n.ProcessInfo.UsedHistoryPrice = detrStruct.UsedHistoryPrice
	}

	// 判断是否允许退票
	var Deduction float64
	// 1.允许退票
	if  !n.ProcessInfo.IsNotPermitted{
		// 票面价是否够扣
		Deduction = n.ProcessInfo.UsedHistoryPrice + n.ProcessInfo.RefundFee + n.ProcessInfo.NoShowFee
	} else {
		// 2.不允许退票, 直接取票面价
		//Deduction = detrStruct.CostInfo.Price
	}
	// 3.燃油税费 无

	detrStruct.UsedFare = Deduction

	n.Log.Print(map[string]interface{}{
		"退票费": n.ProcessInfo.RefundFee,
		"误机费": n.ProcessInfo.NoShowFee,
		"不允许退票": n.ProcessInfo.IsNotPermitted,
		"已使用票面": n.ProcessInfo.UsedHistoryPrice,
		"票面": detrStruct.CostInfo.Price,
		"扣除额": Deduction,
		"应退税":detrStruct.Taxs,
	})
	return nil
}
// DoCalc  GS  JD HU 解析是一样的， 可以使用相同的模板
func (n *ModelNo1) DoCalc() (err error) {
	n.init()

	// 1.匹配退票费
	refundFeeLst, err := n.getRefundFee(n.EtermMap)
	if err != nil {
		return err
	}

	// 2.获取最高退票费
	if len(refundFeeLst) != 0{n.ProcessInfo.RefundFee, n.ProcessInfo.Currency, err = getHighestFee(refundFeeLst)}

	// 3.匹配过程中需要的参数
	if err := n.get16ItemInfo(n.EtermMap); err != nil {
		return  fmt.Errorf("error in ModelNoArrows.DoCalc.get16ItemMap, error:[%s]", err.Error())
	}

	// 4.判断是否误机
	IsMissing, err := ISMissedFlight(n.Passengers.PassengerVoyages[0].TicketNo, n.OfficeNo, n.Passengers.DETR, n.ProcessInfo.NoShowRuleTime)
	if err != nil {
		return fmt.Errorf("error in ModelNo1.DoCalc.ISMissedFlight, error: [%s]", err.Error())
	}
	if !IsMissing {
		n.ProcessInfo.NoShowFee = 0.0
	}
	fmt.Println("误机相关时间规定:", n.ProcessInfo.NoShowRuleTime)
	fmt.Println("是否误机:", IsMissing)
	// 5.计算退票费
	if err := n.calc(); err != nil {
		return  fmt.Errorf("error in ModelNo1.DoCalc.calc, 计算退票费失败, error:[%s]", err.Error())
	}
	return  nil
}
