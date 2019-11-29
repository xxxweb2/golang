package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	pb "study/golang/09grpc/proto"
)

type UserInfoService struct{}

var u = UserInfoService{}

func (s *UserInfoService) GetUserInfo(ctx context.Context, req *pb.UserRequest) (resp *pb.UserResponse, err error) {
	name := req.Name
	if name == "zs" {
		resp = &pb.UserResponse{
			Id:    1,
			Name:  name,
			Age:   12,
			Hobby: []string{"Sing", "Run"},
		}
	}

	return
}

func main() {
	addr := "127.0.0.1:8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("监听异常：%s\n", err)
	}
	fmt.Printf("监听端口：%s\n", addr)
	s := grpc.NewServer()
	pb.RegisterUserInfoServiceServer(s, &u)
	s.Serve(listener)
}
