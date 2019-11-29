package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

func main() {
	cfg, err := ini.Load("my.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	fmt.Println("App Model:", cfg.Section("").Key("app_mode").String())
	fmt.Println("Data Path:", cfg.Section("paths").Key("data").String())
}
