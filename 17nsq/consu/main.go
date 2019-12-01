package main

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"sync"
)

func main() {
	testNSQ()
}

type NSQHandler struct {
}

func (this *NSQHandler) HandleMessage(msg *nsq.Message) error {
	fmt.Println("11111111111111")
	fmt.Println("receive", msg.NSQDAddress, "message:", string(msg.Body))
	return nil
}

func testNSQ() {
	url := "118.25.2.25:4150"
	waiter := sync.WaitGroup{}
	waiter.Add(1)
	go func() {
		defer waiter.Done()
		config := nsq.NewConfig()
		config.MaxInFlight = 9
		for i := 0; i < 10; i++ {
			consumer, err := nsq.NewConsumer("test", "nsq_to_file", config)
			if err != nil {
				fmt.Println("err: ", err)
				return
			}

			fmt.Println("实例成功")

			consumer.AddHandler(&NSQHandler{})
			err = consumer.ConnectToNSQD(url)
			if err != nil {
				fmt.Println("err:", err)
				return
			}

			fmt.Println("连接成功")
		}
		select {}
	}()
	waiter.Wait()
}
