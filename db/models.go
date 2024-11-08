package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID    uint   `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name  string `json:"name" gorm:"type:varchar(32)"`
	Email string `json:"email" gorm:"type:varchar(100)"`
	Age   uint32 `json:"age" gorm:"type:integer"`
}
