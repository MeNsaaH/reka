package controllers

import (
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/provider"
	"github.com/mensaah/reka/resource"
)

var (
	tmpl      *template.Template
	providers map[string]*provider.Provider
)

// GetTemplatesDir get path to templates file
func GetTemplatesDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(workingDir, "web/templates")
}

//LoadTemplates loads templates from views directory
func LoadTemplates(p []*provider.Provider) {
	// Load Providers
	providers = make(map[string]*provider.Provider)
	for _, provider := range p {
		providers[provider.Name] = provider
	}

	tmpl = template.New("").Funcs(template.FuncMap{
		"providerEnabled":      ProviderEnabled,
		"stringIn":             StringIn,
		"getProviderResources": GetProviderResources,
		"styleClass":           StyleClass,
	})

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".tmpl") {
			var err error
			tmpl, err = tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if err := filepath.Walk(GetTemplatesDir(), fn); err != nil {
		panic(err)
	}
}

//GetTemplates returns preloaded templates
func GetTemplates() *template.Template {
	return tmpl
}

// StringIn returns whether a string is in an array
func StringIn(s string, arr []string) bool {
	return true
}

// ProviderEnabled returns whether a provider is enabled or not
func ProviderEnabled(s string) bool {
	return StringIn(s, config.GetProviders())
}

// GetProviderResources returns the resource names supported by the provider
func GetProviderResources(provider string) []string {
	return providers[provider].GetResourceNames()
}

// StyleClass : The css class to represent the state with
func StyleClass(s resource.Status) string {
	if s == resource.Running {
		return "success"
	} else if s == resource.Pending || s == resource.ShuttingDown || s == resource.Stopping {
		return "info"
	} else if s == resource.Stopped {
		return "warning"
	} else {
		return "danger"
	}
}
