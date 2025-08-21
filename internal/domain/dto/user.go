package dto

type UserRequest struct {
	Email string `json:"email" binding:"required" jsonschema:"required,description=User email"`
	Name  string `json:"name" binding:"required" jsonschema:"required,description=User name"`
	Password   string `json:"password" binding:"" jsonschema:"required,description=User password"`
	PhoneNumber string `json:"phone_number" binding:"required" jsonschema:"required,description=User phone number"`
	RoleId     string `json:"role_id" binding:"required" jsonschema:"required,description=User role ID"`
}

type UserData struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Role   *RoleData
}

type UserSuccessResponse struct {
	Response[UserData]
}