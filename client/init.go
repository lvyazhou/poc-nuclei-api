package main

import (
	"context"
	"google.golang.org/grpc"
	"nuclei/server"
)

func main() {
	conn, _ := grpc.Dial("127.0.0.1:9999", grpc.WithInsecure())

	c := server.NewRunPocScanServiceClient(conn)

	req := server.PocScanReq{
		Urls: "http://218.77.53.28:8087",
	}
	reply, _ := c.Poc(context.Background(), &req)

	println(reply.JsonResults)
}
