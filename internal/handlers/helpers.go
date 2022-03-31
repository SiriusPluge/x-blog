package handlers

import (
	"X-Blog/pkg/models"
	"net/http"

	"github.com/alexmolinanasaev/exterr"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Session
const (
	sessionKeyUserID        = "user_id"
	sessionKeyAccessLevel   = "access_level"
	sessionKeyAuthenticated = "authenticated"
)

func Authorization(c *gin.Context, requiredAccess int) exterr.ErrExtender {
	clientAccess := SessionGetAccessLevel(c)

	if clientAccess < requiredAccess {
		return exterr.New("Handler: Authorization access denied").
			SetErrCode(http.StatusForbidden).
			SetAltMsg("Access denied")
	}
	return nil
}

// Сравнивает id пользователя из сессии c передаваемым аргументом ownerUserID
func IsOwner(c *gin.Context, ownerUserID int) bool {
	id, err := SessionGetUserID(c)
	if err != nil {
		return false
	}
	return id == ownerUserID
}

func SessionAuthenticated(c *gin.Context) bool {
	session := sessions.Default(c)

	authenticated := session.Get(sessionKeyAuthenticated)
	if authenticated == nil {
		return false
	}

	return authenticated.(bool)
}

func SessionGetUserID(c *gin.Context) (int, exterr.ErrExtender) {
	if !SessionAuthenticated(c) {
		return -1, exterr.New("Session: SessionGetUserID user unauthorized").
			SetErrCode(http.StatusUnauthorized).
			SetAltMsg(http.StatusText(http.StatusUnauthorized))
	}
	session := sessions.Default(c)
	userID := session.Get(sessionKeyUserID)
	if userID == nil {
		return -1, exterr.New("Session: SessionGetUserID user id is nil").
			SetErrCode(http.StatusUnauthorized).
			SetAltMsg(http.StatusText(http.StatusUnauthorized))
	}
	return userID.(int), nil
}

func SessionGetAccessLevel(c *gin.Context) int {
	if !SessionAuthenticated(c) {
		return models.GuestAccess
	}

	session := sessions.Default(c)

	level := session.Get(sessionKeyAccessLevel)
	if level == nil {
		return models.GuestAccess
	}

	return level.(int)
}

func respError(c *gin.Context, err exterr.ErrExtender) {
	logrus.Error(err)
	logrus.Trace(err.TraceRawString())
	c.JSON(err.GetErrCode(), gin.H{"message": err.GetAltMsg()})
}

func respOk(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{"message": msg})
}
