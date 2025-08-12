package dto

type UserRequest struct {
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Password   string `json:"password" binding:""`
	PhoneNumber string `json:"phone_number" binding:"required"`
	RoleId     string `json:"role_id" binding:"required"`
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