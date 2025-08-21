package models


type RoleModel struct {
	BaseModel
	Name string `json:"name" gorm:"uniqueIndex;not null;size:100"`
	Users []UserModel `json:"users,omitempty" gorm:"foreignKey:RoleId"`
}

func (RoleModel) TableName() string {
	return "roles"
}
