package rules

import (
	_ "embed"
	"fmt"
	"os"
	"vul-parser/internal/domain/models"

	"gopkg.in/yaml.v3"
)

//go:embed builtin.yaml
var builtinRulesData []byte

type yamlRule struct {
	ID             string          `yaml:"id"`
	Name           string          `yaml:"name"`
	Severity       string          `yaml:"severity"`
	Description    string          `yaml:"description"`
	Recommendation string          `yaml:"recommendation"`
	Conditions     []yamlCondition `yaml:"conditions"`
}

type yamlCondition struct {
	Path              string      `yaml:"path"`
	Operator          string      `yaml:"operator"`
	Value             interface{} `yaml:"value"`
	AndValueNotEmpty  bool        `yaml:"and_value_not_empty"`
	ExcludeValueRegex string      `yaml:"exclude_value_regex"`
}

type yamlRuleSet struct {
	Rules []yamlRule `yaml:"rules"`
}

func LoadRules(customPath string) ([]models.Rule, error) {
	var data []byte
	var err error

	if customPath != "" {
		data, err = os.ReadFile(customPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read custom rules: %w", err)
		}
	} else {
		data = builtinRulesData
	}

	var ruleSet yamlRuleSet
	if err := yaml.Unmarshal(data, &ruleSet); err != nil {
		return nil, fmt.Errorf("failed to parse rules: %w", err)
	}

	rules := make([]models.Rule, 0, len(ruleSet.Rules))
	for _, yr := range ruleSet.Rules {
		conditions := make([]models.Condition, 0, len(yr.Conditions))
		for _, yc := range yr.Conditions {
			conditions = append(conditions, models.Condition{
				Path:              yc.Path,
				Operator:          models.ConditionOperator(yc.Operator),
				Value:             yc.Value,
				AndValueNotEmpty:  yc.AndValueNotEmpty,
				ExcludeValueRegex: yc.ExcludeValueRegex,
			})
		}

		rules = append(rules, models.Rule{
			ID:             yr.ID,
			Name:           yr.Name,
			Severity:       models.Severity(yr.Severity),
			Description:    yr.Description,
			Recommendation: yr.Recommendation,
			Conditions:     conditions,
		})
	}

	return rules, nil
}
