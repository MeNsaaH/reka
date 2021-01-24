package config

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// GetLogger returns a logger for logging manager processes
// Logs are also persisted to file
func GetLogger(mgrName, logPath string) *log.Entry {
	logger := log.New()
	if config.Verbose {
		logger.SetLevel(log.DebugLevel)
	}
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(io.MultiWriter(file, os.Stderr))
	} else {
		logger.Infof("Failed to log to file %s, using default stderr: ", err)
	}
	return logger.WithFields(log.Fields{
		"resource": mgrName,
		"provider": mgrName,
	})
}
