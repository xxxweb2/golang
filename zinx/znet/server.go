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

	go func() {

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
		for {
			//如果有客户端连接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//客户端已经建立连接
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err ", err)
						continue
					}

				}
			}()

		}

	}()
}
func (s *Server) Stop() {
	//Todo 将服务器的资源、状态或者一些开启的链接信息 进行停职或回收
}

//运行服务
func (s *Server) Server() {
	s.Start()

	//阻塞状态
	select {}
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
