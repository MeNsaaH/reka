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
	"strings"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider/aws"
	"github.com/mensaah/reka/provider/gcp"
	"github.com/mensaah/reka/provider/types"
	"github.com/mensaah/reka/resource"
	"github.com/mensaah/reka/rules"
	"github.com/mensaah/reka/state"
)

var (
	cfgFile        string
	cfg            *config.Config
	providers      []*types.Provider
	backend        state.Backender
	activeState    *state.State
	verbose        bool
	disableStop    bool
	disableResume  bool
	disableDestroy bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "reka",
	Short: "Run Reka using config file (default $HOME/.reka.yml)",
	Long:  `A Cloud Infrastructure Management Tool to stop, resume, clean and destroy resources based on tags.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Load Config and Defaults
		config.LoadConfig()

		cfg = config.GetConfig()
		err := rules.LoadRules()
		if err != nil {
			log.Fatal(err)
		}

		// Initialize Provider objects
		providers = initProviders()
		backend = state.InitBackend()
		// RefreshResources on every execution
		refreshResources(providers)
		for _, p := range providers {
			res := activeState.Current[p.Name]

			if !disableStop {
				stoppableResources := p.GetStoppableResources(res)
				fmt.Println("Stoppable Resources: ", stoppableResources)
				errs := p.StopResources(stoppableResources)
				logErrors(errs)
			}

			if !disableResume {
				resumableResources := p.GetResumableResources(res)
				fmt.Println("Resumable Resources: ", resumableResources)
				errs := p.ResumeResources(resumableResources)
				logErrors(errs)
			}

			if !disableDestroy {
				destroyableResources := p.GetDestroyableResources(res)
				fmt.Println("Destroyable Resources: ", destroyableResources)
				errs := p.DestroyResources(destroyableResources)
				logErrors(errs)
			}
		}
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
	rootCmd.Flags().BoolP("unused-only", "t", false, "Reap only unused resources")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Output verbose logs (DEBUG)")
	rootCmd.Flags().BoolVar(&disableStop, "disable-stop", false, "Disable stopping of resources")
	rootCmd.Flags().BoolVar(&disableResume, "disable-resume", false, "Disable resuming of resources")
	rootCmd.Flags().BoolVar(&disableDestroy, "disable-destroy", false, "Disable destruction of resources")
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
		log.Debug("Using config file:", viper.ConfigFileUsed())
	}

	if verbose {
		config.SetVerboseLogging()
		log.SetLevel(log.DebugLevel)
	}

}

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
		case gcp.GetName():
			provider, err = gcp.NewProvider()
		}
		if err != nil {
			log.Fatalf("Could not initialize %s Provider: %s", p, err.Error())
		}
		providers = append(providers, provider)
	}
	log.Infof("Loaded Providers: %s", strings.Join(config.GetProviders()[:], ","))
	return providers
}

// Refresh current status of resources from Providers
func refreshResources(providers []*types.Provider) {
	// activeState is the current state stored in backend
	activeState = backend.GetState()

	activeState.Current = make(state.ProvidersState)
	for _, provider := range providers {
		res := provider.GetAllResources()
		activeState.Current[provider.Name] = res
	}

	// Add new resources to desired state if they don't already exists
	// this is to ensure all new resources created are also added to reka's desired state
	for currentProvider, currentProviderResources := range activeState.Current {
		if _, ok := activeState.Desired[currentProvider]; ok {
			for u, w := range currentProviderResources {
				if _, ok := activeState.Desired[currentProvider][u]; ok {
					for _, res := range w {
						if !containsResource(activeState.Desired[currentProvider][u], res) {
							activeState.Desired[currentProvider][u] = append(activeState.Desired[currentProvider][u], res)
						}
					}
				} else {
					activeState.Desired[currentProvider][u] = activeState.Current[currentProvider][u]
				}
			}
		} else {
			activeState.Desired[currentProvider] = currentProviderResources
		}
	}

	backend.WriteState(activeState)
}

func containsResource(res []*resource.Resource, r *resource.Resource) bool {
	for _, rs := range res {
		if rs.UUID == r.UUID {
			return true
		}
	}
	return false
}

func logErrors(errs map[string]error) {
	for k, v := range errs {
		log.Errorf("%s: %s", k, v.Error())
	}
}
