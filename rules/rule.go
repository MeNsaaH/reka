package rules

import (
	"unsafe"

	log "github.com/sirupsen/logrus"

	"github.com/mensaah/reka/config"
	"github.com/mensaah/reka/resource"
)

// Ruler interface for all rules
type Ruler interface {
	validate() error // Checks if the parameters passed are valid for the rule
	CheckResource(*resource.Resource) Action
	String() string
}

// Action : The Action to be applied to the resources matching a particular rule
type Action int

const (
	DoNothing Action = iota
	Stop
	Resume
	Destroy
)

var (
	rules        map[string]Ruler
	excludeRules []*config.ExcludeRule
)

func init() {
	rules = make(map[string]Ruler)
}

// LoadRules loads exclude rules into memory. It is abstracted as a different function because at the time
// of calling init, config values have not been loaded yet
func LoadRules() error {
	var err error
	for _, rule := range config.GetConfig().Rules {
		// Convert Rule in config to rules.Rule type
		r := *((*Rule)(unsafe.Pointer(&rule)))
		r.Tags = resource.Tags(rule.Tags)
		err = ParseRule(r)
		if err != nil {
			return err
		}
	}

	excludeRules = config.GetConfig().Exclude
	return nil
}

// Rule that can be defined for resources
type Rule struct {
	*config.Rule
}

// All checks for resource exclusion are to be done here.
func (r Rule) shouldExcludeResource(res *resource.Resource) bool {
	// Check if resource is included in the rule resource block. Resources not included are to be excluded
	if !hasResourceUri(r.Resources, res) {
		return true
	}

	for _, exRule := range excludeRules {
		if hasResourceUri(exRule.Resources, res) && hasTags(res, exRule.Tags) {
			return true
		}
	}
	return false
}

// ParseRule Get the rule to use for a particular condition
func ParseRule(rule Rule) error {
	r := []string{rule.Condition.ActiveDuration.StartTime, rule.Condition.TerminationDate, rule.Condition.TerminationPolicy}
	if hasMultipleConditions(r) {
		log.Fatalf("Multiple Conditions specified for rule: `%s`", rule.Name)
	}
	if _, ok := rules[rule.Name]; ok {
		log.Fatalf("Rule with name `%s` already exists", rule.Name)
	}

	var activeRule Ruler
	if rule.Condition.ActiveDuration.StartTime != "" || rule.Condition.ActiveDuration.StopTime != "" {
		activeRule = &ActiveDurationRule{Rule: &rule}
	} else if rule.Condition.TerminationDate != "" {
		activeRule = &TerminationDateRule{Rule: &rule}
	} else if rule.Condition.TerminationPolicy != "" {
		activeRule = &TerminationPolicyRule{Rule: &rule}
	} else {
		log.Fatalf("No Conditions specified for rule: `%s`", rule.Name)
	}

	err := activeRule.validate()
	if err != nil {
		return err
	}
	rules[rule.Name] = activeRule
	return nil
}

// Checks if multiple conditions are specified for the same rule
func hasMultipleConditions(rules []string) bool {
	count := 0
	for _, r := range rules {
		if r != "" {
			count++
		}
	}
	if count > 1 {
		return true
	}
	return false
}

// GetRules : Return an array of created Rules
func GetRules() []Ruler {
	var rulers []Ruler
	for _, m := range rules {
		rulers = append(rulers, m)
	}
	return rulers
}

// GetResourceAction get action to be performed on a resource
func GetResourceAction(res *resource.Resource) Action {
	// Returns the first Matching Rule Action for a resource
	for _, r := range rules {
		if action := r.CheckResource(res); action != DoNothing {
			return action
		}
	}
	return DoNothing
}

func hasTags(res *resource.Resource, tags resource.Tags) bool {
	for k, v := range tags {
		if res.Tags[k] != v {
			return false
		}
	}
	return true
}

func hasResourceUri(uris []string, res *resource.Resource) bool {
	if len(uris) > 0 {
		for _, v := range uris {
			if v == res.Uri() {
				return true
			}
		}
		return false
	}
	return true
}
