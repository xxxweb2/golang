package main

import (
	"fmt"
	"time"
)

type Temp struct {
	MChan chan int
}

func main() {
	temp := &Temp{}
	var closedchan = make(chan int)
	temp.MChan = closedchan

	printNum(temp)

	close(temp.MChan)
	time.Sleep(time.Second)
	fmt.Println("全部结束")

}

func printNum(mChan *Temp) {
	numChan := make(chan int)
	//要打印的数字
	n := 0
	go func() {
		for {
			select {
			case a := <-mChan.MChan:
				fmt.Println("程序停止了, num: ", a)
				return
			case numChan <- n:
				fmt.Println("数字: ", n)
				n++
			}
		}

	}()

	time.Sleep(time.Second)
}
