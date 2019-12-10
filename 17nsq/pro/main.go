package main

import (
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
)

func main() {
	for i := 0; i < 10; i++ {
		sendMessage()
	}

	time.Sleep(time.Second * 10)
}

func sendMessage() {
	url := "118.25.2.25:4150"
	producer, err := nsq.NewProducer(url, nsq.NewConfig())
	if err != nil {
		panic(err)
	}
	fmt.Println("实例化成功")
	err = producer.Publish("test", []byte("hello world"))
	if err != nil {
		fmt.Println("推送失败")
		panic(err)
	}
	fmt.Println("推送成功")
	producer.Stop()
}
