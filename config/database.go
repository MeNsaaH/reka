package config

import (
	"fmt"
	"path"
)

// DatabaseConfig Config for Dabatabase
type DatabaseConfig struct {
	Type     string
	Name     string
	Host     string
	User     string
	Password string
}

// GetConnectionString the connection string  for database
func (db *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", db.Host, db.User, db.Password, db.Name)
}

// SqliteDefaultPath the default database path to use for sqlite
func (db *DatabaseConfig) SqliteDefaultPath() string {
	return path.Join(workingDir, "reka.db")
}
