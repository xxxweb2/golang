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

	city := City{
		Username: "xuxinxin",
	}

	tx := db.Begin()
	err = tx.Create(city).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("数据回滚, ", db.Error)
	}
	tx.Commit()

	//原生sql
	err = db.Exec("delete from city where user_id = ?", 5).Error
	fmt.Println(err)
	var tmp Tmp
	db.Raw("select username, sex from city where user_id = ?", 4).Scan(&tmp)
	fmt.Println("res", tmp)
}

type Tmp struct {
	Username string
	Sex      string
}
