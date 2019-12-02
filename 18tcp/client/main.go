package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	SERVER_IP       = "127.0.0.1"
	SERVER_PORT     = 8080
	SERVER_RECE_LEN = 10
)

func main() {
	address := SERVER_IP + ":" + strconv.Itoa(SERVER_PORT)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		lineLen := len(line)

		n := 0
		for writteen := 0; writteen < lineLen; writteen += n {
			var toWrite string
			if lineLen-writteen > SERVER_RECE_LEN {
				toWrite = line[writteen : writteen+SERVER_RECE_LEN]
			} else {
				toWrite = line[writteen:]
			}

			n, err = conn.Write([]byte(toWrite))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Write:", toWrite)
			msg := make([]byte, SERVER_RECE_LEN)
			n, err = conn.Read(msg)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Response:", string(msg))

		}
	}

}
