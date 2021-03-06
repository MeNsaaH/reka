package config

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	config     = &Config{}
	workingDir string
	err        error
)

const (
	appName = "REKA"
)

type ExcludeRule struct {
	Name      string
	Region    string
	Tags      map[string]string
	Resources []string
}

type Rule struct {
	Name      string
	Condition struct {
		ActiveDuration struct {
			StartTime string
			StopTime  string
			StartDay  string
			StopDay   string
		}
		TerminationPolicy string
		TerminationDate   string
	}
	Resources []string
	Region    string
	Tags      map[string]string
}

func (r Rule) String() string {
	return r.Name
}

type Backend struct {
	Type   string
	Path   string
	Bucket string
	Region string
}

// Config : The Config values passed to application
type Config struct {
	Name            string
	Providers       []string
	Timezone        string
	RefreshInterval int32
	LogPath         string
	Verbose         bool

	Database *DatabaseConfig

	Web struct {
		// Authentication Details to login to Reka
		Auth struct {
			Username string
			Password string
		}
	}

	staticPath string // Path to Static File

	// Exclude block prevents certain resources from been tracked or affected by reka.
	Exclude []*ExcludeRule

	// StateBackend is how state is stored (read & write)
	// State files contain details used for infrastructure resumption and history of
	// infrastructural management
	StateBackend *Backend
	// Rules block define how reka should behave given certain resources. These rules
	// usually target resources based on tags/labels which are attached to the resources
	Rules []*Rule
	// AWS Config
	Aws *aws.Config
	// Gcp configuration
	Gcp *Gcp
}

// RemoteBackendTypes allowed remote storage
var RemoteBackendTypes = []string{"s3", "azblob", "gs"}

// LocalBackendTypes local backend storage
var LocalBackendTypes = []string{"local"}

// StateBackendTypes all state backend storage types
var StateBackendTypes = append(RemoteBackendTypes, LocalBackendTypes...)

// LoadConfig load all passed configs and defaults
func LoadConfig() *Config {
	workingDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	viper.SetEnvPrefix(appName)
	viper.AutomaticEnv() // Load Variables from Environment with REKA prefix
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// Defaults
	viper.SetDefault("StaticPath", "web/static")
	viper.SetDefault("StateBackend.Type", "local")
	viper.SetDefault("StateBackend.Path", path.Join(workingDir, "reka-state.json"))
	// viper.SetDefault("DbType", "sqlite") // Default Database type is sqlite
	viper.SetDefault("LogPath", path.Join(workingDir, "logs"))
	viper.SetDefault("RefreshInterval", 4)             // interval between running refresh and checking for resources to updates
	viper.SetDefault("aws.DefaultRegion", "us-east-2") // Default AWS Region for users https://docs.aws.amazon.com/emr/latest/ManagementGuide/emr-plan-region.html

	// Load Config file
	if configPath := viper.GetString("Config"); configPath != "" {
		dir, file := filepath.Split(configPath)
		viper.SetConfigName(file)   // name of config file (without extension)
		viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
		// If REKA_CONFIG_FILE is set load config from that
		viper.AddConfigPath(dir)
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
				log.Fatalf("error: %s. Consider passing the `--config` variable or settings %s_CONFIG environment", err, appName)

			} else {
				// Config file was found but another error was produced
				log.Fatalf("Error: %s", err)
			}
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	if !path.IsAbs(config.staticPath) {
		config.staticPath = path.Join(workingDir, config.staticPath)
	}

	if len(config.Providers) < 1 {
		log.Fatal("No providers specified. Reka needs atleast one provider to monitor")
	}

	if !Contains(StateBackendTypes, config.StateBackend.Type) {
		log.Fatalf("State Backend Type is not permitted, Please use one of %s", StateBackendTypes)

	} else {
		if Contains(RemoteBackendTypes, config.StateBackend.Type) {
			if config.StateBackend.Bucket == "" || config.StateBackend.Path == "" {
				log.Fatalf("State Backend (type: %s) - bucket and path must be set in config", config.StateBackend.Type)
			}
			if config.StateBackend.Type == "s3" && config.StateBackend.Region == "" {
				log.Fatal("State Backend (type: s3) - Must configure Region")
			}
		}
	}

	if !path.IsAbs(viper.GetString("StaticPath")) {
		config.staticPath = path.Join(workingDir, viper.GetString("StaticPath"))
	}

	awsConfig := loadAwsConfig(viper.GetString("aws.AccessKeyID"), viper.GetString("aws.SecretAccessKey"), viper.GetString("aws.DefaultRegion"))
	config.Aws = &awsConfig

	if !path.IsAbs(config.LogPath) {
		config.LogPath = path.Join(workingDir, config.LogPath)
	}
	// Create the Logs directory if it does not exists
	if _, err := os.Stat(config.LogPath); os.IsNotExist(err) {
		err = os.Mkdir(config.LogPath, os.ModePerm)
		if err != nil {
			log.Fatal("Could not create log path: ", err)
		}
	}

	return config
}

// Contains check if array contains value
func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// GetConfig return the config object
func GetConfig() *Config {
	return config
}

// GetDB Return database config
func GetDB() *DatabaseConfig {
	return config.Database
}

// GetAWS Return database config
func GetAWS() *aws.Config {
	return config.Aws
}

//StaticPath returns path to application static folder
func StaticPath() string {
	return config.staticPath
}

// GetProviders returns list of selected providers
func GetProviders() []string {
	return config.Providers
}

// SetVerboseLogging sets logging to be in DEBUG mode
func SetVerboseLogging() {
	config.Verbose = true
}
