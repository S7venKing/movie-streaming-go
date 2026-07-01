package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/database"
	moviesDto "github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/dto/movies"
	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/models"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")

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

func AdminReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieId := c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "MovieId is required"})
			return
		}
		var req struct {
			AdminReview string `json:"admin_review"`
		}

		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		setiment, rankVal, err := GetReviewRanking(req.AdminReview)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting review ranking"})
			return
		}
		filter := bson.M{"imdb_id": movieId}

		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  setiment,
				},
			},
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := movieCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Moview not found"})
			return
		}

		resp.RankingName = setiment
		resp.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, resp)
	}
}

func GetReviewRanking(admin_review string) (string, int, error) {
	rankings, err := GetRankings()
	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""
	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited = sentimentDelimited + ranking.RankingName + ","
		}
	}
	sentimentDelimited = strings.Trim(sentimentDelimited, ",")
	err = godotenv.Load(".env")

	if err != nil {
		log.Println("Warning: .env file not found")
	}
	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	if OpenAiApiKey == "" {
		return "", 0, errors.New("could not read OPENAI_API_KEY")
	}

	llm, err := openai.New(openai.WithToken(OpenAiApiKey))

	if err != nil {
		return "", 0, err
	}

	base_promt_template := os.Getenv("BASE_PROMT_TEMPLATE")
	if base_promt_template == "" {
		return "", 0, errors.New("could not read BASE_PROMT_TEMPLATE")
	}

	base_promt := strings.Replace(base_promt_template, "{rankings}", sentimentDelimited, 1)

	response, err := llm.Call(context.Background(), base_promt+admin_review)

	if err != nil {
		return "", 0, err
	}

	rankVal := 0
	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}
	return response, rankVal, nil
}

func GetRankings() ([]models.Ranking, error) {
	var rankings []models.Ranking
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	cursor, err := rankingCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}
	return rankings, nil
}
