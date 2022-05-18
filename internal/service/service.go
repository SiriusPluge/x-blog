package service

import (
	"voting-app/internal/repository"
	"voting-app/pkg/models"

	"github.com/alexmolinanasaev/exterr"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
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

type HyperLedApi interface {
	AddWallet() exterr.ErrExtender
	BuyTokens(tokenAmount int, address string) ([]byte, exterr.ErrExtender)
}

type Service struct {
	BlogApp
	HyperLedApi
}

func NewService(repos *repository.Repository, contract *gateway.Contract) *Service {
	return &Service{
		BlogApp:     NewGetService(repos.BlogApp),
		HyperLedApi: NewEthService(contract),
	}
}
