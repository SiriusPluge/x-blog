package handlers

import (
	"X-Blog/pkg/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/alexmolinanasaev/exterr"
	"github.com/ansel1/merry"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Регистрация пользователя
func (h *Handler) SignUp(c *gin.Context) {
	requestBody := c.Request.Body
	defer requestBody.Close()

	// Get body
	buf, err := io.ReadAll(requestBody)
	if err != nil {
		respError(c,
			exterr.NewWithErr("Handler: SignUp can't get request body", err).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg(http.StatusText(http.StatusBadRequest)))
		return
	}

	// Unmarshal
	newUser := &models.User{}
	err = json.Unmarshal(buf, newUser)
	if err != nil {
		respError(c,
			exterr.NewWithErr("Handler: SignUp unmarshal error", err).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg(http.StatusText(http.StatusBadRequest)))
		return
	}

	// Вызываем сервис
	id, errExt := h.services.SignUp(newUser)
	if errExt != nil {
		respError(c, exterr.NewWithExtErr("Handler: SignUp error", errExt))
		return
	}

	respOk(c, fmt.Sprintf("Successful registration. Your id: %d", id))

}

// Аунтетификация пользователя
func (h *Handler) SignIn(c *gin.Context) {
	requestBody := c.Request.Body
	defer requestBody.Close()

	// Get body
	buf, err := io.ReadAll(requestBody)
	if err != nil {
		respError(c,
			exterr.NewWithErr("Handler: SignIn can't get request body", err).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg(http.StatusText(http.StatusBadRequest)))
		return
	}

	// Unmarshal
	user := models.User{}
	err = json.Unmarshal(buf, &user)
	if err != nil {
		respError(c,
			exterr.NewWithErr("Handler: SignIn unmarshal error", err).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg(http.StatusText(http.StatusBadRequest)))
		return
	}

	if user.Email == "" || user.Password == "" {
		respError(c,
			exterr.New("Handler: SignIn empty email or password").
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("Empty email or password"))
		return
	}

	// Вызываем сервис
	authenticatedUser, errExt := h.services.SignIn(user.Email, user.Password)
	if errExt != nil {
		respError(c, exterr.NewWithExtErr("Handler: SignIn error", errExt))
		return
	}

	// Сохраняем информацию в сессию
	session := sessions.Default(c)
	session.Set(sessionKeyAuthenticated, true)
	session.Set(sessionKeyUserID, authenticatedUser.ID)
	session.Set(sessionKeyAccessLevel, authenticatedUser.AccessLevel)
	err = session.Save()
	if err != nil {
		respError(c, exterr.NewWithErr("Handler: SignIn failed to save session", err).
			SetErrCode(http.StatusInternalServerError).
			SetAltMsg(http.StatusText(http.StatusInternalServerError)))
		return
	}

	respOk(c, "Successfully authenticated user")
}

// TODO: реализовать
// Главная страница
func (h *Handler) Home(c *gin.Context) {

}

// TODO: реализовать
// Лента статей
func (h *Handler) Posts(c *gin.Context) {

}

// Создать статью
func (h *Handler) CreatePost(c *gin.Context) {
	// Получение id автора
	author, errExt := SessionGetUserID(c)
	if errExt != nil {
		respError(c, exterr.NewWithExtErr("Handler: CreatePost get user id error", errExt))
		return
	}

	// Get body
	requestBody := c.Request.Body
	defer requestBody.Close()

	buf, err := io.ReadAll(requestBody)
	if err != nil {
		respError(c,
			exterr.NewWithErr("Handler: CreatePost can't get request body", err).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg(http.StatusText(http.StatusBadRequest)))
		return
	}

	// Unmarshal
	newPost := &models.Post{}
	err = json.Unmarshal(buf, newPost)
	if err != nil {
		respError(c,
			exterr.NewWithErr("Handler: CreatePost unmarshal error", err).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg(http.StatusText(http.StatusBadRequest)))
		return
	}

	newPost.Author = author

	// Вызов сервиса
	postID, errExt := h.services.CreatePost(newPost)
	if errExt != nil {
		respError(c, exterr.NewWithExtErr("Handler: CreatePost error", errExt))
		return
	}

	respOk(c, fmt.Sprintf("Post created - id=%d", postID))
}

// Посмотреть статью
func (h *Handler) GetPost(c *gin.Context) {
	// Проверка доступа пользователя
	errExt := Authorization(c, models.UserAccess)
	if errExt != nil {
		respError(c, errExt)
		return
	}

	// Получаем id статьи
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		respError(c, exterr.NewWithErr("Handler: GetPost convert id to int", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Post id must be an integer!"))
		return
	}

	// Обращение в базу данных
	post, errExt := h.services.GetPost(id)
	if errExt != nil {
		errExt.SetAltMsg("Post not found")
		errExt.SetErrCode(http.StatusNotFound)
		respError(c, errExt)
		return
	}

	// FIXME: выпвод ответа через respOk
	msg := fmt.Sprintf("Post show owner - id[%d]", post.ID)

	c.JSON(http.StatusOK, gin.H{"message": msg, "post_info": post})
}

// Редактировать статью
func (h *Handler) EditPost(c *gin.Context) {
	// Получаем id статьи, которую нужно изменить
	postIDstr := c.Param("id")
	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		respError(c, exterr.NewWithErr("Handler: EditPost convert id to int", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Post id must be an integer!"))
		return
	}

	// Поиск поста в базе
	authorID, errExt := h.services.GetPostAuthorID(postID)
	if errExt != nil {
		respError(c, exterr.NewWithExtErr("Handler: EditPost post not found", errExt).
			SetErrCode(http.StatusNotFound).
			SetAltMsg("Post not found"))
		return
	}

	// Проверка доступа пользователя
	errNotAdmin := Authorization(c, models.AdminAccess)
	if !IsOwner(c, authorID) && errNotAdmin != nil {
		respError(c, exterr.New("Handler: EditPost access denied").
			SetErrCode(http.StatusForbidden).
			SetAltMsg("Access denied"))
		return
	}

	// Get body
	requestBody := c.Request.Body
	defer requestBody.Close()

	buf, err := io.ReadAll(requestBody)
	if err != nil {
		respError(c,
			exterr.NewWithErr("Handler: EditPost can't get request body", err).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg(http.StatusText(http.StatusBadRequest)))
		return
	}

	// Unmarshal
	newPost := &models.Post{}
	err = json.Unmarshal(buf, newPost)
	if err != nil {
		respError(c,
			exterr.NewWithErr("Handler: EditPost unmarshal error", err).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg(http.StatusText(http.StatusBadRequest)))
		return
	}
	newPost.ID = postID

	// Это передаем в сервис, потому что у пользователя и админа разные права на редактирование
	admin := false
	if errNotAdmin == nil {
		admin = true
	}

	// Вызов сервиса
	errExt = h.services.EditPost(newPost, admin)
	if err != nil {
		respError(c, exterr.NewWithExtErr("Handler: EditPost error", errExt))
		return
	}

	respOk(c, "Post successfully edited")
}

// Удалить статью
func (h *Handler) DeletePost(c *gin.Context) {
	// Получаем id статьи, которую нужно удалить
	idString := c.Param("id")
	postID, err := strconv.Atoi(idString)
	if err != nil {
		respError(c, exterr.NewWithErr("Handler: DeletePost convert id to int", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Post id must be an integer!"))
		return
	}

	// Поиск поста в базе
	authorID, errExt := h.services.GetPostAuthorID(postID)
	if errExt != nil {
		respError(c, exterr.NewWithExtErr("Handler: EditPost post not found", errExt).
			SetErrCode(http.StatusNotFound).
			SetAltMsg("Post not found"))
		return
	}

	// Проверка доступа пользователя
	errNotAdmin := Authorization(c, models.AdminAccess)
	if !IsOwner(c, authorID) && errNotAdmin != nil {
		respError(c, exterr.New("Handler: EditPost access denied").
			SetErrCode(http.StatusForbidden).
			SetAltMsg("Access denied"))
		return
	}

	errExt = h.services.DeletePost(postID)
	if errExt != nil {
		respError(c, exterr.NewWithExtErr("Handler: EditPost error", errExt))
		return
	}

	respOk(c, "Post successfully delited")
}

// Лайкнуть статью
func (h *Handler) LikePost(c *gin.Context) {
	// Проверка доступа пользователя
	errExt := Authorization(c, models.UserAccess)
	if errExt != nil {
		respError(c, errExt)
		return
	}

	// получение ID пользователя из сессии
	userID, errExt := SessionGetUserID(c)
	if errExt != nil {
		respError(c, exterr.NewWithExtErr("Handler: LikePost error", errExt))
		return
	}

	// получение id поста
	idString := c.Param("id")
	postID, err := strconv.Atoi(idString)
	if err != nil {
		respError(c, exterr.NewWithErr("Handler: LikePost convert id to int", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Post id must be an integer!"))
		return
	}

	// Лайкается пост
	errExt = h.services.LikePost(userID, postID)
	if errExt != nil {
		respError(c, errExt)
		return
	}

	respOk(c, fmt.Sprintf("Post liked - id[%d]", postID))
}

// TODO: реализовать
// Убрать лайк
func (h *Handler) UnlikePost(c *gin.Context) {
	// Проверка доступа пользователя
	errExt := Authorization(c, models.UserAccess)
	if errExt != nil {
		respError(c, errExt)
		return
	}

	// получение ID пользователя из сессии
	userID, errExt := SessionGetUserID(c)
	if errExt != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request: ERROR authenticated user!"})
		return
	}

	// получение id поста
	idString := c.Param("id")
	postID, err := strconv.Atoi(idString)
	if err != nil {
		respError(c, exterr.NewWithErr("Handler: LikePost convert id to int", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Post id must be an integer!"))
		return
	}

	// Убираем лайк у поста
	err = h.services.UnlikePost(userID, postID)
	if err != nil {
		c.JSON(merry.HTTPCode(err), gin.H{"message": err.Error()})
		return
	}

	msg := fmt.Sprintf("Post unliked - id[%d]", postID)
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// Добавить в избранное
func (h *Handler) AddFavoritesPost(c *gin.Context) {
	// Проверка доступа пользователя
	errExt := Authorization(c, models.UserAccess)
	if errExt != nil {
		respError(c, errExt)
		return
	}

	// получение ID пользователя из сессии
	userID, errExt := SessionGetUserID(c)
	if errExt != nil {
		respError(c, exterr.NewWithExtErr("Handler: LikePost error", errExt))
		return
	}

	// получение id поста
	idString := c.Param("id")
	postID, err := strconv.Atoi(idString)
	if err != nil {
		respError(c, exterr.NewWithErr("Handler: LikePost convert id to int", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Post id must be an integer!"))
		return
	}

	// Лайкается пост
	errExt = h.services.AddFavoritPost(userID, postID)
	if errExt != nil {
		respError(c, errExt)
		return
	}

	respOk(c, fmt.Sprintf("Post add in favorites - id[%d]", postID))
}

// TODO: реализовать
// Убрать из избранного
func (h *Handler) RemoveFavoritesPost(c *gin.Context) {
	// Проверка доступа пользователя
	errExt := Authorization(c, models.UserAccess)
	if errExt != nil {
		respError(c, errExt)
		return
	}

	// получение ID пользователя из сессии
	userID, errExt := SessionGetUserID(c)
	if errExt != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request: ERROR authenticated user!"})
		return
	}

	// получение id поста
	idString := c.Param("id")
	postID, err := strconv.Atoi(idString)
	if err != nil {
		respError(c, exterr.NewWithErr("Handler: LikePost convert id to int", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("Post id must be an integer!"))
		return
	}

	// Удалить из избранного пользователя пост
	err = h.services.UnfavoritesPost(userID, postID)
	if err != nil {
		c.JSON(merry.HTTPCode(err), gin.H{"message": err.Error()})
		return
	}

	msg := fmt.Sprintf("Post unfavorites - id[%d]", postID)
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// TODO: реализовать
// Показать страницу пользователя
func (h *Handler) GetUser(c *gin.Context) {
	// Проверка доступа пользователя
	errExt := Authorization(c, models.UserAccess)
	if errExt != nil {
		respError(c, errExt)
		return
	}

	// получение id пользователя из url
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		respError(c, exterr.NewWithErr("Handler: GetUser convert id to int", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("User id must be an integer!"))
		return
	}

	// Обращение к базе данных
	userInfo, errExt := h.services.GetUser(id)
	if errExt != nil {
		respError(c, exterr.NewWithExtErr("Handler: GetUser error", errExt))
		return
	}

	// if userID != id { // Если пользовать смотрит чужую страницу
	// }

	// TODO: подумать, как измениь RespOK, чтобы вставлять туда дополниткльную инфу
	msg := fmt.Sprintf("User show owner - id[%d]", id)

	c.JSON(http.StatusOK, gin.H{"message": msg, "user_info": userInfo})
}

// TODO: реализовать
// Редактировать страницу пользователя
func (h *Handler) EditUser(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		respError(c, exterr.NewWithErr("Handler: EditUser convert id to int", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("User id must be an integer!"))
		return
	}

	respOk(c, fmt.Sprintf("Edit User Page %d", id))
}

// Удалить пользователя
func (h *Handler) DeleteUser(c *gin.Context) {
	// Получаем id пользователя, который должен быть удален
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respError(c, exterr.NewWithErr("Handler: DeleteUser convert id to int", err).
			SetErrCode(http.StatusBadRequest).
			SetAltMsg("User id must be an integer!"))
		return
	}

	// Проверка доступа пользователя
	errNotAdmin := Authorization(c, models.AdminAccess)
	if !IsOwner(c, id) && errNotAdmin != nil {
		respError(c, exterr.New("Handler: DeleteUser access denied").
			SetErrCode(http.StatusForbidden).
			SetAltMsg("Access denied"))
		return
	}

	errExt := h.services.DeleteUser(id)
	if errExt != nil {
		respError(c, exterr.NewWithExtErr("Handler: DeleteUser error", errExt))
		return
	}

	respOk(c, fmt.Sprintf("User %d deleted", id))
}
