package service

import (
	"X-Blog/internal/repository"
	"X-Blog/pkg/models"

	"github.com/alexmolinanasaev/exterr"
)

type BlogApp interface {
	// user method`s
	SignUp(user *models.User) (int, exterr.ErrExtender)
	SignIn(email string, password string) (*models.User, exterr.ErrExtender)
	GetUser(id int) (*models.User, exterr.ErrExtender)
	DeleteUser(id int) exterr.ErrExtender

	// post method`s
	CreatePost(incomingPost *models.Post) (int, exterr.ErrExtender)
	GetPost(id int) (*models.Post, exterr.ErrExtender)
	GetPostAuthorID(id int) (int, exterr.ErrExtender)
	EditPost(newPost *models.Post, admin bool) exterr.ErrExtender
	DeletePost(id int) exterr.ErrExtender
	LikePost(userID, postID int) exterr.ErrExtender
	UnlikePost(userID, postID int) exterr.ErrExtender
	AddFavoritPost(userID, postID int) exterr.ErrExtender
	UnfavoritesPost(userID, postID int) exterr.ErrExtender
}

type Service struct {
	BlogApp
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		BlogApp: NewGetService(repos.BlogApp),
	}
}
