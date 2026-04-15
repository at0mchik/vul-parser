package models

type ConditionOperator string

const (
	OpEq        ConditionOperator = "eq"
	OpContains  ConditionOperator = "contains"
	OpRegexKey  ConditionOperator = "regex_key"
	OpRegexVal  ConditionOperator = "regex_value"
)

type Condition struct {
    Path             string
    Operator         ConditionOperator
    Value            interface{}
    AndValueNotEmpty bool
    ExcludeValueRegex string
}

type Rule struct {
	ID             string
	Name           string
	Severity       Severity
	Description    string
	Recommendation string
	Conditions     []Condition
}