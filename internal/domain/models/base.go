package models

import (
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	ID string `gorm:"primaryKey;type:varchar(36);default:(uuid_generate_v4()::text)"`
}
