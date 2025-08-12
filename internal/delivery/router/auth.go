package router

import (
	"lms-bootcamp/internal/delivery/handler"
	"lms-bootcamp/internal/di"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthRouter struct {
	router *gin.RouterGroup
	handler  *handler.AuthHandler
}

func NewAuthRouter(r *gin.RouterGroup, db *gorm.DB) *AuthRouter {
	di := di.NewAuthDI(db).Init()
	return &AuthRouter{router: r, handler: di}
}

func (r *AuthRouter) GetRouter() *gin.RouterGroup {
	r.router.POST("/login", r.handler.Login)
	r.router.POST("/forgot-password", r.handler.ForgotPassword)
	r.router.POST("/reset-password", r.handler.ResetPassword)

	return r.router
}