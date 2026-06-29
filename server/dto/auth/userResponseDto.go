package dto

import "github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/models"

type UserResponse struct {
	ID              string                  `json:"id"`
	UserID          string                  `json:"user_id"`
	FirstName       string                  `json:"first_name"`
	LastName        string                  `json:"last_name"`
	Email           string                  `json:"email"`
	Role            string                  `json:"role"`
	FavouriteGenres []models.Genre `json:"favourite_genres"`
	CreatedAt       string                  `json:"created_at"`
	UpdatedAt       string                  `json:"updated_at"`
}