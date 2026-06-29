package main

import (
	"fmt"

	controller "github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, Magic Stream Movie")
	})

	movie := router.Group("/movies")
	{
		movie.GET("", controller.GetMovies())
		movie.GET("/:id", controller.GetMovieByID())

		movie.POST("", controller.CreateMovie())

		movie.PUT("/:id", controller.UpdateMovie())

		movie.DELETE("/:id", controller.DeleteMovie())
	}

	if err := router.Run(":9080"); err != nil {
		fmt.Println("Failed to start server")
	}

}
