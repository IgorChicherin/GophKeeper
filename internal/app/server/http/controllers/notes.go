package controllers

import (
	"github.com/IgorChicherin/gophkeeper/internal/app/server/http/middlewares"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/http/models"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getNodeURI struct {
	NodeID int `uri:"nodeID"`
}

type NotesController struct {
	UserUseCase  usecases.UserUseCase
	NotesUseCase usecases.NotesUseCase
}

func (nc NotesController) Route(api *gin.RouterGroup) {
	middleware := middlewares.AuthMiddleware(nc.UserUseCase)
	notes := api.Group("/notes").Use(middleware)
	{
		notes.POST("/create", nc.CreateNote)
		notes.GET("/:nodeID", nc.GetNote)
		notes.GET("", nc.GetNotes)
	}
}

// @BasePath /api
// login godoc
// @Summary create note
// @Schemes
// @Description create user note
// @Tags notes
// @Accept json
// @Produce json
// @Success 204
// @Success 200 {object} models.Note
// @Failure 400,401 {object} models.DefaultErrorResponse
// @Router /notes/create [post]
func (nc NotesController) CreateNote(c *gin.Context) {
	user, err := GetUser(c, nc.UserUseCase)
	if err != nil {
		controllerLog(c).WithError(err).Errorln("can't get user")
		return
	}

	var req models.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		controllerLog(c).WithError(err).Errorln("can't parse data")
		c.AbortWithStatusJSON(http.StatusBadRequest, models.DefaultErrorResponse{
			Error: err.Error(),
		})
		return
	}

	note, err := nc.NotesUseCase.CreateUserNote(user, req)
	if err != nil {
		controllerLog(c).WithError(err).Errorln("can't create note")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.DefaultErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, note)
}

// @BasePath /api
// login godoc
// @Summary get note
// @Schemes
// @Description get user note
// @Tags notes
// @Accept json
// @Produce json
// @Success 200 {object} models.Note
// @Success 204
// @Router /notes/:noteID [get]
func (nc NotesController) GetNote(c *gin.Context) {
	user, err := GetUser(c, nc.UserUseCase)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.DefaultErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var nodeURIParams getNodeURI

	err = c.BindUri(&nodeURIParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.DefaultErrorResponse{
			Error: err.Error(),
		})
		return
	}

	node, err := nc.NotesUseCase.GetNote(user, nodeURIParams.NodeID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.DefaultErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, node)
}

// @BasePath /api
// login godoc
// @Summary get all notes
// @Schemes
// @Description get user all notes
// @Tags notes
// @Accept json
// @Produce json
// @Success 200 {json} []models.Note
// @Success 204
// @Router /notes [get]
func (nc NotesController) GetNotes(c *gin.Context) {
	user, err := GetUser(c, nc.UserUseCase)
	if err != nil {
		controllerLog(c).WithError(err).Errorln("can't get user")
		return
	}

	notes, err := nc.NotesUseCase.GetUserNotes(user)
	if err != nil {
		controllerLog(c).WithError(err).Errorln("get notes error")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.DefaultErrorResponse{
			Error: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, notes)
}
