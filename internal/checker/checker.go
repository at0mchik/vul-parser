package checker

import (
	"fmt"
	"regexp"
	"strings"
	"vul-parser/internal/domain/models"
)

type Checker struct {
	rules []models.Rule
}

func NewChecker(rules []models.Rule) *Checker {
	return &Checker{rules: rules}
}

func (c *Checker) Check(data interface{}) []models.Vulnerability {
	var vulnerabilities []models.Vulnerability
	
	for _, rule := range c.rules {
		for _, cond := range rule.Conditions {
			vulns := c.evaluateCondition(data, cond, rule)
			vulnerabilities = append(vulnerabilities, vulns...)
		}
	}
	
	return c.deduplicate(vulnerabilities)
}

func (c *Checker) evaluateCondition(data interface{}, cond models.Condition, rule models.Rule) []models.Vulnerability {
	var results []models.Vulnerability
	
	c.walk(data, "", func(path string, value interface{}) {
		if c.matchesCondition(path, value, cond) {
			results = append(results, models.Vulnerability{
				RuleID:       rule.ID,
				Severity:     rule.Severity,
				Description:  rule.Description,
				Recommendation: rule.Recommendation,
				Path:         path,
				Value:        value,
			})
		}
	})
	
	return results
}

func (c *Checker) matchesCondition(path string, value interface{}, cond models.Condition) bool {
	if cond.Path != "*" && !c.pathMatches(path, cond.Path) {
		return false
	}
	
	if cond.AndValueNotEmpty && c.isEmptyValue(value) {
		return false
	}
	
	// новое: исключить значения, подходящие под regex
	if cond.ExcludeValueRegex != "" {
		if c.regexMatch(fmt.Sprintf("%v", value), cond.ExcludeValueRegex) {
			return false
		}
	}
	
	switch cond.Operator {
	case models.OpEq:
		return c.equal(value, cond.Value)
	case models.OpContains:
		return c.contains(value, cond.Value)
	case models.OpRegexKey:
		return c.regexMatch(path, cond.Value)
	case models.OpRegexVal:
		return c.regexMatch(fmt.Sprintf("%v", value), cond.Value)
	default:
		return false
	}
}

func (c *Checker) pathMatches(path, pattern string) bool {
	if pattern == "$" {
		return path == ""
	}
	
	pattern = strings.TrimPrefix(pattern, "$.")
	if pattern == "" {
		return path == ""
	}
	
	return path == pattern || strings.HasPrefix(path, pattern+".")
}

func (c *Checker) equal(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func (c *Checker) contains(a, b interface{}) bool {
	return strings.Contains(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
}

func (c *Checker) regexMatch(s string, pattern interface{}) bool {
	re, err := regexp.Compile(fmt.Sprintf("%v", pattern))
	if err != nil {
		return false
	}
	return re.MatchString(s)
}

func (c *Checker) isEmptyValue(v interface{}) bool {
	if v == nil {
		return true
	}
	s, ok := v.(string)
	if ok && s == "" {
		return true
	}
	return false
}

func (c *Checker) walk(data interface{}, prefix string, fn func(string, interface{})) {
	if data == nil {
		return
	}
	
	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			newPath := key
			if prefix != "" {
				newPath = prefix + "." + key
			}
			fn(newPath, val)
			c.walk(val, newPath, fn)
		}
	case []interface{}:
		for i, val := range v {
			newPath := fmt.Sprintf("%s[%d]", prefix, i)
			fn(newPath, val)
			c.walk(val, newPath, fn)
		}
	default:
		// leaf value already handled by parent
	}
}

func (c *Checker) deduplicate(vulns []models.Vulnerability) []models.Vulnerability {
	seen := make(map[string]bool)
	var result []models.Vulnerability
	
	for _, v := range vulns {
		key := fmt.Sprintf("%s:%s", v.RuleID, v.Path)
		if !seen[key] {
			seen[key] = true
			result = append(result, v)
		}
	}
	
	return result
}