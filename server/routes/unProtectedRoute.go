package routes

import (
	controller "github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUnProtectedRoutes(router *gin.Engine) {

	router.GET("/movies", controller.GetMovies())
	router.GET("/movies/:id", controller.GetMovieByID())

	router.POST("/auth/register", controller.Register())
	router.POST("/auth/login", controller.Login())
}
