package auxiliary

import (
	"encoding/json"
	"fmt"
	"github.com/wikensmith/toLogCenter"
	"github.com/wikensmith/parse16Items/modules"
	"github.com/wikensmith/parse16Items/structs"
	"strings"
)

var ParseModels = func(string) (modules.Calc, error){return nil, nil}

type MatchModel struct {
	Passengers        *structs.Passengers_
	Origin16ItemsData string
	OfficeNo string
	etermMap          map[string]string
	Log *toLogCenter.Logger
}

func (m *MatchModel)Name() string {
	return "MatchModel"

}

func (m *MatchModel) parseOriginData() ([]string, error) {
	params2 := new(structs.ResGetRefRules16)
	if err := json.Unmarshal([]byte(m.Origin16ItemsData), params2); err != nil {
		return nil, fmt.Errorf("error in DoParamsGetRefRules16.Unmarshal, error: [%s]", err.Error())
	}
	if params2.Error != 0 {
		return nil, fmt.Errorf("error in DoParamsGetRefRules16, error: [%s]", params2.Message)
	}
	if len(params2.EtermStr) == 0 {
		return nil, fmt.Errorf("error in DoParamsGetRefRules16, error: [%s]", params2.Message)
	}
	//m.etermLst = params2.EtermStr
	return params2.EtermStr, nil
}

// getFXGAndResultMap 获取字典 {"xs fxg01//16":"result1","xs fxg02//":"result2"}
func (m *MatchModel) getFXGAndResultMap(etermLst []string) {

	m.etermMap = make(map[string]string, 0)
	tempK := ""
	for k, d := range etermLst {
		if strings.HasPrefix(d, ">XS FSG") ||
			strings.HasPrefix(d, ">xs fsg") ||
			strings.HasPrefix(d, ">xs fxg") ||
			strings.HasPrefix(d, ">XS FXG") {
			// 731-3391534281
			if k+2 <= len(etermLst)-1 {
				if etermLst[k+2] == ">XSFSPN" {
					tempK = d
					m.etermMap[d] = etermLst[k+1]
				}
			}else {
				tempK = d
				m.etermMap[d] = etermLst[k+1]}
		}
		if d == ">XSFSPN" {
			m.etermMap[tempK] += etermLst[k+1]
		}
	}
}

var parts = map[string]string{
	"CHANGE":  "FOR CHANGE",
	//"NOSHOW":  "FOR NOSHOW",  有时没有这一项
	"REFUNDS": "FOR REFUND",
}

func (m *MatchModel) matchModule() (modules.Calc, error) {
	for _, v := range m.etermMap {
		obj, err := ParseModels(v)
		if err != nil {
			return nil, fmt.Errorf("error in matchModule.ParseModels, error: [%s]", err.Error())
		}
		model := obj.New(m.Passengers, m.etermMap, m.OfficeNo, m.Log)
		return model, nil


		//// 1.匹配ModelNo1模板
		//reg := regexp.MustCompile(`CANCELLATION\s*<<.*?CHARGE ([A-Z]{3})(\d+)(\.*\d*).*?FOR REFUND/CANCEL`)
		//subLst01 := reg.FindAllStringSubmatch(v, -1)
		//if len(subLst01) != 0 {
		//	return &modules.ModelNo1{
		//		Passengers: m.Passengers,
		//		EtermMap:   m.etermMap,
		//		OfficeNo:   m.OfficeNo,
		//		Log:        m.Log}, nil
		//}
	}
	return nil, fmt.Errorf("error in matchModule, error: [%s]", "没有匹配到模板")
}

func (m *MatchModel) Do() (modules.Calc, error) {
	// ["指令1", "result1","指令2","result2"]
	etermLst, err := m.parseOriginData()
	if err != nil {
		return nil, fmt.Errorf("error in MatchModel.parseOriginData,解析16项原始数据失败, error: [%s]", err.Error())
	}
	// etermLst {"fxg01//16": "result1", "fxg02//16": "result2"}
	m.getFXGAndResultMap(etermLst)
	return m.matchModule()
}
