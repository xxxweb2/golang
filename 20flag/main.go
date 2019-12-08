package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var ip = flag.Int("flagname", 1234, "help message for flagname")
	var age = flag.Int("age", 1234, "help message for flagname")
	flag.Parse()
	fmt.Println("ip: ", *ip, age)
	fmt.Println("os:",os.Args)
}
