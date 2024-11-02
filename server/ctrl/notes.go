package ctrl

import (
	"errors"
	"net/http"

	"quicknotes/internal/files"
	"quicknotes/internal/httputil"
	"quicknotes/model"

	"github.com/gin-gonic/gin"
)

type NotesController struct{}

var fileHandler = files.FileHandler{}

func (c NotesController) Register(r *gin.RouterGroup) {
	group := r.Group("/notes")

	group.GET("/", c.list)
	group.POST("/", c.create)
	group.GET("/:path", c.getNote)
	group.PUT("/:path", c.updateNote)
	group.GET("/tags", c.getTags)
}

// list godoc
//
//	@Summary		List notes
//	@Description	get notes
//	@Tags			notes
//	@Produce		json
//	@Param			query		query	model.ListNotesInput	false	"Query"
//	@Success		200			{array}	model.Note
//	@Router			/notes	[get]
func (c NotesController) list(ctx *gin.Context) {
	query := model.ListNotesInput{}
	err := ctx.BindQuery(&query)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	notes, err := fileHandler.SearchNotes(query.Term, query.Tag, query.Sort, query.Order, query.Limit)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, notes)
}

// getNote godoc
//
//	@Summary		Get note
//	@Description	get note
//	@Tags			notes
//	@Produce		json
//	@Param			path		path	string	true	"File path"
//	@Success		200			{array}	model.Note
//	@Router			/notes	[get]
func (c NotesController) getNote(ctx *gin.Context) {
	path := ctx.Param("path")
	note, err := fileHandler.GetNote(path)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, note)
}

// updateNote godoc
//
//	@Summary		Update note
//	@Description	update note
//	@Tags			notes
//	@Produce		json
//	@Param			path		path	string				true	"File path"
//	@Param			note		body	model.AddNoteInput	true	"Add note"
//	@Success		200
//	@Router			/notes	[get]
func (c NotesController) updateNote(ctx *gin.Context) {
	path := ctx.Param("path")
	if path == "" {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("path parameter is required"))
		return
	}
	input := model.PatchNoteInput{}
	err := ctx.BindQuery(&input)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	err = fileHandler.UpdateNote(path, input)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// create godoc
//
//	@Summary		Create note
//	@Description	create note
//	@Tags			notes
//	@Produce		json
//	@Param			note		body		model.AddNoteInput	true	"Add note"
//	@Success		200			{array}		model.Note
//	@Failure		400			{object}	httputil.HTTPError
//	@Router			/notes		[get]
func (c NotesController) create(ctx *gin.Context) {
	query := model.AddNoteInput{}
	err := ctx.BindQuery(&query)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	err = fileHandler.CreateNote(query)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

// getTags godoc
//
//	@Summary		Get note tags
//	@Description	get note tags
//	@Tags			notes
//	@Produce		json
//	@Success		200			{array}	string
//	@Router			/notes/tags	[get]
func (c NotesController) getTags(ctx *gin.Context) {
	tags, err := fileHandler.GetTags()
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, tags)
}
