package parseClient

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type Sheet struct {
	TicketNo string
	Status int
	Msg string
	Deduction float64
	FarePrice float64
	usedFarePrice float64
	NoShowFee float64
	TAxes float64
}

func Read(path, airlineName string) ([]string, error ){
	xlsx, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error in Read.OpenFile, error:[%s]", err.Error())
	}
	//cell := xlsx.GetCellValue(airlineName, "A2")
	//fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows := xlsx.GetRows(airlineName)
	resultLst := make([]string, 0)
	for k, row := range rows {
		if k == 0 {
			continue
		}
		if len(row) == 0 {
			continue
		}
		resultLst = append(resultLst, row[0])
	}
	return resultLst, nil
}

func Write(xlsx *excelize.File, sheetName string, row string, val ...string) ( *excelize.File) {
	lst := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N"}
	for i, v := range val{
		xlsx.SetCellValue(sheetName, lst[i] + row, v)
	}
	return xlsx
}


