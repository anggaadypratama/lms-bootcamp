package handler

import (
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/pkg/utils"
	"lms-bootcamp/internal/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	u *usecase.AuthUseCase
	utils *utils.Utils
}

func NewAuthHandler(usecase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		u: usecase,
		utils: utils.NewUtils(),
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
        var errMsg []string
        if strings.Contains(err.Error(), "validation") {
            errMsg = h.utils.ParseValidationError(err)
        } else {
            errMsg = []string{err.Error()}
        }
        c.JSON(400, dto.NewResponse("Invalid request payload", 400, nil, &errMsg))
        return
	}

	loginData, err := h.u.Login(c.Request.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Login failed", 500, nil, &errMsg))
		return
	}

	if loginData == nil {
		errMsg := "Invalid credentials"
		c.JSON(401, dto.NewResponse("Login failed", 401, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse("Login successful", 200, loginData, nil))
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var forgotPasswordRequest dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&forgotPasswordRequest); err != nil {
		var errMsg []string
		if strings.Contains(err.Error(), "validation") {
			errMsg = h.utils.ParseValidationError(err)
		} else {
			errMsg = []string{err.Error()}
		}
		c.JSON(400, dto.NewResponse("Invalid request payload", 400, nil, &errMsg))
		return
	}

	if err := h.u.ForgotPassword(c.Request.Context(), forgotPasswordRequest.Email); err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Forgot password failed", 500, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse[string]("Forgot password email sent", 200, nil, nil))
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var resetPasswordRequest dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&resetPasswordRequest); err != nil {
        var errMsg []string
        if strings.Contains(err.Error(), "validation") {
            errMsg = h.utils.ParseValidationError(err)
        } else {
            errMsg = []string{err.Error()}
        }
        c.JSON(400, dto.NewResponse("Invalid request payload", 400, nil, &errMsg))
        return
	}

	if err := h.u.ResetPassword(c.Request.Context(), resetPasswordRequest.Token, resetPasswordRequest.NewPassword, resetPasswordRequest.ConfirmPassword); err != nil {
		errMsg := err.Error()
		c.JSON(500, dto.NewResponse("Reset password failed", 500, nil, &errMsg))
		return
	}

	c.JSON(200, dto.NewResponse[string]("Password reset successfully", 200, nil, nil))
}