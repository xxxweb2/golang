package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "study/golang/shippy/consignment-service/proto/consignment"
)

const (
	PORT = ":50051"
)

type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error)
}

type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

type service struct {
	repo Repository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	resp := &pb.Response{Created: true, Consignment: consignment}
	fmt.Println("请求")

	return resp, nil
}

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("listen on: %s\n", PORT)

	server := grpc.NewServer()
	repo := Repository{}

	pb.RegisterShippingServiceServer(server, &service{repo: repo})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
