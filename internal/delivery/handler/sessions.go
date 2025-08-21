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


type SessionHandler struct {
	u usecase.SessionUseCase
	utils *utils.Utils
}

func NewSessionHandler(usecase usecase.SessionUseCase) *SessionHandler {

	return &SessionHandler{
		u: usecase,
		utils: utils.NewUtils(),
	}
}

func (h *SessionHandler) GetSessions(c *gin.Context)  {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
	defer cancel()

	var page, pageSize int
	if p := c.Query("page"); p != "" {
		pageValue, err := strconv.Atoi(p)
		if err != nil || pageValue <= 0 {
			page = 1
		} else {
			page = pageValue
		}
	} else {
		page = 1
	}

	if ps := c.Query("per_page"); ps != "" {
		pageSizeValue, err := strconv.Atoi(ps)
		if err != nil || pageSizeValue <= 0 {
			pageSize = 10
		} else {
			pageSize = pageSizeValue
		}
	} else {
		pageSize = 0
	}

	search := c.Query("search")
	
	sessions, err := h.u.GetAllSessions(ctx, search, page, pageSize)

	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Failed to retrieve sessions", 500, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("Sessions retrieved successfully", 200, &sessions, nil))
}

func (h *SessionHandler) GetSessionByID(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
	defer cancel()

	sessionID := c.Param("id")

	if sessionID == "" || len(sessionID) == 0 {
		c.JSON(400, dto.NewResponse[*string]("Invalid session ID", 400, nil, nil))
		return
	}

	session, err := h.u.GetSessionByID(c.Request.Context(), sessionID)
	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Failed to retrieve session", 500, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("Session retrieved successfully", 200, session, nil))
}

func (h *SessionHandler) CreateSession(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
	defer cancel()

	var session *dto.SessionRequest

    if err := c.ShouldBindJSON(&session); err != nil {
        var errMsg []string
        if strings.Contains(err.Error(), "validation") {
            errMsg = h.utils.ParseValidationError(err)
        } else {
            errMsg = []string{err.Error()}
        }
        c.JSON(400, dto.NewResponse("Invalid request payload", 400, nil, &errMsg))
        return
    }

	_, err := h.u.CreateSession(c.Request.Context(), session)
	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Failed to create session", 500, nil, &errMsg))
		return
	}

	c.JSON(201, dto.NewResponse[string]("Session created successfully", 201, nil, nil))
}

func (h *SessionHandler) UpdateSession(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	sessionID := c.Param("id")

	if sessionID == "" {
		c.JSON(400, dto.NewResponse[*string]("Invalid session ID", 400, nil, nil))
		return
	}

	var session *dto.SessionRequest

	if err := c.ShouldBindJSON(&session); err != nil {
        var errMsg []string
        if strings.Contains(err.Error(), "validation") {
            errMsg = h.utils.ParseValidationError(err)
        } else {
            errMsg = []string{err.Error()}
        }
        c.JSON(400, dto.NewResponse("Invalid request payload", 400, nil, &errMsg))
        return
	}

	err := h.u.UpdateSession(c.Request.Context(), sessionID, session)
	if err != nil {
		errMsg := err.Error()
		c.JSON(400, dto.NewResponse("Failed to update session", 400, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("Session updated successfully", 200, &session, nil))
}

func (h *SessionHandler) DeleteSession(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	sessionID := c.Param("id")

	if sessionID == "" && len(sessionID) == 0 {
		c.JSON(400, dto.NewResponse[*string]("Invalid session ID", 400, nil, nil))
		return
	}

	err := h.u.DeleteSession(c.Request.Context(), sessionID)
	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Failed to delete session", 500, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse[*string]("Session deleted successfully", 200, nil, nil))
}