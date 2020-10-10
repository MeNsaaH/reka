package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/viper"
)

var (
	config     *Config
	workingDir string
	err        error
)

const (
	appName = "REKA"
)

// AWSConfig Related Configurations
type AWSConfig struct {
	// AWS Configs
	Config        aws.Config
	DefaultRegion string
}

// DatabaseConfig Config for Dabatabase
type DatabaseConfig struct {
	Host     string
	Name     string //database name
	User     string
	Password string
	Type     string
}

// GetConnectionString the connection string  for database
func (db *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.db.Host, config.db.User, config.db.Password, config.db.Name)
}

// SqliteDefaultPath the default database path to use for sqlite
func (db *DatabaseConfig) SqliteDefaultPath() string {
	return path.Join(workingDir, "reka.db")
}

// Config : The Config values passed to application
type Config struct {
	// A list of supported providers to be enabled
	Providers  []string
	staticPath string

	Aws *AWSConfig
	db  *DatabaseConfig

	// Authentication details set from config
	Username string
	Password string

	// RefreshInterval in hours
	RefreshInterval int32
	// LogPath
	LogPath string
}

// LoadConfig load all passed configs and defaults
func LoadConfig() *Config {
	workingDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	viper.SetEnvPrefix(appName) // will be uppercased automatically
	viper.SetConfigName("reka") // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name

	// If REKA_CONFIG_FILE is set load config from that
	viper.AddConfigPath(viper.GetString("ConfigFile"))
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		// panic(fmt.Errorf("Fatal error config file: %s", err))
		fmt.Errorf("Fatal error config file: %s", err)
	}

	// Defaults
	viper.SetDefault("StaticPath", "web/static")
	viper.SetDefault("DbType", "sqlite") // Default Database type is sqlite
	viper.SetDefault("LogPath", path.Join(workingDir, "logs"))
	viper.SetDefault("RefreshInterval", 4) // interval between running refresh and checking for resources to updates

	staticPath := viper.GetString("StaticPath")

	config = &Config{}

	if !path.IsAbs(staticPath) {
		config.staticPath = path.Join(workingDir, staticPath)
	}

	// Load Configuration
	config.Providers = []string{"aws"} // TODO Remove Test providers init with aws

	config.db = &DatabaseConfig{
		Type:     viper.GetString("DBType"),
		Host:     viper.GetString("DbHost"),
		Name:     viper.GetString("DbName"),
		User:     viper.GetString("DbUser"),
		Password: viper.GetString("DbPassword"),
	}
	config.Aws = &AWSConfig{}
	config.RefreshInterval = viper.GetInt32("RefreshInterval")
	config.LogPath = viper.GetString("LogPath")
	if _, err := os.Stat(config.LogPath); os.IsNotExist(err) {
		err = os.Mkdir(config.LogPath, os.ModePerm)
		if err != nil {
			log.Fatal("Could not create log path: ", err)
		}
	}
	return config
}

// GetConfig return the config object
func GetConfig() *Config {
	return config
}

// GetDB Return database config
func GetDB() *DatabaseConfig {
	return config.db
}

// GetAWS Return database config
func GetAWS() *AWSConfig {
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
