package mcp

import (
	"lms-bootcamp/internal/domain/dto"

	server "github.com/ckanthony/gin-mcp"
)


type McpServer struct {
	server *server.GinMCP
}

func NewMcpServer(mcp *server.GinMCP) *McpServer {
	return &McpServer{
		server: mcp,
	}
}

func (m *McpServer) RegisterRoutes() *server.GinMCP{
	m.server.RegisterSchema("POST", "/users", nil, &dto.UserRequest{})
	m.server.RegisterSchema("POST", "/courses", nil, &dto.CourseRequest{})
	m.server.RegisterSchema("POST", "/roles", nil, &dto.RoleRequest{})
	m.server.RegisterSchema("POST", "/session", nil, &dto.SessionRequest{})

	return m.server
}
