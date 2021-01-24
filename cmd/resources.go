/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"log"
	"os"
	"text/template"

	"github.com/mensaah/reka/provider/aws"
	"github.com/mensaah/reka/provider/gcp"
	"github.com/mensaah/reka/provider/types"
	"github.com/spf13/cobra"
)

var (
	listResource bool
	outputFormat string
)
var markdownTemplate = `
# Supported Resources
{{- range .Providers }}
## {{ .Name }}
| Resource | Destroyable| Stoppable|
| ---------|:----------:| --------:|
{{- range .Managers}}
| {{ .Name }}      | true | {{ .IsStoppable }} |
{{- end}}
{{- end }}
`

// resourcesCmd represents the resources command
var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "Shows a list of resources currently supported by reka",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		type pr struct {
			Providers []*types.Provider
		}
		p := pr{Providers: []*types.Provider{}}
		awsProvider, _ := aws.NewProvider()
		gcpProvider, _ := gcp.NewProvider()

		p.Providers = append(p.Providers, awsProvider, gcpProvider)

		// Create a new template and parse the letter into it.
		t := template.Must(template.New("markdown").Parse(markdownTemplate))

		// Execute the template for each recipient.
		err := t.Execute(os.Stdout, p)
		if err != nil {
			log.Println("executing template:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(resourcesCmd)
	resourcesCmd.Flags().StringVarP(&outputFormat, "output", "o", "markdown", "Set output type (default Markdown)")
}

func formatProviderDetails(p *types.Provider) {
}
