package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	reQQEmail = `(\d+)@qq.com`
)

func GetEmail() {
	//1.去网站拿数据
	resp, err := http.Get("https://tieba.baidu.com/p/6051076813?red_tag=1573533731")
	HandleError(err, "http.get url")

	defer resp.Body.Close()

	//2.读取页面内容
	pageBytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "ioUtil.ReadAll")

	//字节转字符串
	pageStr := string(pageBytes)
	//3.过滤数据
	re := regexp.MustCompile(reQQEmail)
	results := re.FindAllStringSubmatch(pageStr, -1)
	for _, result := range results {
		fmt.Println("email: ", result[0])
		fmt.Println("qq: ", result[1])
	}
}

func HandleError(err error, why string) {
	if err != nil {
		fmt.Println(err, why)
	}
}

func main() {
	GetEmail()
}
