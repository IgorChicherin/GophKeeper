package router

import (
	docs "github.com/IgorChicherin/gophkeeper/api"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/http/controllers"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/http/middlewares"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/http/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/usecases"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/authlib"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	conn *pgx.Conn,
	authService authlib.AuthService,
) *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = false

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.Use(middlewares.LoggerMiddleware())

	docs.SwaggerInfo.BasePath = "/api/"

	userRepo := repositories.NewUserRepository(conn, authService)

	userUseCase := usecases.NewUserUseCase(authService, userRepo)
	auth := controllers.AuthController{UserUseCase: userUseCase}

	api := router.Group("/api")
	{
		auth.Route(api)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
