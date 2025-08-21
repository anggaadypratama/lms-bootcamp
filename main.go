package main

import (
	"log"
	"net/http"

	"lms-bootcamp/config"
	"lms-bootcamp/internal/delivery/mcp"
	"lms-bootcamp/internal/delivery/router"
	"lms-bootcamp/internal/domain/models"
	"lms-bootcamp/internal/pkg/database"

	server "github.com/ckanthony/gin-mcp"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	conn, err := database.NewDatabaseConnection(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	conn.Migrate(
		&models.UserModel{},
		&models.CourseModel{},
		&models.RoleModel{},
		&models.SessionModel{},
		&models.AssignmentModel{},
	)

	router := router.NewBaseRouter(r, conn.DB)
	router.GetRouter()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello dari LMS Bootcamp! 🚀",
			"status":  "success",
		})
	})

	mcpServer := server.New(r, &server.Config{
		Name:        "LMS Bootcamp",
		Description: "An example API automatically exposed via MCP.",
		BaseURL: "http://localhost:" + cfg.Port,
	})

	mcp.NewMcpServer(mcpServer).RegisterRoutes()
	mcpServer.Mount("/mcp")

	log.Fatal(r.Run(":" + cfg.Port))
}