package main

import (
	"fmt"

	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, Magic Stream Movie")
	})

	routes.SetupUnProtectedRoutes(router)
	routes.SetupProtectedRoutes(router)

	if err := router.Run(":9080"); err != nil {
		fmt.Println("Failed to start server")
	}

}
