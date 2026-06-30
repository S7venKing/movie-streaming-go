package routes

import (
	controller "github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/controllers"
	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/middleware"
	"github.com/gin-gonic/gin"
)

func SetupUnProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleWare())

	movie := router.Group("/movies")
	{
		movie.GET("", controller.GetMovies())
		movie.GET("/:id", controller.GetMovieByID())
	}

	auth := router.Group("/auth")
	{		auth.POST("/register", controller.Register())
		auth.POST("/login", controller.Login())

	}
}
