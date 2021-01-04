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
	"fmt"
	"os"
	"unsafe"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider"
	"github.com/mensaah/reka/provider/aws"
	"github.com/mensaah/reka/provider/gcp"
	"github.com/mensaah/reka/resource"
	"github.com/mensaah/reka/rules"
	"github.com/mensaah/reka/state"
)

var (
	cfgFile      string
	cfg          *config.Config
	providers    []*provider.Provider
	backend      state.Backender
	currentState *state.State
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "reka",
	Short: "Run Reka using config file (default $HOME/.reka.yml)",
	Long:  `A Cloud Infrastructure Management Tool to stop, resume, clean and destroy resources based on tags.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		refreshResources(providers)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.reka.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".reka" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".reka")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}

	// Load Config and Defaults
	config.LoadConfig()
	cfg = config.GetConfig()

	for _, rule := range cfg.Rules {
		// Convert Rule in config to rules.Rule type
		r := *((*rules.Rule)(unsafe.Pointer(&rule)))
		r.Tags = resource.Tags(rule.Tags)
		rules.ParseRule(r)
	}

	// Initialize Provider objects
	providers = initProviders()
	backend = state.InitBackend()
}

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
		case gcp.GetName():
			provider, err = gcp.NewProvider()
		}
		if err != nil {
			log.Fatalf("Could not initialize %s Provider: %s", p, err.Error())
		}
		// TODO Config providers
		providers = append(providers, provider)
	}
	return providers
}

// Refresh current status of resources from Providers
// TODO Reconcile state so that new resources are added to desired states and former resources removed
func refreshResources(providers []*provider.Provider) {
	// currentState is the state stored in backend
	currentState = backend.GetState()

	currentState.Current = make(state.ProvidersState)
	for _, provider := range providers {
		allResources := provider.GetAllResources()
		currentState.Current[provider.Name] = allResources
	}

	// Add new resources to desired state if they don't already exists
	for k, v := range currentState.Current {
		if m, ok := currentState.Desired[k]; ok || currentState.Desired[k] == nil {
			log.Error(m)
			// TODO Return difference between two Resources object
			continue
		}
		currentState.Desired[k] = v
	}

	backend.WriteState(currentState)
}

func reapResources() {
	// stoppableResources := provider.GetStoppableResources(allResources)
	// fmt.Println("Stoppable Resources: ", stoppableResources)
	// errs := provider.StopResources(stoppableResources)
	// fmt.Println("Errors Stopping Resources: ", errs)

	// resumableResources := provider.GetResumableResources(allResources)
	// fmt.Println("Resumable Resources: ", resumableResources)
	// errs = provider.ResumeResources(resumableResources)
	// fmt.Println("Errors Resuming Resources: ", errs)

	// destroyableResources := provider.GetDestroyableResources(allResources)
	// fmt.Println("Destroyable Resources: ", destroyableResources)
	// errs = provider.DestroyResources(destroyableResources)
	// fmt.Println("Errors Destroying Resources: ", errs)
}
