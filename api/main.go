/*
Copyright © 2020 Mmadu Manasseh <mmadumanasseh@gmail.com>

Implements the API endpoints for interacting with `reka` via CLI and UI
*/
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func listResources(c *gin.Context) {
	name := c.Param("provider")
	c.JSON(http.StatusOK, gin.H{
		"message":  []string{"ec2", "s3"},
		"provider": name,
	})
}

func listProvider(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": []string{"aws", "gcp"},
	})
}

func main() {
	r := gin.Default()
	r.GET("/", listProvider)
	r.GET("/:provider/", listResources)
	log.Fatal(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
