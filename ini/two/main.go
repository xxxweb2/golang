package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func main() {

	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("test page")
	})

	http.ListenAndServe(":8080", nil)

	//f, _ := os.Create("pro.txt")
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	//cfg, err := ini.Load("my.ini")
	//if err != nil {
	//	fmt.Printf("Fail to read file: %v", err)
	//	os.Exit(1)
	//}
	//
	//add(1, 2)
	//// 典型读取操作，默认分区可以使用空字符串表示
	//fmt.Println("App Mode:", cfg.Section("").Key("app_mode").String())
	//fmt.Println("Data Path:", cfg.Section("paths").Key("data").String())
	//
	//// 我们可以做一些候选值限制的操作
	//fmt.Println("Server Protocol:",
	//	cfg.Section("server").Key("protocol").In("http", []string{"http", "https"}))
	//// 如果读取的值不在候选列表内，则会回退使用提供的默认值
	//fmt.Println("Email Protocol:",
	//	cfg.Section("server").Key("protocol").In("smtp", []string{"imap", "smtp"}))
	//
	//// 试一试自动类型转换
	//fmt.Printf("Port Number: (%[1]T) %[1]d\n", cfg.Section("server").Key("http_port").MustInt(9999))
	//fmt.Printf("Enforce Domain: (%[1]T) %[1]v\n", cfg.Section("server").Key("enforce_domain").MustBool(false))
	//
	//// 差不多了，修改某个值然后进行保存
	//cfg.Section("").Key("app_mode").SetValue("production")
	//cfg.SaveTo("my.ini.local")
}

//func add(a, b int) {
//	fmt.Println(a + b)
//}
