package router

import (
	docs "github.com/IgorChicherin/gophkeeper/api"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/http/controllers"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/http/middlewares"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/repositories"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/usecases"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/authlib"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/crypto/crypto509"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	conn *pgx.Conn,
	authService authlib.AuthService,
	publicKey []byte,
	decrypter crypto509.Decrypter,
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
	auth := controllers.AuthController{UserUseCase: userUseCase, PublicCert: string(publicKey)}

	notesRepo := repositories.NewNotesRepository(conn)
	notesUseCase := usecases.NewNotesUseCase(notesRepo, decrypter)
	notes := controllers.NotesController{NotesUseCase: notesUseCase, UserUseCase: userUseCase}

	api := router.Group("/api")
	{
		auth.Route(api)
		notes.Route(api)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
