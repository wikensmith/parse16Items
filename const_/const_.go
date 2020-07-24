package const_

var Param = &HostParam{
	PrivateKey: "cqyshkgsetermsdk",
	BaseURL: "http://eterm.70168.com:8081",
	AirLineCode: "GS",
	MysqlURI: "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
	GRPCHost: "0.0.0.0",
	GRPCPort: "8031",
	Project:  "RefundRulesCalc",
	Module:   "INT",
	User:     "7921",
	UWURL: "http://192.168.0.212:8058/admin_api/uwing",
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
	// UW 网站查询过期pnr
	UWURL string  // 调用查询接口地址

}

