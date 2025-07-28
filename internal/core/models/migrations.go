package models

import "gorm.io/gorm"

type Migrations struct {
	gorm.Model
	Name string `gorm:"uniqueIndex" json:"name"`
}
