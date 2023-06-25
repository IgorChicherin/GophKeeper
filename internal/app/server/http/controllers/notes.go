package controllers

import (
	"github.com/IgorChicherin/gophkeeper/internal/app/server/http/middlewares"
	"github.com/IgorChicherin/gophkeeper/internal/app/server/usecases"
	"github.com/IgorChicherin/gophkeeper/internal/shared/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getNoteURI struct {
	NoteID int `uri:"noteID"`
}

type NotesController struct {
	UserUseCase  usecases.UserUseCase
	NotesUseCase usecases.NotesUseCase
}

func (nc NotesController) Route(api *gin.RouterGroup) {
	middleware := middlewares.AuthMiddleware(nc.UserUseCase)
	notes := api.Group("/notes").Use(middleware)
	{
		notes.POST("/create", nc.createNote)
		notes.PUT("/update/:noteID", nc.updateNote)
		notes.GET("/:noteID", nc.getNote)
		notes.GET("/delete/:noteID", nc.deleteNote)
		notes.GET("", nc.getNotes)
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
// @Success 200 {object} models.CreateNoteRequest
// @Failure 400,401 {object} models.DefaultErrorResponse
// @Router /notes/create [post]
func (nc NotesController) createNote(c *gin.Context) {
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
// @Success 200 {object} models.DecodedNoteResponse
// @Success 204
// @Router /notes/:noteID [get]
func (nc NotesController) getNote(c *gin.Context) {
	user, err := GetUser(c, nc.UserUseCase)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.DefaultErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var nodeURIParams getNoteURI

	err = c.BindUri(&nodeURIParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.DefaultErrorResponse{
			Error: err.Error(),
		})
		return
	}

	node, err := nc.NotesUseCase.GetNote(user, nodeURIParams.NoteID)
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
// @Summary delete note
// @Schemes
// @Description delete user note
// @Tags notes
// @Accept json
// @Produce json
// @Success 200 204
// @Router /notes/delete/:noteID [get]
func (nc NotesController) deleteNote(c *gin.Context) {
	user, err := GetUser(c, nc.UserUseCase)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.DefaultErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var nodeURIParams getNoteURI

	err = c.BindUri(&nodeURIParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.DefaultErrorResponse{
			Error: err.Error(),
		})
		return
	}

	err = nc.NotesUseCase.DeleteUserNote(user, nodeURIParams.NoteID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.DefaultErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

// @BasePath /api
// login godoc
// @Summary get all notes
// @Schemes
// @Description get user all notes
// @Tags notes
// @Accept json
// @Produce json
// @Success 200 {json} []models.DecodedNoteResponse
// @Success 204
// @Router /notes [get]
func (nc NotesController) getNotes(c *gin.Context) {
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

// @BasePath /api
// login godoc
// @Summary update note
// @Schemes
// @Description update user note
// @Tags notes
// @Accept json
// @Produce json
// @Success 204
// @Success 200 {object} models.CreateNoteRequest
// @Failure 400,401 {object} models.DefaultErrorResponse
// @Router /notes/update/:noteID [put]
func (nc NotesController) updateNote(c *gin.Context) {
	user, err := GetUser(c, nc.UserUseCase)
	if err != nil {
		controllerLog(c).WithError(err).Errorln("can't get user")
		return
	}

	var nodeURIParams getNoteURI

	err = c.BindUri(&nodeURIParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.DefaultErrorResponse{
			Error: err.Error(),
		})
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

	note, err := nc.NotesUseCase.UpdateUserNote(user, nodeURIParams.NoteID, req)
	if err != nil {
		controllerLog(c).WithError(err).Errorln("can't create note")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.DefaultErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, note)
}
