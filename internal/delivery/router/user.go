package router

import (
	"lms-bootcamp/internal/delivery/handler"
	"lms-bootcamp/internal/di"
	"lms-bootcamp/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRouter struct {
	router *gin.RouterGroup
	handler  *handler.UserHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewUserRouter(r *gin.RouterGroup, db *gorm.DB) *UserRouter {
	di := di.NewUserDI(db).Init()
	authMiddleware := middleware.NewAuthMiddleware()
	return &UserRouter{router: r, handler: di, authMiddleware: authMiddleware}
}

func (r *UserRouter) GetRouter() *gin.RouterGroup {
	// r.router.Use(r.authMiddleware.ValidateRole(dto.RoleAdmin, dto.RoleMentor))
	// {
		r.router.GET("/user", r.handler.GetUsers)
		r.router.GET("/user/:id", r.handler.GetUserByID)
		r.router.PUT("/user/:id", r.handler.UpdateUser)
		r.router.POST("/user", r.handler.CreateUser)
		r.router.DELETE("/user/:id", r.handler.DeleteUser)
	// }

	return r.router
}