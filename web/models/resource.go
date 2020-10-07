package models

import (
	"gorm.io/gorm"

	"github.com/mensaah/reka/types"
)

//Resource model to store details about Resource
type Resource struct {
	gorm.Model
	types.Resource
}

//BeforeSave gorm hook
func (res *Resource) BeforeSave(db *gorm.DB) (err error) {
	return
}
