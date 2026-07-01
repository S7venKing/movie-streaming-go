package routes

import (
	controller "github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/controllers"
	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	movies := router.Group("/movies")
	movies.Use(middleware.AuthMiddleWare())
	{
		movies.POST("/addmovie", controller.CreateMovie())
		movies.PUT("/:id", controller.UpdateMovie())
		movies.DELETE("/:id", controller.DeleteMovie())
	}

	auth := router.Group("/auth")
	auth.Use(middleware.AuthMiddleWare())
	{
		auth.POST("/logout", controller.Logout())
		auth.GET("/me", controller.Me())
	}
}
