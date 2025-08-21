package dto


type LoginRequest struct {
	Email    string `json:"email" binding:"required" jsonschema:"required,description=User email"`
	Password string `json:"password" binding:"required" jsonschema:"required,description=User password"`
}

type LoginData struct {
	Token string `json:"token"`
	User  *UserData `json:"user"`
}

type LoginResponse struct {
	Response[LoginData]
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" jsonschema:"required,description=User email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required" jsonschema:"required,description=Password reset token"`
	NewPassword string `json:"new_password" binding:"required,min=6" jsonschema:"required,description=New password"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword" jsonschema:"required,description=Confirm password"`
}
