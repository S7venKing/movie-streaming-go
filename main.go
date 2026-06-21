package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, Magic Stream Movie")
	})

	if err := router.Run(":9080"); err != nil {
		fmt.Println("Failed to start server")
	}

}
