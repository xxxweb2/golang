package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

type Temp struct {
	MChan chan int
}

type Person struct {
	Name string
	age int
	sex int
}

func main() {
	fmt.Println(runtime.GOROOT())
	fmt.Println(runtime.NumCPU())
	fmt.Println(runtime.Caller(0))
	fmt.Println(os.Hostname())
	fmt.Println(os.Getenv("path"))
	fmt.Println(os.Geteuid())
	//realPath := path.Clean("a/../../../")
	//fmt.Println(realPath)

	//var age = flag.Int("age", 18, "age age")
	//var height = flag.Int("height", 20, "height height")
	//var b = flag.Bool("b", false, "height height")
	//
	//flag.Parse()
	////flag.Usage()
	//fmt.Printf("age=%d  height=%d bool=%t", *age, *height, *b)
	//
	//
	//fmt.Println(os.Args)
}
func Exts(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		fmt.Println("exit:",path[i])
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}
//temp := &Temp{}
//var closedchan = make(chan int)
//temp.MChan = closedchan
//
//printNum(temp)
//
//close(temp.MChan)
//time.Sleep(time.Second)
//fmt.Println("全部结束")
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
