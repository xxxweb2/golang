package main

import "github.com/gin-gonic/gin"

func main() {
	c := gin.Default()

	c.GET("/test", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "hello world",
		})
	})
	c.GET("/tes/:name/c", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "hello name",
		})
	})
	c.Run()

}
