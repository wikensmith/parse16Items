package parse

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wikensmith/toLogCenter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"github.com/wikensmith/parse16Items/auxiliary"
	"github.com/wikensmith/parse16Items/config"
	"github.com/wikensmith/parse16Items/db"
	parseCalcServer "github.com/wikensmith/parse16Items/message"
	"github.com/wikensmith/parse16Items/structs"
	"time"
)
var Param = &structs.HostParam{
	AirLineCode: "GS",
	MysqlURI: "",
	GRPCHost: "0.0.0.0",
	GRPCPort: "8031",
	Project:  "RefundRulesCalc",
	Module:   "INT",
	User:     "7921",
}


func init() {
	config.InitConfig("./config.toml")
	db.InitDB(config.GetConfig().Mysql.URI)
}

type server struct{}

func (s *server) ParseCalc(ctx context.Context, in *parseCalcServer.Req) (*parseCalcServer.Res, error) {
	fmt.Println("--------------- start " + time.Now().Format("2006-01-02 15:04:05") + "---------------")
	// 定义日志属性
	l := toLogCenter.Logger{
		Project:  Param.Project,
		Module:  Param.Module,
		User:    Param.User,
	}
	executor, err := new(auxiliary.Executor).New(in.Data, l)
	if err != nil {
		return &parseCalcServer.Res{Data: "", Status: 1, Message: err.Error()}, nil
	}
	resData, err := executor.Do()
	defer func() {
		executor.Log.PrintReturn(resData)
		executor.Log.Send()
	}()
	if err != nil {
		executor.Log.Print(fmt.Sprintf("error in ParseCalc, error: [%s]", err.Error()))
		return &parseCalcServer.Res{Data: "", Status: 1, Message: fmt.Sprintf("error in ParseCalc, error: [%s]", err.Error())}, nil
	}
	executor.Log.Level("info")
	data, err := json.Marshal(resData)
	fmt.Println("返回数据...", "Data:", resData, "err:", err)
	fmt.Println("--------------- end " + time.Now().Format("2006-01-02 15:04:05") + "---------------")
	return &parseCalcServer.Res{Data: string(data), Status: 0, Message: "success"}, nil
}

func Start() {
	defer func() {
		db.DBClose()
	}()

	fmt.Println(fmt.Sprintf("GS %s:%s grpc...", Param.GRPCHost, Param.GRPCPort))
	lis, err := net.Listen("tcp", Param.GRPCHost + ":" + Param.GRPCPort)
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	s := grpc.NewServer()
	parseCalcServer.RegisterWaiterServer(s, &server{})
	reflection.Register(s)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
