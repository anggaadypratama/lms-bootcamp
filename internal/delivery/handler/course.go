package handler

import (
	"context"
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	usecase *usecase.CourseUseCase
}

func NewCourseHandler(uc *usecase.CourseUseCase) *CourseHandler {
	return &CourseHandler{
		usecase: uc,
	}
}


func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var course *dto.CourseRequest
	if err := c.ShouldBindJSON(&course); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusBadRequest, dto.NewResponse("Invalid request payload", 400, nil, &errMsg))
		return
	}


	if err := h.usecase.CreateCourse(c.Request.Context(), course); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusInternalServerError, dto.NewResponse("Failed to create course", 400, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("Courses created successfully", 201, &course, nil))
}


func (h *CourseHandler) GetCourseByIDOrSlug(c *gin.Context) {
	idOrSlug := c.Param("id_or_slug")
	course, err := h.usecase.GetCourseByIDOrSlug(c.Request.Context(), idOrSlug)
	if err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusInternalServerError, dto.NewResponse("Failed to retrieve course", 400, nil, &errMsg))
		return
	}
	c.JSON(200, dto.NewResponse("Courses retrieved successfully", 200, &course, nil))
}

func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
	defer cancel()

	var filter dto.CourseFilter
	var page, pageSize int
	
	if f := c.Query("filter"); f != "" {
		filter.Title = f
	}
	

	if err := c.ShouldBindQuery(&filter); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusInternalServerError,dto.NewResponse("Invalid filter", 400, nil, &errMsg))
		return
	}

	if pageParam := c.Query("page"); pageParam != "" {
		p, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			errMsg := err.Error()
			c.JSON(http.StatusInternalServerError,dto.NewResponse("Invalid Page", 400, nil, &errMsg))
			return
		}

		page = p

		if pageSizeParam := c.Query("per_page"); pageSizeParam != "" {
			ps, err := strconv.Atoi(c.Query("per_page"))
			if err != nil {
				errMsg := err.Error()
				c.JSON(http.StatusInternalServerError,dto.NewResponse("Invalid Page Size", 400, nil, &errMsg))
				return
			}
			pageSize = ps
		}
	}

	

	courses, err := h.usecase.GetAllCourses(c.Request.Context(), &filter, &page, &pageSize)
	if err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusInternalServerError,dto.NewResponse("Failed to retrieve courses", 400, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("Courses retrieved successfully", 200, &courses, nil))
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteCourse(c.Request.Context(), id); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusInternalServerError, dto.NewResponse("Failed to delete course", 400, nil, &errMsg))
		return
	}
	c.JSON(200, dto.NewResponse[[]string]("Course deleted successfully", 200, nil, nil))
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	var course *dto.CourseRequest
	if err := c.ShouldBindJSON(&course); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusBadRequest, dto.NewResponse("Invalid request payload", 400, nil, &errMsg))
		return
	}

	id := c.Param("id")
	if err := h.usecase.UpdateCourse(c.Request.Context(), id, course); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusInternalServerError, dto.NewResponse("Failed to update course", 400, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("Course updated successfully", 200, &course, nil))
}

func (h *CourseHandler) GetUser(c *gin.Context) {
	id := c.Param("course_id")
	role := c.Query("role")

	if id == "" {
		errMsg := "Invalid ID"
		c.JSON(http.StatusBadRequest, dto.NewResponse("Failed to retrieve user", 400, nil, &errMsg))
		return
	}
	

	user, err := h.usecase.GetCoursesByUser(c.Request.Context(), id, role)
	if err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusInternalServerError, dto.NewResponse("Failed to retrieve user", 400, nil, &errMsg))
		return
	}
	c.JSON(200, dto.NewResponse("User retrieved successfully", 200, &user, nil))
}
