package main

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"study/golang/21nsq-request"
)

func main() {
	url := "118.25.2.25:4150"
	producer, err := xunsq.NewProducer(url, xunsq.NewConfig())
	if err != nil {
		panic(err)
	}

	err = producer.Publish("test",[]byte("hello world"))
	if err != nil {
		fmt.Println("推送失败")
		panic(err)
	}

	fmt.Println("推送成功")

	producer.Stop()
}
