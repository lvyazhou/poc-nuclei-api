package server

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func Init() {
	// 开启服务端口监听
	s, _ := net.Listen("tcp", ":9999")
	grpcServer := grpc.NewServer()

	// 注册rpc服务
	pocService := PocService{}
	RegisterRunPocScanServiceServer(grpcServer, &pocService)

	fmt.Printf("start server poc-osint-server %s-%s on %d success...\n", "127.0.0.1", "v1", 9999)
	// 启动rpc server
	err := grpcServer.Serve(s)
	if err != nil {
		fmt.Println("start server poc-osint-server failed ", err)
	}
}
