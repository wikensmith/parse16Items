package structs

type HostParam struct {
	AirLineCode string // 解析航司二字码
	MysqlURI string // mysql 数据库uri
	GRPCHost string
	GRPCPort string
	// 日志中心
	Project string
	Module string
	User string  // 推送至日志中心的用户名: 工号
}
