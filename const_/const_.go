package const_

var Param = &HostParam{
	PrivateKey: "cqyshkgsetermsdk",
	BaseURL: "http://eterm.70168.com:8081",
	AirLineCode: "GS",
	MysqlURI: "",
	GRPCHost: "0.0.0.0",
	GRPCPort: "8031",
	Project:  "RefundRulesCalc",
	Module:   "INT",
	User:     "7921",
}


type HostParam struct {
	BaseURL string // 大师兄接口基础路由
	PrivateKey string // 大师兄接口私钥
	AirLineCode string // 解析航司二字码
	MysqlURI string // mysql 数据库uri
	GRPCHost string
	GRPCPort string
	// 日志中心
	Project string
	Module string
	User string  // 推送至日志中心的用户名: 工号
}

