package main

import (
	"context"
	"google.golang.org/grpc"
	pb "nuclei/pb/server"
)

func main() {
	conn, _ := grpc.Dial("127.0.0.1:9999", grpc.WithInsecure())

	c := pb.NewRunPocServiceClient(conn)

	req := pb.PocReq{
		Urls: "http://218.77.53.28:8087",
	}
	reply, _ := c.Poc(context.Background(), &req)

	println(reply.JsonResults)
}
