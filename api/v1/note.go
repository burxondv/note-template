package v1

import (
	"net/http"
	"strconv"

	"github.com/burxondv/note-template/api/models"
	"github.com/burxondv/note-template/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Router /notes [post]
// @Summary Create a note
// @Description Create a note
// @Tags note
// @Accept json
// @Produce json
// @Param note body models.CreateNoteRequest true "Note"
// @Success 201 {object} models.Note
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateNote(c *gin.Context) {
	var (
		req models.CreateNoteRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Note().Create(&repo.Note{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseNoteModel(resp))
}

// @Router /notes/{id} [get]
// @Summary Get note by id
// @Description Get note by id
// @Tags note
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Note
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetNote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Note().Get(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseNoteModel(resp))
}

// @Router /notes [get]
// @Summary Get all notes
// @Description Get all notes
// @Tags note
// @Accept json
// @Produce json
// @Param filter query models.GetAllNotesParams false "Filter"
// @Success 200 {object} models.GetAllNotesResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllNotes(c *gin.Context) {
	req, err := validateGetAllNoteParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.Note().GetAll(&repo.GetAllNotesParams{
		Page:       req.Page,
		Limit:      req.Limit,
		Search:     req.Search,
		UserID:     req.UserID,
		SortByData: req.SortByData,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, getNoteResponse(result))
}

// @Router /notes/{id} [put]
// @Summary Update a note
// @Description Update a note
// @Tags note
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param note body models.UpdateNote true "Note"
// @Success 200 {object} models.Note
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateNote(c *gin.Context) {
	var req repo.Note

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req.ID = int64(id)

	noteData, err := h.storage.Note().Get(req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	updated, err := h.storage.Note().Update(&repo.Note{
		ID:          noteData.ID,
		UserID:      noteData.UserID,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   noteData.CreatedAt,
		UpdatedAt:   noteData.UpdatedAt,
		DeletedAt:   noteData.DeletedAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseNoteModel(updated))
}

// @Router /notes/{id} [delete]
// @Summary Delete a note
// @Description Delete a note
// @Tags note
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteNote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = h.storage.Note().Delete(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted",
	})
}

func getNoteResponse(data *repo.GetAllNotesResult) *models.GetAllNotesResponse {
	response := models.GetAllNotesResponse{
		Notes: make([]*models.Note, 0),
		Count: data.Count,
	}

	for _, note := range data.Notes {
		u := parseNoteModel(note)
		response.Notes = append(response.Notes, &u)
	}

	return &response
}

func validateGetAllNoteParams(c *gin.Context) (*models.GetAllNotesParams, error) {
	var (
		limit      int = 10
		page       int = 1
		err        error
		userID     int
		sortByDate string = "desc"
	)

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("user_id") != "" {
		userID, err = strconv.Atoi(c.Query("user_id"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("sort_by_date") != "" &&
		(c.Query("sort_by_date") == "desc" || c.Query("sort_by_date") == "asc") {
		sortByDate = c.Query("sort_by_date")
	}

	return &models.GetAllNotesParams{
		Limit:      int32(limit),
		Page:       int32(page),
		Search:     c.Query("search"),
		UserID:     int64(userID),
		SortByData: sortByDate,
	}, nil
}

func parseNoteModel(note *repo.Note) models.Note {
	return models.Note{
		ID:          note.ID,
		UserID:      note.UserID,
		Title:       note.Title,
		Description: note.Description,
		CreatedAt:   note.CreatedAt,
		UpdatedAt:   note.UpdatedAt,
		DeletedAt:   note.DeletedAt,
	}
}
