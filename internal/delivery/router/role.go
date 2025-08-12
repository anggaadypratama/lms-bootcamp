package router

import (
	"lms-bootcamp/internal/delivery/handler"
	"lms-bootcamp/internal/di"
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleRouter struct {
	router *gin.RouterGroup
	handler  *handler.RoleHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewRoleRouter(r *gin.RouterGroup, db *gorm.DB) *RoleRouter {
	di := di.NewRoleDI(db).Init()
	authMiddleware := middleware.NewAuthMiddleware()
	return &RoleRouter{router: r, handler: di, authMiddleware: authMiddleware}
}

func (r *RoleRouter) GetRouter() *gin.RouterGroup {
	r.router.Use(r.authMiddleware.ValidateRole(dto.RoleAdmin))
	{
		r.router.GET("/role", r.handler.GetRoles)
		r.router.GET("/role/:id", r.handler.GetRoleByID)
		r.router.PUT("/role/:id", r.handler.UpdateRole)
		r.router.POST("/role", r.handler.CreateRole)
		r.router.DELETE("/role/:id", r.handler.DeleteRole)
	}

	return r.router
}