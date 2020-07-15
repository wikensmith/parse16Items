## <center>退票中心16项解析</center>


### 使用方法
#### 1. 指定自定义参数
```go
func main(){
// 指定mysql数据库uri
  const_.Param.MysqlURI ="root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
  // 指定grpc访问端口
  const_.Param.GRPCPort = "8088"
  // 指定发往日志中心的用户工号
  const_.Param.User = "7921"
  // 指定模板解析函数
  auxiliary.ParseModels = ParseModel  // 注入自定义解析函数
  Start() // 启动程序
}


```
#### 2. 模板解析函数示例ParseModel
```go
// ParseModel模板解析函数， etermStr为黑屏16项文本
func ParseModel(etermStr string) (modules.Calc, error){
	l := []string{"a", "b"}
	for _, v := range l {
	    // 符合特定条件的时候使用特定的模板, 该模板为自己定义
		if v == "b"{
			return new(ModelNo1), nil
		}
		if v == "a" {
			return new(ModelNo2), nil
		}
	}
	return nil, fmt.Errorf("error in ParseModel,没有匹配到模板, error: [%s]", "")
}
```
#### 3. 模板函数示例
```go 
/*
该类实现16项文件解析及退票费，税金的计算。 把计算结果赋值给指针属性Passegner中的DETR实现返回值
*/
//ModelNo1 << -------的模板
type ModelNo1 struct {
	Passengers      *structs.Passengers_
	EtermMap        map[string]string
	ticketNoLst     []string
	OfficeNo        string
	ProcessInfo     *structs.ProcessInfo
	Log             *toLogCenter.Logger
}
// 接口函数
func (m *ModelNo1)Name() string {
	return "ModelNo1"
}
// 接口函数
func (m *ModelNo1)New(Passengers *structs.Passengers_,EtermMap map[string]string,OfficeNo string,Log *toLogCenter.Logger) modules.Calc {
	m.Passengers = Passengers
	m.EtermMap = EtermMap
	m.OfficeNo = OfficeNo
	m.Log = Log
	return m
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
		n.ProcessInfo.NoShowFee, n.ProcessInfo.Currency, err = GetHighestFee(resultLst)
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
		if Deduction > detrStruct.CostInfo.Price {
			Deduction = detrStruct.CostInfo.Price
		}
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
// DoCalc  接口函数 （入口函数）
func (n *ModelNo1) DoCalc() (err error) {
	n.init()

	// 1.匹配退票费
	refundFeeLst, err := n.getRefundFee(n.EtermMap)
	if err != nil {
		return err
	}

	// 2.获取最高退票费
	if len(refundFeeLst) != 0{n.ProcessInfo.RefundFee, n.ProcessInfo.Currency, err = GetHighestFee(refundFeeLst)}

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

```
#### 4. 测试
```go
func DoTest()  {
    // GS 生成测试xlsx文件的sheet名称。 测试文件地址, grpc 端口
	parseClient.DoLstTest("GS", "./test.xlsx", "8088")
}
```

