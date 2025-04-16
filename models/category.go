package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name          string     `json:"name"`
	Desciption    string     `json:"description"`
	ParentID      *uint      `json:"parent_id,omitempty"`
	Parent        *Category  `gorm:"foreignKey:ParentID", json:"parent,omitempty"`
	Subcategories []Category `gorm:"foreignKey:ParentID" json:"subcategories,omitempty"`
}
