package rules

import (
	log "github.com/sirupsen/logrus"

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

var rules map[string]Ruler

func init() {
	rules = make(map[string]Ruler)
}

// Rule that can be defined for resources
type Rule struct {
	Name      string
	Condition struct {
		ActiveDuration struct {
			StartTime string
			StopTime  string
			StartDay  string
			StopDay   string
		}
		TerminationPolicy string
		TerminationDate   string
	}
	Tags   resource.Tags
	Region string
}

func (r Rule) String() string {
	return r.Name
}

// ExcludeRule : Excludes resources that matches this rules
type ExcludeRule struct {
	Name      string
	Region    string
	Tags      resource.Tags
	Resources []string

	Validate func() error // Checks if the parameters passed are valid for the rule
}

// ParseRule Get the rule to use for a particular condition
func ParseRule(rule Rule) Ruler {
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
		log.Fatal(err)
	}
	rules[rule.Name] = activeRule
	return activeRule
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
