package parseClient

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/wikensmith/parse16Items/message"
	"github.com/wikensmith/parse16Items/structs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"strings"
	"time"
)


func main2(ticketNo, officeNo string) (*parseCalcServer.Res, error){
	// 建立连接到gRPC服务
	conn, err := grpc.Dial("127.0.0.1:8088", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 函数结束时关闭连接
	defer conn.Close()
	// 创建Waiter服务的客户端
	t := parseCalcServer.NewWaiterClient(conn)
	detr, err := DoQuery(ticketNo, officeNo)
	fmt.Println("detr:", detr)
	if err != nil {
		fmt.Println("error in DoQuery:", err.Error())
		return nil, err
	}
	param := new(QueueStruct)

	if err := json.Unmarshal([]byte(Data), param); err != nil {
		fmt.Println("error in Unmarshal:", err.Error())
		return nil, err
	}

	param.BuyOrders[0].Passengers[0].RefundCenterDETR = detr
	param.BuyOrders[0].Passengers[0].PassengerVoyages[0].TicketNo = ticketNo
	param.BuyOrders[0].OfficeNo = officeNo

	res, err :=  json.Marshal(param)
	if err != nil {
		fmt.Println("error in Marshal:", err.Error())
		return nil, err
	}


	// 调用gRPC接口
	tr, err := t.ParseCalc(context.Background(), &parseCalcServer.Req{Data: string(res)})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("服务端响应: %v", tr.Status)
	log.Printf("服务端响应: %v", tr.Message)
	log.Printf("服务端响应: %v", tr.Data)
	return tr, nil
}

func DoLstTest(airlineName string, pathName string)  {
	lst, err := Read(pathName)
	//lst, err := Read("./解析客票.xlsx")
	if err != nil {
		fmt.Println(err)
	}
	//lst := []string{"826-3393458399"}

	officeNo1 := "CKG177"
	officeNo2 := "CKG262"
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet(airlineName)
	name := time.Now().Format("result2006-01-02T15-04-05")
	defer func() {
		xlsx.SetActiveSheet(index)
		if err := xlsx.SaveAs(name+".xlsx"); err != nil {
			fmt.Println(fmt.Errorf("error in Write.SaveAs, error: [%s]", err.Error()).Error())
		}
	}()
	xlsx = Write(xlsx, airlineName, "1", "票号", "是否成功", "票面价", "扣减", "税金")
	for i, v := range lst {
		fmt.Printf("总计: %d 条, 已完成: %d 条\n", len(lst), i)
		if v == ""{
			continue
		}
		res, err := main2(v, officeNo1)
		if err != nil {
			fmt.Println("error in main2", i, v)
			continue
		}
		if strings.Contains(res.Message, "没有权限获取票号数据") {
			res, err = main2(v, officeNo2)
			if err != nil {
				fmt.Println("error in main2", i, v)
				continue
			}
		}
		if res.Data == "" {
			if err := Write(xlsx,airlineName, strconv.Itoa(i+2), v, "解析失败",res.Message); err != nil {
				fmt.Println("error in write", err)
			}
			continue
		}
		tempStruct := new(structs.ComingData)

		_ = json.Unmarshal([]byte(res.Data), tempStruct)
		detrStruct := new(structs.DETRStruct)
		err = json.Unmarshal([]byte(tempStruct.BuyOrders[0].Passengers[0].RefundCenterDETR), detrStruct)
		fmt.Println(err)

		if err := Write(xlsx, airlineName, strconv.Itoa(i+2), v, "成功",
			strconv.FormatFloat(detrStruct.CostInfo.Price,'f', 2, 64),
			strconv.FormatFloat(detrStruct.UsedFare,'f', 2, 64),
			fmt.Sprintf("%#v\n", detrStruct.Taxs),
			); err != nil {
			fmt.Println(err)
		}
		if i %5 == 0 {
			xlsx.SetActiveSheet(index)
			if err := xlsx.SaveAs(name+".xlsx"); err != nil {
				fmt.Println(fmt.Errorf("error in Write.SaveAs, error: [%s]", err.Error()).Error())
			}
		}
	}
}

func DoSingleTest()  {

}