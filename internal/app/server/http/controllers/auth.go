package controllers

import (
	"errors"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/http/models"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserUseCase usecases.UserUseCase
	PublicCert  string
}

func (ac AuthController) Route(api *gin.RouterGroup) {
	auth := api.Group("/user")
	{
		auth.POST("/login", ac.login)
		auth.POST("/register", ac.register)
	}
}

// @BasePath /api
// login godoc
// @Summary login
// @Schemes
// @Description user login
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.LoginRequest true "login"
// @Success 200 {object} models.LoginResponse
// @Failure 400,401,500
// @Router /user/login [post]
func (ac AuthController) login(c *gin.Context) {
	var userData models.LoginRequest

	if err := c.ShouldBind(&userData); err != nil {
		controllerLog(c).WithError(err).Errorln("can't parse request")
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	token, err := ac.UserUseCase.Login(userData.Login, userData.Password)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("user repository error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if errors.Is(err, usecases.ErrUserNotFound) || errors.Is(err, usecases.ErrUnauthorized) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Header("Authorization", token)
	c.AbortWithStatusJSON(http.StatusOK, models.LoginResponse{Cert: ac.PublicCert, Token: token})
}

// @BasePath /api
// login godoc
// @Summary register
// @Schemes
// @Description user registration
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.UserRequest true "user account"
// @Success 200
// @Failure 404,500
// @Failure 400,409 {object} models.DefaultErrorResponse
// @Router /user/register [post]
func (ac AuthController) register(c *gin.Context) {
	var userData models.UserRequest

	if err := c.ShouldBind(&userData); err != nil {
		controllerLog(c).WithError(err).Errorln("can't parse request")
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	token, err := ac.UserUseCase.Register(userData.Login, userData.Password)

	if errors.Is(err, usecases.ErrUserHasBeenRegistered) {
		c.AbortWithStatusJSON(http.StatusConflict, map[string]string{"error": "user with this login has been created"})
		return
	}

	if err != nil {
		controllerLog(c).WithError(err).Errorln("user repository error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Authorization", token)
	c.AbortWithStatus(http.StatusOK)
}
