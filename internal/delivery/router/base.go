package router

import (
	"lms-bootcamp/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BaseRouter struct {
	router *gin.Engine
	db *gorm.DB
	authMiddleware *middleware.AuthMiddleware
}

func NewBaseRouter(r *gin.Engine, db *gorm.DB) *BaseRouter {
	authMiddleware := middleware.NewAuthMiddleware()
	return &BaseRouter{router: r, db: db, authMiddleware: authMiddleware}
}

func (r *BaseRouter) GetRouter() *gin.Engine {
	v1Public := r.router.Group("/api/v1")
	v1Protected := r.router.Group("/api/v1/")
	r.router.Use(gin.Recovery())

	v1Protected.Use(r.authMiddleware.Authenticate())
	{
		NewRoleRouter(v1Protected, r.db).GetRouter()
		NewUserRouter(v1Protected, r.db).GetRouter()
	}

	NewAuthRouter(v1Public, r.db).GetRouter()

	return r.router
}