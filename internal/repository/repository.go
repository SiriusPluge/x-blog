package repository

import (
	"X-Blog/pkg/models"

	"github.com/alexmolinanasaev/exterr"
)

type BlogApp interface {
	// user method`s
	SignUp(user *models.User) (int, exterr.ErrExtender)
	DeleteUser(id int) exterr.ErrExtender
	GetUserByEmail(email string) (*models.User, exterr.ErrExtender)
	GetUserByID(id int) (*models.User, exterr.ErrExtender)

	// post method`s
	CreatePost(*models.Post) (int, exterr.ErrExtender)
	GetPost(id int) (*models.Post, exterr.ErrExtender)
	GetPostAuthorID(id int) (int, exterr.ErrExtender)
	UpdatePost(post *models.Post) exterr.ErrExtender
	DeletePost(id int) exterr.ErrExtender
	LikePost(userID, postID, amountLikes int) exterr.ErrExtender
	GetCheckLikePost(userID, postID int) exterr.ErrExtender
	GetAmountLikePost(postID int) (int, exterr.ErrExtender)
	GetCheckFavoritesPost(userID, postID int) exterr.ErrExtender
	GetAmountFavoritesPost(postID int) (int, exterr.ErrExtender)
	UnlikePost(userID, postID, amountLikes int) exterr.ErrExtender
	FavoritesPost(userID, postID, amountFavorites int) exterr.ErrExtender
	UnfavoritesPost(userID, postID, amountFavorites int) exterr.ErrExtender
}

type Repository struct {
	BlogApp
}

func NewRepository(db *PostgresDB) *Repository {
	return &Repository{
		BlogApp: NewGetItemPostgres(db),
	}
}
