package router

import (
	"lms-bootcamp/internal/delivery/handler"
	"lms-bootcamp/internal/di"
	"lms-bootcamp/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SessionRouter struct {
	router *gin.RouterGroup
	handler  *handler.SessionHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewSessionRouter(r *gin.RouterGroup, db *gorm.DB) *SessionRouter {
	di := di.NewSessionDI(db).Init()
	authMiddleware := middleware.NewAuthMiddleware()
	return &SessionRouter{router: r, handler: di, authMiddleware: authMiddleware}
}

func (r *SessionRouter) GetRouter() *gin.RouterGroup {
	sessionGroup := r.router.Group("/session")
	sessionGroup.POST("/", r.handler.CreateSession)
	sessionGroup.PUT("/:id", r.handler.UpdateSession)
	sessionGroup.DELETE("/:id", r.handler.DeleteSession)
	sessionGroup.GET("/:id", r.handler.GetSessionByID)
	sessionGroup.GET("/", r.handler.GetSessions)
	// sessionGroup.GET("/", r.handler.GetAllSessions)
	// sessionGroup.GET("/:id", r.handler.GetSessionByID)
	// sessionGroup.PUT("/:id", r.handler.UpdateSession)
	// sessionGroup.DELETE("/:id", r.handler.DeleteSession)

	
	

	// r.router.Use(r.authMiddleware.ValidateRole(dto.RoleAdmin, dto.RoleMentor))
	// {
	// 	// r.router.GET("/course", r.handler.GetCourses)
	// 	// r.router.GET("/course/:id", r.handler.GetCourseByID)
	// 	// r.router.PUT("/course/:id", r.handler.UpdateCourse)
	// 	// r.router.DELETE("/course/:id", r.handler.DeleteCourse)
	// }

	return sessionGroup
}