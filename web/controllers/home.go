package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//HomeGet handles GET / route
func HomeGet(c *gin.Context) {
	h := DefaultH(c)
	c.HTML(http.StatusOK, "index", h)
}
