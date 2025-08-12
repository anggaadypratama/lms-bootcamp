package handler

import (
	"context"
	"time"

	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)


type RoleHandler struct {
	u usecase.RoleUseCase
}

func NewRoleHandler(usecase usecase.RoleUseCase) *RoleHandler {
	return &RoleHandler{
		u: usecase,
	}
}

func (h *RoleHandler) GetRoles(c *gin.Context)  {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
	defer cancel()

	roles, err := h.u.GetAllRoles(ctx)
	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Failed to retrieve roles", 500, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("Roles retrieved successfully", 200, &roles, nil))
}

func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
	defer cancel()

	roleID := c.Param("id")

	role, err := h.u.GetRoleByID(c.Request.Context(), roleID)
	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Failed to retrieve role", 500, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("Role retrieved successfully", 200, role, nil))
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
	defer cancel()

	var role dto.RoleRequest

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := h.u.CreateRole(c.Request.Context(), role.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, dto.NewResponse[string]("Role created successfully", 201, nil, nil))
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
	defer cancel()

	roleID := c.Param("id")

	var role dto.RoleRequest

	if err := c.ShouldBindJSON(&role); err != nil {
		errorMsg := err.Error()
		c.JSON(400, dto.NewResponse("Invalid request payload", 400, nil, &errorMsg))
		return
	}

	err := h.u.UpdateRole(c.Request.Context(), roleID, role.Name)
	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Failed to update role", 500, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("Role updated successfully", 200, &role, nil))
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
	defer cancel()

	roleID := c.Param("id")

	err := h.u.DeleteRole(c.Request.Context(), roleID)
	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Failed to delete role", 500, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse[string]("Role deleted successfully", 200, nil, nil))
}