package types

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// GetLogger returns a logger for logging tasks to file
// TODO Create logger per Job with unique filename
func GetLogger(mgrName, logPath string) *log.Entry {
	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(io.MultiWriter(file, os.Stderr))
	} else {
		logger.Infof("Failed to log to file %s, using default stderr: ", err)
	}
	return logger.WithFields(log.Fields{
		"resource": mgrName,
	})
}
