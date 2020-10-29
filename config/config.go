package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	log "github.com/sirupsen/logrus"
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

// AwsConfig Related Configurations
type AwsConfig struct {
	// AWS Configs
	Config          aws.Config
	AccessKey       string `yaml:"accessKey"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	DefaultRegion   string `yaml:"defaultRegion"`
}

// DatabaseConfig Config for Dabatabase
type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// GetConnectionString the connection string  for database
func (db *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", db.Host, db.User, db.Password, db.Name)
}

// SqliteDefaultPath the default database path to use for sqlite
func (db *DatabaseConfig) SqliteDefaultPath() string {
	return path.Join(workingDir, "reka.db")
}

// Config : The Config values passed to application
type Config struct {
	Name            string          `yaml:"name"`
	Providers       []string        `yaml:"providers"`
	Database        *DatabaseConfig `yaml:"database"`
	Aws             *AwsConfig      `yaml:"aws"`
	RefreshInterval int32           `yaml:"refreshInterval"`
	LogPath         string          `yaml:"logPath"`

	Auth struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"auth"`

	staticPath string // Path to Static File
}

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
	// viper.SetDefault("DbType", "sqlite") // Default Database type is sqlite
	viper.SetDefault("LogPath", path.Join(workingDir, "logs"))
	viper.SetDefault("RefreshInterval", 4) // interval between running refresh and checking for resources to updates

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

	config = &Config{}

	staticPath := viper.GetString("StaticPath")
	if !path.IsAbs(staticPath) {
		config.staticPath = path.Join(workingDir, staticPath)
	}

	config.Providers = viper.GetStringSlice("Providers")
	if len(config.Providers) < 1 {
		log.Fatal("No providers specified. Reka needs atleast one provider to track")
	}

	config.Auth.Username = viper.GetString("Auth.Username")
	config.Auth.Password = viper.GetString("Auth.Password")

	config.Database = &DatabaseConfig{
		Type:     viper.GetString("Database.Type"),
		Host:     viper.GetString("Database.Host"),
		Name:     viper.GetString("Database.Name"),
		User:     viper.GetString("Database.User"),
		Password: viper.GetString("Database.Password"),
	}
	config.Aws = &AwsConfig{}
	config.RefreshInterval = viper.GetInt32("RefreshInterval")

	config.LogPath = viper.GetString("LogPath")
	if !path.IsAbs(config.LogPath) {
		config.LogPath = path.Join(workingDir, staticPath)
	}
	// Create the Logs directory if it does not exists
	if _, err := os.Stat(config.LogPath); os.IsNotExist(err) {
		err = os.Mkdir(config.LogPath, os.ModePerm)
		if err != nil {
			log.Fatal("Could not create log path: ", err)
		}
	}

	fmt.Println(config.Auth)
	return config
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
func GetAWS() *AwsConfig {
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
