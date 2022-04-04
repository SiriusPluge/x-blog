package handlers

import (
	"X-Blog/internal/service"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

var sessionLifeTime = int((10 * time.Hour).Seconds()) // Время жизни сессии в секундах

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Setup the cookie store for session management
	store := sessions.NewCookieStore([]byte("secret"))       // Передаем ключи аутентификации и шифрования
	store.Options(sessions.Options{MaxAge: sessionLifeTime}) // Установим время жизни сессии
	router.Use(sessions.Sessions("mysession", store))

	api := router.Group("/api")
	{
		api.POST("/sign-up", h.SignUp) // Регистрация пользователя
		api.POST("/sign-in", h.SignIn) // Аунтетификация пользователя

		api.POST("/whoami", h.Whoami) // Получение публчиных данных о пользователе

		api.GET("/home", h.Home)   // Главная страница
		api.GET("/posts", h.Posts) // Лента статей

		postGroup := api.Group("/post")
		{
			postGroup.POST("/create", h.CreatePost) // Создать статью
			postGroup.GET("/:id", h.GetPost)        // Посмотреть статью
			postGroup.PUT("/:id", h.EditPost)       // Редактировать статью
			postGroup.DELETE("/:id", h.DeletePost)  // Удалить статью

			postGroup.PUT("/like/add/:id", h.LikePost)      // Лайкнуть статью
			postGroup.PUT("/like/remove/:id", h.UnlikePost) // Убрать лайк

			postGroup.PUT("/favorites/add/:id", h.AddFavoritesPost)       // Добавить в избранное
			postGroup.PUT("/favorites/remove/:id", h.RemoveFavoritesPost) // Убрать из избранного
		}

		userGroup := api.Group("/user")
		{
			userGroup.GET("/:id", h.GetUser)       // Показать страницу пользователя
			userGroup.PUT("/:id", h.EditUser)      // Редактировать страницу пользователя
			userGroup.DELETE("/:id", h.DeleteUser) // Удалить пользователя

			userGroup.POST("/wallet/add", h.AddWallet) // Добавить кошелёк
			userGroup.POST("/token/buy", h.BuyToken)

			userGroup.POST("/gratitude", h.Gratitude) // Отправить благодарность автору статьи (криптовалютой)
		}

	}

	return router
}
