package dto

type RoleRequest struct {
	Name string `json:"name" binding:"required"`
}

type RoleData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoleSuccessResponse struct {
	Response[RoleData]
}