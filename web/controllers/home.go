package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"

	"github.com/mensaah/reka/types"
	"github.com/mensaah/reka/web/models"
)

//HomeGet handles GET / route
func HomeGet(c *gin.Context) {
	var r types.Resource

	db := models.GetDB()
	h := DefaultH(c)
	// TODO Specify Provider to query from
	// provider := strings.ToLower(c.Param("provider")) || config.GetProviders()[0]
	// Filter where resource Manager is for a particular provider

	h["TrackedResourceCount"] = db.Not("state = ?", types.Destroyed).Find(&r).RowsAffected
	h["RunningResourceCount"] = db.Where("state = ?", types.Running).Find(&r).RowsAffected
	h["StoppedResourceCount"] = db.Where("state = ?", types.Stopped).Find(&r).RowsAffected
	h["DestroyedResourceCount"] = db.Where("state = ?", types.Destroyed).Find(&r).RowsAffected

	// Get Recent Resource Updates
	var recentResourceUpdates []types.Resource
	if err := db.Order("updated_at desc").Limit(10).Find(&recentResourceUpdates); err != nil {
		log.Error("DB Error: ", err.Error)
	}
	h["RecentResourceUpdates"] = recentResourceUpdates

	c.HTML(http.StatusOK, "index", h)
}
