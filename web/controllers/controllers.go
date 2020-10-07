package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/mensaah/reka/config"
)

//DefaultH returns common to all pages template data
func DefaultH(c *gin.Context) gin.H {
	return gin.H{
		"Providers": config.GetProviders(),
		"Context":   c,
	}
}
