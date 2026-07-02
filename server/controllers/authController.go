package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/database"
	dto "github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/dto/auth"
	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/models"
	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/utils"
	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection("users")

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var request dto.RegisterRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{
			"email": request.Email,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database error",
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Email already exists",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(request.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to hash password",
			})
			return
		}

		user := models.User{
			ID:              bson.NewObjectID(),
			UserID:          uuid.NewString(),
			FirstName:       request.FirstName,
			LastName:        request.LastName,
			Email:           request.Email,
			Password:        string(hashedPassword),
			Role:            "USER",
			Token:           "",
			RefreshToken:    "",
			FavouriteGenres: []models.Genre{},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		_, err = userCollection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to register",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Register successfully",
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var request dto.LoginRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{
			"email": request.Email,
		}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid email or password",
			})
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(request.Password),
		)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid email or password",
			})
			return
		}

		accessToken, err := utils.GenerateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Cannot generate access token",
			})
			return
		}

		refreshToken := uuid.NewString()

		_, err = userCollection.UpdateOne(
			ctx,
			bson.M{
				"user_id": user.UserID,
			},
			bson.M{
				"$set": bson.M{
					"token":         accessToken,
					"refresh_token": refreshToken,
					"updated_at":    time.Now(),
				},
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update token",
			})
			return
		}

		response := dto.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User: dto.UserResponse{
				ID:              user.ID.Hex(),
				UserID:          user.UserID,
				FirstName:       user.FirstName,
				LastName:        user.LastName,
				Email:           user.Email,
				Role:            user.Role,
				FavouriteGenres: user.FavouriteGenres,
				CreatedAt:       user.CreatedAt.Format(time.RFC3339),
				UpdatedAt:       user.UpdatedAt.Format(time.RFC3339),
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Lấy access token từ Authorization header
		tokenString, err := utils.GetAccessToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		// Validate JWT
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid user id",
			})
			return
		}

		result, err := userCollection.UpdateOne(
			ctx,
			bson.M{
				"user_id": claims.UserID,
			},
			bson.M{
				"$set": bson.M{
					"token":         "",
					"refresh_token": "",
					"updated_at":    time.Now(),
				},
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Logout failed",
			})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successfully",
		})
	}
}

func Me() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := utils.GetAccessToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			return
		}

		// User ID trong JWT
		userID := claims.UserID

		var user models.User

		err = userCollection.FindOne(c, bson.M{
			"_id": userID,
		}).Decode(&user)

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database error",
			})
			return
		}

		response := dto.UserResponse{
			ID:              user.ID.Hex(),
			UserID:          user.UserID,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			Role:            user.Role,
			FavouriteGenres: user.FavouriteGenres,
			CreatedAt:       user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       user.UpdatedAt.Format(time.RFC3339),
		}

		c.JSON(http.StatusOK, response)
	}
}
