package znet

import (
	"fmt"
	"net"
	"study/golang/zinx/ziface"
)

//iServer
type Server struct {
	//服务器名称
	Name string
	//ip版本
	IPVersion string
	//监听的ip
	IP string
	//端口
	Port int
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP: %s, Port: %d, is starting\n", s.IP, s.Port)
	// 1. 获取一个tcp的addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error:", err)
		return
	}
	//2.监听地址
	listenner, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen ", s.IPVersion, " err ", err)
		return
	}
	fmt.Println("start Zinx server succ, ", s.Name, " succ, Listenning...")

	//阻塞的等待客户端连接，处理客户端业务
	for  {
		//如果有客户端连接过来，阻塞会返回
		listenner.AcceptTCP()
	}
}
func (s *Server) Stop() {

}
func (s *Server) Server() {
	s.Start()
}

//初始化Server模块
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}

	return s
}
