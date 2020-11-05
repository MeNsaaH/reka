package resource

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/now"
	log "github.com/sirupsen/logrus"
)

// Tags : A tag objects
type Tags map[string]string

const (
	rekaNamespace      = "reka"
	destructionDateTag = "destruction-date"
	destructionPolicy  = "destruction-policy"
	activeDurationTag  = "active-duration"
	include            = "include"
)

// Check if Tag is a valid Reka tag i.e begins with REKA_*
func isValidRekaTag(key string) bool {
	k := strings.TrimSpace(strings.ToLower(key))
	return strings.HasPrefix(k, rekaNamespace)
}

// Check if Tag is a timer tag i.e. Tags that contain resource time control e.g destruction-date
func isTimerTag(key string) bool {
	// Timer Tags:
	// destructionDate, activeDuration
	timerTags := []string{destructionDateTag, activeDurationTag}
	for _, v := range timerTags {
		if strings.HasSuffix(strings.TrimSpace(strings.ToLower(key)), v) {
			return true
		}
	}
	return false
}

// Gets the tag with `reka-` namespace prefix
func getCleanTag(key string) string {
	cleanedTag := strings.TrimSpace(strings.ToLower(key))
	// Get substring `reka-` index after
	cleanedTag = cleanedTag[len(fmt.Sprintf("%s-", rekaNamespace)):]
	return cleanedTag
}

// Tags
// using `reka` namespace
// - destruction-date: The date when all resources should be reaped e.g `YYYY-MM-DD HH:MM`; 10h (relative to starting time)
// - destruction-policy: stop, destroy (Whether to just stop the resources or to destroy them)
// - active-duration : Timeframe within which the Resource should be active e.g 10:30-18:00
// - include: Use to explicitly add a resource to be tracked by reka
// - exclude-all: Reka starts excluding all resources and they must be manually included to using `reka-include`
// - include-all: Reka starts tracking all the resources

// ShouldInitiateDestruction checks if the tags to initiate destruction are valid at the time of execution
func ShouldInitiateDestruction(tags Tags) bool {
	// For destructionDate, check if current date is/surpassed destruction date
	for k, v := range tags {
		if isValidRekaTag(k) {
			cleanedTag := getCleanTag(k)
			if cleanedTag == destructionDateTag {
				date, err := now.Parse(v)
				if err != nil {
					log.Errorf("Invalid Tag for destruction-date: %v, %s", v, err)
					continue
				}
				currentDate := time.Now()
				if int(date.Sub(currentDate)) <= 0 {
					return true
				}
			}
		}
	}
	return false
}

func getCleanActiveDurationValue(value string) (time.Time, time.Time, error) {
	v := strings.Split(value, "-")
	if len(v) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("Invalid value for active-duration, format should be HH:MM-HH:MM: %v", v)
	}
	startTime, err := now.Parse(v[0])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("Invalid Start Time for active-duration: %v", v)

	}
	stopTime, err := now.Parse(v[1])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("Invalid Stop Time for active-duration: %v", v)
	}
	return startTime, stopTime, nil
}

// ShouldInitiateStopping checks if the tags to initiate destruction are valid at the time of execution
func ShouldInitiateStopping(tag Tags) bool {
	for k, v := range tag {
		if isValidRekaTag(k) {
			cleanedTag := getCleanTag(k)
			// For activeDuration, check if current time is within active duration
			// activeDuration specified does not destroy the instance. It justs stops the instance, if the instance
			// implements the ResourceStopResumer Interface
			if activeDurationTag == cleanedTag {
				startTime, stopTime, err := getCleanActiveDurationValue(v)
				if err != nil {
					log.Errorf(err.Error())
					continue
				}
				currentTime := time.Now()

				// Initialize Stopping only if currentTime is not within active duration
				if !(currentTime.Sub(startTime) > 0 && currentTime.Sub(stopTime) <= 0) {
					return true
				}
			}
		}
	}
	return false
}

// ShouldInitiateResumption checks if the tags to initiate destruction are valid at the time of execution
func ShouldInitiateResumption(tag Tags) bool {
	for k, v := range tag {
		if isValidRekaTag(k) {
			cleanedTag := getCleanTag(k)
			// For activeDuration, check if current time is within active duration
			// activeDuration specified does not destroy the instance. It justs stops the instance, if the instance
			// implements the ResourceStopResumer Interface
			if activeDurationTag == cleanedTag {
				startTime, stopTime, err := getCleanActiveDurationValue(v)
				if err != nil {
					log.Errorf(err.Error())
					continue
				}
				currentTime := time.Now()

				// Initialize Resumption only if currentTime is within active duration
				if currentTime.Sub(startTime) > 0 && currentTime.Sub(stopTime) <= 0 {
					return true
				}
			}
		}
	}
	return false
}
