package provider

import (
	"strings"
	"time"
)

// ResourceTags : A tag objects
type ResourceTags map[string]string

const (
	rekaNamespace     = "reka"
	terminationDate   = "termination-date"
	terminationPolicy = "termination-policy"
	activeDuration    = "active-duration"
	include           = "include"
)

// Tags
// using `reka` namespace
// - termination-date: The date when all resources should be reaped e.g 10-10-10 7pm; 10h (relative to starting time)
// - termination-policy: stop, destroy (Whether to just stop the resources or to destroy them)
// - active-duration : Timeframe within which the Resource should be active e.g 7am-8pm,
// - include: Use to explicitly add a resource to be tracked by reka
// - exclude-all: Reka starts excluding all resources and they must be manually included to using `reka-include`
// - include-all: Reka starts tracking all the resources

func getDurationToDestroyTimeFromTag(key, value string) time.Duration {

	switch key {
	case terminationDate:
		v := strings.TrimSpace(strings.ToLower(value))
		// If it can be in format "300m", "1.5h30m"
		duration, err := time.ParseDuration(strings.ReplaceAll(v, " ", ""))
		if err == nil {
			return duration
		}
		// For time in format YYYY-mm-dd
		_, err = time.Parse("2006-01-02", v)
		if err == nil {
			return duration
		}
	}
	return time.Duration(0)
}

// Checks if the tags to initiate destruction are valid at the time of execution
func destroyTagValid(tag ResourceTags) bool {
	destroyTagSet := false
	for k, _ := range tag {
		k := strings.TrimSpace(strings.ToLower(k))
		if !strings.HasPrefix(k, rekaNamespace) {
			continue
		}
	}
	return destroyTagSet
}
