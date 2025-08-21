package handler

import (
	"context"
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/pkg/utils"
	"lms-bootcamp/internal/usecase"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)


type UserHandler struct {
	u usecase.UserUseCase
	utils *utils.Utils
}

func NewUserHandler(usecase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		u: usecase,
		utils: utils.NewUtils(),
	}
}

func (h *UserHandler) GetUsers(c *gin.Context)  {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
	defer cancel()

	var page, pageSize *int
	var filter *map[string]interface{}
	var filterBy string

	if fb := c.Query("filter_by"); fb != "" {
		filterBy = fb
		
		if f := c.Query("filter"); f != "" {

			validFields := map[string]string{
				"name":  "name",
				"email": "email",
				"phone": "phone_number",
				"role":  "role_id",
			}

			if dbField, exists := validFields[filterBy]; exists {
				filter = &map[string]interface{}{dbField: f}
			}
		}
	}
	
	if p := c.Query("page"); p != "" {
		pageValue, err := strconv.Atoi(p)
		if err != nil || pageValue <= 0 {
			pageValue = 1
		}
		page = &pageValue
	}

	if ps := c.Query("per_page"); ps != "" {
		pageSizeValue, _ := strconv.Atoi(ps)
		pageSize = &pageSizeValue
	}

	users, err := h.u.GetAllUsers(ctx, filter, page, pageSize)

	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Failed to retrieve users", 500, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("Users retrieved successfully", 200, &users, nil))
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
	defer cancel()

	userID := c.Param("id")

	user, err := h.u.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Failed to retrieve user", 500, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("User retrieved successfully", 200, user, nil))
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
	defer cancel()

	var user dto.UserRequest

	if err := c.ShouldBindJSON(&user); err != nil {
        var errMsg []string
        if strings.Contains(err.Error(), "validation") {
            errMsg = h.utils.ParseValidationError(err)
        } else {
            errMsg = []string{err.Error()}
        }
        c.JSON(400, dto.NewResponse("Invalid request payload", 400, nil, &errMsg))
        return
	}

	_, err := h.u.CreateUser(c.Request.Context(), &user)
	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Failed to create user", 500, nil, &errMsg))
		return
	}

	c.JSON(201, dto.NewResponse[string]("User created successfully", 201, nil, nil))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.Param("id")

	var user dto.UserRequest

	if err := c.ShouldBindJSON(&user); err != nil {
        var errMsg []string
        if strings.Contains(err.Error(), "validation") {
            errMsg = h.utils.ParseValidationError(err)
        } else {
            errMsg = []string{err.Error()}
        }
        c.JSON(400, dto.NewResponse("Invalid request payload", 400, nil, &errMsg))
        return
	}

	err := h.u.UpdateUser(c.Request.Context(), userID, &user)
	if err != nil {
		errMsg := err.Error()
		c.JSON(400, dto.NewResponse("Failed to update user", 400, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("User updated successfully", 200, &user, nil))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.Param("id")

	err := h.u.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Failed to delete user", 500, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse[string]("User deleted successfully", 200, nil, nil))
}