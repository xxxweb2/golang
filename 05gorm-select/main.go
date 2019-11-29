package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type City struct {
	UserId   int `db:"user_id"`
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
	city = City{}
	db.Last(&city)
	fmt.Println("Last", city)
	//所有记录
	var citys = make([]City, 0, 0)
	db.Find(&citys)
	fmt.Println("all: ", citys)

	//where
	city = City{}
	db.Where("Username = ?", "stu001").First(&city)
	fmt.Println("where 获取第一个匹配记录", city)

	//where in
	city = City{}
	db.Where("Username in (?)", []string{"stu001", "stu002"}).Find(&citys)
	fmt.Println("where in ", citys)

	//select
	city = City{}
	db.Select("user_id, username, sex").Where("user_id = ?", 3).First(&city)
	fmt.Println("where select ", city)

	//order
	city = City{}
	db.Order("user_id desc,user_id asc").Offset(1).Limit(2).First(&city)
	fmt.Println("where order", city)

	tmps := make([]Tmp, 0, 0)
	db = db.Table("city").Select("user_id as uid,sex as se").Group("user_id").Scan(&tmps)
	err = db.Error
	fmt.Println("error: ", err)

	fmt.Println(&tmps)

}

type Tmp struct {
	Uid int
	Se  string
}
