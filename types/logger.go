package types

import (
	"io"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
)

// GetLogger returns a logger for logging tasks to file
// TODO Create logger per Job with unique filename
func GetLogger(mgrName string) *log.Entry {
	logger := log.New()
	file, err := os.OpenFile(path.Join(config.GetConfig().LogPath, "logrus.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(io.MultiWriter(file, os.Stderr))
	} else {
		logger.Infof("Failed to log to file %s, using default stderr: ", err)
	}
	return logger.WithFields(log.Fields{
		"resource": mgrName,
	})
}
