package parse

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wikensmith/parse16Items/auxiliary"
	"github.com/wikensmith/parse16Items/const_"
	"github.com/wikensmith/parse16Items/db"
	parseCalcServer "github.com/wikensmith/parse16Items/message"
	"github.com/wikensmith/toLogCenter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)


func init() {
	db.InitDB(const_.Param.MysqlURI)
}

type server struct{}

func (s *server) ParseCalc(ctx context.Context, in *parseCalcServer.Req) (*parseCalcServer.Res, error) {
	fmt.Println("--------------- start " + time.Now().Format("2006-01-02 15:04:05") + "---------------")
	// 定义日志属性
	l := toLogCenter.Logger{
		Project:  const_.Param.Project,
		Module:  const_.Param.Module,
		User:    const_.Param.User,
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

	fmt.Println(fmt.Sprintf("GS %s:%s grpc...", const_.Param.GRPCHost, const_.Param.GRPCPort))
	lis, err := net.Listen("tcp", const_.Param.GRPCHost + ":" + const_.Param.GRPCPort)
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
