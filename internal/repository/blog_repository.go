package repository

import (
	"X-Blog/pkg/models"
	"database/sql"
	"fmt"

	"github.com/alexmolinanasaev/exterr"
)

type GetItemPostgres struct {
	db *PostgresDB
}

func NewGetItemPostgres(db *PostgresDB) *GetItemPostgres {
	return &GetItemPostgres{
		db: db,
	}
}

func (p *GetItemPostgres) SignUp(u *models.User) (int, exterr.ErrExtender) {
	var id int

	// FIXME: лишняя логика в репо; перенести эту прверку в сервис
	userBool := p.checkUserByEmail(u.Email)
	if userBool {
		return 0, exterr.New("Repository user already exist")
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (
			first_name,
			last_name,
			brithday,
			gender,
			position,
			email,
			password,
			access_level,
			registration_date
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
		ON CONFLICT (email) DO NOTHING
		RETURNING id`, usersTable)

	row := p.db.DB.QueryRow(query, u.FirstName, u.LastName, u.Brithday, u.Gender, u.Position, u.Email, u.Password, models.UserAccess)
	if err := row.Scan(&id); err != nil {
		return 0, exterr.NewWithErr("User not created", err)
	}

	return id, nil
}

// Найти пользователя в БД по почте
func (p *GetItemPostgres) GetUserByEmail(email string) (*models.User, exterr.ErrExtender) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE email=$1`, usersTable)

	row := p.db.DB.QueryRow(query, email)
	user := &models.User{}
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Brithday,
		&user.Gender,
		&user.Position,
		&user.Email,
		&user.Password,
		&user.RegistrationDate,
		&user.AccessLevel,
		&user.Deleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exterr.NewWithErr("Repository: GetUserByEmail not found", err)
		}
		return nil, exterr.NewWithErr("Repository: GetUserByEmail error", err)
	}

	return user, nil
}

// Найти пользователя в БД по id
func (p *GetItemPostgres) GetUserByID(id int) (*models.User, exterr.ErrExtender) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id=$1`, usersTable)

	row := p.db.DB.QueryRow(query, id)
	user := &models.User{}
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Brithday,
		&user.Gender,
		&user.Position,
		&user.Email,
		&user.Password,
		&user.RegistrationDate,
		&user.AccessLevel,
		&user.Deleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exterr.NewWithErr("Repository: GetUserByID not found", err)
		}
		return nil, exterr.NewWithErr("Repository: GetUserByID error", err)
	}

	return user, nil
}

func (p *GetItemPostgres) DeleteUser(id int) exterr.ErrExtender {
	query := fmt.Sprintf(`
		UPDATE %s SET deleted=true WHERE id=$1`, usersTable)

	_, err := p.db.DB.Exec(query, id)
	if err != nil {
		return exterr.NewWithErr("Repository: DeleteUser error", err)
	}
	return nil
}

func (p *GetItemPostgres) RecoverUser(id int) exterr.ErrExtender {
	query := fmt.Sprintf(`
		UPDATE %s SET deleted=false WHERE id=$1`, usersTable)

	_, err := p.db.DB.Exec(query, id)
	if err != nil {
		return exterr.NewWithErr("Repository: RecoverUser error", err)
	}
	return nil
}

func (p *GetItemPostgres) CreatePost(post *models.Post) (int, exterr.ErrExtender) {
	postID := 0
	query := fmt.Sprintf(`
	INSERT INTO %s 
		(title, author_id, content_text, amount_likes, access_level, date_creation, last_change) 
	VALUES 
		($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
	RETURNING id`, postTable)

	row := p.db.DB.QueryRow(query, post.Title, post.Author, post.Content, post.AmountLikes, post.AccessLevel)
	if err := row.Scan(&postID); err != nil {
		return postID, exterr.NewWithErr("Repository: CreatePost error", err)
	}

	return postID, nil
}

func (p *GetItemPostgres) GetPost(id int) (*models.Post, exterr.ErrExtender) {
	post := &models.Post{}

	query := fmt.Sprintf(`
	SELECT
		id,
		title,
		author_id,
		content_text,
		date_creation,
		last_change,
		amount_likes,
		amount_favorites,
		access_level
	FROM %s
	WHERE id=$1`, postTable)

	row := p.db.DB.QueryRow(query, id)
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Author,
		&post.Content,
		&post.DateCreation,
		&post.LastChange,
		&post.AmountLikes,
		&post.AmountFavorites,
		&post.AccessLevel,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exterr.NewWithErr("[Repo: GetPost] not found", err)
		}
		return nil, exterr.NewWithErr("[Repo: GetPost] scan error", err)
	}

	return post, nil
}

func (p *GetItemPostgres) GetPostAuthorID(id int) (int, exterr.ErrExtender) {
	query := fmt.Sprintf(`SELECT author_id FROM %s WHERE id=$1`, postTable)

	row := p.db.DB.QueryRow(query, id)
	authorID := 0
	err := row.Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, exterr.NewWithErr("[Repo: GetPost] not found", err)
		}
		return 0, exterr.NewWithErr("[Repo: GetPost] scan error", err)
	}

	return authorID, nil
}

func (p *GetItemPostgres) DeletePost(id int) exterr.ErrExtender {
	query := fmt.Sprintf(`
		DELETE FROM %s WHERE id=$1`, postTable)

	_, err := p.db.DB.Exec(query, id)
	if err != nil {
		return exterr.NewWithErr("Repository: DeletePost error", err)
	}
	return nil
}

func (p *GetItemPostgres) LikePost(userID, postID, amountLikes int) exterr.ErrExtender {
	// Добавление лайка к посту
	qUpdatePost := fmt.Sprintf(`
		UPDATE %s SET amount_likes=$1 WHERE id=$2`, postTable)

	_, err := p.db.DB.Exec(qUpdatePost, amountLikes+1, postID)
	if err != nil {
		return exterr.NewWithErr("Repository: LikePost UPDATE post table error", err)
	}

	// Обновление связующей таблицы (Идентифицирующей кто поставил лайк)
	qUpdateLikedPost := fmt.Sprintf(`
		INSERT INTO %s
			(author_id,
			post_id)
		VALUES($1, $2)`, likedPostTable)

	_, err = p.db.DB.Exec(qUpdateLikedPost, userID, postID)
	if err != nil {
		return exterr.NewWithErr("Repository: LikePost INSERT INTO liked_post table error", err)
	}

	return nil
}

// FIXME: исправить название
// Получение информации ставил ли лайк к посту пользователь
func (p *GetItemPostgres) GetCheckLikePost(userID, postID int) exterr.ErrExtender {
	var post models.Post

	qLikeUser := fmt.Sprintf("SELECT id FROM %s WHERE author_id=$1 AND post_id=$2", likedPostTable)
	row := p.db.DB.QueryRow(qLikeUser, userID, postID)
	err := row.Scan(
		&post.ID,
	)
	// Если имеется ошибка, то соответсвенно лайк поставлен не был и возвращаем error
	if err != nil {
		return exterr.New("there are no matches in the database")
	}

	// Если ошибки не имеется, то соответсвенно лайк был поставлен и возращаем nil
	return nil
}

func (p *GetItemPostgres) GetAmountLikePost(postID int) (int, exterr.ErrExtender) {
	var post models.Post

	// Получение количества лайков у поста
	qAmountLike := fmt.Sprintf("SELECT amount_likes FROM %s WHERE id=$1", postTable)

	row := p.db.DB.QueryRow(qAmountLike, postID)
	err := row.Scan(
		&post.AmountLikes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, exterr.NewWithErr("Repository: GetAmountLikePost post not found", err)
		}
		return 0, exterr.NewWithErr("Repository: GetAmountLikePost error", err)
	}

	return post.AmountLikes, nil
}

func (p *GetItemPostgres) UnlikePost(userID, postID, amountLikes int) exterr.ErrExtender {
	var id int

	amountLikes -= 1
	if amountLikes < 0 {
		amountLikes = 0
	}

	// Снятие лайка к посту
	qUpdatePost := fmt.Sprintf(`
		UPDATE %s SET amount_likes=$1 WHERE id=$2`, postTable)

	_, err := p.db.DB.Exec(qUpdatePost, amountLikes, postID)
	if err != nil {
		return exterr.NewWithErr("update error like in post", err)
	}

	// Получение id из связующей таблицы (для последующего удаления)
	queryID := fmt.Sprintf(`
		SELECT id
		FROM %s	
		WHERE author_id=$1 AND post_id=$2`, likedPostTable)

	row := p.db.DB.QueryRow(queryID, userID, postID)
	err = row.Scan(
		&id,
	)
	if err != nil {
		return exterr.NewWithErr("error update(unliked) liked_post Table", err)
	}

	// Обновление связующей таблицы (Идентифицирующей кто поставил лайк)
	query := fmt.Sprintf(`
		DELETE FROM %s WHERE id=$1`, likedPostTable)

	_, err = p.db.DB.Exec(query, id)
	if err != nil {
		return exterr.NewWithErr("delete error likes for post", err)
	}

	return nil
}

func (p *GetItemPostgres) UpdatePost(post *models.Post) exterr.ErrExtender {
	query := fmt.Sprintf(`
	UPDATE %s SET
		title = $2,
		content_text = $3,  
		last_change = CURRENT_TIMESTAMP,
		access_level = $4
	WHERE id = $1`, postTable)

	_, err := p.db.DB.Exec(query, post.ID, post.Title, post.Content, post.AccessLevel)
	if err != nil {
		return exterr.NewWithErr("Repository: UpdatePost error", err)
	}

	return nil
}

func (p *GetItemPostgres) GetCheckFavoritesPost(userID, postID int) exterr.ErrExtender {
	var id int

	qLikeUser := fmt.Sprintf("SELECT id FROM %s WHERE author_id=$1 AND post_id=$2", favoritPostTable)
	row := p.db.DB.QueryRow(qLikeUser, userID, postID)
	err := row.Scan(
		&id,
	)
	// Если имеется ошибка, то соответсвенно в избранное не добавлено и возвращаем error
	if err != nil {
		return exterr.NewWithErr("there are no matches in the database", err)
	}

	// Если ошибки не имеется, то соответсвенно лайк был поставлен и возращаем nil
	return nil
}

func (p *GetItemPostgres) GetAmountFavoritesPost(postID int) (int, exterr.ErrExtender) {
	var post models.Post

	// Получение количества лайков у поста
	qAmountLike := fmt.Sprintf("SELECT amount_favorites FROM %s WHERE id=$1", postTable)

	row := p.db.DB.QueryRow(qAmountLike, postID)
	err := row.Scan(
		&post.AmountLikes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, exterr.NewWithErr("Repository: GetAmountFavoritesPost post not found", err)
		}
		return 0, exterr.NewWithErr("Repository: GetAmountFavoritPost error", err)
	}

	return post.AmountFavorites, nil
}

func (p *GetItemPostgres) FavoritesPost(userID, postID, amountFavorites int) exterr.ErrExtender {
	// Добавление количества избранных к посту
	qUpdatePost := fmt.Sprintf(`
		UPDATE %s SET amount_favorites=$1 WHERE id=$2`, postTable)

	_, err := p.db.DB.Exec(qUpdatePost, amountFavorites+1, postID)
	if err != nil {
		return exterr.NewWithErr("Repository: FavoritesPost UPDATE post table error", err)
	}

	// Обновление связующей таблицы (Идентифицирующей кто добавил в избранное)
	qUpdateLikedPost := fmt.Sprintf(`
		INSERT INTO %s
			(author_id,
			post_id)
		VALUES($1, $2)`, favoritPostTable)

	_, err = p.db.DB.Exec(qUpdateLikedPost, userID, postID)
	if err != nil {
		return exterr.NewWithErr("Repository: FavoritesPost INSERT INTO favorit_post table error", err)
	}

	return nil
}

func (p *GetItemPostgres) UnfavoritesPost(userID, postID, amountFavorites int) exterr.ErrExtender {
	var id int
	amountFavorites -= 1
	if amountFavorites < 0 {
		amountFavorites = 0
	}

	// Убрать пост из избранного
	qUpdatePost := fmt.Sprintf(`
		UPDATE %s SET amount_favorites=$1 WHERE id=$2`, postTable)

	_, err := p.db.DB.Exec(qUpdatePost, amountFavorites, postID)
	if err != nil {
		return exterr.NewWithErr("update error unfavorites post", err)
	}

	// Получение id из связующей таблицы (для последующего удаления)
	queryID := fmt.Sprintf(`
		SELECT id
		FROM %s	
		WHERE author_id=$1 AND post_id=$2`, favoritPostTable)

	row := p.db.DB.QueryRow(queryID, userID, postID)
	err = row.Scan(
		&id,
	)
	if err != nil {
		return exterr.NewWithErr("error update(unfavorites) favorit_post Table", err)
	}

	// Обновление связующей таблицы (Идентифицирующей кто добавил избранное)
	query := fmt.Sprintf(`
		DELETE FROM %s WHERE id=$1`, favoritPostTable)

	_, err = p.db.DB.Exec(query, id)
	if err != nil {
		return exterr.NewWithErr("delete error unfavorites post", err)
	}

	return nil
}
