package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func main() {
	fmt.Printf("consumer_test\n")

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"198.13.34.80:9092"}, config)
	if err != nil {
		fmt.Printf("consumer_test create consumer error %s\n", err.Error())
		return
	}

	defer consumer.Close()
	partition_consumer, err := consumer.ConsumePartition("kafka_go_test", 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Printf("try create partition_consumer error %s\n", err.Error())
		return
	}

	defer partition_consumer.Close()

	for {
		select {
		case msg := <-partition_consumer.Messages():
			fmt.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s \n",
				msg.Offset, msg.Partition, msg.Timestamp, string(msg.Value))
		case err := <-partition_consumer.Errors():
			fmt.Printf("err : %s\n", err.Error())
		}
	}
}
