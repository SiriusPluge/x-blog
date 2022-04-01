package service

import (
	"X-Blog/internal/repository"
	"X-Blog/pkg/models"

	"github.com/alexmolinanasaev/exterr"
	"github.com/ethereum/go-ethereum/ethclient"
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

type EthApi interface {
	
}

type Service struct {
	BlogApp
	EthApi
}

func NewService(repos *repository.Repository, ethClient *ethclient.Client) *Service {
	return &Service{
		BlogApp: NewGetService(repos.BlogApp),
		EthApi: NewEthService(ethClient),
	}
}
