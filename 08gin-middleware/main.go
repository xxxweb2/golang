package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

//定义中间件
func MiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		t := time.Now()
		fmt.Println("中间件开始执行了")
		context.Set("request", "中间件")
		status := context.Writer.Status()
		fmt.Println("中间件执行完毕", status)
		t2 := time.Since(t)
		fmt.Println("time", t2)
	}
}

func main() {
	r := gin.Default()
	//r.Use(MiddleWare())

	r.GET("/ce", func(c *gin.Context) {
		req, _ := c.Get("request")
		fmt.Println("request:", req)
		c.JSON(200, gin.H{"request": req})
	})

	r.Run()
}
