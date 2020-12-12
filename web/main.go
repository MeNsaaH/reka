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
	"fmt"
	"net/http"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider"
	"github.com/mensaah/reka/provider/aws"
	"github.com/mensaah/reka/resource"
	"github.com/mensaah/reka/rules"
	"github.com/mensaah/reka/web/controllers"
	"github.com/mensaah/reka/web/models"
)

var (
	providers []*provider.Provider
	scheduler *gocron.Scheduler
)

func main() {
	// Parse Command Line Arguments
	pflag.String("config", "", "Path to config file")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// Load Config and Defaults
	config.LoadConfig()
	models.SetDB(config.GetDB())

	cfg := config.GetConfig()

	for _, rule := range cfg.Rules {
		// Convert Rule in config to rules.Rule type
		r := *((*rules.Rule)(unsafe.Pointer(&rule)))
		r.Tags = resource.Tags(rule.Tags)
		rules.ParseRule(r)
	}

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
	router.GET("/provider/:provider", controllers.HomeGet)
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
func initProviders() []*provider.Provider {
	var providers []*provider.Provider
	for _, p := range config.GetProviders() {
		var (
			provider *provider.Provider
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
	_, err := scheduler.Every(uint64(frequency)).Hours().StartImmediately().Do(refreshResources, providers)
	// _, err := scheduler.Every(uint64(frequency)).Hour().Do(refreshResources, providers)
	if err != nil {
		log.Errorf("error creating job: %v", err)
	}
}

// Refresh current status of resources from Providers
func refreshResources(providers []*provider.Provider) {
	for _, provider := range providers {
		allResources := provider.GetAllResources()

		for _, resources := range allResources {
			if len(resources) > 0 {
				models.CreateOrUpdateResources(resources)
			}
		}
		stoppableResources := provider.GetStoppableResources(allResources)
		fmt.Println("Stoppable Resources: ", stoppableResources)
		errs := provider.StopResources(stoppableResources)
		fmt.Println("Errors Stopping Resources: ", errs)

		resumableResources := provider.GetResumableResources(allResources)
		fmt.Println("Resumable Resources: ", resumableResources)
		errs = provider.ResumeResources(resumableResources)
		fmt.Println("Errors Resuming Resources: ", errs)

		destroyableResources := provider.GetDestroyableResources(allResources)
		fmt.Println("Destroyable Resources: ", destroyableResources)
		errs = provider.DestroyResources(destroyableResources)
		fmt.Println("Errors Destroying Resources: ", errs)
	}
}
