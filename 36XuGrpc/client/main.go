package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "study/golang/36XuGrpc/proto"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("链接异常: %s\n", err)
		return
	}
	defer conn.Close()

	client := pb.NewUserInfoServiceClient(conn)
	req := new(pb.UserRequest)
	req.Name = "zs"
	response, err := client.GetUserInfo(context.Background(), req)

	if err != nil {
		fmt.Printf("响应异常 %s\n", err)
		return
	}
	fmt.Printf("响应结果： %v\n", response)
}
