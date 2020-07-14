package parseClient


type QueueInfo interface {
	GetIntlFlag() bool
	GetYATPRefundNo() interface{}
	GetIsVoluntary() bool
}

// 队列信息结构体， for channelName
type QueueStruct struct {
	YATPRefundNo         interface{} `json:"Id" binding:"required"` // YATP退票单Id, 实际的平台订单号
	IsVoluntary          bool        `json:"IsVoluntary" binding:"required"`
	RefundReason         string      `json:"Reason" binding:"required"`
	PassengersName       string      //退票所有乘机人姓名，多个逗号隔开
	PassengerInfoList    []queuePassengerInfo
	Voyages              []Voyages_
	BuyOrders            []BuyOrder
	CreateDepartment     string  //退票申请部门
	OrganizationUnitName string  //退票所属部门
	RefundOrderNo        string  //YATP退票单号
	Reason               string  //退票理由
	RefundAmount         float64 //销售实退金额
	BuyAmount            float64 // 采购订单退总金额(应收)
	CreationTime         string  //退票申请时间 2020-02-25T18:22:52.6402835+08:00
	Status               string  //退票状态 （申请中：Applying，已取消：Cancel，待审核 Recheck，已完成：Finish）
	ApiInfo              apiInfo // 接口帐号信息
	Remark               string  // 传入附件信息 地址
	TicketNos            string  `json:"TicketNo"`
	IntlFlag             bool
	CreatorName          string // 退票申请人
	AuditorName          string
	Attachment           []Attachment_
	ProfitCenters        ProfitCenters_
	PNR                  PNR `json:"Pnr"` // pnr 信息

}

// 销售单pnr信息
type PNR struct {
	GdsName    string
	OfficeNo   string
	PnrCode    string
	BigPnrCode string
	Source     string
	SourceNo   string
	Id         interface{}
}

type ProfitCenters_ struct {
	Params Params_
}
type Params_ struct {
	UserId     string
	AsmsAcount string
	AsmsPwd    string
	Xm         string
	Sj         string
}

func (q *QueueStruct) GetIntlFlag() bool {
	return q.IntlFlag
}
func (q *QueueStruct) GetYATPRefundNo() interface{} {
	return q.YATPRefundNo
}
func (q *QueueStruct) GetIsVoluntary() bool {
	return q.IsVoluntary
}

type Attachment_ struct {
	File File_
}
type File_ struct {
	FileName string
	FilePath string
	Hash     string
}

type apiInfo struct {
	ApiAccount  string
	ApiPassword string
	Key         string
	Extra       string
}
type Voyages_ struct {
	Id         interface{}
	DepAirport string // 出发地三字码
	ArrAirport string
	AirLine    string // 航司二字码
	Cabin      string
	DepTime    string
	ArrTime    string
}
type queuePassengerInfo struct {
	Name     string
	Segment  string
	CertType string // 证件类型
	CertNo   string
	Voyages  []VoyagesId
}
type VoyagesId struct {
	TicketNo string
	VoyageId interface{}
}
type BuyOrder struct {
	Submitter        string //提交者，格式为: "20200610150203:YATP^首次提交@@20200610160203:web^再次提交"
	IsVoluntary      bool
	RefundBuyOrderNo string
	OuterBuyOderId   string
	BuyChannel       BuyChannel_
	Pnr              map[string]string
	GdsName          string
	OfficeNo         string
	Currency         string
	ExchangeRate     float64
	BuyAmount        float64
	BuyForeignAmount float64
	ReceivedAmount   float64
	RefundType       string
	SubmitStatus     string
	RefundStatus     string
	Passengers       []Passengers_
	Paylogs          string
	CreationTime     string
	Id               interface{} //采购订单Id
	ExtensionData    string
}
type RefundCenterMsg struct {
	IsExchange         bool     // 是否换开
	IsUsePartly        bool     //  是否部份使用
	ContinuousTicketNo []string // 连续票号列表，如果
}

type PassengerVoyages_ struct {
	TicketNo string // 票号
	Id       interface{}
	VoyageId interface{}
}

type RefundCenterFeeDetail struct {
	TicketNo string                 // 票号
	Taxes    map[string]interface{} // 如果是部份退票的税金明细

}

type Passengers_ struct {
	PassengerId          interface{}
	ChangePassengerId    string
	PassengerVoyages     []PassengerVoyages_
	RefundCenterDETR     string // 退票中心执行DETR查询结果json字符串  {"blackScreen": "data", "webInput":"data2"}
	Name                 string
	PassengerType        string
	CertType             string
	CertNo               string
	SaleTicketPrice      float64
	SaleTicketTax        float64
	SaleInsurePrice      float64
	SaleChangePrice      float64
	SaleAgencyFee        float64
	SaleRebate           float64
	SaleCashBack         float64
	SaleAdjustPrice      float64
	SaleServiceFee       float64
	SaleAdjustServiceFee float64
	SaleDelayFee         float64
	SaleOtherFee         float64
	SaleSum              float64
	BuyTicketPrice       float64
	BuyTicketTax         float64
	BuyAgencyFee         float64
	SaleRefundFee        float64
	BuyRebate            float64
	BuyCashBack          float64
	BuyInsurePrice       float64
	BuyChangePrice       float64
	BuyAdjustPrice       float64
	BuyServiceFee        float64
	BuyAdjustServiceFee  float64
	BuyDelayFee          float64
	BuyRefundFee         float64
	BuyOtherFee          float64
	BuySum               float64
	BuyTaxDetail         string
	Id                   interface{} // 乘机人ID
}
type BuyChannel_ struct {
	Account         string
	AccountAbbr     string
	ChannelInt      string
	ChannelType     string // 渠道类型
	ChannelTypeName string // 渠道名称
	SaleBuyCapacity string
	OrderCapacity   string
	Name            string // 渠道名称
	NameAbbr        string
	Remark          string
	Organization    string
	OrganizationId  string
	PaymentId       int64
	PaymentName     string
	Id              int64
}

//退票返回参数结构体
type RefundResponse struct {
	Status             int     // 0 成功 1 请求接口异常 2 接口返回异常  3 参数错误 4 其他异常
	Msg                string  // 退票消息,  申请退票成功（退票状态，如可知退票状态）; status非0, 返回具体异常信息
	YatpRefundId       string  // YATP 退票单号
	PlatformRefundId   string  // 平台退票单号
	IntlFlag           bool    //是否国际 true:国际，false：国内
	ShouldRefundPrice  float64 // 应退金额, 采购方应退金额
	ActualRefundPrice  float64 // 实退金额 没有就不填
	RefundBuyPassenger []*RefundBuyPassenger_
	BuyOrders          []map[string]interface{}
	ElapsedTime        int64  // 耗时ms
	ExecutorName       string // 执退票接口名称
}

func (r *RefundResponse) GetIntlFlag() bool {
	return r.IntlFlag
}
func (r *RefundResponse) GetYATPRefundNo() interface{} {
	return r.YatpRefundId
}

type RefundBuyPassenger_ struct {
	BuyOrderId   interface{}      // YATP 采购单ID
	PassengerId  interface{}      // YATP乘客ID
	TickNoAction []*TickNoAction_ // 退票票号
}
type TickNoAction_ struct {
	PlatformRefundId string // 平台退票订单号
	TicketNo         string
	Message          string
	Remark           string
	Status           int
}

// parse for return
type BuyOrderForReturn struct {
	BuyOrders []map[string]interface{}
}

// LogMsg 日志中心存储信息
type LogMsg struct {
	Result         *RefundResponse
	InputParameter map[string]interface{}
	QueueName      string
}

// 黑屏解析结果
type DETRParseResult struct {
	IsExchange           string
	IsUsePartly          string
	IsContinuousTicketNo string
	Exception            string
}

type RefundCenterDETR struct {
	BlackScreen string // 黑屏信息
	WebInput    string // 前端返回信息
}





type TripInfos struct {
	TripNo         string // 航段
	TicketNoStatus string // 航段票号状态
	TicketNo       string // 票号
}

type ConstInfo struct {
	CONJTKT string // 连续票号 如 "1/2"
	EXCH    string // 是否换开, 如果换开了， 这里是原来的票号
}
