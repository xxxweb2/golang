package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	pb "study/golang/36XuGrpc/proto"
)

type UserInfoService struct {
}

var u = UserInfoService{}

func (c *UserInfoService) GetUserInfo(ctx context.Context, req *pb.UserRequest) (resp *pb.UserResponse, err error) {
	name := req.Name
	if name == "zs" {
		resp = &pb.UserResponse{
			Id:    1,
			Name:  name,
			Age:   12,
			Hobby: []string{"Sing", "run"},
		}
	}

	return
}

func main() {
	addr := "127.0.0.1:8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("监听异常：%s\n", err)
		return
	}
	fmt.Printf("监听端口：%s\n", addr)
	s := grpc.NewServer()
	pb.RegisterUserInfoServiceServer(s, &u)

	s.Serve(listener)
}
