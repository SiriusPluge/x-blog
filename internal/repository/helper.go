package repository

import (
	"X-Blog/pkg/models"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	usersTable       = "author"
	postTable        = "post"
	likedPostTable   = "liked_post"
	favoritPostTable = "favorit_post"
	imageTable       = "image"
	bookTable        = "books"
)

func (p *GetItemPostgres) checkUserByEmail(email string) bool {
	var user models.User
	query := fmt.Sprintf(`
		SELECT
			id,
			first_name,
			last_name,
			brithday,
			gender,
			position,
			email,
			password,
			registration_date,
			access_level
		FROM %s
		WHERE email=$1`, usersTable)

	row := p.db.DB.QueryRow(query, email)
	err := row.Scan(
		user.ID,
		user.FirstName,
		user.LastName,
		user.Brithday,
		user.Gender,
		user.Position,
		user.Email,
		user.Password,
		user.RegistrationDate,
		user.AccessLevel,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("NOT FOUND USER BY EMAIL", err.Error())
			return false
		}

		logrus.Errorf("ERROR FOR SCAN USER BY EMAIL", err.Error())
		return false
	}

	return true
}
