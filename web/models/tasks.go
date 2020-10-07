package models

import "gorm.io/gorm"

//Task model stores details about all tasks executed by Reka.
type Task struct {
	gorm.Model
	// Status of Execution: SUCCESS, RUNNING, ERROR
	Status string
	// Logfile stores path to logs of the execution
	LogFile string
}
