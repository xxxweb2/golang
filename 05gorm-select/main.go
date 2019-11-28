package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type City struct {
	UserId   int
	Username string
	Sex      string
	Email    string
}

func (c *City) TableName() string {
	return "city"
}

func main() {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// 全局禁用表名复数
	db.SingularTable(true)

	var city City

	db.First(&city)

	//第一条记录
	fmt.Println("First: ", city)
	//最后一条记录
	db.Last(&city)
	fmt.Println("Last", city)
	//所有记录
	var citys = make([]City, 0, 0)
	db.Find(&citys)
	fmt.Println("all: ", citys)


}
