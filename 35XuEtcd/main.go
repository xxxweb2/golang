package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	config := clientv3.Config{
		Endpoints:   []string{"118.25.2.25:2379"},
		DialTimeout: 10 * time.Second,
	}
	client, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}

	kv := clientv3.NewKV(client)
	ctx,_ := context.WithTimeout(context.Background(), 3*time.Second)
	putResp, err := kv.Put(ctx, "/job/v3", "test")
	if err != nil {
		panic(err)
	}
	fmt.Println("putRes:", putResp.PrevKv)

	defer client.Close()
}
