package controllers

import (
	"errors"
	"fmt"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/http/models"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/usecases"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func controllerLog(c *gin.Context) *log.Entry {
	entryRaw, ok := c.Get("logger")
	if !ok {
		return log.NewEntry(log.StandardLogger())
	}

	entry, ok := entryRaw.(*log.Entry)
	if !ok {
		return log.NewEntry(log.StandardLogger())
	}

	return entry
}

func GetUser(c *gin.Context, userUseCase usecases.UserUseCase) (models.User, error) {
	token := c.GetHeader("Authorization")

	if token == "" {
		controllerLog(c).Errorln("unauthorized")
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized user"})
		return models.User{}, errors.New("unauthorized")
	}

	user, err := userUseCase.GetUser(token)

	if err != nil {
		controllerLog(c).WithError(err).Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return models.User{}, errors.New(fmt.Sprintf("can't decode token: %s", err))
	}

	return user, nil
}
