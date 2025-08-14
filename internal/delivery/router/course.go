package router

import (
	"lms-bootcamp/internal/delivery/handler"
	"lms-bootcamp/internal/di"
	"lms-bootcamp/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CourseRouter struct {
	router *gin.RouterGroup
	handler  *handler.CourseHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewCourseRouter(r *gin.RouterGroup, db *gorm.DB) *CourseRouter {
	di := di.NewCourseDI(db).Init()
	authMiddleware := middleware.NewAuthMiddleware()
	return &CourseRouter{router: r, handler: di, authMiddleware: authMiddleware}
}

func (r *CourseRouter) GetRouter() *gin.RouterGroup {
	courseGroup := r.router.Group("/course")
	courseGroup.POST("/", r.handler.CreateCourse)
	courseGroup.GET("/", r.handler.GetAllCourses)
	courseGroup.GET("/:id_or_slug", r.handler.GetCourseByIDOrSlug)
	courseGroup.DELETE("/:id", r.handler.DeleteCourse)
	courseGroup.PUT("/:id", r.handler.UpdateCourse)
	courseGroup.GET("/user/:course_id", r.handler.GetUser)

	// r.router.Use(r.authMiddleware.ValidateRole(dto.RoleAdmin, dto.RoleMentor))
	// {
	// 	// r.router.GET("/course", r.handler.GetCourses)
	// 	// r.router.GET("/course/:id", r.handler.GetCourseByID)
	// 	// r.router.PUT("/course/:id", r.handler.UpdateCourse)
	// 	// r.router.DELETE("/course/:id", r.handler.DeleteCourse)
	// }

	return r.router
}