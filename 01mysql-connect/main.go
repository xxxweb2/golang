package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	database, err := sqlx.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	fmt.Println(database, err)
}
