package main

import (
	"context"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	pb "study/golang/shippy/consignment-service/proto/consignment"
)

const (
	ADDRESS           = "localhost:50051"
	DEFAULT_INFO_FILE = "consignment.json"
)

func parseFile(fileName string) (*pb.Consignment, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var consignment *pb.Consignment
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, errors.New("consignment.json file err")
	}

	return consignment, err
}

func main() {
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect error: %b", err)
	}
	defer conn.Close()

	client := pb.NewShippingServiceClient(conn)
	infoFile := DEFAULT_INFO_FILE
	if len(os.Args) > 1 {
		infoFile = os.Args[1]
	}

	consignment, err := parseFile(infoFile)
	if err != nil {
		log.Fatalf("parse info file error: %v", err)
	}

	resp, err := client.CreateConsignment(context.Background(), consignment)

	if err != nil {
		log.Fatalf("create consignment error: %v", err)
	}

	// 新货物是否托运成功
	log.Printf("created: %t", resp.Created)
}
