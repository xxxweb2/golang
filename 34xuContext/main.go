package main

import (
	"context"
	"fmt"
	"math"
	"study/golang/34xuContext/xuContext"
	"time"
)

func main() {
	ctx := context.WithValue(context.Background(), "logId", 123123)
	printNum3(ctx)
}

func printNum3(ctx context.Context) {
	logId, ok := ctx.Value("logId").(int)
	if !ok {
		fmt.Println("err", ok)
		return
	}

	fmt.Println("logId=", logId)
}
func printNum2(ctx xuContext.Context) {

	//要打印的数字
	n := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("结束")
			return
		default:
			fmt.Println("数字: ", n)
			n++
		}
		time.Sleep(time.Second)
	}
}

//ctx, cancel := xuContext.WithCancel(context.Background())
//numChan := printNum(ctx)
//for n := range numChan {
//	if n > 2 {
//		cancel()
//		break
//	}
//}
//time.Sleep(time.Second)
//fmt.Println("全部结束")
//func printNum(ctx xuContext.Context) chan int {
//	numChan := make(chan int)
//	//要打印的数字
//	n := 0
//	go func() {
//		for {
//			select {
//			case <-ctx.Done():
//				fmt.Println("程序停止了, num: ", n)
//				return
//			case numChan <- n:
//				fmt.Println("数字: ", n)
//				n++
//			}
//		}
//	}()
//
//	return numChan
//}
