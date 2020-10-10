package models

import (
	"gorm.io/gorm/clause"

	"github.com/mensaah/reka/types"
)

// CreateOrUpdateResources : Creates resource if it does not exists else updates it
func CreateOrUpdateResources(resources []*types.Resource) error {
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

// CreateAllResourceManagers : Creates all resource managers
func CreateAllResourceManagers(providers []*types.Provider) error {
	// Update columns to new value on `id` conflict

	var err error
	for _, provider := range providers {
		for _, rmgr := range provider.ResourceManagers {
			err = db.Where(types.ResourceManager{Name: rmgr.Name}).FirstOrCreate(&rmgr).Error
		}
	}
	return err
}
