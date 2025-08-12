package middleware

import (
	"fmt"
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)



type AuthMiddleware struct {
	secretKey string
	utils *utils.Utils
}

func NewAuthMiddleware() *AuthMiddleware {
	cfg := utils.NewUtils()

	return &AuthMiddleware{
		secretKey: cfg.JwtSecret,
		utils: cfg,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

        tokenParts := strings.Split(token, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            errMsg := "Invalid authorization header format"
            c.JSON(http.StatusUnauthorized, dto.NewResponse("Unauthorized", 401, nil, &errMsg))
            c.Abort()
            return
        }

        tokenString := tokenParts[1]
		res, err := m.utils.ValidateJWTToken(tokenString)

		if err != nil {
			errMsg := err.Error()
			c.JSON(401, dto.NewResponse("Unauthorized", 401, nil, &errMsg))
			c.Abort()
			return
		}

		responseData, ok := res["data"].(map[string]interface{})
		if !ok {
			errMsg := "Invalid token data structure"
			c.JSON(401, dto.NewResponse("Unauthorized", 401, nil, &errMsg))
			c.Abort()
			return
		}

		fmt.Println("Authenticated user:", responseData)

		roleData, ok := responseData["Role"].(map[string]interface{})
		if !ok {
			errMsg := "Invalid role data structure"
			c.JSON(401, dto.NewResponse("Unauthorized", 401, nil, &errMsg))
			c.Abort()
			return
		}

		c.Set("email", responseData["email"])
		c.Set("user_id", responseData["id"])
		c.Set("role", roleData)
		c.Next()
	}
}

func (m *AuthMiddleware) ValidateRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
		userRoleVal, exists := c.Get("role")
		if !exists {
			errMsg := "User role not found"
			c.JSON(http.StatusForbidden, dto.NewResponse("Forbidden", 403, nil, &errMsg))
			c.Abort()
			return
		}

		userRole, ok := userRoleVal.(map[string]interface{})
		if !ok {
			errMsg := "Invalid role data structure"
			c.JSON(http.StatusForbidden, dto.NewResponse("Forbidden", 403, nil, &errMsg))
			c.Abort()
			return
		}

		roleName, ok := userRole["name"].(string)
		if !ok {
			errMsg := "Role name not found"
			c.JSON(http.StatusForbidden, dto.NewResponse("Forbidden", 403, nil, &errMsg))
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if roleName == allowedRole {
				c.Next()
				return
			}
		}

        errMsg := "Insufficient permissions"
        c.JSON(http.StatusForbidden, dto.NewResponse("Forbidden", 403, nil, &errMsg))
        c.Abort()
    }
}