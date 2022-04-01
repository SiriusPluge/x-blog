package repository

import (
	"fmt"
	"strings"

	"github.com/ansel1/merry"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var schemaUser = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id SERIAL NOT NULL PRIMARY KEY,
    first_name  VARCHAR(30) NOT NULL,
    last_name  VARCHAR(50) NOT NULL,
    brithday VARCHAR(30),
    gender VARCHAR(10),
	position VARCHAR(100),
	email VARCHAR(100) NOT NULL UNIQUE,
	password VARCHAR(250) NOT NULL,
	registration_date TIMESTAMP NOT NULL,
	access_level INT NOT NULL DEFAULT 0,
	deleted BOOL DEFAULT false
);`, usersTable)

var schemaPost = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id SERIAL NOT NULL UNIQUE PRIMARY KEY,
	title VARCHAR(255),
    author_id INT REFERENCES author (id),
	content_text TEXT,
	date_creation TIMESTAMP NOT NULL,
	last_change TIMESTAMP,
	amount_likes INT DEFAULT 0,
	amount_favorites INT DEFAULT 0,
	access_level INT NOT NULL DEFAULT 0
);`, postTable)

var schemaLikedPost = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
	id SERIAL NOT NULL UNIQUE PRIMARY KEY,
	author_id INT REFERENCES author (id),
	post_id INT REFERENCES post (id)
);`, likedPostTable)

var schemaFavoritePost = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
	id SERIAL NOT NULL UNIQUE PRIMARY KEY,
	author_id INT REFERENCES author (id),
	post_id INT REFERENCES post (id)
);`, favoritPostTable)

var schemaImage = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id  SERIAL NOT NULL UNIQUE PRIMARY KEY,
    link_image TEXT NOT NULL
);`, imageTable)

var schemaBook = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id  SERIAL NOT NULL UNIQUE PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
	year_publication INT,
	link_book TEXT NOT NULL
)`, bookTable)

var schema = strings.Join([]string{schemaUser, schemaPost, schemaLikedPost, schemaFavoritePost, schemaImage, schemaBook}, "\n")

type PostgresDB struct {
	DB *sqlx.DB
}

type ConfigPostgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConnectPostgresDB is used to connect the Postgres Database
func NewConnectionPostgresDB(cfg *ConfigPostgres) *PostgresDB {

	logrus.Info(fmt.Sprintf("Initializing the connection to the database in port: %s \n", cfg.Port))

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		logrus.Fatalf("error open postgresDB: %s", merry.Details(err))
	}
	// defer db.Close()

	errConnDB := db.Ping()
	if errConnDB != nil {
		logrus.Fatalf("error connection postgresDB: %s", errConnDB.Error())
	}

	db.MustExec(schema)

	return &PostgresDB{DB: db}
}
