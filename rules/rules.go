package rules

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
	"github.com/mensaah/reka/resource"
)

// TerminationDateRule defines rule that sets when a resource should be terminated based on termination-date rule
// set in config file. This rule checks if the set termination-date is past and activates destroy on the instance
type TerminationDateRule struct {
	*Rule
	Date time.Time
}

func (r *TerminationDateRule) validate() error {
	date, err := now.Parse(r.Condition.TerminationDate)
	if err != nil {
		return fmt.Errorf("Error parsing conditon.terminationDate: %s", err)
	}
	r.Date = date
	return nil
}

// CheckResource Returns a list of resources whose termination Date is exceeed
func (r TerminationDateRule) CheckResource(res *resource.Resource) Action {
	if r.shouldExcludeResource(res) {
		return DoNothing
	}
	if hasTags(res, r.Tags) {
		currentDate := time.Now()
		if int(r.Date.Sub(currentDate)) <= 0 {
			return Destroy
		}
	}
	return DoNothing
}

// ActiveDurationRule defines the period within which a resource should be active. Any time
// Outside the duration, the resource is stopped and resumed again when period is active
type ActiveDurationRule struct {
	*Rule
	StartTime time.Time
	StopTime  time.Time
}

func (r *ActiveDurationRule) validate() error {
	startTime, err := now.Parse(r.Condition.ActiveDuration.StartTime)
	if err != nil {
		return fmt.Errorf("Invalid Start Time for condition.activeDuration: %s", r.Condition.ActiveDuration.StartTime)
	}
	r.StartTime = startTime

	stopTime, err := now.Parse(r.Condition.ActiveDuration.StopTime)
	if err != nil {
		return fmt.Errorf("Invalid Stop Time for condition.activeDuration: %s", r.Condition.ActiveDuration.StopTime)
	}
	r.StopTime = stopTime
	return nil
}

// CheckResource Returns a list of resources activeDuration is valid
func (r ActiveDurationRule) CheckResource(res *resource.Resource) Action {
	if r.shouldExcludeResource(res) {
		return DoNothing
	}
	if hasTags(res, r.Tags) {
		currentTime := time.Now()
		// Initialize Stopping only if currentTime is not within active duration
		if !(currentTime.Sub(r.StartTime) > 0 && currentTime.Sub(r.StopTime) <= 0) && res.IsActive() {
			return Stop
		}
		// Initialize Resumption only if currentTime is within active duration
		if currentTime.Sub(r.StartTime) > 0 && currentTime.Sub(r.StopTime) <= 0 && res.IsStopped() {
			return Resume
		}
	}
	return DoNothing
}

// TerminationPolicyRule defines rule that sets when a resource should be terminated.
type TerminationPolicyRule struct {
	*Rule
	Policy string
}

func (r *TerminationPolicyRule) validate() error {
	validPolicies := []string{"unused"}
	for _, p := range validPolicies {
		if p == r.Condition.TerminationPolicy {
			return nil
		}
	}
	r.Policy = r.Condition.TerminationPolicy
	return fmt.Errorf("Error parsing condition.terminationPolicy: Invalid Policy %s", r.Policy)
}

// CheckResource Returns a list of resources whose termination Date is exceeed
func (r TerminationPolicyRule) CheckResource(res *resource.Resource) Action {
	if r.shouldExcludeResource(res) {
		return DoNothing
	}
	if res.IsUnused() {
		return Destroy
	}
	return DoNothing
}
