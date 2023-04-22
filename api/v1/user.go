package v1

import (
	"net/http"
	"strconv"

	"github.com/burxondv/note-template/api/models"
	"github.com/burxondv/note-template/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Router /users [post]
// @Summary Create a user
// @Description Create a user
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User"
// @Success 201 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		req models.CreateUserRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.User().Create(&repo.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		ImageURL:    req.ImageURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseUserModel(resp))
}

// @Router /users/{id} [get]
// @Summary Get user by id
// @Description Get user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.User().Get(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseUserModel(resp))
}

// @Router /users [get]
// @Summary Get all users
// @Description Get all users
// @Tags user
// @Accept json
// @Produce json
// @Param filter query models.GetAllUserParams false "Filter"
// @Success 200 {object} models.GetAllUsersResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	req, err := validateGetAllUserParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.User().GetAll(&repo.GetAllUsersParams{
		Page:   req.Page,
		Limit:  req.Limit,
		Search: req.Search,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, getUsersResponse(result))
}

// @Router /users/{id} [put]
// @Summary Update a user
// @Description Update a user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param user body models.CreateUserRequest true "User"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var req repo.User

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

	updated, err := h.storage.User().Update(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseUserModel(updated))
}

// @Router /users/{id} [delete]
// @Summary Delete a user
// @Description Delete a user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = h.storage.User().Delete(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted",
	})
}

func getUsersResponse(data *repo.GetAllUsersResult) *models.GetAllUsersResponse {
	response := models.GetAllUsersResponse{
		Users: make([]*models.User, 0),
		Count: int(data.Count),
	}

	for _, user := range data.Users {
		u := parseUserModel(user)
		response.Users = append(response.Users, &u)
	}

	return &response
}

func validateGetAllUserParams(c *gin.Context) (*models.GetAllUserParams, error) {
	var (
		limit int = 10
		page  int = 1
		err   error
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

	return &models.GetAllUserParams{
		Limit:  int32(limit),
		Page:   int32(page),
		Search: c.Query("search"),
	}, nil
}

func parseUserModel(user *repo.User) models.User {
	return models.User{
		ID:          int(user.ID),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		ImageURL:    user.ImageURL,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}
}
