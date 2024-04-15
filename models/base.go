package models

import (
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

func NotInTrash(db *gorm.DB) *gorm.DB {
	return db.Where("status != ?", entities.InTrash.Value())
}