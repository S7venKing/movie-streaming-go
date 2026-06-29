package dto

import (
	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/models"
)

type CreateMovieRequest struct {
	ImdbID      string         `json:"imdb_id" validate:"required"`
	Title       string         `json:"title" validate:"required"`
	PosterPath  string         `json:"poster_path"`
	YouTubeID   string         `json:"youtube_id"`
	Genre       []models.Genre `json:"genre"`
	AdminReview string         `json:"admin_review"`
	Ranking     models.Ranking `json:"ranking"`
}

type UpdateMovieRequest struct {
	Title       string         `json:"title"`
	PosterPath  string         `json:"poster_path"`
	YouTubeID   string         `json:"youtube_id"`
	Genre       []models.Genre `json:"genre"`
	AdminReview string         `json:"admin_review"`
	Ranking     models.Ranking `json:"ranking"`
}

type MovieResponse struct {
	ID          string         `json:"id"`
	ImdbID      string         `json:"imdb_id"`
	Title       string         `json:"title"`
	PosterPath  string         `json:"poster_path"`
	YouTubeID   string         `json:"youtube_id"`
	Genre       []models.Genre `json:"genre"`
	AdminReview string         `json:"admin_review"`
	Ranking     models.Ranking `json:"ranking"`
}
