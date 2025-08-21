package dto

type RoleRequest struct {
	Name string `json:"name" binding:"required" jsonschema:"required,description=Role Name"`
}

type RoleData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoleSuccessResponse struct {
	Response[RoleData]
}