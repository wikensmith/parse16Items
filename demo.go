package parse

import (
	"fmt"
	"github.com/wikensmith/parse16Items/auxiliary"
	"github.com/wikensmith/parse16Items/const_"
	"github.com/wikensmith/parse16Items/modules"
	"github.com/wikensmith/parse16Items/structs"
	"github.com/wikensmith/toLogCenter"
)

func ParseModel(etermStr string) (modules.Calc, error){
	l := []string{"a", "b"}
	for _, v := range l {
		if v == "b"{
			return &ModelNo1{}, nil
		}
		if v == "a" {
			return &ModelNo2{}, nil
		}
	}
	return nil, fmt.Errorf("error in ParseModel,没有匹配到模板, error: [%s]", "")
}

type ModelNo1 struct {
	Passengers      *structs.Passengers_
	EtermMap        map[string]string
	ticketNoLst     []string
	OfficeNo        string
	ProcessInfo     *structs.ProcessInfo
	Log             *toLogCenter.Logger
}
func (m *ModelNo1)New(Passengers *structs.Passengers_,EtermMap map[string]string,OfficeNo string,Log *toLogCenter.Logger) modules.Calc {
	m.Passengers = Passengers
	m.EtermMap = EtermMap
	m.OfficeNo = OfficeNo
	m.Log = Log
	return m
}
func (m *ModelNo1)DoCalc() error {
	fmt.Println("func111111111111111111")
	return nil
}
func (m *ModelNo1)Name() string {
	return "ModelNo1"
}

type ModelNo2 struct {
	Passengers      *structs.Passengers_
	EtermMap        map[string]string
	ticketNoLst     []string
	OfficeNo        string
	ProcessInfo     *structs.ProcessInfo
	Log             *toLogCenter.Logger
}
func (m *ModelNo2)New(Passengers *structs.Passengers_,EtermMap map[string]string,OfficeNo string,Log *toLogCenter.Logger) modules.Calc {
	m.Passengers = Passengers
	m.EtermMap = EtermMap
	m.OfficeNo = OfficeNo
	m.Log = Log
	return m
}
func (m *ModelNo2)DoCalc() error {
	fmt.Println("\nfunc22222222222222222222")
	return nil
}
func (m *ModelNo2)Name() string {
	return "ModelNo2"
}
// 开始
func Demo() {
	const_.Param.MysqlURI = "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	const_.Param.GRPCHost = "0.0.0.0"
	const_.Param.GRPCPort = "8088"
	const_.Param.User = "7921"
	auxiliary.ParseModels = ParseModel
	Start()
}
// 测试
func ForTest()  {

}
