package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/database"
	moviesDto "github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/dto/movies"
	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")

func mapMovieResponse(movie models.Movie) moviesDto.MovieResponse {
	return moviesDto.MovieResponse{
		ID:          movie.ID.Hex(),
		ImdbID:      movie.ImdbID,
		Title:       movie.Title,
		PosterPath:  movie.PosterPath,
		YouTubeID:   movie.YouTubeID,
		Genre:       movie.Genre,
		AdminReview: movie.AdminReview,
		Ranking:     movie.Ranking,
	}
}

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to fetch movies",
			})
			return
		}
		defer cursor.Close(ctx)

		var movies []models.Movie

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to decode movies",
			})
			return
		}

		responses := make([]moviesDto.MovieResponse, 0, len(movies))

		for _, movie := range movies {
			responses = append(responses, mapMovieResponse(movie))
		}

		c.JSON(http.StatusOK, responses)
	}
}

func GetMovieByID() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id := c.Param("id")

		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid movie id",
			})
			return
		}

		var movie models.Movie

		err = movieCollection.FindOne(ctx, bson.M{
			"_id": objectID,
		}).Decode(&movie)

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Movie not found",
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, mapMovieResponse(movie))
	}
}

func CreateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var request moviesDto.CreateMovieRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		movie := models.Movie{
			ID:          bson.NewObjectID(),
			ImdbID:      request.ImdbID,
			Title:       request.Title,
			PosterPath:  request.PosterPath,
			YouTubeID:   request.YouTubeID,
			Genre:       request.Genre,
			AdminReview: request.AdminReview,
			Ranking:     request.Ranking,
		}

		_, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create movie",
			})
			return
		}

		c.JSON(http.StatusCreated, mapMovieResponse(movie))
	}
}

func UpdateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id := c.Param("id")

		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid movie id",
			})
			return
		}

		var request moviesDto.UpdateMovieRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		update := bson.M{
			"$set": bson.M{
				"title":        request.Title,
				"poster_path":  request.PosterPath,
				"youtube_id":   request.YouTubeID,
				"genre":        request.Genre,
				"admin_review": request.AdminReview,
				"ranking":      request.Ranking,
			},
		}

		result, err := movieCollection.UpdateOne(
			ctx,
			bson.M{"_id": objectID},
			update,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update movie",
			})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Movie not found",
			})
			return
		}

		var movie models.Movie

		_ = movieCollection.FindOne(ctx, bson.M{
			"_id": objectID,
		}).Decode(&movie)

		c.JSON(http.StatusOK, mapMovieResponse(movie))
	}
}

func DeleteMovie() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id := c.Param("id")

		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid movie id",
			})
			return
		}

		result, err := movieCollection.DeleteOne(ctx, bson.M{
			"_id": objectID,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to delete movie",
			})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Movie not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Movie deleted successfully",
		})
	}
}
