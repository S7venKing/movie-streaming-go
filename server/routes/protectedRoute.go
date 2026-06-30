package routes

import (
	controller "github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/controllers"
	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleWare())
	movie := router.Group("/movies")
	{
		movie.POST("addmovie", controller.CreateMovie())

		movie.PUT("/:id", controller.UpdateMovie())

		movie.DELETE("/:id", controller.DeleteMovie())
	}

	auth := router.Group("/auth")
	{
		auth.POST("/logout", controller.Logout())
		auth.GET("/me", controller.Me())
	}
}
