/*
Copyright © 2020 Mmadu Manasseh <mmadumanasseh@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider/aws"
	"github.com/mensaah/reka/types"
	"github.com/mensaah/reka/web/controllers"
	"github.com/mensaah/reka/web/models"
)

var (
	providers []*types.Provider
	scheduler *gocron.Scheduler
)

func main() {

	config.LoadConfig()
	models.SetDB(config.GetDB())

	// Initialize Provider objects
	providers = initProviders()

	err := models.AutoMigrate(providers)
	if err != nil {
		log.Fatal("Database Migration Error: ", err)
	}
	initCronJob(config.GetConfig().RefreshInterval)

	// Load Templates
	controllers.LoadTemplates(providers)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.SetHTMLTemplate(controllers.GetTemplates())
	router.StaticFS("/static", http.Dir(config.StaticPath()))
	// router.Use(controllers.ContextData())

	router.GET("/", controllers.HomeGet)
	// router.NoRoute(controllers.NotFound)
	// router.NoMethod(controllers.MethodNotAllowed)

	// Start cron
	scheduler.StartAsync()

	router.Run(":8080")
}

//initLogger initializes logrus logger with some defaults
func initLogger() {
	log.SetFormatter(&log.TextFormatter{})
	//logrus.SetOutput(os.Stderr)
	if gin.Mode() == gin.DebugMode {
		log.SetLevel(log.DebugLevel)
	}
}

// TODO Add logger to Providers during configuration
func initProviders() []*types.Provider {
	var providers []*types.Provider
	for _, p := range config.GetProviders() {
		var (
			provider *types.Provider
			err      error
		)
		switch p {
		case aws.GetName():
			provider, err = aws.NewProvider()
			if err != nil {
				log.Fatal("Could not initialize AWS Provider: ", err)
			}
		}
		// TODO Config providers
		providers = append(providers, provider)
	}
	return providers
}

func initCronJob(frequency int32) {
	scheduler = gocron.NewScheduler(time.UTC)

	// Periodic tasks
	_, err := scheduler.Every(uint64(frequency)).Hour().StartImmediately().Do(refreshResources, providers)
	// _, err := scheduler.Every(uint64(frequency)).Hour().Do(refreshResources, providers)
	if err != nil {
		log.Errorf("error creating job: %v", err)
	}
}

// Refresh current status of resources from Providers
func refreshResources(providers []*types.Provider) {
	for _, provider := range providers {
		allResources := provider.GetAllResources()

		for _, resources := range allResources {
			if len(resources) > 0 {
				models.CreateOrUpdateResources(resources)
			}
		}
		// destroyableResources := provider.GetDestroyableResources(allResources)
		// save(destroyableResources)
		// stoppableResources := provider.GetStoppableResources(allResources)
		// save(stoppableResources)
		// resumableResources := provider.GetResumableResources(allResources)
		// save(resumableResources)
	}
}
