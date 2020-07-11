package utils

import (
	"strings"
)

func IsTagValid(tag string) bool {
	var Tags = []string{"lifeduration", "terminationdate"}
	cleanedTag := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(tag, " ", ""), "-", ""), "_", ""))
	for _, t := range Tags {
		if t == cleanedTag {
			return true
		}
	}
	return false
}

// TODO Compare Tag values
func IsEligibleForDestruction(value string) bool {
	return true
}
