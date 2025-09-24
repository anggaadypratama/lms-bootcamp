package router

import (
	"lms-bootcamp/internal/delivery/handler"
	"lms-bootcamp/internal/di"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssignmentRoute struct {
	router *gin.RouterGroup
	handler  *handler.AssignmentHandler
}

func NewAssignmentRouter(r *gin.RouterGroup, db *gorm.DB) *AssignmentRoute {
	di := di.NewAssignmentDI(db).Init()
	return &AssignmentRoute{router: r, handler: di}
}

func (r *AssignmentRoute) GetRouter() *gin.RouterGroup {
	r.router.POST("/create", r.handler.CreateAssignment)
	r.router.GET("/list", r.handler.ListAssignments)
	r.router.GET("/detail/:id", r.handler.GetAssignmentDetail)
	r.router.PUT("/update/:id", r.handler.UpdateAssignment)
	r.router.DELETE("/delete/:id", r.handler.DeleteAssignment)

	return r.router
}