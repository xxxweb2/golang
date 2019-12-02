package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	SERVER_IP       = "127.0.0.1"
	SERVER_PORT     = 8080
	SERVER_RECE_LEN = 10
)

func main() {
	fmt.Println("server has been start ====>")
	address := SERVER_IP + ":" + strconv.Itoa(SERVER_PORT)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		handleConn(conn)

		conn.Close()
	}

}

func handleConn(conn net.Conn) {
	var err error
	for {
		data := make([]byte, SERVER_RECE_LEN)
		_, err = conn.Read(data)
		if err != nil {
			fmt.Println(err)
			break
		}
		strData := string(data)
		fmt.Println("Received:", strData)
		upper := strings.ToUpper(strData)
		_, err = conn.Write([]byte(upper))
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Send:", upper)
	}
}
