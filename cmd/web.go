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
package cmd

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/web/controllers"
	"github.com/mensaah/reka/web/models"
)

var (
	scheduler *gocron.Scheduler
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		models.SetDB(config.GetDB())
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
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}

//initLogger initializes logrus logger with some defaults
func initLogger() {
	log.SetFormatter(&log.TextFormatter{})
	//logrus.SetOutput(os.Stderr)
	if gin.Mode() == gin.DebugMode {
		log.SetLevel(log.DebugLevel)
	}
}

func initCronJob(frequency int32) {
	scheduler = gocron.NewScheduler(time.UTC)

	// Periodic tasks
	// _, err := scheduler.Every(uint64(frequency)).Hours().StartImmediately().Do(refreshResources, providers)
	_, err := scheduler.Every(uint64(frequency)).Hour().Do(refreshResources, providers)
	if err != nil {
		log.Errorf("error creating job: %v", err)
	}
}
