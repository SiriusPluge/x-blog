package models

import (
	"time"
)

const (
	GuestAccess int = 0 + iota // 0
	UserAccess                 // 1
	AdminAccess                // 2
)

type User struct {
	ID               int       `json:"id"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Brithday         string    `json:"brithday"`
	Gender           string    `json:"gender"`
	Position         string    `json:"position"`
	RegistrationDate time.Time `json:"registration_date"`
	AccessLevel      int       `json:"access_level"`
	Deleted          bool      `json:"deleted"`
}

type UserPublic struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Brithday  string `json:"brithday"`
	Gender    string `json:"gender"`
	Position  string `json:"position"`
}

type Post struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Content         string    `json:"content_text"`
	Author          int       `json:"author"`
	DateCreation    time.Time `json:"date_creation"`
	LastChange      time.Time `json:"last_change"`
	AmountLikes     int       `json:"amount_likes"`
	AmountFavorites int       `json:"amount_favorites"`
	AccessLevel     int       `json:"access_level"`
}

type Image struct {
	ID        int    `json:"id"`
	LinkImage string `json:"link_image"`
}

type Books struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	YearPublication string `json:"year_publication"`
	LinkBook        string `json:"link_book"`
}
