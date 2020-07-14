package structs

type ProcessInfo struct {
	IsNotPermitted bool // 是否允许退票
	IsPartlyUsed bool  // 是否部份使用
	NoShowFee float64 // 误机费
	NoShowRuleTime int //
	UsedFare float64 // 已使用票面价
	Nonfuel string // 不允许退的燃油税费
	Currency string // 币种
	item16RefundFee float64 // 16项中的退票费
	UsedHistoryPrice float64 // 自愿退票已使用票面价
	RefundFee float64 // 16项匹配出来的退票费

}
