package main

import "study/golang/zinx/znet"

func main() {
	//创建
	s := znet.NewServer("[zinx V0.1]")
	s.Server()
}
