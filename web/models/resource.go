package models

import (
	"gorm.io/gorm/clause"

	"github.com/mensaah/reka/provider"
	"github.com/mensaah/reka/resource"
)

// CreateOrUpdateResources : Creates resource if it does not exists else updates it
func CreateOrUpdateResources(resources []*resource.Resource) error {
	// Update columns to new value on `id` conflict
	var err error
	for _, r := range resources {
		err = db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "uuid"}},
			DoUpdates: clause.AssignmentColumns([]string{"region", "state"}),
		}).Create(&r).Error
	}
	return err
}

// CreateAllManagers : Creates all resource managers
func CreateAllManagers(providers []*provider.Provider) error {
	// Update columns to new value on `id` conflict

	var err error
	for _, provider := range providers {
		for _, rmgr := range provider.Managers {
			err = db.Where(resource.Manager{Name: rmgr.Name}).FirstOrCreate(&rmgr).Error
		}
	}
	return err
}
