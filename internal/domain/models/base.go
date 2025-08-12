package models

import "gorm.io/gorm"

type BaseModel struct {
	gorm.Model
   	ID       string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
}