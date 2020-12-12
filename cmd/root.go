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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider"
	"github.com/mensaah/reka/provider/aws"
	"github.com/mensaah/reka/resource"
	"github.com/mensaah/reka/rules"
)

var (
	cfgFile   string
	providers []*provider.Provider
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
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Load Config and Defaults
	config.LoadConfig()
	cfg := config.GetConfig()

	for _, rule := range cfg.Rules {
		// Convert Rule in config to rules.Rule type
		r := *((*rules.Rule)(unsafe.Pointer(&rule)))
		r.Tags = resource.Tags(rule.Tags)
		rules.ParseRule(r)
	}

	// Initialize Provider objects
	providers = initProviders()

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

// Refresh current status of resources from Providers
func refreshResources(providers []*provider.Provider) {
	for _, provider := range providers {
		allResources := provider.GetAllResources()
		s, err := json.MarshalIndent(allResources, "", "\t")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(s))

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
}
