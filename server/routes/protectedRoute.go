package routes

import (
	controller "github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/controllers"
	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {

	router.Use(middleware.AuthMiddleWare())

	router.POST("/movies/addmovie", controller.CreateMovie())
	router.PUT("/movies/:id", controller.UpdateMovie())
	router.DELETE("/movies/:id", controller.DeleteMovie())

	router.POST("/auth/logout", controller.Logout())
	router.GET("/auth/me", controller.Me())
}
