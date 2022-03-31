package service

import (
	"X-Blog/internal/repository"
	"X-Blog/pkg/models"
	"net/http"

	"github.com/alexmolinanasaev/exterr"
)

const (
	salt = "hjqrhjqw124617ajfhajs"
)

type GetService struct {
	repo repository.BlogApp
}

func NewGetService(repo repository.BlogApp) *GetService {
	return &GetService{
		repo: repo,
	}
}

// TODO: Добавить валидацию почты и проверку надёжности пароля
// Регистрация
func (s *GetService) SignUp(u *models.User) (int, exterr.ErrExtender) {
	if u.FirstName == "" || u.LastName == "" || u.Email == "" || u.Password == "" {
		return 0, exterr.New("Empty user info").
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Empty user info (FirstName, LastName, Email or Password)")

	}

	// Проверка, существует ли пользователь
	_, errExt := s.repo.GetUserByEmail(u.Email)
	if errExt == nil {
		return 0, exterr.New("Service: user found, can not create new user").
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Email already taken")
	}

	u.AccessLevel = models.UserAccess             // Уровень доступа при регистрации
	u.Password = GeneratePasswordHash(u.Password) // Записываем хешированный пароль

	// Запись в базу
	id, errExt := s.repo.SignUp(u)
	if errExt != nil {
		return 0, exterr.NewWithExtErr("Service can not add new user in repository", errExt).
			SetErrCode(http.StatusInternalServerError).
			SetAltMsg(http.StatusText(http.StatusInternalServerError))
	}

	return id, nil
}

// Аутентификация(вход)
func (s *GetService) SignIn(email string, password string) (*models.User, exterr.ErrExtender) {
	hashedPassword := GeneratePasswordHash(password)

	user, errExt := s.repo.GetUserByEmail(email)
	if errExt != nil {
		return nil, exterr.NewWithExtErr("Service: user not found ", errExt).
			SetErrCode(http.StatusNotFound).
			SetAltMsg("User not found")
	}
	if user.Deleted {
		return nil, exterr.New("Service: GetUser error").
			SetErrCode(http.StatusUnauthorized).
			SetAltMsg("User deleted")
	}

	if hashedPassword != user.Password {
		return nil, exterr.New("Wrong password").
			SetErrCode(http.StatusUnauthorized).
			SetAltMsg("Wrong password")
	}

	return user, nil
}

// Получение информации о пользователе
func (s *GetService) GetUser(id int) (*models.User, exterr.ErrExtender) {
	// Проверка, существует ли пользователь
	user, errExt := s.repo.GetUserByID(id)
	if errExt != nil {
		return nil, exterr.NewWithExtErr("Service: GetUser error", errExt).
			SetErrCode(http.StatusNotFound).
			SetAltMsg("User not found")
	}
	if user.Deleted {
		return nil, exterr.New("Service: GetUser error").
			SetErrCode(http.StatusConflict).
			SetAltMsg("User deleted")
	}

	return user, nil
}

// Удаление пользователя
func (s *GetService) DeleteUser(id int) exterr.ErrExtender {
	// Проверка, существует ли пользователь
	user, errExt := s.repo.GetUserByID(id)
	if errExt != nil {
		return exterr.NewWithExtErr("Service: DeleteUser error", errExt).
			SetErrCode(http.StatusNotFound).
			SetAltMsg("User not found")
	}
	if user.Deleted {
		return exterr.New("Service: DeleteUser error").
			SetErrCode(http.StatusConflict).
			SetAltMsg("User already deleted")
	}

	// Запись в базу
	errExt = s.repo.DeleteUser(id)
	if errExt != nil {
		return exterr.NewWithExtErr("Service: DeleteUser error", errExt).
			SetErrCode(http.StatusInternalServerError).
			SetAltMsg(http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

//   Post
//  Создание статьи
func (s *GetService) CreatePost(incomingPost *models.Post) (int, exterr.ErrExtender) {
	// Проверка заголовка статьи
	if incomingPost.Title == "" {
		return -1, exterr.New("Service: CreatePost empty title").
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Empty title")
	}

	// Проверка автора статьи
	if incomingPost.Author <= 0 {
		return -1, exterr.New("Service: CreatePost empty author").
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Empty author")
	}

	post := &models.Post{
		Title:       incomingPost.Title,
		Author:      incomingPost.Author,
		Content:     incomingPost.Content,
		AccessLevel: models.UserAccess,
	}

	// Запись в базу
	postID, errExt := s.repo.CreatePost(post)
	if errExt != nil {
		return -1, exterr.NewWithExtErr("Service: CreatePost error", errExt).
			SetErrCode(http.StatusInternalServerError).
			SetAltMsg(http.StatusText(http.StatusInternalServerError))
	}

	return postID, nil
}

func (s *GetService) GetPost(id int) (*models.Post, exterr.ErrExtender) {
	// Проверка id
	if id == 0 {
		return nil, exterr.New("Service: GetPost empty id").
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Empty post id")
	}

	// Взаимодействие с базой данных
	post, errExt := s.repo.GetPost(id)
	if errExt != nil {
		return nil, exterr.NewWithExtErr("Service: GetPost error", errExt).
			SetErrCode(http.StatusInternalServerError).
			SetAltMsg(http.StatusText(http.StatusInternalServerError))
	}

	return post, nil
}

func (s *GetService) GetPostAuthorID(id int) (int, exterr.ErrExtender) {
	// Проверка id
	if id == 0 {
		return 0, exterr.New("Service: GetPostAuthor empty id").
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Empty post id")
	}

	// Взаимодействие с базой данных
	authorID, errExt := s.repo.GetPostAuthorID(id)
	if errExt != nil {
		return 0, exterr.NewWithExtErr("Service: GetPostAuthor error", errExt).
			SetErrCode(http.StatusInternalServerError).
			SetAltMsg(http.StatusText(http.StatusInternalServerError))
	}

	return authorID, nil
}

func (s *GetService) EditPost(newPost *models.Post, admin bool) exterr.ErrExtender {
	// Получаем статью из базы
	oldPost, errExt := s.repo.GetPost(newPost.ID)
	if errExt != nil {
		return exterr.NewWithExtErr("Service: EditPost post not found", errExt).
			SetErrCode(http.StatusNotFound).
			SetAltMsg("Post not found")
	}

	// Заголовок, если в новом пустой, то оставить старый
	if newPost.Title == "" {
		newPost.Title = oldPost.Title
	}

	// Контент, если в новом пустой, то оставить старый
	if newPost.Content == "" {
		newPost.Content = oldPost.Content
	}

	// Уровень доступа - если изменения вносит не администратор, то отдать ошибку
	if newPost.AccessLevel != 0 && !admin {
		return exterr.New("Service: access denied").
			SetErrCode(http.StatusForbidden).
			SetAltMsg("Only anministrator allowed to change post access level")
	}

	// Запись в базу данных
	errExt = s.repo.UpdatePost(newPost)
	if errExt != nil {
		return exterr.NewWithExtErr("Service: EditPost error", errExt).
			SetErrCode(http.StatusInternalServerError).
			SetAltMsg(http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (s *GetService) DeletePost(id int) exterr.ErrExtender {
	errExt := s.repo.DeletePost(id)
	if errExt != nil {
		return exterr.NewWithExtErr("Service: DeletePost error", errExt).
			SetErrCode(http.StatusInternalServerError).
			SetAltMsg(http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func (s *GetService) LikePost(userID, postID int) exterr.ErrExtender {
	// Проверка поставлен ли ранее лайк пользователем
	errExt := s.repo.GetCheckLikePost(userID, postID)
	if errExt == nil {
		return exterr.New("error for checkLikePost").
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("User liked has already put a like")
	}
	// Получение количества лайков под постом
	amountLike, errExt := s.repo.GetAmountLikePost(postID)
	if errExt != nil {
		return exterr.NewWithExtErr("get error amount likes post", errExt).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("get error amount likes post")
	}

	// Запись лайка под пост
	errExt = s.repo.LikePost(userID, postID, amountLike)
	if errExt != nil {
		return exterr.NewWithExtErr("get error for like post", errExt).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("get error for like post")
	}

	return nil
}

func (s *GetService) UnlikePost(userID, postID int) exterr.ErrExtender {
	// Проверка поставлен ли ранее лайк пользователем
	errExt := s.repo.GetCheckLikePost(userID, postID)
	if errExt != nil {
		return exterr.NewWithExtErr("error for checkLikePost", errExt).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("the user did not like the post")
	}

	// Получение количества лайков под постом
	amountLike, err := s.repo.GetAmountLikePost(postID)
	if err != nil {
		return exterr.NewWithExtErr("get error amount likes post", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("get error amount likes post")
	}

	// Запись лайка под пост
	err = s.repo.UnlikePost(userID, postID, amountLike)
	if err != nil {
		return exterr.NewWithExtErr("get error for like post", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("get error for like post")
	}

	return nil
}

func (s *GetService) AddFavoritPost(userID, postID int) exterr.ErrExtender {
	// Проверка не добавлен ли пост в избарнное пользователем
	errExt := s.repo.GetCheckFavoritesPost(userID, postID)
	if errExt == nil {
		return exterr.New("error for FavoritesLikePost").
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("user favorit has already put a like")
	}
	// Получение количества добавленного в избранное под постом
	amountFavorites, errExt := s.repo.GetAmountFavoritesPost(postID)
	if errExt != nil {
		return exterr.NewWithExtErr("get error amount favorites post", errExt).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("get error amount favorites post")
	}

	// Запись добавления в избранное под пост
	errExt = s.repo.FavoritesPost(userID, postID, amountFavorites)
	if errExt != nil {
		return exterr.NewWithExtErr("get error for favorit post", errExt).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("get error for favorit post")
	}

	return nil
}

func (s *GetService) UnfavoritesPost(userID, postID int) exterr.ErrExtender {
	// Проверка поставлен ли ранее лайк пользователем
	errExt := s.repo.GetCheckFavoritesPost(userID, postID)
	if errExt != nil {
		return exterr.NewWithExtErr("error for checkFavoritesPost", errExt).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("the user did not Favorites the post")
	}

	// Получение количества лайков под постом
	amountFavorites, err := s.repo.GetAmountFavoritesPost(postID)
	if err != nil {
		return exterr.NewWithExtErr("get error amount Favorites post", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("get error amount Favorites post")
	}

	// Запись лайка под пост
	err = s.repo.UnfavoritesPost(userID, postID, amountFavorites)
	if err != nil {
		return exterr.NewWithExtErr("get error for Favorites post", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("get error for Favorites post")
	}

	return nil
}
